package types

import "github.com/google/uuid"

type UserResponse struct {
	Username    string    `json:"username"`
	Email       string    `json:"email"`
	ID          uuid.UUID `json:"id"`
	AccessToken string    `json:"access_token"`
	RefeshToken string    `json:"refresh_token"`
}
