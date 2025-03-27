package main

import (
	"context"
	"encoding/json"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/mahesh-yadav/go-recipes-api/config"
	"github.com/mahesh-yadav/go-recipes-api/database"
	_ "github.com/mahesh-yadav/go-recipes-api/docs"
	"github.com/mahesh-yadav/go-recipes-api/utils"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
	"go.mongodb.org/mongo-driver/v2/bson"
)

func initDb(config *config.Config) {
	if config.InitializeDB {
		recipes := make([]Recipe, 0)
		file, _ := os.ReadFile("recipes.json")
		err := json.Unmarshal(file, &recipes)
		if err != nil {
			log.Fatal("error unmarshalling: ", err)
		}

		var listOfRecipes []interface{}
		for _, recipe := range recipes {
			listOfRecipes = append(listOfRecipes, recipe)
		}

		collection := database.GetCollection(config, "recipes")

		insertManyResult, err := collection.InsertMany(context.Background(), listOfRecipes)
		if err != nil {
			log.Fatal("error inserting: ", err)
		}

		log.Println("Inserted recipes: ", len(insertManyResult.InsertedIDs))
	} else {
		log.Println("Collection already exists. Skipping initialization.")
	}
}

// Data Models
type Recipe struct {
	Name         string    `json:"name" bson:"name" binding:"required" example:"Chocolate Chip Cookies"`
	Tags         []string  `json:"tags" bson:"tags" binding:"required" example:"dessert,snack"`
	Ingredients  []string  `json:"ingredients" bson:"ingredients" binding:"required" example:"2 1/4 cups all-purpose flour,1 tsp baking soda,1 cup butter,3/4 cup granulated sugar,3/4 cup brown sugar,2 large eggs,2 cups semi-sweet chocolate chips"`
	Instructions []string  `json:"instructions" bson:"instructions" binding:"required" example:"Preheat oven to 375°F (190°C),Mix dry ingredients,Cream butter and sugars,Beat in eggs,Stir in chocolate chips,Drop spoonfuls onto baking sheets,Bake for 9 to 11 minutes"`
	PublishedAt  time.Time `json:"published_at" bson:"published_at" example:"2023-03-10T15:04:05Z"`
}

type ViewRecipe struct {
	ID           bson.ObjectID `json:"id" bson:"_id" example:"c0283p3d0cvuglq85log"`
	Name         string        `json:"name" bson:"name" example:"Chocolate Chip Cookies"`
	Tags         []string      `json:"tags" bson:"tags" example:"dessert,snack"`
	Ingredients  []string      `json:"ingredients" bson:"ingredients" example:"2 1/4 cups all-purpose flour,1 tsp baking soda,1 cup butter,3/4 cup granulated sugar,3/4 cup brown sugar,2 large eggs,2 cups semi-sweet chocolate chips"`
	Instructions []string      `json:"instructions" bson:"instructions" example:"Preheat oven to 375°F (190°C),Mix dry ingredients,Cream butter and sugars,Beat in eggs,Stir in chocolate chips,Drop spoonfuls onto baking sheets,Bake for 9 to 11 minutes"`
	PublishedAt  time.Time     `json:"published_at" bson:"published_at" example:"2023-03-10T15:04:05Z"`
}

type AddUpdateRecipe struct {
	Name         string   `json:"name" binding:"required" example:"Chocolate Chip Cookies"`
	Tags         []string `json:"tags" binding:"required" example:"dessert,snack"`
	Ingredients  []string `json:"ingredients" binding:"required" example:"2 1/4 cups all-purpose flour,1 tsp baking soda,1 cup butter,3/4 cup granulated sugar,3/4 cup brown sugar,2 large eggs,2 cups semi-sweet chocolate chips"`
	Instructions []string `json:"instructions" binding:"required" example:"Preheat oven to 375°F (190°C),Mix dry ingredients,Cream butter and sugars,Beat in eggs,Stir in chocolate chips,Drop spoonfuls onto baking sheets,Bake for 9 to 11 minutes"`
}

type ListRecipes struct {
	Count int          `json:"count"`
	Data  []ViewRecipe `json:"data"`
}

type SearchParams struct {
	Tag string `form:"tag" binding:"required,min=3"`
}

// NewRecipeHandler godoc
//
//	@Summary		Create a new recipe
//	@Description	Create a new recipe
//	@Tags			recipes
//	@Accept			json
//	@Produce		json
//	@Param			recipe	body	AddUpdateRecipe	true	"Add Recipe"
//	@Success		201
//	@Failure		400	{object}	utils.HTTPError
//	@Failure		500	{object}	utils.HTTPError
//	@Router			/recipes [post]
func NewRecipeHandler(c *gin.Context) {
	var addRecipe AddUpdateRecipe
	if err := c.ShouldBindJSON(&addRecipe); err != nil {
		utils.NewError(c, http.StatusBadRequest, err)
		return
	}

	recipe := Recipe{
		Name:         addRecipe.Name,
		Tags:         addRecipe.Tags,
		Ingredients:  addRecipe.Ingredients,
		Instructions: addRecipe.Instructions,
	}
	recipe.PublishedAt = time.Now()

	config := config.GetConfig()
	collection := database.GetCollection(config, "recipes")

	result, err := collection.InsertOne(context.TODO(), recipe)
	if err != nil {
		log.Println("error inserting documents: ", err)
		utils.NewError(c, http.StatusInternalServerError, err)
		return
	}

	c.JSON(http.StatusCreated, result)

}

// ListRecipesHandler godoc
//
//	@Summary		List Recipes
//	@Description	List Recipes
//	@Tags			recipes
//	@Accept			json
//	@Produce		json
//	@Success		200	{object}	ListRecipes
//	@Failure		500	{object}	utils.HTTPError
//	@Router			/recipes [get]
func ListRecipesHandler(c *gin.Context) {
	config := config.GetConfig()
	collection := database.GetCollection(config, "recipes")

	cursor, err := collection.Find(context.TODO(), bson.D{})
	if err != nil {
		log.Println("Error finding documents: ", err)
		utils.NewError(c, http.StatusInternalServerError, err)
		return
	}
	defer cursor.Close(context.TODO())

	recipes := make([]ViewRecipe, 0)
	for cursor.Next(context.TODO()) {
		var recipe ViewRecipe
		if err := cursor.Decode(&recipe); err != nil {
			log.Println("Error decoding document: ", err)
			utils.NewError(c, http.StatusInternalServerError, err)
			return
		}
		recipes = append(recipes, recipe)
	}

	c.JSON(http.StatusOK, ListRecipes{
		Count: len(recipes),
		Data:  recipes,
	})
}

// UpdateRecipeHandler godoc
//
//	@Summary		Update a recipe
//	@Description	Update a recipe
//	@Tags			recipes
//	@Accept			json
//	@Produce		json
//	@Param			id		path	string			true	"Recipe ID"
//	@Param			recipe	body	AddUpdateRecipe	true	"Update Recipe"
//	@Success		200
//	@Failure		400	{object}	utils.HTTPError
//	@Failure		500	{object}	utils.HTTPError
//	@Router			/recipes/{id} [put]
func UpdateRecipeHandler(c *gin.Context) {
	id := c.Param("id")

	objectID, err := bson.ObjectIDFromHex(id)
	if err != nil {
		utils.NewError(c, http.StatusBadRequest, err)
		return
	}

	var updateRecipe AddUpdateRecipe
	if err := c.ShouldBindJSON(&updateRecipe); err != nil {
		utils.NewError(c, http.StatusBadRequest, err)
		return
	}

	filter := bson.D{{Key: "_id", Value: objectID}}
	updateDoc := bson.D{{
		Key: "$set",
		Value: bson.D{
			{Key: "name", Value: updateRecipe.Name},
			{Key: "tags", Value: updateRecipe.Tags},
			{Key: "ingredients", Value: updateRecipe.Ingredients},
			{Key: "instructions", Value: updateRecipe.Instructions},
			{Key: "published_at", Value: time.Now()},
		}}}

	config := config.GetConfig()
	collection := database.GetCollection(config, "recipes")

	result, err := collection.UpdateOne(context.TODO(), filter, updateDoc)
	if err != nil {
		log.Println("Error updating document: ", err)
		utils.NewError(c, http.StatusInternalServerError, err)
		return
	}

	c.JSON(http.StatusOK, result)

}

// DeleteRecipeHandler godoc
//
//	@Summary		Delete a recipe
//	@Description	Delete a recipe
//	@Tags			recipes
//	@Accept			json
//	@Produce		json
//	@Param			id	path	string	true	"Recipe ID"
//	@Success		204
//	@Failure		400	{object}	utils.HTTPError
//	@Failure		500	{object}	utils.HTTPError
//	@Router			/recipes/{id} [delete]
func DeleteRecipeHandler(c *gin.Context) {
	id := c.Param("id")
	objectID, err := bson.ObjectIDFromHex(id)
	if err != nil {
		utils.NewError(c, http.StatusBadRequest, err)
		return
	}

	filter := bson.D{{Key: "_id", Value: objectID}}
	config := config.GetConfig()
	collection := database.GetCollection(config, "recipes")

	result, err := collection.DeleteOne(context.TODO(), filter)
	if err != nil {
		log.Println("Error deleting document: ", err)
		utils.NewError(c, http.StatusInternalServerError, err)
		return
	}

	c.JSON(http.StatusOK, result)
}

// SearchRecipeHandler godoc
//
//	@Summary		Search recipes by tag
//	@Description	Search recipes by tag
//	@Tags			recipes
//	@Accept			json
//	@Produce		json
//	@Param			tag	query		string	true	"Tag to search for"
//	@Success		200	{object}	ListRecipes
//	@Failure		400	{object}	utils.HTTPError
//	@Failure		500	{object}	utils.HTTPError
//	@Router			/recipes/search [get]
func SearchRecipeHandler(c *gin.Context) {
	var searchParams SearchParams

	if err := c.ShouldBindQuery(&searchParams); err != nil {
		utils.NewError(c, http.StatusBadRequest, err)
		return
	}

	filter := bson.D{{Key: "tags", Value: bson.D{{Key: "$in", Value: []string{searchParams.Tag}}}}}
	config := config.GetConfig()
	collection := database.GetCollection(config, "recipes")

	cursor, err := collection.Find(context.TODO(), filter)
	if err != nil {
		log.Println("Error finding documents: ", err)
		utils.NewError(c, http.StatusInternalServerError, err)
		return
	}
	defer cursor.Close(context.TODO())

	recipes := make([]ViewRecipe, 0)
	for cursor.Next(context.TODO()) {
		var recipe ViewRecipe
		if err := cursor.Decode(&recipe); err != nil {
			log.Println("Error decoding document: ", err)
			utils.NewError(c, http.StatusInternalServerError, err)
			return
		}
		recipes = append(recipes, recipe)
	}

	c.JSON(http.StatusOK, ListRecipes{
		Count: len(recipes),
		Data:  recipes,
	})

}

//	@Title			Recipes API
//	@Version		1.0
//	@Description	This is a simple API for managing recipes.
//
//	@Host			localhost:8080
//	@BasePath		/
func main() {
	config := config.GetConfig()

	gin.SetMode(config.GinMode)

	database.ConnectToDB(config)
	initDb(config)

	router := gin.Default()

	router.POST("/recipes", NewRecipeHandler)
	router.GET("/recipes", ListRecipesHandler)
	router.PUT("/recipes/:id", UpdateRecipeHandler)
	router.DELETE("/recipes/:id", DeleteRecipeHandler)
	router.GET("/recipes/search", SearchRecipeHandler)

	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	router.Run()
}
