package uploader

import (
	"bbs-go/pkg/bbsurls"
	"bytes"
	"context"
	"fmt"
	"github.com/mlogclub/simple/common/strs"
	"net/http"
	"net/url"
	"sync"
	"time"

	"bbs-go/pkg/config"
	"github.com/tencentyun/cos-go-sdk-v5"
)

// 腾讯云oss
type txyCosUploader struct {
	once   sync.Once
	object *cos.ObjectService
}

func (txy *txyCosUploader) PutImage(data []byte, contentType string) (string, error) {
	if strs.IsBlank(contentType) {
		contentType = "image/jpeg"
	}
	key := generateImageKey(data, contentType)
	return txy.PutObject(key, data, contentType)
}

func (txy *txyCosUploader) PutObject(key string, data []byte, contentType string) (string, error) {
	object := txy.getObject()
	reader := bytes.NewReader(data)
	opt := &cos.ObjectPutOptions{
		ObjectPutHeaderOptions: &cos.ObjectPutHeaderOptions{
			ContentType: contentType,
		},
		ACLHeaderOptions: &cos.ACLHeaderOptions{
			XCosACL: "public-read",
		},
	}
	if _, err := object.Put(context.Background(), key, reader, opt); err != nil {
		return "", err
	}
	c := config.Instance.Uploader.TxyCos
	return bbsurls.UrlJoin(c.Host, key), nil
}

func (txy *txyCosUploader) CopyImage(originUrl string) (string, error) {
	data, contentType, err := download(originUrl)
	if err != nil {
		return "", err
	}
	return txy.PutImage(data, contentType)
}

func (txy *txyCosUploader) getObject() *cos.ObjectService {
	txy.once.Do(func() {
		c := config.Instance.Uploader.TxyCos
		urls := fmt.Sprintf("https://%s.cos.%s.myqcloud.com", c.Bucket, c.Region)
		u, _ := url.Parse(urls)
		b := &cos.BaseURL{BucketURL: u}
		client := cos.NewClient(b, &http.Client{
			Timeout: 100 * time.Second,
			Transport: &cos.AuthorizationTransport{
				SecretID:  c.SecretID,
				SecretKey: c.SecretKey,
			},
		})
		txy.object = client.Object

	})
	return txy.object
}
