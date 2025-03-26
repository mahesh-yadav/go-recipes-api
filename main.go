package main

import (
	"encoding/json"
	"errors"
	"net/http"
	"os"
	"time"

	"slices"

	"github.com/gin-gonic/gin"
	_ "github.com/mahesh-yadav/go-recipes-api/docs"
	"github.com/mahesh-yadav/go-recipes-api/httputil"
	"github.com/rs/xid"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
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
	ID           string    `json:"id" example:"c0283p3d0cvuglq85log"`
	Name         string    `json:"name" binding:"required" example:"Chocolate Chip Cookies"`
	Tags         []string  `json:"tags" example:"dessert,snack"`
	Ingredients  []string  `json:"ingredients" binding:"required" example:"2 1/4 cups all-purpose flour,1 tsp baking soda,1 cup butter,3/4 cup granulated sugar,3/4 cup brown sugar,2 large eggs,2 cups semi-sweet chocolate chips"`
	Instructions []string  `json:"instructions" example:"Preheat oven to 375°F (190°C),Mix dry ingredients,Cream butter and sugars,Beat in eggs,Stir in chocolate chips,Drop spoonfuls onto baking sheets,Bake for 9 to 11 minutes"`
	PublishedAt  time.Time `json:"published_at" example:"2023-03-10T15:04:05Z"`
}

type AddUpdateRecipe struct {
	Name         string   `json:"name" binding:"required" example:"Chocolate Chip Cookies"`
	Tags         []string `json:"tags" example:"dessert,snack"`
	Ingredients  []string `json:"ingredients" binding:"required" example:"2 1/4 cups all-purpose flour,1 tsp baking soda,1 cup butter,3/4 cup granulated sugar,3/4 cup brown sugar,2 large eggs,2 cups semi-sweet chocolate chips"`
	Instructions []string `json:"instructions" example:"Preheat oven to 375°F (190°C),Mix dry ingredients,Cream butter and sugars,Beat in eggs,Stir in chocolate chips,Drop spoonfuls onto baking sheets,Bake for 9 to 11 minutes"`
}

type ListRecipes struct {
	Count int      `json:"count"`
	Data  []Recipe `json:"data"`
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
//	@Param			recipe	body		AddUpdateRecipe	true	"Add Recipe"
//	@Success		201		{object}	Recipe
//	@Failure		400		{object}	httputil.HTTPError
//	@Router			/recipes [post]
func NewRecipeHandler(c *gin.Context) {
	var addRecipe AddUpdateRecipe
	if err := c.ShouldBindJSON(&addRecipe); err != nil {
		httputil.NewError(c, http.StatusBadRequest, err)
		return
	}

	recipe := Recipe{
		Name:         addRecipe.Name,
		Tags:         addRecipe.Tags,
		Ingredients:  addRecipe.Ingredients,
		Instructions: addRecipe.Instructions,
	}

	recipe.ID = xid.New().String()
	recipe.PublishedAt = time.Now()

	recipes = append(recipes, recipe)
	c.JSON(http.StatusCreated, recipe)
}

// ListRecipesHandler godoc
//
//	@Summary		List Recipes
//	@Description	List Recipes
//	@Tags			recipes
//	@Accept			json
//	@Produce		json
//	@Success		200	{object}	ListRecipes
//	@Router			/recipes [get]
func ListRecipesHandler(c *gin.Context) {
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
//	@Param			id	path		string	true	"Recipe ID"
//	@Param			recipe	body		AddUpdateRecipe	true	"Update Recipe"
//	@Success		200		{object}	Recipe
//	@Failure		400		{object}	httputil.HTTPError
//	@Router			/recipes/{id} [put]
func UpdateRecipeHandler(c *gin.Context) {
	id := c.Param("id")

	var updateRecipe AddUpdateRecipe
	if err := c.ShouldBindJSON(&updateRecipe); err != nil {
		httputil.NewError(c, http.StatusBadRequest, err)
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
		httputil.NewError(c, http.StatusBadRequest, errors.New("Recipe not found"))
		return
	}

	recipe := Recipe{
		Name:         updateRecipe.Name,
		Tags:         updateRecipe.Tags,
		Ingredients:  updateRecipe.Ingredients,
		Instructions: updateRecipe.Instructions,
	}

	recipe.ID = id
	recipe.PublishedAt = time.Now()
	recipes[index] = recipe
	c.JSON(http.StatusOK, recipe)
}

// DeleteRecipeHandler godoc
//
//	@Summary		Delete a recipe
//	@Description	Delete a recipe
//	@Tags			recipes
//	@Accept			json
//	@Produce		json
//	@Param			id	path		string	true	"Recipe ID"
//	@Success		204		{object}	Recipe
//	@Failure		400		{object}	httputil.HTTPError
//	@Router			/recipes/{id} [delete]
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
		httputil.NewError(c, http.StatusBadRequest, errors.New("Recipe not found"))
		return
	}

	recipes = append(recipes[:index], recipes[index+1:]...)
	c.JSON(http.StatusNoContent, nil)
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
//	@Failure		400	{object}	httputil.HTTPError
//	@Router			/recipes/search [get]
func SearchRecipeHandler(c *gin.Context) {
	var searchParams SearchParams

	if err := c.ShouldBindQuery(&searchParams); err != nil {
		httputil.NewError(c, http.StatusBadRequest, err)
		return
	}

	listOfRecipes := make([]Recipe, 0)

	for _, recipe := range recipes {
		if slices.Contains(recipe.Tags, searchParams.Tag) {
			listOfRecipes = append(listOfRecipes, recipe)
		}
	}

	c.JSON(http.StatusOK, ListRecipes{
		Count: len(listOfRecipes),
		Data:  listOfRecipes,
	})
}

// @Title			Recipes API
// @Version		1.0
// @Description	This is a simple API for managing recipes.
//
// @Host			localhost:8080
// @BasePath		/
func main() {
	initDb()
	router := gin.Default()

	router.POST("/recipes", NewRecipeHandler)
	router.GET("/recipes", ListRecipesHandler)
	router.PUT("/recipes/:id", UpdateRecipeHandler)
	router.DELETE("/recipes/:id", DeleteRecipeHandler)
	router.GET("/recipes/search", SearchRecipeHandler)

	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	router.Run()
}
