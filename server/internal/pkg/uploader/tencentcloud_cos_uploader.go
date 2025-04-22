package uploader

import (
	"bytes"
	"context"
	"fmt"
	"net/http"
	"net/url"
	"sync"

	"github.com/tencentyun/cos-go-sdk-v5"

	"github.com/mlogclub/simple/common/strs"

	"bbs-go/internal/pkg/config"
)

// 腾讯云cos
type tencentCloudCosUploader struct {
	once   sync.Once
	client *cos.Client
}

func (cosUploader *tencentCloudCosUploader) PutImage(data []byte, contentType string) (string, error) {
	if strs.IsBlank(contentType) {
		contentType = "image/jpeg"
	}
	key := generateImageKey(data, contentType)
	return tencentCos.PutObject(key, data, contentType)
}

func (cosUploader *tencentCloudCosUploader) PutObject(key string, data []byte, contentType string) (string, error) {
	cosUploader.initClient()

	opt := &cos.ObjectPutOptions{
		ObjectPutHeaderOptions: &cos.ObjectPutHeaderOptions{
			ContentType: contentType,
		},
	}
	if _, err := cosUploader.client.Object.Put(context.Background(), key, bytes.NewReader(data), opt); err != nil {
		return "", err
	}

	return fmt.Sprintf("%s/%s", cosUploader.client.BaseURL.BucketURL, key), nil
}

func (cosUploader *tencentCloudCosUploader) CopyImage(originUrl string) (string, error) {
	data, contentType, err := download(originUrl)
	if err != nil {
		return "", err
	}
	return tencentCos.PutImage(data, contentType)
}

func (cosUploader *tencentCloudCosUploader) initClient() {
	tencentCos.once.Do(func() {
		c := config.Instance.Uploader.TencentCos

		u, _ := url.Parse(fmt.Sprintf("https://%s.cos.%s.myqcloud.com", c.Bucket, c.Region))
		su, _ := url.Parse(fmt.Sprintf("https://cos.%s.myqcloud.com", c.Region))
		b := &cos.BaseURL{BucketURL: u, ServiceURL: su}
		tencentCos.client = cos.NewClient(b, &http.Client{
			Transport: &cos.AuthorizationTransport{
				SecretID:  c.SecretId,
				SecretKey: c.SecretKey,
			},
		})
	})
}
