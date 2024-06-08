package types

import (
	"time"

	"github.com/google/uuid"
)

type NewSnippetRequest struct {
	Title       string `json:"title"`
	Description string `json:"description"`
	Code        string `json:"code"`
}
type NewSnippetResponse struct {
	ID          uuid.UUID `json:"id"`
	Title       string    `json:"title"`
	Description string    `json:"description"`
	Code        string    `json:"code"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
	UserID      uuid.UUID `json:"user_id"`
	CreatedBy   string    `json:"created_by"`
}
