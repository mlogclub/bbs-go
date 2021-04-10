package uploader

import (
	"bbs-go/package/urls"
	"io/ioutil"
	"os"
	"path/filepath"

	"bbs-go/package/config"
)

// 本地文件系统
type localUploader struct{}

func (local *localUploader) PutImage(data []byte) (string, error) {
	key := generateImageKey(data)
	return local.PutObject(key, data)
}

func (local *localUploader) PutObject(key string, data []byte) (string, error) {
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
	return urls.UrlJoin(c.Host, key), nil
}

func (local *localUploader) CopyImage(originUrl string) (string, error) {
	data, err := download(originUrl)
	if err != nil {
		return "", err
	}
	return local.PutImage(data)
}
