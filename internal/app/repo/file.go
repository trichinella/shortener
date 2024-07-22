package repo

import (
	"bufio"
	"context"
	"fmt"
	"github.com/google/uuid"
	"github.com/mailru/easyjson"
	"os"
	"shortener/internal/app/config"
	"shortener/internal/app/entity"
	"shortener/internal/app/handler/inout"
	"shortener/internal/app/human"
	"shortener/internal/app/logging"
	"shortener/internal/app/random"
	"shortener/internal/app/service/authentification"
	"time"
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

// GetShortcutByShortURL Получить ссылку на основе URL
func (r *FileRepository) GetShortcutByShortURL(ctx context.Context, shortURL string) (*entity.Shortcut, error) {
	shortcut, ok := r.Shortcuts[shortURL]

	if ok {
		return &shortcut, nil
	}

	return nil, fmt.Errorf("unknown short url")
}

// GetShortcutByOriginalURL Получить ссылку на основе URL
func (r *FileRepository) GetShortcutByOriginalURL(ctx context.Context, originalURL string) (*entity.Shortcut, error) {
	for _, shortcut := range r.Shortcuts {
		if shortcut.OriginalURL == originalURL {
			return &shortcut, nil
		}
	}

	return nil, nil
}

// CreateShortcut Создать ссылку - пока будем хранить в мапе
func (r *FileRepository) CreateShortcut(ctx context.Context, originalURL string) (*entity.Shortcut, error) {
	shortcut, err := r.GetShortcutByOriginalURL(ctx, originalURL)

	if shortcut != nil {
		return shortcut, NewDuplicateShortcutError(err, shortcut)
	}

	var hash string
	for {
		hash = random.GenerateRandomString(7)
		if !HasShortcut(ctx, r, hash) {
			break
		}
	}

	id, err := uuid.NewUUID()
	if err != nil {
		return nil, err
	}

	userID, ok := ctx.Value(authentification.ContextUserID).(uuid.UUID)
	if !ok {
		userID = uuid.Nil
	}
	shortcut = &entity.Shortcut{
		ID:          id,
		OriginalURL: originalURL,
		ShortURL:    hash,
		UserID:      userID,
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

	r.Shortcuts[hash] = *shortcut

	return shortcut, nil
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

func (r *FileRepository) CreateBatch(ctx context.Context, batchInput inout.ExternalBatchInput) (result inout.ExternalBatchOutput, err error) {
	//нормальное поведение
	if len(batchInput) == 0 {
		return result, nil
	}

	for _, input := range batchInput {
		shortcut, err := r.CreateShortcut(ctx, input.OriginalURL)
		if err != nil {
			return nil, err
		}

		externalOutput := inout.ExternalOutput{}
		externalOutput.ExternalID = input.ExternalID
		externalOutput.ShortURL = human.GetFullShortURL(shortcut)

		result = append(result, externalOutput)
	}

	return result, nil
}

func (r *FileRepository) GetShortcutsByUserID(ctx context.Context, userID uuid.UUID) ([]entity.Shortcut, error) {
	var shortcuts []entity.Shortcut
	for _, shortcut := range r.Shortcuts {
		if shortcut.UserID == userID {
			shortcuts = append(shortcuts, shortcut)
		}
	}

	return shortcuts, nil
}

func (r *FileRepository) DeleteList(ctx context.Context, userID uuid.UUID, list inout.ShortURLList) error {
	for _, shortURL := range list {
		shortcut, err := r.GetShortcutByShortURL(ctx, shortURL)
		if err != nil {
			continue
		}
		if userID != shortcut.UserID {
			continue
		}

		now := time.Now()
		shortcut.DeletedDate = &now
		r.Shortcuts[shortURL] = *shortcut
	}

	return nil
}
