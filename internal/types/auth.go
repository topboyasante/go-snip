package types

import "github.com/google/uuid"

type UserLoginRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type UserSignUpRequest struct {
	Username string `json:"username"`
	Email    string `json:"email"`
	Password string `json:"password"`
}

type ActivateAccountRequest struct {
	Email     string `json:"email"`
	AuthToken int    `json:"auth_token"`
}
type ForgotPasswordRequest struct {
	Email     string `json:"email"`
}

type UserResponse struct {
	Username    string    `json:"username"`
	Email       string    `json:"email"`
	ID          uuid.UUID `json:"id"`
	AccessToken string    `json:"access_token"`
	RefeshToken string    `json:"refresh_token"`
}
