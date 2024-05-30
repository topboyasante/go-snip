package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/topboyasante/go-snip/api/v1/controllers"
	"github.com/topboyasante/go-snip/api/v1/middleware"
)

func SnippetRoutes(r *gin.RouterGroup) {
	snippetRoutes := r.Group("/snippets")
	snippetRoutes.Use(middleware.RequireAuth)
	
	snippetRoutes.GET("", controllers.GetSnippets)
	snippetRoutes.GET("/:id", controllers.GetSnippet)
	snippetRoutes.POST("/create", controllers.CreateSnippet)
	snippetRoutes.DELETE("/:id", controllers.DeleteSnippet)
}
