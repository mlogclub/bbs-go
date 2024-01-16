package uploader

import (
	"bbs-go/internal/pkg/bbsurls"
	"io/ioutil"
	"os"
	"path/filepath"

	"bbs-go/internal/pkg/config"
)

// 本地文件系统
type localUploader struct{}

func (local *localUploader) PutImage(data []byte, contentType string) (string, error) {
	key := generateImageKey(data, contentType)
	return local.PutObject(key, data, contentType)
}

func (local *localUploader) PutObject(key string, data []byte, contentType string) (string, error) {
	if err := os.MkdirAll("/", os.ModeDir); err != nil {
		return "", err
	}
	c := config.Instance.Uploader.Local
	filename := filepath.Join(c.Path, key)
	if err := os.MkdirAll(filepath.Dir(filename), os.ModePerm); err != nil {
		return "", err
	}
	if err := ioutil.WriteFile(filename, data, os.ModePerm); err != nil {
		return "", err
	}
	return bbsurls.UrlJoin(c.Host, key), nil
}

func (local *localUploader) CopyImage(originUrl string) (string, error) {
	data, contentType, err := download(originUrl)
	if err != nil {
		return "", err
	}
	return local.PutImage(data, contentType)
}
