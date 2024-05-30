package models

import (
	"github.com/google/uuid"
	"github.com/topboyasante/go-snip/internal/database"
)

type Snippet struct {
	BaseModel
	Title       string    `json:"title"`
	Description string    `json:"description"`
	Code        string    `json:"code"`
	UserID      uuid.UUID `json:"user_id"`
	User        User      `json:"user" gorm:"foreignKey:UserID"`
}

func (snippet *Snippet) Create() (*Snippet, error) {
	err := database.DB.Create(&snippet).Error
	if err != nil {
		return &Snippet{}, err
	}
	return snippet, nil
}

func GetSnippets(userId uuid.UUID) ([]Snippet, error) {
	var snippets []Snippet
	err := database.DB.Preload("User").Where("user_id = ?", userId).Find(&snippets).Error
	if err != nil {
		return []Snippet{}, err
	}
	return snippets, nil
}

func GetSnippet(userId uuid.UUID, snippetId string) (Snippet, error) {
	var snippet Snippet
	err := database.DB.Where("user_id = ?", userId).Where("id = ?", snippetId).Find(&snippet).Error
	if err != nil {
		return Snippet{}, err
	}
	return snippet, nil
}

func (snippet *Snippet) Delete() error {
	err := database.DB.Where("id = ?", snippet.ID).Delete(&snippet).Error
	if err != nil {
		return err
	}
	return nil
}
