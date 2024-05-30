package auth

import (
	"math/rand"

	"github.com/topboyasante/go-snip/api/v1/models"
	"github.com/topboyasante/go-snip/internal/database"
	"gorm.io/gorm"
)

func GenerateAuthToken() int {
	// Generate a random number between 1000 and 9999
	randomNumber := rand.Intn(9000) + 1000

	return randomNumber
}

func IsEmailUnique(email string) bool {
	var user models.User
	result := database.DB.Where("email = ?", email).First(&user)
	return result.Error == gorm.ErrRecordNotFound
}

func IsUsernameUnique(username string) bool {
	var user models.User
	result := database.DB.Where("username = ?", username).First(&user)
	return result.Error == gorm.ErrRecordNotFound
}
