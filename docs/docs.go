// Package docs Code generated by swaggo/swag. DO NOT EDIT
package docs

import "github.com/swaggo/swag"

const docTemplate = `{
    "schemes": {{ marshal .Schemes }},
    "swagger": "2.0",
    "info": {
        "description": "{{escape .Description}}",
        "title": "{{.Title}}",
        "contact": {},
        "version": "{{.Version}}"
    },
    "host": "{{.Host}}",
    "basePath": "{{.BasePath}}",
    "paths": {
        "/auth/refresh": {
            "post": {
                "description": "Refresh an existing JWT token and return a new one",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "auth"
                ],
                "summary": "Refresh JWT token",
                "parameters": [
                    {
                        "type": "string",
                        "description": "{token}",
                        "name": "Authorization",
                        "in": "header",
                        "required": true
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "$ref": "#/definitions/handlers.JWTOutput"
                        }
                    },
                    "401": {
                        "description": "Unauthorized",
                        "schema": {
                            "$ref": "#/definitions/models.ErrorResponse"
                        }
                    },
                    "500": {
                        "description": "Internal Server Error",
                        "schema": {
                            "$ref": "#/definitions/models.ErrorResponse"
                        }
                    }
                }
            }
        },
        "/auth/signin": {
            "post": {
                "description": "Authenticate a user and return a JWT token",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "auth"
                ],
                "summary": "Sign in a user",
                "parameters": [
                    {
                        "description": "User Credentials",
                        "name": "user",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "$ref": "#/definitions/models.User"
                        }
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "$ref": "#/definitions/handlers.JWTOutput"
                        }
                    },
                    "400": {
                        "description": "Bad Request",
                        "schema": {
                            "$ref": "#/definitions/models.ErrorResponse"
                        }
                    },
                    "401": {
                        "description": "Unauthorized",
                        "schema": {
                            "$ref": "#/definitions/models.ErrorResponse"
                        }
                    },
                    "500": {
                        "description": "Internal Server Error",
                        "schema": {
                            "$ref": "#/definitions/models.ErrorResponse"
                        }
                    }
                }
            }
        },
        "/auth/signup": {
            "post": {
                "description": "Create a new user account",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "auth"
                ],
                "summary": "Sign up a new user",
                "parameters": [
                    {
                        "description": "User Sign Up",
                        "name": "user",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "$ref": "#/definitions/models.User"
                        }
                    }
                ],
                "responses": {
                    "201": {
                        "description": "Created"
                    },
                    "400": {
                        "description": "Bad Request",
                        "schema": {
                            "$ref": "#/definitions/models.ErrorResponse"
                        }
                    },
                    "500": {
                        "description": "Internal Server Error",
                        "schema": {
                            "$ref": "#/definitions/models.ErrorResponse"
                        }
                    }
                }
            }
        },
        "/recipes": {
            "get": {
                "description": "Get a list of all recipes",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "recipes"
                ],
                "summary": "List all recipes",
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "$ref": "#/definitions/models.ListRecipes"
                        }
                    },
                    "500": {
                        "description": "Internal Server Error",
                        "schema": {
                            "$ref": "#/definitions/models.ErrorResponse"
                        }
                    }
                }
            },
            "post": {
                "description": "Create a new recipe",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "recipes"
                ],
                "summary": "Create a new recipe",
                "parameters": [
                    {
                        "description": "Add Recipe",
                        "name": "recipe",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "$ref": "#/definitions/models.AddUpdateRecipe"
                        }
                    }
                ],
                "responses": {
                    "201": {
                        "description": "Created"
                    },
                    "400": {
                        "description": "Bad Request",
                        "schema": {
                            "$ref": "#/definitions/models.ErrorResponse"
                        }
                    },
                    "500": {
                        "description": "Internal Server Error",
                        "schema": {
                            "$ref": "#/definitions/models.ErrorResponse"
                        }
                    }
                }
            }
        },
        "/recipes/search": {
            "get": {
                "description": "Search recipes by tag",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "recipes"
                ],
                "summary": "Search recipes by tag",
                "parameters": [
                    {
                        "type": "string",
                        "description": "Tag to search for",
                        "name": "tag",
                        "in": "query",
                        "required": true
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "$ref": "#/definitions/models.ListRecipes"
                        }
                    },
                    "400": {
                        "description": "Bad Request",
                        "schema": {
                            "$ref": "#/definitions/models.ErrorResponse"
                        }
                    },
                    "500": {
                        "description": "Internal Server Error",
                        "schema": {
                            "$ref": "#/definitions/models.ErrorResponse"
                        }
                    }
                }
            }
        },
        "/recipes/{id}": {
            "get": {
                "description": "Get details of a specific recipe by its ID",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "recipes"
                ],
                "summary": "Get a recipe by ID",
                "parameters": [
                    {
                        "type": "string",
                        "description": "Recipe ID",
                        "name": "id",
                        "in": "path",
                        "required": true
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "$ref": "#/definitions/models.ViewRecipe"
                        }
                    },
                    "400": {
                        "description": "Bad Request",
                        "schema": {
                            "$ref": "#/definitions/models.ErrorResponse"
                        }
                    },
                    "404": {
                        "description": "Not Found",
                        "schema": {
                            "$ref": "#/definitions/models.ErrorResponse"
                        }
                    },
                    "500": {
                        "description": "Internal Server Error",
                        "schema": {
                            "$ref": "#/definitions/models.ErrorResponse"
                        }
                    }
                }
            },
            "put": {
                "description": "Update a recipe",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "recipes"
                ],
                "summary": "Update a recipe",
                "parameters": [
                    {
                        "type": "string",
                        "description": "Recipe ID",
                        "name": "id",
                        "in": "path",
                        "required": true
                    },
                    {
                        "description": "Update Recipe",
                        "name": "recipe",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "$ref": "#/definitions/models.AddUpdateRecipe"
                        }
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK"
                    },
                    "400": {
                        "description": "Bad Request",
                        "schema": {
                            "$ref": "#/definitions/models.ErrorResponse"
                        }
                    },
                    "500": {
                        "description": "Internal Server Error",
                        "schema": {
                            "$ref": "#/definitions/models.ErrorResponse"
                        }
                    }
                }
            },
            "delete": {
                "description": "Delete a recipe",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "recipes"
                ],
                "summary": "Delete a recipe",
                "parameters": [
                    {
                        "type": "string",
                        "description": "Recipe ID",
                        "name": "id",
                        "in": "path",
                        "required": true
                    }
                ],
                "responses": {
                    "204": {
                        "description": "No Content"
                    },
                    "400": {
                        "description": "Bad Request",
                        "schema": {
                            "$ref": "#/definitions/models.ErrorResponse"
                        }
                    },
                    "500": {
                        "description": "Internal Server Error",
                        "schema": {
                            "$ref": "#/definitions/models.ErrorResponse"
                        }
                    }
                }
            }
        }
    },
    "definitions": {
        "handlers.JWTOutput": {
            "type": "object",
            "properties": {
                "expires": {
                    "type": "string"
                },
                "token": {
                    "type": "string"
                }
            }
        },
        "models.AddUpdateRecipe": {
            "type": "object",
            "required": [
                "ingredients",
                "instructions",
                "name",
                "tags"
            ],
            "properties": {
                "ingredients": {
                    "type": "array",
                    "items": {
                        "type": "string"
                    },
                    "example": [
                        "2 1/4 cups all-purpose flour",
                        "1 tsp baking soda",
                        "1 cup butter",
                        "3/4 cup granulated sugar",
                        "3/4 cup brown sugar",
                        "2 large eggs",
                        "2 cups semi-sweet chocolate chips"
                    ]
                },
                "instructions": {
                    "type": "array",
                    "items": {
                        "type": "string"
                    },
                    "example": [
                        "Preheat oven to 375°F (190°C)",
                        "Mix dry ingredients",
                        "Cream butter and sugars",
                        "Beat in eggs",
                        "Stir in chocolate chips",
                        "Drop spoonfuls onto baking sheets",
                        "Bake for 9 to 11 minutes"
                    ]
                },
                "name": {
                    "type": "string",
                    "example": "Chocolate Chip Cookies"
                },
                "tags": {
                    "type": "array",
                    "items": {
                        "type": "string"
                    },
                    "example": [
                        "dessert",
                        "snack"
                    ]
                }
            }
        },
        "models.ErrorResponse": {
            "type": "object",
            "properties": {
                "code": {
                    "type": "integer"
                },
                "message": {
                    "type": "string"
                }
            }
        },
        "models.ListRecipes": {
            "type": "object",
            "properties": {
                "count": {
                    "type": "integer"
                },
                "data": {
                    "type": "array",
                    "items": {
                        "$ref": "#/definitions/models.ViewRecipe"
                    }
                }
            }
        },
        "models.User": {
            "type": "object",
            "required": [
                "password",
                "username"
            ],
            "properties": {
                "password": {
                    "type": "string"
                },
                "username": {
                    "type": "string"
                }
            }
        },
        "models.ViewRecipe": {
            "type": "object",
            "properties": {
                "id": {
                    "type": "string",
                    "example": "c0283p3d0cvuglq85log"
                },
                "ingredients": {
                    "type": "array",
                    "items": {
                        "type": "string"
                    },
                    "example": [
                        "2 1/4 cups all-purpose flour",
                        "1 tsp baking soda",
                        "1 cup butter",
                        "3/4 cup granulated sugar",
                        "3/4 cup brown sugar",
                        "2 large eggs",
                        "2 cups semi-sweet chocolate chips"
                    ]
                },
                "instructions": {
                    "type": "array",
                    "items": {
                        "type": "string"
                    },
                    "example": [
                        "Preheat oven to 375°F (190°C)",
                        "Mix dry ingredients",
                        "Cream butter and sugars",
                        "Beat in eggs",
                        "Stir in chocolate chips",
                        "Drop spoonfuls onto baking sheets",
                        "Bake for 9 to 11 minutes"
                    ]
                },
                "name": {
                    "type": "string",
                    "example": "Chocolate Chip Cookies"
                },
                "published_at": {
                    "type": "string",
                    "example": "2023-03-10T15:04:05Z"
                },
                "tags": {
                    "type": "array",
                    "items": {
                        "type": "string"
                    },
                    "example": [
                        "dessert",
                        "snack"
                    ]
                }
            }
        }
    }
}`

// SwaggerInfo holds exported Swagger Info so clients can modify it
var SwaggerInfo = &swag.Spec{
	Version:          "1.0",
	Host:             "localhost:8080",
	BasePath:         "/",
	Schemes:          []string{},
	Title:            "Recipes API",
	Description:      "This is a simple API for managing recipes.",
	InfoInstanceName: "swagger",
	SwaggerTemplate:  docTemplate,
	LeftDelim:        "{{",
	RightDelim:       "}}",
}

func init() {
	swag.Register(SwaggerInfo.InstanceName(), SwaggerInfo)
}
