package main

import (
	"encoding/json"
	"net/http"
	"os"
	"time"

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
	Name         string    `json:"name"`
	Tags         []string  `json:"tags"`
	Ingredients  []string  `json:"ingredients"`
	Instructions []string  `json:"instructions"`
	PublishedAt  time.Time `json:"published_at"`
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

func main() {
	initDb()
	router := gin.Default()

	router.POST("/recipes", NewRecipeHandler)
	router.GET("/recipes", ListRecipesHandler)
	router.Run()
}
