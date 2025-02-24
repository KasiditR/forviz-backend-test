package routes

import "github.com/gin-gonic/gin"

func MainRoutes(routes *gin.RouterGroup) {
	authRoute := routes.Group("/auth")
	authRoutes(authRoute)

	bookRoute := routes.Group("/book")
	bookRoutes(bookRoute)
}
