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

	"bbs-go/internal/models/dto"
)

// 腾讯云cos
type TencentCosUploader struct {
	m          sync.Mutex
	client     *cos.Client
	currentCfg dto.UploadConfig
}

func (u *TencentCosUploader) PutImage(cfg dto.UploadConfig, data []byte, contentType string) (string, error) {
	if strs.IsBlank(contentType) {
		contentType = "image/jpeg"
	}
	key := generateImageKey(data, contentType)
	return u.PutObject(cfg, key, data, contentType)
}

func (u *TencentCosUploader) PutObject(cfg dto.UploadConfig, key string, data []byte, contentType string) (string, error) {
	if err := u.initClient(cfg); err != nil {
		return "", nil
	}

	opt := &cos.ObjectPutOptions{
		ObjectPutHeaderOptions: &cos.ObjectPutHeaderOptions{
			ContentType: contentType,
		},
	}
	if _, err := u.client.Object.Put(context.Background(), key, bytes.NewReader(data), opt); err != nil {
		return "", err
	}

	return fmt.Sprintf("%s/%s", u.client.BaseURL.BucketURL, key), nil
}

func (u *TencentCosUploader) CopyImage(cfg dto.UploadConfig, originUrl string) (string, error) {
	data, contentType, err := download(originUrl)
	if err != nil {
		return "", err
	}
	return u.PutImage(cfg, data, contentType)
}

func (u *TencentCosUploader) initClient(cfg dto.UploadConfig) error {
	if !u.isCfgChange(cfg) {
		return nil
	}

	u.m.Lock()
	defer u.m.Unlock()

	bucketURL, _ := url.Parse(fmt.Sprintf("https://%s.cos.%s.myqcloud.com", cfg.TencentCos.Bucket, cfg.TencentCos.Region))
	serviceURL, _ := url.Parse(fmt.Sprintf("https://cos.%s.myqcloud.com", cfg.TencentCos.Region))
	baseURL := &cos.BaseURL{BucketURL: bucketURL, ServiceURL: serviceURL}
	u.client = cos.NewClient(baseURL, &http.Client{
		Transport: &cos.AuthorizationTransport{
			SecretID:  cfg.TencentCos.SecretId,
			SecretKey: cfg.TencentCos.SecretKey,
		},
	})
	u.currentCfg = cfg

	return nil
}

func (u *TencentCosUploader) isCfgChange(cfg dto.UploadConfig) bool {
	if u.client == nil {
		return true
	}

	if u.currentCfg.TencentCos.Bucket != cfg.TencentCos.Bucket ||
		u.currentCfg.TencentCos.Region != cfg.TencentCos.Region ||
		u.currentCfg.TencentCos.SecretId != cfg.TencentCos.SecretId ||
		u.currentCfg.TencentCos.SecretKey != cfg.TencentCos.SecretKey {
		return true
	}

	return false
}
