package uploader

import (
	"bbs-go/pkg/urls"
	"bytes"
	"github.com/mlogclub/simple"
	"sync"

	"github.com/aliyun/aliyun-oss-go-sdk/oss"
	"github.com/sirupsen/logrus"

	"bbs-go/pkg/config"
)

// 阿里云oss
type aliyunOssUploader struct {
	once   sync.Once
	bucket *oss.Bucket
}

func (aliyun *aliyunOssUploader) PutImage(data []byte, contentType string) (string, error) {
	key := generateImageKey(data, contentType)
	if simple.IsBlank(contentType) {
		contentType = "image/jpeg"
	}
	return aliyun.PutObject(key, data, contentType)
}

func (aliyun *aliyunOssUploader) PutObject(key string, data []byte, contentType string) (string, error) {
	bucket := aliyun.getBucket()
	var options []oss.Option
	if simple.IsNotBlank(contentType) {
		options = append(options, oss.ContentType(contentType))
	}
	if err := bucket.PutObject(key, bytes.NewReader(data), options...); err != nil {
		return "", err
	}
	c := config.Instance.Uploader.AliyunOss
	return urls.UrlJoin(c.Host, key), nil
}

func (aliyun *aliyunOssUploader) CopyImage(originUrl string) (string, error) {
	data, contentType, err := download(originUrl)
	if err != nil {
		return "", err
	}
	return aliyun.PutImage(data, contentType)
}

func (aliyun *aliyunOssUploader) getBucket() *oss.Bucket {
	aliyun.once.Do(func() {
		c := config.Instance.Uploader.AliyunOss
		if client, err := oss.New(c.Endpoint, c.AccessId, c.AccessSecret); err != nil {
			logrus.Error(err)
		} else if aliyun.bucket, err = client.Bucket(c.Bucket); err != nil {
			logrus.Error(err)
		}
	})
	return aliyun.bucket
}
