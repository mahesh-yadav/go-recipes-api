package main

import (
	"encoding/json"
	"net/http"
	"os"
	"time"

	"slices"

	"github.com/gin-gonic/gin"
	"github.com/rs/xid"
)

// In-Memory DB
var recipes []Recipe

func initDb() {
	recipes = make([]Recipe, 0)

	file, _ := os.ReadFile("recipes.json")

	json.Unmarshal(file, &recipes)
}

// Data Models
type Recipe struct {
	ID           string    `json:"id"`
	Name         string    `json:"name" binding:"required"`
	Tags         []string  `json:"tags"`
	Ingredients  []string  `json:"ingredients" binding:"required"`
	Instructions []string  `json:"instructions"`
	PublishedAt  time.Time `json:"published_at"`
}

type SearchParams struct {
	Tag string `form:"tag" binding:"required,min=3"`
}

// Route Handlers
func NewRecipeHandler(c *gin.Context) {
	var recipe Recipe
	if err := c.ShouldBindJSON(&recipe); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	recipe.ID = xid.New().String()
	recipe.PublishedAt = time.Now()

	recipes = append(recipes, recipe)
	c.JSON(http.StatusCreated, recipe)
}

func ListRecipesHandler(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"count":   len(recipes),
		"recipes": recipes,
	})
}

func UpdateRecipeHandler(c *gin.Context) {
	id := c.Param("id")

	var updatedRecipe Recipe
	if err := c.ShouldBindJSON(&updatedRecipe); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	index := -1
	for i := 0; i < len(recipes); i++ {
		if recipes[i].ID == id {
			index = i
			break
		}
	}

	if index == -1 {
		c.JSON(http.StatusNotFound, gin.H{
			"error": "Recipe not found",
		})
		return
	}

	updatedRecipe.ID = id
	updatedRecipe.PublishedAt = time.Now()
	recipes[index] = updatedRecipe
	c.JSON(http.StatusOK, updatedRecipe)
}

func DeleteRecipeHandler(c *gin.Context) {
	id := c.Param("id")

	index := -1
	for i := 0; i < len(recipes); i++ {
		if recipes[i].ID == id {
			index = i
			break
		}
	}

	if index == -1 {
		c.JSON(http.StatusNotFound, gin.H{
			"error": "Recipe not found",
		})
		return
	}

	recipes = append(recipes[:index], recipes[index+1:]...)
	c.JSON(http.StatusOK, gin.H{
		"message": "Recipe deleted successfully",
	})
}

func SearchRecipeHandler(c *gin.Context) {
	var searchParams SearchParams

	if err := c.ShouldBindQuery(&searchParams); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	listOfRecipes := make([]Recipe, 0)

	for _, recipe := range recipes {
		if slices.Contains(recipe.Tags, searchParams.Tag) {
			listOfRecipes = append(listOfRecipes, recipe)
		}
	}

	c.JSON(http.StatusOK, gin.H{
		"count":   len(listOfRecipes),
		"recipes": listOfRecipes,
	})
}

func main() {
	initDb()
	router := gin.Default()

	router.POST("/recipes", NewRecipeHandler)
	router.GET("/recipes", ListRecipesHandler)
	router.PUT("/recipes/:id", UpdateRecipeHandler)
	router.DELETE("/recipes/:id", DeleteRecipeHandler)
	router.GET("/recipes/search", SearchRecipeHandler)
	router.Run()
}
