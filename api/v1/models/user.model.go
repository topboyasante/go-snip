package models

import (
	"github.com/google/uuid"
	"github.com/topboyasante/go-snip/internal/database"
	"golang.org/x/crypto/bcrypt"
)

type User struct {
	BaseModel
	Username  string `json:"username" gorm:"unique"`
	Email     string `json:"email" gorm:"unique"`
	Password  string `json:"password"`
	IsActive  bool   `json:"is_active"`
	AuthToken int    `json:"auth_token"` //use this token for restting passwords, and verifying accounts
}

func (user *User) Create() (*User, error) {
	err := database.DB.Create(&user).Error
	if err != nil {
		return &User{}, err
	}
	return user, nil
}

func (user *User) VerifyPassword(pw string) error {
	return bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(pw))
}

func GetUserByEmail(email string) (User, error) {
	var user User
	err := database.DB.Where("email", email).Find(&user).Error
	if err != nil {
		return User{}, err
	}
	return user, nil
}

func GetUserByUsername(username string) (User, error) {
	var user User
	err := database.DB.Where("username", username).Find(&user).Error
	if err != nil {
		return User{}, err
	}
	return user, nil
}

func GetUserById(id uuid.UUID) (User, error) {
	var user User
	err := database.DB.Where("ID=?", id).Find(&user).Error
	if err != nil {
		return User{}, err
	}
	return user, nil
}

func HashPassword(pw string) ([]byte, error) {
	hash, err := bcrypt.GenerateFromPassword([]byte(pw), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}
	return hash, nil
}

func (user *User) UpdatePassword(input string) (*User, error) {
	hash, err := bcrypt.GenerateFromPassword([]byte(input), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}
	if err := database.DB.Model(&user).Update("password", hash); err != nil {
		return &User{}, nil
	}
	return user, nil
}
