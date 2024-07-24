package mapper

import (
	"shortener/internal/app/entity"
	"shortener/internal/app/handler/inout"
	"shortener/internal/app/human"
)

func GetBaseShortcutListFromShortcuts(shortcuts []entity.Shortcut) inout.BaseShortcutList {
	list := inout.BaseShortcutList{}

	if shortcuts == nil {
		return list
	}

	for _, shortcut := range shortcuts {
		baseShortcut := inout.BaseShortcut{
			OriginalURL: shortcut.OriginalURL,
			ShortURL:    human.GetFullShortURL(&shortcut),
		}

		list = append(list, baseShortcut)
	}

	return list
}
