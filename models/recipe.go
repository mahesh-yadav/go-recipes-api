package models

import (
	"time"

	"go.mongodb.org/mongo-driver/v2/bson"
)

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

type RecipeTagSearchParams struct {
	Tag string `form:"tag" binding:"required,min=3"`
}
