package uploader

import (
	"bbs-go/common/urls"
	"bbs-go/config"
	"bytes"
	"context"
	"github.com/sirupsen/logrus"
	"github.com/tencentyun/cos-go-sdk-v5"
	"net/http"
	"net/url"
	"sync"
)

type tencentCosUploader struct {
	once   sync.Once
	client *cos.Client
}

func (cos *tencentCosUploader) PutImage(data []byte) (string, error) {
	key := generateImageKey(data)
	return cos.PutObject(key, data)
}

func (cos *tencentCosUploader) PutObject(key string, data []byte) (string, error) {
	client := cos.getClient()
	f := bytes.NewReader(data)
	if _, err := client.Object.Put(context.Background(), key, f, nil); err != nil {
		return "", err
	}
	c := config.Instance.Uploader.ObjectStorage
	return urls.UrlJoin(c.Host, key), nil
}

func (cos *tencentCosUploader) CopyImage(originUrl string) (string, error) {
	data, err := download(originUrl)
	if err != nil {
		return "", err
	}
	return cos.PutImage(data)
}

func (client *tencentCosUploader) getClient() *cos.Client {
	client.once.Do(func() {
		con := config.Instance.Uploader.ObjectStorage
		cosurl, err := url.Parse(con.Endpoint)
		if err != nil {
			logrus.Error(err)
		}
		baseUrl := &cos.BaseURL{BucketURL: cosurl}
		client.client = cos.NewClient(baseUrl, &http.Client{
			Transport: &cos.AuthorizationTransport{
				SecretID:  con.AccessId,
				SecretKey: con.AccessSecret,
			},
		})
	},
	)
	return client.client
}
