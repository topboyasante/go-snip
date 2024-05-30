package controllers

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/topboyasante/go-snip/api/v1/models"
	"github.com/topboyasante/go-snip/internal/types"
	"github.com/topboyasante/go-snip/pkg/validators"
	"gorm.io/gorm"
)

func GetSnippets(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized request"})
		return
	}

	// Type Assertion
	uID, ok := userID.(uuid.UUID)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized request"})
		return
	}

	user, err := models.GetUserById(uID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": "no user exists with the provided user ID",
			})
			return
		}
		return
	}

	if !user.IsActive {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "account is not activated",
		})
		return
	}

	snippets, err := models.GetSnippets(user.ID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "could not retrieve snippets",
		})
		return
	}

	c.JSON(200, gin.H{
		"snippets": snippets,
	})
}

func GetSnippet(c *gin.Context) {
	id := c.Param("id")

	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized request"})
		return
	}

	// Type Assertion
	uID, ok := userID.(uuid.UUID)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized request"})
		return
	}

	user, err := models.GetUserById(uID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": "no user exists with the provided user ID",
			})
			return
		}
		return
	}

	if !user.IsActive {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "account is not activated",
		})
		return
	}

	snippet, err := models.GetSnippet(user.ID, id)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "could not retrieve snippet",
		})
		return
	}

	c.JSON(200, gin.H{
		"snippet": snippet,
	})
}

func CreateSnippet(c *gin.Context) {
	var body struct {
		Title       string `json:"title"`
		Description string `json:"description"`
		Code        string `json:"code"`
	}

	if c.Bind(&body) != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "failed to read request body",
		})
		return
	}

	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized request"})
		return
	}

	// Type Assertion
	uID, ok := userID.(uuid.UUID)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized request"})
		return
	}

	if !validators.NotBlank(body.Title) ||
		!validators.NotBlank(body.Code) {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "some fields are empty",
		})
		return
	}

	user, err := models.GetUserById(uID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": "no user exists with the provided user ID",
			})
			return
		}
		return
	}

	if !user.IsActive {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "account is not activated",
		})
		return
	}

	newSnippet := &models.Snippet{
		Title:       body.Title,
		Description: body.Description,
		Code:        body.Code,
		UserID:      user.ID,
		User:        user,
	}

	newSnippet.ID = uuid.New()
	res, err := newSnippet.Create()
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "unable to create snippet",
		})
		return
	}

	snippetRes := &types.NewSnippetResponse{
		ID:          res.ID,
		Title:       res.Title,
		Description: res.Description,
		Code:        res.Code,
		UserID:      res.UserID,
		CreatedAt:   res.CreatedAt,
		UpdatedAt:   res.UpdatedAt,
		CreatedBy:   res.User.Username,
	}

	c.JSON(http.StatusOK, gin.H{
		"snippet": snippetRes,
	})
}

func DeleteSnippet(c *gin.Context) {
	id := c.Param("id")

	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized request"})
		return
	}

	// Type Assertion
	uID, ok := userID.(uuid.UUID)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized request"})
		return
	}

	user, err := models.GetUserById(uID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": "no user exists with the provided user ID",
			})
			return
		}
		return
	}

	if !user.IsActive {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "account is not activated",
		})
		return
	}

	snippet, err := models.GetSnippet(user.ID, id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": "no snippet exists with the provided queries",
			})
			return
		}
		return
	}

	if snippet.UserID != uID {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "you are not authorized to delete this snippet"})
		return
	}

	err = snippet.Delete()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "unable to delete snippet",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"successs": "snippet deleted",
	})
}
