package repo

import (
	"bufio"
	"fmt"
	"github.com/mailru/easyjson"
	"os"
	"shortener/internal/app/config"
	"shortener/internal/app/entity"
	"shortener/internal/app/random"
)

// FileRepository Основная структура
type FileRepository struct {
	Shortcuts map[string]entity.Shortcut
	Config    *config.MainConfig
}

func CreateFileRepository(config *config.MainConfig) (*FileRepository, error) {
	fileRepo := &FileRepository{
		Shortcuts: map[string]entity.Shortcut{},
		Config:    config,
	}

	err := fileRepo.init()
	if err != nil {
		return nil, err
	}

	return fileRepo, nil
}

// GetShortcut Получить ссылку на основе URL
func (r *FileRepository) GetShortcut(shortURL string) (*entity.Shortcut, error) {
	shortcut, ok := r.Shortcuts[shortURL]

	if ok {
		return &shortcut, nil
	}

	return nil, fmt.Errorf("unknown short url")
}

func (r *FileRepository) HasShortcut(shortURL string) bool {
	_, err := r.GetShortcut(shortURL)

	return err == nil
}

// CreateShortcut Создать ссылку - пока будем хранить в мапе
func (r *FileRepository) CreateShortcut(originalURL string) (*entity.Shortcut, error) {
	var hash string
	for {
		hash = random.GenerateRandomString(7)
		if !r.HasShortcut(hash) {
			break
		}
	}

	shortcut := entity.Shortcut{
		OriginalURL: originalURL,
		ShortURL:    hash,
	}

	data, err := easyjson.Marshal(shortcut)

	if err != nil {
		return nil, err
	}

	data = append(data, []byte("\n")...)

	file, err := os.OpenFile(r.Config.FileStoragePath, os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		return nil, err
	}

	defer func() {
		err := file.Close()
		if err != nil {
			panic(err)
		}
	}()

	_, err = file.Write(data)
	if err != nil {
		return nil, err
	}

	r.Shortcuts[hash] = shortcut

	return &shortcut, nil
}

func (r *FileRepository) init() error {
	file, err := os.OpenFile(r.Config.FileStoragePath, os.O_RDONLY|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		return err
	}

	defer func() {
		err := file.Close()
		if err != nil {
			panic(err)
		}
	}()

	sc := bufio.NewScanner(file)
	for sc.Scan() {
		data := sc.Text()

		shortcut := entity.Shortcut{}
		err = easyjson.Unmarshal([]byte(data), &shortcut)

		if err != nil {
			return err
		}

		r.Shortcuts[shortcut.ShortURL] = shortcut
	}

	return nil
}
