package storage

import (
	"fmt"
	"io"
	"log"
	"mime/multipart"
	"os"
	"path/filepath"
	"strings"

	"github.com/badersalis/gidana_backend/internal/config"
	"github.com/google/uuid"
)

var allowedImageExts = map[string]bool{
	".jpg": true, ".jpeg": true, ".png": true, ".gif": true, ".webp": true,
}

type FileStorage interface {
	Save(fh *multipart.FileHeader) (string, error)
	Delete(url string) error
}

type fileStorage struct{}

func NewFileStorage() FileStorage { return &fileStorage{} }

func (fs *fileStorage) Save(fh *multipart.FileHeader) (string, error) {
	ext := strings.ToLower(filepath.Ext(fh.Filename))
	if !allowedImageExts[ext] {
		return "", fmt.Errorf("invalid file type: only jpg, jpeg, png, gif, webp are allowed")
	}

	if config.App.SupabaseURL != "" {
		return UploadFile(fh, "properties")
	}

	filename := uuid.New().String() + ext
	savePath := filepath.Join(config.App.UploadDir, filename)

	src, err := fh.Open()
	if err != nil {
		return "", fmt.Errorf("failed to open file: %w", err)
	}
	defer src.Close()

	data, err := io.ReadAll(src)
	if err != nil {
		return "", fmt.Errorf("failed to read file: %w", err)
	}

	if err := os.WriteFile(savePath, data, 0644); err != nil {
		return "", fmt.Errorf("failed to save file locally: %w", err)
	}
	return "/uploads/properties/" + filename, nil
}

func (fs *fileStorage) Delete(url string) error {
	if url == "" {
		return nil
	}
	if config.App.SupabaseURL != "" {
		if err := DeleteFile(url); err != nil {
			log.Printf("storage: failed to delete %s: %v", url, err)
		}
		return nil
	}
	localPath := filepath.Join(".", url)
	if err := os.Remove(localPath); err != nil && !os.IsNotExist(err) {
		log.Printf("storage: failed to delete local file %s: %v", localPath, err)
	}
	return nil
}
