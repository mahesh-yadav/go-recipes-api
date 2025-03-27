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
	collection := database.GetMongoCollection(config, "recipes")

	if config.InitializeDB {
		database.InitDB(collection)
	}

	database.ConnectToRedis(config)
	redisClient := database.GetRedisClient(config)

	recipesHandler := handlers.NewRecipeHandler(ctx, collection, redisClient, config)

	router := gin.Default()

	router.GET("/recipes", recipesHandler.ListRecipesHandler)

	authorized := router.Group("/")
	authorized.Use(middleware.AuthMiddlewareAPIKey())
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
