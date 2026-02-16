package uploader

import (
	"bbs-go/internal/models/dto"
	"os"
	"path/filepath"
	"strings"
)

type LocalUploader struct{}

func (u *LocalUploader) PutImage(cfg dto.UploadConfig, data []byte, contentType string) (string, error) {
	if strings.TrimSpace(contentType) == "" {
		contentType = "image/jpeg"
	}
	key := generateImageKeyByPrefix("", data, contentType)
	return u.PutObject(cfg, key, data, contentType)
}

func (u *LocalUploader) PutObject(_ dto.UploadConfig, key string, data []byte, _ string) (string, error) {
	cleanKey := strings.TrimPrefix(filepath.ToSlash(filepath.Clean(key)), "/")
	relativePath := filepath.FromSlash(cleanKey)
	fullPath := filepath.Join(".", "res", "uploads", relativePath)

	if err := os.MkdirAll(filepath.Dir(fullPath), 0o755); err != nil {
		return "", err
	}
	if err := os.WriteFile(fullPath, data, 0o644); err != nil {
		return "", err
	}

	return "/res/uploads/" + cleanKey, nil
}

func (u *LocalUploader) CopyImage(cfg dto.UploadConfig, originUrl string) (string, error) {
	data, contentType, err := download(originUrl)
	if err != nil {
		return "", err
	}
	return u.PutImage(cfg, data, contentType)
}
