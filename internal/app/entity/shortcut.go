package entity

import (
	"github.com/google/uuid"
	"time"
)

type Shortcut struct {
	ID          uuid.UUID `json:"uuid"`
	ShortURL    string    `json:"short_url"`
	OriginalURL string    `json:"original_url"`
	CreatedDate time.Time `json:"-"`
}
