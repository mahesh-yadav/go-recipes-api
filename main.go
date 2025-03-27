package main

import (
	"context"

	"github.com/gin-gonic/gin"
	"github.com/mahesh-yadav/go-recipes-api/config"
	"github.com/mahesh-yadav/go-recipes-api/database"
	_ "github.com/mahesh-yadav/go-recipes-api/docs"
	"github.com/mahesh-yadav/go-recipes-api/handlers"
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

	gin.SetMode(config.GinMode)

	database.ConnectToDB(config)

	ctx := context.Background()
	collection := database.GetCollection(config, "recipes")

	if config.InitializeDB {
		database.InitDB(collection)
	}

	recipesHandler := handlers.NewRecipeHandler(ctx, collection)

	router := gin.Default()

	router.POST("/recipes", recipesHandler.CreateRecipeHandler)
	router.GET("/recipes", recipesHandler.ListRecipesHandler)
	router.GET("/recipes/:id", recipesHandler.GetRecipeHandler)
	router.PUT("/recipes/:id", recipesHandler.UpdateRecipeHandler)
	router.DELETE("/recipes/:id", recipesHandler.DeleteRecipeHandler)
	router.GET("/recipes/search", recipesHandler.SearchRecipeHandler)

	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	router.Run()
}
