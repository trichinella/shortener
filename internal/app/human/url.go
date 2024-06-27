package human

import (
	"shortener/internal/app/config"
	"shortener/internal/app/entity"
)

func GetFullShortUrl(cfg *config.MainConfig, contraction *entity.Contraction) string {
	return cfg.DisplayLink + "/" + contraction.ShortUrl
}
