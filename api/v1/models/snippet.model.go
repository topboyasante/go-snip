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

func GetSnippets() ([]Snippet, error) {
	var snippets []Snippet
	err := database.DB.Preload("User").Find(&snippets).Error
	if err != nil {
		return []Snippet{}, err
	}
	return snippets, nil
}

func GetSnippetsWithUser(userId uuid.UUID) ([]Snippet, error) {
	var snippets []Snippet
	err := database.DB.Preload("User").Where("user_id = ?", userId).Find(&snippets).Error
	if err != nil {
		return []Snippet{}, err
	}
	return snippets, nil
}

func GetSnippet(snippetId string) (Snippet, error) {
	var snippet Snippet
	err :=  database.DB.Preload("User").Where("id = ?", snippetId).Find(&snippet).Error
	if err != nil {
		return Snippet{}, err
	}
	return snippet, nil
}

func GetSnippetWithUser(userId uuid.UUID, snippetId string) (Snippet, error) {
	var snippet Snippet
	err :=  database.DB.Preload("User").Where("user_id = ?", userId).Where("id = ?", snippetId).Find(&snippet).Error
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

func (snippet *Snippet) Update(title, description, code string) error {
	err := database.DB.Model(&snippet).Updates(Snippet{
		Title:       title,
		Description: description,
		Code:        code,
	}).Error
	if err != nil {
		return err
	}
	return nil
}
