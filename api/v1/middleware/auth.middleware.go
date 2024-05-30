package middleware

import (
	"fmt"
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"github.com/topboyasante/go-snip/api/v1/models"
	"github.com/topboyasante/go-snip/internal/database"
	"github.com/topboyasante/go-snip/pkg/config"
)

func RequireAuth(c *gin.Context) {
	// Get accessToken from Header
	tokenStr := c.GetHeader("Authorization")

	// Return if no accessToken was provided
	if tokenStr == "" {
		c.AbortWithStatus(http.StatusUnauthorized)
		return
	}

	tokenStr = strings.TrimPrefix(tokenStr, "Bearer ")

	// Parse the accessToken and check if the correct signing method was used
	token, err := jwt.Parse(tokenStr, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}

		return []byte(config.ENV.JWTKey), nil
	})

	if err != nil {
		log.Println(err)

		c.AbortWithStatus(http.StatusUnauthorized)
	}

	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		if float64(time.Now().Unix()) > claims["exp"].(float64) {
			c.AbortWithStatus(http.StatusUnauthorized)
		}

		var user models.User
		database.DB.Where("id = ?", claims["sub"]).First(&user)

		if user.ID == uuid.Nil {
			c.AbortWithStatus(http.StatusUnauthorized)
			return
		}

		c.Set("user_id", user.ID)
	} else {
		c.AbortWithStatus(http.StatusUnauthorized)
	}
	c.Next()
}
