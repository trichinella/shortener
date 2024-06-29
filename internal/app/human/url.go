package human

import (
	"shortener/internal/app/config"
	"shortener/internal/app/entity"
)

func GetFullShortURL(cfg *config.MainConfig, shortcut *entity.Shortcut) string {
	return cfg.DisplayLink + "/" + shortcut.ShortURL
}
