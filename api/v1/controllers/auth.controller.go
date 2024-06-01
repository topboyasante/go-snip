package controllers

import (
	"fmt"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"github.com/topboyasante/go-snip/api/v1/models"
	"github.com/topboyasante/go-snip/internal/database"
	"github.com/topboyasante/go-snip/internal/types"
	"github.com/topboyasante/go-snip/pkg/auth"
	"github.com/topboyasante/go-snip/pkg/config"
	"github.com/topboyasante/go-snip/pkg/email"
	"github.com/topboyasante/go-snip/pkg/validators"
	"golang.org/x/crypto/bcrypt"
)

func SignIn(c *gin.Context) {
	// Create a struct to hold the request body
	var body struct {
		Username string `json:"username"`
		Password string `json:"password"`
	}

	// Parse the request body and store it in the body struct
	if c.Bind(&body) != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "failed to read request body",
		})
		return
	}

	// Validations to check for empty fields
	if !validators.NotBlank(body.Username) || !validators.NotBlank(body.Password) {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "some fields are empty",
		})
		return
	}

	// Find the user with the provided email
	user, err := models.GetUserByUsername(body.Username)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "user does not exist",
		})
		return
	}

	// Return if the account has not been activated
	if !user.IsActive {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "account is not activated",
		})
		return
	}

	err = user.VerifyPassword(body.Password)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "invalid password",
		})
		return
	}

	access_token, refesh_token, err := auth.CreateJWTTokens(user.ID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "unable to create accessToken",
		})
		return
	}

	// Generate a new Auth Token on Sign up
	user.AuthToken = auth.GenerateAuthToken()
	database.DB.Save(&user)

	// Response to be sent to the user
	userRes := &types.UserResponse{
		Username:    user.Username,
		Email:       user.Email,
		ID:          user.ID,
		AccessToken: access_token,
		RefeshToken: refesh_token,
	}

	c.JSON(http.StatusOK, gin.H{
		"success": "logged in",
		"user":    userRes,
	})
}

func SignUp(c *gin.Context) {
	// Create a struct to hold the request body
	var body struct {
		Username string `json:"username"`
		Email    string `json:"email"`
		Password string `json:"password"`
	}

	// Parse the request body and store it in the body struct
	if c.Bind(&body) != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "failed to read request body",
		})
		return
	}

	// Validations to check for empty fields
	if !validators.NotBlank(body.Username) || !validators.NotBlank(body.Email) || !validators.NotBlank(body.Password) {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "some fields are empty",
		})
		return
	}

	// Validations to check for a correct email
	if !validators.Matches(body.Email, validators.EmailRX) {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "invalid email",
		})
		return
	}

	//Check if a user exists with that email or username
	if !auth.IsEmailUnique(body.Email) {
		c.JSON(http.StatusBadRequest, gin.H{"error": "user with provided email exists"})
		return
	}
	if !auth.IsUsernameUnique(body.Username) {
		c.JSON(http.StatusBadRequest, gin.H{"error": "user with provided username exists"})
		return
	}

	// Hash the password from the request body
	hash, err := bcrypt.GenerateFromPassword([]byte(body.Password), 10)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "failed to hash password",
		})
		return
	}

	// Create a new models.User instance which contains the values from the request body, and the hashed password
	user := models.User{
		Username: body.Username,
		Email:    body.Email,
		Password: string(hash),
	}

	// set the user ID to a new UUID
	user.ID = uuid.New()

	//Generate a new auth token
	user.AuthToken = auth.GenerateAuthToken()

	// All new accounts are not active by default
	user.IsActive = false

	// Send an email to the user to activate their account
	email.SendMailWithSMTP(
		email.EmailConfig,
		"Activate your account",
		"web/activate-account.html",
		struct {
			Name      string
			AuthToken int
		}{Name: user.Username, AuthToken: user.AuthToken},
		[]string{body.Email},
	)

	// Insert the user in the DB
	newUser, err := user.Create()
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "failed to create user",
		})
		return
	}

	// Return a 200 status when everything was successful
	c.JSON(http.StatusOK, gin.H{
		"success": "account created. please activate your account",
		"user":    newUser,
	})
}

func ActivateAccount(c *gin.Context) {
	// Create an instance of models.User to hold the existing user data
	var user models.User

	// Create a struct to hold the request body
	var body struct {
		Email     string `json:"email"`
		AuthToken int    `json:"auth_token"`
	}

	// Parse the request body and store it in the body struct
	if c.Bind(&body) != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "failed to read request body",
		})
		return
	}

	// Validations to check for empty fields
	if !validators.NotBlank(body.Email) || !validators.NotZero(body.AuthToken) {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "some fields are empty",
		})
		return
	}

	// Validations to check for a correct email
	if !validators.Matches(body.Email, validators.EmailRX) {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "invalid email",
		})
		return
	}

	// Find the user with the provided email and store the user details in the user variable
	user, err := models.GetUserByEmail(body.Email)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "user does not exist",
		})
		return
	}

	// Check if the token is valid
	if body.AuthToken != user.AuthToken {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "token is invalid",
		})
		return
	}

	// Return if the account has already been activated, and activate the user account if it has not
	if user.IsActive {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "account has already been activated",
		})
		return
	}

	// Activate user account
	user.IsActive = true

	// Generate a new auth token on account activation
	user.AuthToken = auth.GenerateAuthToken()
	database.DB.Save(&user)

	c.JSON(http.StatusOK, gin.H{
		"success": "account has been activated",
	})
}

func ForgotPassword(c *gin.Context) {
	var user models.User

	// Create a struct to hold the request body
	var body struct {
		Email string `json:"email"`
	}

	// Parse the request body and store it in the body struct
	if c.Bind(&body) != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "failed to read request body",
		})
		return
	}

	// Validations to check for empty fields
	if !validators.NotBlank(body.Email) {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "the email field is empty",
		})
		return
	}

	// Validations to check for a correct email
	if !validators.Matches(body.Email, validators.EmailRX) {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "invalid email",
		})
		return
	}

	// Find the user with the provided email and store the user details in the user variable
	// Find the user with the provided email
	user, err := models.GetUserByEmail(body.Email)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "user does not exist",
		})
		return
	}

	// Send an email with the auth token to the user
	email.SendMailWithSMTP(
		email.EmailConfig,
		"Activate your account",
		"web/reset-password.html",
		struct {
			Name      string
			AuthToken int
		}{Name: user.Username, AuthToken: user.AuthToken},
		[]string{body.Email},
	)

	c.JSON(http.StatusOK, gin.H{
		"success": "reset code sent",
	})
}

func ResetPassword(c *gin.Context) {
	// Create an instance of models.User to hold the existing user data
	var user models.User

	// Create a struct to hold the request body
	var body struct {
		Email       string `json:"email"`
		AuthToken   int    `json:"auth_token"`
		NewPassword string `json:"new_password"`
	}

	// Parse the request body and store it in the body struct
	if c.Bind(&body) != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "failed to read request body",
		})
		return
	}

	// Validations to check for empty fields
	if !validators.NotBlank(body.Email) || !validators.NotZero(body.AuthToken) || !validators.NotBlank(body.NewPassword) {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "some fields are empty",
		})
		return
	}

	// Validations to check for a correct email
	if !validators.Matches(body.Email, validators.EmailRX) {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "invalid email",
		})
		return
	}

	// Check if the token is valid
	if body.AuthToken != user.AuthToken {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "token is invalid",
		})
		return
	}

	hash, err := bcrypt.GenerateFromPassword([]byte(body.NewPassword), 10)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "failed to hash password",
		})
		return
	}

	user.Password = string(hash)
	database.DB.Save(&user)

	c.JSON(http.StatusOK, gin.H{
		"success": "password has been reset",
	})
}

func RefreshAccessToken(c *gin.Context) {
	var body struct {
		RefreshToken string `json:"refresh_token"`
	}

	if err := c.Bind(&body); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "failed to read request body"})
		return
	}

	token, err := jwt.Parse(body.RefreshToken, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(config.ENV.JWTKey), nil
	})

	if err != nil || !token.Valid {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid refresh token"})
		return
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok || float64(time.Now().Unix()) > claims["exp"].(float64) {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "expired refresh token"})
		return
	}

	userID := claims["sub"].(string)
	newAccessToken, newRefreshToken, err := auth.CreateJWTTokens(userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate tokens"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"access_token":  newAccessToken,
		"refresh_token": newRefreshToken,
	})
}
