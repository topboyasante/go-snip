package main

import (
	"log"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/topboyasante/go-snip/api/v1/models"
	"github.com/topboyasante/go-snip/api/v1/routes"
	"github.com/topboyasante/go-snip/internal/database"
	"github.com/topboyasante/go-snip/pkg/config"

	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
	_ "github.com/topboyasante/go-snip/docs"
)

//	@title			Go-Snip API
//	@version		1.0
//	@description	API Documentation for Go-Snip.
//	@termsOfService	http://swagger.io/terms/

//	@contact.name	Nana Kwasi Asante
//	@contact.url	https://www.nkasante.com
//	@contact.email	asantekwasi101@gmail.com

//	@license.name	Apache 2.0
//	@license.url	http://www.apache.org/licenses/LICENSE-2.0.html

//	@host		localhost:4000
//	@BasePath	/api/v1


//	@securityDefinitions.apikey	JWT
//	@in							header
//	@name						Authorization

//	@externalDocs.description	OpenAPI
//	@externalDocs.url			https://swagger.io/resources/open-api/

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

	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	r.Run(config.ENV.ServerPort)
}
