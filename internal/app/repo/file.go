package repo

import (
	"bufio"
	"encoding/json"
	"fmt"
	"os"
	"shortener/internal/app/config"
	"shortener/internal/app/entity"
	"shortener/internal/app/random"
)

// FileRepository Основная структура
type FileRepository struct {
	Contractions []*entity.Contraction
	Config       *config.MainConfig
}

func CreateFileRepository(config *config.MainConfig) *FileRepository {
	fileRepo := &FileRepository{
		Contractions: []*entity.Contraction{},
		Config:       config,
	}

	err := fileRepo.init()
	if err != nil {
		panic(err)
	}

	return fileRepo
}

// GetContraction Получить ссылку на основе URL
func (r *FileRepository) GetContraction(shortURL string) (*entity.Contraction, error) {
	for _, contraction := range r.Contractions {
		if contraction.ShortURL == shortURL {
			return contraction, nil
		}
	}

	return nil, fmt.Errorf("unknown short url")
}

func (r *FileRepository) HasContraction(shortURL string) bool {
	_, err := r.GetContraction(shortURL)

	return err == nil
}

// CreateContraction Создать ссылку - пока будем хранить в мапе
func (r *FileRepository) CreateContraction(originalURL string) *entity.Contraction {
	var hash string
	for {
		hash = random.GenerateRandomString(7)
		if !r.HasContraction(hash) {
			break
		}
	}

	contraction := &entity.Contraction{
		OriginalURL: originalURL,
		ShortURL:    hash,
	}

	data, err := json.Marshal(contraction)
	if err != nil {
		panic(err)
	}

	data = append(data, []byte("\n")...)

	file, err := os.OpenFile(r.Config.FileStoragePath, os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		panic(err)
	}

	defer func() {
		err := file.Close()
		if err != nil {
			panic(err)
		}
	}()

	_, err = file.Write(data)
	if err != nil {
		panic(err)
	}

	r.Contractions = append(r.Contractions, contraction)

	return contraction
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

		contraction := entity.Contraction{}
		err = json.Unmarshal([]byte(data), &contraction)

		if err != nil {
			return err
		}

		r.Contractions = append(r.Contractions, &contraction)
	}

	return nil
}
