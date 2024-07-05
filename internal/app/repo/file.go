package repo

import (
	"bufio"
	"context"
	"fmt"
	"github.com/mailru/easyjson"
	"os"
	"shortener/internal/app/config"
	"shortener/internal/app/entity"
	"shortener/internal/app/logging"
	"shortener/internal/app/random"
)

// FileRepository репозиторий на основе хранения в файле
type FileRepository struct {
	Shortcuts map[string]entity.Shortcut
}

func CreateFileRepository() (*FileRepository, error) {
	fileRepo := &FileRepository{
		Shortcuts: map[string]entity.Shortcut{},
	}

	err := fileRepo.init()
	if err != nil {
		return nil, err
	}

	return fileRepo, nil
}

// GetShortcut Получить ссылку на основе URL
func (r *FileRepository) GetShortcut(ctx context.Context, shortURL string) (*entity.Shortcut, error) {
	shortcut, ok := r.Shortcuts[shortURL]

	if ok {
		return &shortcut, nil
	}

	return nil, fmt.Errorf("unknown short url")
}

// CreateShortcut Создать ссылку - пока будем хранить в мапе
func (r *FileRepository) CreateShortcut(ctx context.Context, originalURL string) (*entity.Shortcut, error) {
	var hash string
	for {
		hash = random.GenerateRandomString(7)
		if !HasShortcut(ctx, r, hash) {
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

	file, err := os.OpenFile(config.State().FileStoragePath, os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		return nil, err
	}

	defer func() {
		err := file.Close()
		if err != nil {
			logging.Sugar.Fatal(err)
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
	file, err := os.OpenFile(config.State().FileStoragePath, os.O_RDONLY|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		return err
	}

	defer func() {
		err := file.Close()
		if err != nil {
			logging.Sugar.Fatal(err)
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
