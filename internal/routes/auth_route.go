package routes

import (
	"github.com/KasiditR/forviz-backend-api-test/internal/handlers"
	"github.com/gin-gonic/gin"
)

func authRoutes(routes *gin.RouterGroup) {
	routes.POST("/register", handlers.Register())
	routes.POST("/login", handlers.Login())
}
