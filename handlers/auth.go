package handlers

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"github.com/mahesh-yadav/go-recipes-api/config"
	"github.com/mahesh-yadav/go-recipes-api/models"
	"github.com/rs/zerolog/log"
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

// SignUpHandler godoc
//
//	@Summary		Sign up a new user
//	@Description	Create a new user account
//	@Tags			auth
//	@Accept			json
//	@Produce		json
//	@Param			user	body	models.User	true	"User Sign Up"
//	@Success		201
//	@Failure		400	{object}	models.ErrorResponse
//	@Failure		500	{object}	models.ErrorResponse
//	@Router			/auth/signup [post]
func (handler *AuthHandler) SignUpHandler(c *gin.Context) {
	var user models.User
	if err := c.ShouldBindJSON(&user); err != nil {
		c.JSON(http.StatusBadRequest, models.ErrorResponse{
			Code:    http.StatusBadRequest,
			Message: "Invalid User data",
		})
		return
	}

	h := sha256.New()
	h.Write([]byte(user.Password))
	user.Password = hex.EncodeToString(h.Sum(nil))
	result, err := handler.collection.InsertOne(handler.ctx, user)
	if err != nil {
		log.Panic().Msg("Error inserting user into MongoDB")
		return
	}

	c.JSON(http.StatusCreated, result)
}

// SignInHandler godoc
//
//	@Summary		Sign in a user
//	@Description	Authenticate a user and return a JWT token
//	@Tags			auth
//	@Accept			json
//	@Produce		json
//	@Param			user	body		models.User	true	"User Credentials"
//	@Success		200		{object}	JWTOutput
//	@Failure		400		{object}	models.ErrorResponse
//	@Failure		401		{object}	models.ErrorResponse
//	@Failure		500		{object}	models.ErrorResponse
//	@Router			/auth/signin [post]
func (handler *AuthHandler) SignInHandler(c *gin.Context) {
	var user models.User
	if err := c.ShouldBindJSON(&user); err != nil {
		c.JSON(http.StatusBadRequest, models.ErrorResponse{
			Code:    http.StatusBadRequest,
			Message: "Invalid User data",
		})
		return
	}

	h := sha256.New()
	h.Write([]byte(user.Password))
	hashedPassword := hex.EncodeToString(h.Sum(nil))
	filter := bson.M{"username": user.Username, "password": hashedPassword}

	cursor := handler.collection.FindOne(handler.ctx, filter)
	if cursor.Err() != nil {
		c.JSON(http.StatusUnauthorized, models.ErrorResponse{
			Code:    http.StatusUnauthorized,
			Message: "invalid credentials",
		})
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
		log.Panic().Msg("Error creating JWT token")
		return
	}
	jwtOutput := JWTOutput{
		Token:   tokenString,
		Expires: expirationTime,
	}

	c.JSON(http.StatusOK, jwtOutput)
}

// RefreshTokenHandler godoc
//
//	@Summary		Refresh JWT token
//	@Description	Refresh an existing JWT token and return a new one
//	@Tags			auth
//	@Accept			json
//	@Produce		json
//	@Param			Authorization	header		string	true	"{token}"
//	@Success		200				{object}	JWTOutput
//	@Failure		401				{object}	models.ErrorResponse
//	@Failure		500				{object}	models.ErrorResponse
//	@Router			/auth/refresh [post]
func (handler *AuthHandler) RefreshTokenHandler(c *gin.Context) {
	tokenStr := c.GetHeader("Authorization")
	claims := &Claims{}

	token, err := jwt.ParseWithClaims(tokenStr, claims, func(token *jwt.Token) (interface{}, error) {
		return []byte(handler.config.JWTSecret), nil
	})

	if err != nil {
		c.JSON(http.StatusUnauthorized, models.ErrorResponse{
			Code:    http.StatusUnauthorized,
			Message: "invalid token",
		})
		return
	}

	if token == nil || !token.Valid {
		c.JSON(http.StatusUnauthorized, models.ErrorResponse{
			Code:    http.StatusUnauthorized,
			Message: "invalid token",
		})
		return
	}

	expirationTime := time.Now().Add(time.Duration(handler.config.JWTExpirationTimeSeconds) * time.Second)
	claims.ExpiresAt = jwt.NewNumericDate(expirationTime)
	newToken := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	newTokenString, err := newToken.SignedString([]byte(handler.config.JWTSecret))
	if err != nil {
		log.Panic().Msg("Error creating JWT token")
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
