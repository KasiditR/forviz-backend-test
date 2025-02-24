package routes

import (
	"github.com/KasiditR/forviz-backend-api-test/internal/handlers"
	"github.com/KasiditR/forviz-backend-api-test/internal/middlewares"
	"github.com/gin-gonic/gin"
)

func bookRoutes(routes *gin.RouterGroup) {
	routes.Use(middlewares.Authentication())
	routes.POST("/add", handlers.AddBook())
	routes.DELETE("/delete/:id", handlers.DeleteBook())
	routes.PUT("/edit", handlers.EditBook())
	routes.GET("/detail/:id", handlers.GetBook())
	routes.GET("/search", handlers.SearchBook())
	routes.POST("/borrow", handlers.BorrowBook())
	routes.POST("/return", handlers.ReturnBook())
	routes.GET("/most-borrow", handlers.GetMostBorrowedBooks())
}
