package main

import (
	"github.com/KasiditR/forviz-backend-api-test/internal/config"
	"github.com/KasiditR/forviz-backend-api-test/internal/database"
	"github.com/KasiditR/forviz-backend-api-test/internal/routes"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
)

func main() {
	router := gin.New()
	router.Use(cors.New(cors.Config{
		AllowOriginFunc: func(origin string) bool {
			return true
		},
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Content-Type", "Authorization"},
		AllowCredentials: true,
	}))
	router.Use((gin.Logger()))
	router.GET("/", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"message": "Hello, World!"})
	})
	apiRoute := router.Group("/api/v1")
	routes.MainRoutes(apiRoute)
	database.ConnectDatabase()
	log.Fatal(router.Run(":" + config.LoadConfig().Port))
}
