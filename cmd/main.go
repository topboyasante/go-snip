package main

import (
	"log"

	"github.com/gin-gonic/gin"
	"github.com/gin-contrib/cors"
	"github.com/topboyasante/go-snip/api/v1/models"
	"github.com/topboyasante/go-snip/api/v1/routes"
	"github.com/topboyasante/go-snip/internal/database"
	"github.com/topboyasante/go-snip/pkg/config"
)

func init() {
	database.ConnectToDB()
}

func main() {
	//Run AutoMigrations
	err := database.DB.AutoMigrate(&models.User{}, &models.Snippet{})
	if err != nil {
		log.Fatal(err)
	}

	r := gin.Default()
	r.Use(cors.Default())

	v1 := r.Group("/api/v1")
	{
		routes.AuthRoutes(v1)
		routes.SnippetRoutes(v1)
	}

	r.Run(config.ENV.ServerPort)
}
