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

// Get All Snippets godoc
//
//	@Summary		Get Snippets
//	@Description	Get all snippets
//	@Tags			Snippets
//	@Accept			json
//	@Produce		json
//	@Success		200	{object}	types.APISuccessMessage
//	@Failure		400	{object}	types.APIErrorMessage
//	@Failure		500	{object}	types.APIErrorMessage
//	@Router			/snippets [get]
func GetSnippets(c *gin.Context) {
	snippets, err := models.GetSnippets()
	if err != nil {
		c.JSON(http.StatusBadRequest, types.APIErrorMessage{
			ErrorMessage: "could not retrieve snippets",
		})
		return
	}

	c.JSON(200, types.APISuccessMessage{
		Data: snippets,
	})
}

// Get One Snippets godoc
//
//	@Summary		Get Snippet
//	@Description	Get a snippet
//	@Tags			Snippets
//	@Accept			json
//	@Produce		json
//	@Param			id	path		string	true	"Snippet ID"
//	@Success		200	{object}	types.APISuccessMessage
//	@Failure		400	{object}	types.APIErrorMessage
//	@Failure		500	{object}	types.APIErrorMessage
//	@Router			/snippets/{id} [get]
func GetSnippet(c *gin.Context) {
	id := c.Param("id")

	snippet, err := models.GetSnippet(id)
	if err != nil {
		c.JSON(http.StatusBadRequest, types.APIErrorMessage{
			ErrorMessage: "could not retrieve snippet",
		})
		return
	}

	c.JSON(200, types.APISuccessMessage{
		Data: snippet,
	})
}

// Create Snippet godoc
//
//	@Summary		Create a Snippet
//	@Description	Create a Snippet
//	@Tags			Snippets
//	@Accept			json
//	@Produce		json
//	@Security		ApiKeyAuth
//	@Param			Snippet	body		types.NewSnippetRequest	true	"snippet"
//	@Success		200		{object}	types.APISuccessMessage
//	@Failure		400		{object}	types.APIErrorMessage
//	@Failure		500		{object}	types.APIErrorMessage
//	@Router			/snippets/create [post]
func CreateSnippet(c *gin.Context) {
	var body struct {
		Title       string `json:"title"`
		Description string `json:"description"`
		Code        string `json:"code"`
	}

	if c.Bind(&body) != nil {
		c.JSON(http.StatusBadRequest, types.APIErrorMessage{
			ErrorMessage: "failed to read request body",
		})
		return
	}

	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, types.APIErrorMessage{ErrorMessage: "unauthorized request"})
		return
	}

	// Type Assertion
	uID, ok := userID.(uuid.UUID)
	if !ok {
		c.JSON(http.StatusUnauthorized, types.APIErrorMessage{ErrorMessage: "unauthorized request"})
		return
	}

	if !validators.NotBlank(body.Title) ||
		!validators.NotBlank(body.Code) {
		c.JSON(http.StatusBadRequest, types.APIErrorMessage{
			ErrorMessage: "some fields are empty",
		})
		return
	}

	user, err := models.GetUserById(uID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			c.JSON(http.StatusBadRequest, types.APIErrorMessage{
				ErrorMessage: "no user exists with the provided user ID",
			})
			return
		}
		return
	}

	if !user.IsActive {
		c.JSON(http.StatusBadRequest, types.APIErrorMessage{
			ErrorMessage: "account is not activated",
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
		c.JSON(http.StatusBadRequest, types.APIErrorMessage{
			ErrorMessage: "unable to create snippet",
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

	c.JSON(http.StatusOK, types.APISuccessMessage{
		Data: snippetRes,
	})
}


// Delete snippet godoc
//
//	@Summary		Delete snippet
//	@Description	Delete snippet
//	@Tags			Snippets
//	@Accept			json
//	@Produce		json
//	@Security		ApiKeyAuth
//	@Param			id	path		string	true	"Snippet ID"
//	@Success		200		{object}	types.APISuccessMessage
//	@Failure		400		{object}	types.APIErrorMessage
//	@Failure		500		{object}	types.APIErrorMessage
//	@Router			/snippets/{id} [delete]
func DeleteSnippet(c *gin.Context) {
	id := c.Param("id")

	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, types.APIErrorMessage{ErrorMessage: "unauthorized request"})
		return
	}

	// Type Assertion
	uID, ok := userID.(uuid.UUID)
	if !ok {
		c.JSON(http.StatusUnauthorized, types.APIErrorMessage{ErrorMessage: "unauthorized request"})
		return
	}

	user, err := models.GetUserById(uID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			c.JSON(http.StatusBadRequest, types.APIErrorMessage{
				ErrorMessage: "no user exists with the provided user ID",
			})
			return
		}
		return
	}

	if !user.IsActive {
		c.JSON(http.StatusBadRequest, types.APIErrorMessage{
			ErrorMessage: "account is not activated",
		})
		return
	}

	snippet, err := models.GetSnippetWithUser(user.ID, id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			c.JSON(http.StatusBadRequest, types.APIErrorMessage{
				ErrorMessage: "no snippet exists with the provided queries",
			})
			return
		}
		return
	}

	if snippet.UserID != uID {
		c.JSON(http.StatusUnauthorized, types.APIErrorMessage{ErrorMessage: "you are not authorized to delete this snippet"})
		return
	}

	err = snippet.Delete()
	if err != nil {
		c.JSON(http.StatusInternalServerError, types.APIErrorMessage{
			ErrorMessage: "unable to delete snippet",
		})
		return
	}

	c.JSON(http.StatusOK, types.APISuccessMessage{
		Data: "snippet deleted",
	})
}

// Update Snippet godoc
//
//	@Summary		Update Snippet
//	@Description	Update Snippet
//	@Tags			Snippets
//	@Accept			json
//	@Produce		json
//	@Security		ApiKeyAuth
//	@Param			id	path		string	true	"Snippet ID"
//	@Param			Snippet	body		types.NewSnippetRequest	true	"snippet"
//	@Success		200		{object}	types.APISuccessMessage
//	@Failure		400		{object}	types.APIErrorMessage
//	@Failure		500		{object}	types.APIErrorMessage
//	@Router			/snippets/{id} [put]
func UpdateSnippet(c *gin.Context) {
	var body struct {
		Title       string `json:"title"`
		Description string `json:"description"`
		Code        string `json:"code"`
	}

	c.Bind(&body)

	id := c.Param("id")

	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, types.APIErrorMessage{ErrorMessage: "unauthorized request"})
		return
	}

	// Type Assertion
	uID, ok := userID.(uuid.UUID)
	if !ok {
		c.JSON(http.StatusUnauthorized, types.APIErrorMessage{ErrorMessage: "unauthorized request"})
		return
	}

	snippet, err := models.GetSnippetWithUser(uID, id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			c.JSON(http.StatusBadRequest, types.APIErrorMessage{
				ErrorMessage: "no snippet exists with the provided queries",
			})
			return
		}
		return
	}

	if snippet.UserID != uID {
		c.JSON(http.StatusUnauthorized, types.APIErrorMessage{ErrorMessage: "you are not authorized to update this snippet"})
		return
	}

	err = snippet.Update(body.Title, body.Description, body.Code)
	if err != nil {
		c.JSON(http.StatusInternalServerError, types.APIErrorMessage{
			ErrorMessage: "unable to delete snippet",
		})
		return
	}

	c.JSON(200, types.APISuccessMessage{
		Data: "snippet updated",
	})
}
