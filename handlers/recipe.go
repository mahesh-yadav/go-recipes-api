package handlers

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/rs/zerolog/log"

	"github.com/gin-gonic/gin"
	"github.com/mahesh-yadav/go-recipes-api/config"
	"github.com/mahesh-yadav/go-recipes-api/models"
	"github.com/redis/go-redis/v9"
	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
)

type RecipeHandler struct {
	collection  *mongo.Collection
	ctx         context.Context
	redisClient *redis.Client
	config      *config.Config
}

func NewRecipeHandler(ctx context.Context, collection *mongo.Collection, redisClient *redis.Client, config *config.Config) *RecipeHandler {
	return &RecipeHandler{
		collection:  collection,
		ctx:         ctx,
		redisClient: redisClient,
		config:      config,
	}
}

// ListRecipesHandler godoc
//
//	@Summary		List all recipes
//	@Description	Get a list of all recipes
//	@Tags			recipes
//	@Accept			json
//	@Produce		json
//	@Success		200	{object}	models.ListRecipes
//	@Failure		500	{object}	models.ErrorResponse
//	@Router			/recipes [get]
func (handler *RecipeHandler) ListRecipesHandler(c *gin.Context) {
	var redisResults string
	var err error = nil
	if handler.config.EnableRedisCache {
		redisResults, err = handler.redisClient.Get(handler.ctx, "recipes").Result()
	}
	if !config.GetConfig().EnableRedisCache || err == redis.Nil {
		log.Info().Msg("Fetching from MongoDB...")

		cursor, err := handler.collection.Find(context.TODO(), bson.D{})
		if err != nil {
			log.Panic().Msg("Error fetching recipes from MongoDB")
		}
		defer cursor.Close(handler.ctx)

		recipes := make([]models.ViewRecipe, 0)
		for cursor.Next(context.TODO()) {
			var recipe models.ViewRecipe
			if err := cursor.Decode(&recipe); err != nil {
				log.Panic().Msg("Error decoding recipe from MongoDB")
			}
			recipes = append(recipes, recipe)
		}

		if handler.config.EnableRedisCache {
			data, err := json.Marshal(recipes)
			if err != nil {
				log.Panic().Msg("Error marshalling recipies to JSON")
			}
			handler.redisClient.Set(handler.ctx, "recipes", string(data), 0)
		}

		c.JSON(http.StatusOK, models.ListRecipes{
			Count: len(recipes),
			Data:  recipes,
		})
	} else if err != nil {
		log.Panic().Msg("Error fetching recipies from Redis cache")
	} else {
		log.Info().Msg("Retrieved from Redis cache...")
		recipes := make([]models.ViewRecipe, 0)
		err := json.Unmarshal([]byte(redisResults), &recipes)
		if err != nil {
			log.Panic().Msg("Error unmarshalling recipies from Redis cache to JSON")
		}

		c.JSON(http.StatusOK, models.ListRecipes{
			Count: len(recipes),
			Data:  recipes,
		})
	}
}

// GetRecipeHandler godoc
//
//	@Summary		Get a recipe by ID
//	@Description	Get details of a specific recipe by its ID
//	@Tags			recipes
//	@Accept			json
//	@Produce		json
//	@Param			id	path		string	true	"Recipe ID"
//	@Success		200	{object}	models.ViewRecipe
//	@Failure		400	{object}	models.ErrorResponse
//	@Failure		404	{object}	models.ErrorResponse
//	@Failure		500	{object}	models.ErrorResponse
//	@Router			/recipes/{id} [get]
func (handler *RecipeHandler) GetRecipeHandler(c *gin.Context) {
	id := c.Param("id")

	objectID, err := bson.ObjectIDFromHex(id)
	if err != nil {
		log.Error().Str("ID", id).Msg("Invalid Recipe ID")
		c.JSON(http.StatusBadRequest, models.ErrorResponse{
			Code:    http.StatusBadRequest,
			Message: "Invalid Recipe ID",
		})
		return
	}

	filter := bson.D{{Key: "_id", Value: objectID}}

	var recipe models.ViewRecipe
	err = handler.collection.FindOne(handler.ctx, filter).Decode(&recipe)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			c.JSON(http.StatusNotFound, models.ErrorResponse{
				Code:    http.StatusNotFound,
				Message: fmt.Sprintf("Recipe not found with ID: %s", id),
			})
			return
		}

		log.Panic().Msg("Error fetching recipe from MongoDB")
	}

	c.JSON(http.StatusOK, recipe)
}

// CreateRecipeHandler godoc
//
//	@Summary		Create a new recipe
//	@Description	Create a new recipe
//	@Tags			recipes
//	@Accept			json
//	@Produce		json
//	@Param			recipe	body	models.AddUpdateRecipe	true	"Add Recipe"
//	@Success		201
//	@Failure		400	{object}	models.ErrorResponse
//	@Failure		500	{object}	models.ErrorResponse
//	@Router			/recipes [post]
func (handler *RecipeHandler) CreateRecipeHandler(c *gin.Context) {
	var addRecipe models.AddUpdateRecipe
	if err := c.ShouldBindJSON(&addRecipe); err != nil {
		c.JSON(http.StatusBadRequest, models.ErrorResponse{
			Code:    http.StatusBadRequest,
			Message: "Invalid Recipe data",
		})
		return
	}

	recipe := models.Recipe{
		Name:         addRecipe.Name,
		Tags:         addRecipe.Tags,
		Ingredients:  addRecipe.Ingredients,
		Instructions: addRecipe.Instructions,
	}
	recipe.PublishedAt = time.Now()

	result, err := handler.collection.InsertOne(handler.ctx, recipe)
	if err != nil {
		log.Panic().Msg("Error inserting recipe into MongoDB")
	}

	if handler.config.EnableRedisCache {
		log.Info().Msg("Removing recipes from Redis cache...")
		handler.redisClient.Del(handler.ctx, "recipes")
	}
	c.JSON(http.StatusCreated, result)
}

// UpdateRecipeHandler godoc
//
//	@Summary		Update a recipe
//	@Description	Update a recipe
//	@Tags			recipes
//	@Accept			json
//	@Produce		json
//	@Param			id		path	string					true	"Recipe ID"
//	@Param			recipe	body	models.AddUpdateRecipe	true	"Update Recipe"
//	@Success		200
//	@Failure		400	{object}	models.ErrorResponse
//	@Failure		500	{object}	models.ErrorResponse
//	@Router			/recipes/{id} [put]
func (handler *RecipeHandler) UpdateRecipeHandler(c *gin.Context) {
	id := c.Param("id")

	objectID, err := bson.ObjectIDFromHex(id)
	if err != nil {
		log.Error().Str("ID", id).Msg("Invalid Recipe ID")
		c.JSON(http.StatusBadRequest, models.ErrorResponse{
			Code:    http.StatusBadRequest,
			Message: "Invalid Recipe ID",
		})
		return
	}

	var updateRecipe models.AddUpdateRecipe
	if err := c.ShouldBindJSON(&updateRecipe); err != nil {
		c.JSON(http.StatusBadRequest, models.ErrorResponse{
			Code:    http.StatusBadRequest,
			Message: "Invalid Recipe data",
		})
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

	result, err := handler.collection.UpdateOne(handler.ctx, filter, updateDoc)
	if err != nil {
		log.Panic().Msg("Error updating recipe in MongoDB")
		return
	}

	if handler.config.EnableRedisCache {
		log.Info().Msg("Removing recipes from Redis cache...")
		handler.redisClient.Del(handler.ctx, "recipes")
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
//	@Failure		400	{object}	models.ErrorResponse
//	@Failure		500	{object}	models.ErrorResponse
//	@Router			/recipes/{id} [delete]
func (handler *RecipeHandler) DeleteRecipeHandler(c *gin.Context) {
	id := c.Param("id")
	objectID, err := bson.ObjectIDFromHex(id)
	if err != nil {
		log.Error().Str("ID", id).Msg("Invalid Recipe ID")
		c.JSON(http.StatusBadRequest, models.ErrorResponse{
			Code:    http.StatusBadRequest,
			Message: "Invalid Recipe ID",
		})
		return
	}

	filter := bson.D{{Key: "_id", Value: objectID}}

	result, err := handler.collection.DeleteOne(handler.ctx, filter)
	if err != nil {
		log.Panic().Msg("Error deleting recipe in MongoDB")
		return
	}

	if handler.config.EnableRedisCache {
		log.Info().Msg("Removing recipes from Redis cache...")
		handler.redisClient.Del(handler.ctx, "recipes")
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
//	@Success		200	{object}	models.ListRecipes
//	@Failure		400	{object}	models.ErrorResponse
//	@Failure		500	{object}	models.ErrorResponse
//	@Router			/recipes/search [get]
func (handler *RecipeHandler) SearchRecipeHandler(c *gin.Context) {
	var searchParams models.RecipeTagSearchParams

	if err := c.ShouldBindQuery(&searchParams); err != nil {
		c.JSON(http.StatusBadRequest, models.ErrorResponse{
			Code:    http.StatusBadRequest,
			Message: "Invalid Recipe search parameters",
		})
		return
	}

	filter := bson.D{{Key: "tags", Value: bson.D{{Key: "$in", Value: []string{searchParams.Tag}}}}}

	cursor, err := handler.collection.Find(handler.ctx, filter)
	if err != nil {
		log.Panic().Msg("Error searching recipes in MongoDB")
		return
	}
	defer cursor.Close(handler.ctx)

	recipes := make([]models.ViewRecipe, 0)
	for cursor.Next(handler.ctx) {
		var recipe models.ViewRecipe
		if err := cursor.Decode(&recipe); err != nil {
			log.Panic().Msg("Error decoding recipe from MongoDB")
			return
		}
		recipes = append(recipes, recipe)
	}

	c.JSON(http.StatusOK, models.ListRecipes{
		Count: len(recipes),
		Data:  recipes,
	})
}
