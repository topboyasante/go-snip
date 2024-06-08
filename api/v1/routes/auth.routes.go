package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/topboyasante/go-snip/api/v1/controllers"
)

func AuthRoutes(r *gin.RouterGroup) {
	authRoutes := r.Group("/auth")
	authRoutes.POST("/sign-in/", controllers.SignIn)
	authRoutes.POST("/sign-up/", controllers.SignUp)
	authRoutes.POST("/activate-account", controllers.ActivateAccount)
	authRoutes.POST("/forgot-password", controllers.ForgotPassword)
	authRoutes.POST("/reset-password", controllers.ResetPassword)
}
