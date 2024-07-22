package entity

import (
	"github.com/google/uuid"
	"time"
)

//go:generate easyjson -disallow_unknown_fields -all ./shortcut.go

type Shortcut struct {
	ID          uuid.UUID  `json:"uuid"`
	UserID      uuid.UUID  `json:"user_id"`
	ShortURL    string     `json:"short_url"`
	OriginalURL string     `json:"original_url"`
	CreatedDate *time.Time `json:"-"`
	DeletedDate *time.Time `json:"deleted_date"`
}
