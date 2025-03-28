package main

import (
	"context"

	"github.com/gin-gonic/gin"
	"github.com/mahesh-yadav/go-recipes-api/config"
	"github.com/mahesh-yadav/go-recipes-api/database"
	_ "github.com/mahesh-yadav/go-recipes-api/docs"
	"github.com/mahesh-yadav/go-recipes-api/handlers"
	"github.com/mahesh-yadav/go-recipes-api/logger"
	"github.com/mahesh-yadav/go-recipes-api/middleware"
	"github.com/rs/zerolog/log"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

// @Title			Recipes API
// @Version		1.0
// @Description	This is a simple API for managing recipes.
//
// @Host			localhost:8080
// @BasePath		/
func main() {
	config := config.GetConfig()

	log.Logger = logger.SetupLogger(config)

	gin.SetMode(config.GinMode)

	database.ConnectToMongoDB(config)

	ctx := context.Background()
	recipeCollection := database.GetMongoCollection(config, "recipes")
	if config.InitializeDB {
		database.InitDB(recipeCollection)
	}
	database.ConnectToRedis(config)
	redisClient := database.GetRedisClient(config)

	recipesHandler := handlers.NewRecipeHandler(ctx, recipeCollection, redisClient, config)

	userCollection := database.GetMongoCollection(config, "users")
	authHandler := handlers.NewAuthHandler(ctx, config, userCollection)

	router := gin.New()
	router.Use(gin.Logger(), middleware.GlobalErrorMiddleware())

	router.GET("/recipes", recipesHandler.ListRecipesHandler)
	router.POST("/auth/signup", authHandler.SignUpHandler)
	router.POST("/auth/signin", authHandler.SignInHandler)
	router.POST("/auth/refresh", authHandler.AuthMiddlewareJWT(), authHandler.RefreshTokenHandler)

	authorized := router.Group("/")
	authorized.Use(authHandler.AuthMiddlewareJWT())
	{
		authorized.POST("/recipes", recipesHandler.CreateRecipeHandler)
		authorized.GET("/recipes/:id", recipesHandler.GetRecipeHandler)
		authorized.PUT("/recipes/:id", recipesHandler.UpdateRecipeHandler)
		authorized.DELETE("/recipes/:id", recipesHandler.DeleteRecipeHandler)
		authorized.GET("/recipes/search", recipesHandler.SearchRecipeHandler)
	}

	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	router.Run()
}
