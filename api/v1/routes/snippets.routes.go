package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/topboyasante/go-snip/api/v1/controllers"
	"github.com/topboyasante/go-snip/api/v1/middleware"
)

func SnippetRoutes(r *gin.RouterGroup) {
	snippetRoutes := r.Group("/snippets")
	
	snippetRoutes.GET("", controllers.GetSnippets)
	snippetRoutes.GET("/:id", controllers.GetSnippet)
	
	snippetRoutes.Use(middleware.RequireAuth)
	
	snippetRoutes.POST("/create", controllers.CreateSnippet)
	snippetRoutes.PUT("/:id", controllers.UpdateSnippet)
	snippetRoutes.DELETE("/:id", controllers.DeleteSnippet)
}
