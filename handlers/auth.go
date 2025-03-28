package handlers

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"github.com/mahesh-yadav/go-recipes-api/config"
	"github.com/mahesh-yadav/go-recipes-api/models"
	"github.com/mahesh-yadav/go-recipes-api/utils"
	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
)

type AuthHandler struct {
	ctx        context.Context
	config     *config.Config
	collection *mongo.Collection
}

func NewAuthHandler(ctx context.Context, config *config.Config, collection *mongo.Collection) *AuthHandler {
	return &AuthHandler{
		ctx:        ctx,
		config:     config,
		collection: collection,
	}
}

type Claims struct {
	Username string `json:"username"`
	jwt.RegisteredClaims
}

type JWTOutput struct {
	Token   string    `json:"token"`
	Expires time.Time `json:"expires"`
}

func (handler *AuthHandler) SignUpHandler(c *gin.Context) {
	var user models.User
	if err := c.ShouldBindJSON(&user); err != nil {
		utils.NewError(c, http.StatusBadRequest, err)
		return
	}

	h := sha256.New()
	h.Write([]byte(user.Password))
	user.Password = hex.EncodeToString(h.Sum(nil))
	result, err := handler.collection.InsertOne(handler.ctx, user)
	if err != nil {
		utils.NewError(c, http.StatusInternalServerError, err)
		return
	}

	c.JSON(http.StatusCreated, result)
}

func (handler *AuthHandler) SignInHandler(c *gin.Context) {
	var user models.User
	if err := c.ShouldBindJSON(&user); err != nil {
		utils.NewError(c, http.StatusBadRequest, err)
		return
	}

	h := sha256.New()
	h.Write([]byte(user.Password))
	hashedPassword := hex.EncodeToString(h.Sum(nil))
	filter := bson.M{"username": user.Username, "password": hashedPassword}

	cursor := handler.collection.FindOne(handler.ctx, filter)
	if cursor.Err() != nil {
		utils.NewError(c, http.StatusUnauthorized, errors.New("invalid credentials"))
		return
	}

	expirationTime := time.Now().Add(time.Duration(handler.config.JWTExpirationTimeSeconds) * time.Second)
	claims := &Claims{
		user.Username,
		jwt.RegisteredClaims{
			Issuer:    "recipes-api",
			ExpiresAt: jwt.NewNumericDate(expirationTime),
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString([]byte(handler.config.JWTSecret))
	if err != nil {
		utils.NewError(c, http.StatusInternalServerError, err)
		return
	}
	jwtOutput := JWTOutput{
		Token:   tokenString,
		Expires: expirationTime,
	}

	c.JSON(http.StatusOK, jwtOutput)
}

func (handler *AuthHandler) RefreshTokenHandler(c *gin.Context) {
	tokenStr := c.GetHeader("Authorization")
	claims := &Claims{}

	token, err := jwt.ParseWithClaims(tokenStr, claims, func(token *jwt.Token) (interface{}, error) {
		return []byte(handler.config.JWTSecret), nil
	})

	if err != nil {
		c.AbortWithStatus(http.StatusUnauthorized)
	}

	if token == nil || !token.Valid {
		c.AbortWithStatus(http.StatusUnauthorized)
	}

	expirationTime := time.Now().Add(time.Duration(handler.config.JWTExpirationTimeSeconds) * time.Second)
	claims.ExpiresAt = jwt.NewNumericDate(expirationTime)
	newToken := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	newTokenString, err := newToken.SignedString([]byte(handler.config.JWTSecret))
	if err != nil {
		utils.NewError(c, http.StatusInternalServerError, err)
		return
	}
	jwtOutput := JWTOutput{
		Token:   newTokenString,
		Expires: expirationTime,
	}
	c.JSON(http.StatusOK, jwtOutput)
}

func (handler *AuthHandler) AuthMiddlewareJWT() gin.HandlerFunc {
	return func(c *gin.Context) {
		tokenStr := c.GetHeader("Authorization")

		claims := &Claims{}

		token, err := jwt.ParseWithClaims(tokenStr, claims, func(token *jwt.Token) (interface{}, error) {
			return []byte(handler.config.JWTSecret), nil
		})

		if err != nil {
			c.AbortWithStatus(http.StatusUnauthorized)
		}

		if token == nil || !token.Valid {
			c.AbortWithStatus(http.StatusUnauthorized)
		}

		c.Next()
	}
}
