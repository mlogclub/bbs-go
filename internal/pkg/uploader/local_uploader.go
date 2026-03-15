package uploader

import (
	"bytes"
	"io"
	"os"
	"path/filepath"
	"strings"

	"bbs-go/internal/models/dto"
	"bbs-go/internal/pkg/respath"
)

type LocalUploader struct{}

func (u *LocalUploader) PutObject(_ dto.UploadConfig, key string, body io.Reader, opts *PutOptions) (string, error) {
	cleanKey := strings.TrimPrefix(filepath.ToSlash(filepath.Clean(key)), "/")
	fullPath := respath.UploadsPath(filepath.FromSlash(cleanKey))
	if err := os.MkdirAll(filepath.Dir(fullPath), 0o755); err != nil {
		return "", err
	}
	dest, err := os.Create(fullPath)
	if err != nil {
		return "", err
	}
	defer dest.Close()
	src := body
	if opts != nil && opts.ContentLength > 0 {
		src = io.LimitReader(body, opts.ContentLength)
	}
	if _, err := io.Copy(dest, src); err != nil {
		_ = os.Remove(fullPath)
		return "", err
	}
	return respath.UploadsURLPrefix + cleanKey, nil
}

func (u *LocalUploader) CopyImage(cfg dto.UploadConfig, originUrl string) (string, error) {
	data, ct, err := download(originUrl)
	if err != nil {
		return "", err
	}
	ct = NormalizeImageContentType(ct)
	key := GenerateImageKey(data, ct)
	opts := &PutOptions{ContentType: ct, ContentLength: int64(len(data))}
	return u.PutObject(cfg, key, bytes.NewReader(data), opts)
}
