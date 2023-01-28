package uploader

import (
	"bytes"
	"github.com/sirupsen/logrus"
	"io/ioutil"
	"mime/multipart"
	"net/http"
	"server/pkg/bbsurls"
	"server/pkg/config"
)

// 本地文件系统
type serverUploader struct{}

func (server *serverUploader) PutImage(data []byte, contentType string) (string, error) {
	key := generateImageKey(data, contentType)
	return server.PutObject(key, data, contentType)
}

func (server *serverUploader) PutObject(key string, data []byte, contentType string) (string, error) {
	c := config.Instance.Uploader.Server
	bodyBuf := &bytes.Buffer{}
	bodyWriter := multipart.NewWriter(bodyBuf)
	fileWriter, err := bodyWriter.CreateFormFile("file", key)
	if err != nil {
		logrus.Errorf("error wriTing to buffer")
		return "", err
	}
	_, err = fileWriter.Write(data)
	if err != nil {
		return "", err
	}
	contentType = bodyWriter.FormDataContentType()
	bodyWriter.Close()
	resp, err := http.Post(c.URL, contentType, bodyBuf)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()
	resp_body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}
	return bbsurls.UrlJoin(c.Host, string(resp_body)), nil
}

func (server *serverUploader) CopyImage(originUrl string) (string, error) {
	data, contentType, err := download(originUrl)
	if err != nil {
		return "", err
	}
	return server.PutImage(data, contentType)
}
