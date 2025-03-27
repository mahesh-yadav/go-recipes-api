package handlers

import (
	"context"
	"encoding/json"
	"log"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/mahesh-yadav/go-recipes-api/config"
	"github.com/mahesh-yadav/go-recipes-api/models"
	"github.com/mahesh-yadav/go-recipes-api/utils"
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
//	@Summary		List Recipes
//	@Description	Returns list of recipes
//	@Tags			recipes
//	@Accept			json
//	@Produce		json
//	@Success		200	{object}	ListRecipes
//	@Failure		500	{object}	utils.HTTPError
//	@Router			/recipes [get]
func (handler *RecipeHandler) ListRecipesHandler(c *gin.Context) {
	var redisResults string
	var err error = nil
	if handler.config.EnableRedisCache {
		redisResults, err = handler.redisClient.Get(handler.ctx, "recipes").Result()
	}
	if !config.GetConfig().EnableRedisCache || err == redis.Nil {
		log.Println("Fetching from MongoDB...")

		cursor, err := handler.collection.Find(context.TODO(), bson.D{})
		if err != nil {
			utils.NewError(c, http.StatusInternalServerError, err)
			return
		}
		defer cursor.Close(handler.ctx)

		recipes := make([]models.ViewRecipe, 0)
		for cursor.Next(context.TODO()) {
			var recipe models.ViewRecipe
			if err := cursor.Decode(&recipe); err != nil {
				utils.NewError(c, http.StatusInternalServerError, err)
				return
			}
			recipes = append(recipes, recipe)
		}

		if handler.config.EnableRedisCache {
			data, err := json.Marshal(recipes)
			if err != nil {
				utils.NewError(c, http.StatusInternalServerError, err)
				return
			}
			handler.redisClient.Set(handler.ctx, "recipes", string(data), 0)
		}

		c.JSON(http.StatusOK, models.ListRecipes{
			Count: len(recipes),
			Data:  recipes,
		})
	} else if err != nil {
		utils.NewError(c, http.StatusInternalServerError, err)
		return
	} else {
		log.Println("Retrieved from Redis cache...")
		recipes := make([]models.ViewRecipe, 0)
		err := json.Unmarshal([]byte(redisResults), &recipes)
		if err != nil {
			utils.NewError(c, http.StatusInternalServerError, err)
			return
		}

		c.JSON(http.StatusOK, models.ListRecipes{
			Count: len(recipes),
			Data:  recipes,
		})
	}
}

// GetRecipeHandler godoc
//
//	@Summary		Get a recipe
//	@Description	Returns a single recipe
//	@Tags			recipes
//	@Accept			json
//	@Produce		json
//	@Param			id		path	string			true	"Recipe ID"
//	@Success		200
//	@Failure		400	{object}	utils.HTTPError
//	@Failure		404	{object}	utils.HTTPError
//	@Failure		500	{object}	utils.HTTPError
//	@Router			/recipes/{id} [get]
func (handler *RecipeHandler) GetRecipeHandler(c *gin.Context) {
	id := c.Param("id")

	objectID, err := bson.ObjectIDFromHex(id)
	if err != nil {
		utils.NewError(c, http.StatusBadRequest, err)
		return
	}

	filter := bson.D{{Key: "_id", Value: objectID}}

	var recipe models.ViewRecipe
	err = handler.collection.FindOne(handler.ctx, filter).Decode(&recipe)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			utils.NewError(c, http.StatusNotFound, err)
			return
		}

		utils.NewError(c, http.StatusInternalServerError, err)
		return
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
//	@Param			recipe	body	AddUpdateRecipe	true	"Add Recipe"
//	@Success		201
//	@Failure		400	{object}	utils.HTTPError
//	@Failure		500	{object}	utils.HTTPError
//	@Router			/recipes [post]
func (handler *RecipeHandler) CreateRecipeHandler(c *gin.Context) {
	var addRecipe models.AddUpdateRecipe
	if err := c.ShouldBindJSON(&addRecipe); err != nil {
		utils.NewError(c, http.StatusBadRequest, err)
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
		utils.NewError(c, http.StatusInternalServerError, err)
		return
	}

	if handler.config.EnableRedisCache {
		log.Println("Removing recipes from Redis cache...")
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
//	@Param			id		path	string			true	"Recipe ID"
//	@Param			recipe	body	AddUpdateRecipe	true	"Update Recipe"
//	@Success		200
//	@Failure		400	{object}	utils.HTTPError
//	@Failure		500	{object}	utils.HTTPError
//	@Router			/recipes/{id} [put]
func (handler *RecipeHandler) UpdateRecipeHandler(c *gin.Context) {
	id := c.Param("id")

	objectID, err := bson.ObjectIDFromHex(id)
	if err != nil {
		utils.NewError(c, http.StatusBadRequest, err)
		return
	}

	var updateRecipe models.AddUpdateRecipe
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

	result, err := handler.collection.UpdateOne(handler.ctx, filter, updateDoc)
	if err != nil {
		utils.NewError(c, http.StatusInternalServerError, err)
		return
	}

	if handler.config.EnableRedisCache {
		log.Println("Removing recipes from Redis cache...")
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
//	@Failure		400	{object}	utils.HTTPError
//	@Failure		500	{object}	utils.HTTPError
//	@Router			/recipes/{id} [delete]
func (handler *RecipeHandler) DeleteRecipeHandler(c *gin.Context) {
	id := c.Param("id")
	objectID, err := bson.ObjectIDFromHex(id)
	if err != nil {
		utils.NewError(c, http.StatusBadRequest, err)
		return
	}

	filter := bson.D{{Key: "_id", Value: objectID}}

	result, err := handler.collection.DeleteOne(handler.ctx, filter)
	if err != nil {
		utils.NewError(c, http.StatusInternalServerError, err)
		return
	}

	if handler.config.EnableRedisCache {
		log.Println("Removing recipes from Redis cache...")
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
//	@Success		200	{object}	ListRecipes
//	@Failure		400	{object}	utils.HTTPError
//	@Failure		500	{object}	utils.HTTPError
//	@Router			/recipes/search [get]
func (handler *RecipeHandler) SearchRecipeHandler(c *gin.Context) {
	var searchParams models.RecipeTagSearchParams

	if err := c.ShouldBindQuery(&searchParams); err != nil {
		utils.NewError(c, http.StatusBadRequest, err)
		return
	}

	filter := bson.D{{Key: "tags", Value: bson.D{{Key: "$in", Value: []string{searchParams.Tag}}}}}

	cursor, err := handler.collection.Find(handler.ctx, filter)
	if err != nil {
		utils.NewError(c, http.StatusInternalServerError, err)
		return
	}
	defer cursor.Close(handler.ctx)

	recipes := make([]models.ViewRecipe, 0)
	for cursor.Next(handler.ctx) {
		var recipe models.ViewRecipe
		if err := cursor.Decode(&recipe); err != nil {
			utils.NewError(c, http.StatusInternalServerError, err)
			return
		}
		recipes = append(recipes, recipe)
	}

	c.JSON(http.StatusOK, models.ListRecipes{
		Count: len(recipes),
		Data:  recipes,
	})
}
