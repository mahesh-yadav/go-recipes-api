package main

import (
	"time"

	"github.com/gin-gonic/gin"
)

// Data Models
type Recipe struct {
	Name         string    `json:"name"`
	Tags         []string  `json:"tags"`
	Ingredients  []string  `json:"ingredients"`
	Instructions []string  `json:"instructions"`
	PublishedAt  time.Time `json:"published_at"`
}

func main() {
	router := gin.Default()
	router.Run()
}
