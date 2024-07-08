package human

import (
	"shortener/internal/app/config"
	"shortener/internal/app/entity"
)

func GetFullShortURL(shortcut *entity.Shortcut) string {
	return config.State().DisplayLink + "/" + shortcut.ShortURL
}
