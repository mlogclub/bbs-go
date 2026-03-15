package uploader

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"sync"

	"github.com/mlogclub/simple/common/strs"
	"github.com/tencentyun/cos-go-sdk-v5"

	"bbs-go/internal/models/dto"
)

// 腾讯云cos
type TencentCosUploader struct {
	m          sync.Mutex
	client     *cos.Client
	currentCfg dto.UploadConfig
}

func (u *TencentCosUploader) PutObject(cfg dto.UploadConfig, key string, body io.Reader, opts *PutOptions) (string, error) {
	if err := u.initClient(cfg); err != nil {
		return "", err
	}
	headerOpts := &cos.ObjectPutHeaderOptions{}
	if opts != nil {
		headerOpts.ContentType = opts.ContentType
		headerOpts.ContentDisposition = opts.ContentDisposition
	}
	opt := &cos.ObjectPutOptions{ObjectPutHeaderOptions: headerOpts}
	if _, err := u.client.Object.Put(context.Background(), key, body, opt); err != nil {
		return "", err
	}
	return fmt.Sprintf("%s/%s", u.client.BaseURL.BucketURL, key), nil
}

func (u *TencentCosUploader) CopyImage(cfg dto.UploadConfig, originUrl string) (string, error) {
	data, ct, err := download(originUrl)
	if err != nil {
		return "", err
	}
	ct = NormalizeImageContentType(ct)
	key := GenerateImageKey(data, ct)
	opts := &PutOptions{ContentType: ct, ContentLength: int64(len(data))}
	return u.PutObject(cfg, key, bytes.NewReader(data), opts)
}

func (u *TencentCosUploader) initClient(cfg dto.UploadConfig) error {
	if !u.isCfgChange(cfg) {
		return nil
	}

	u.m.Lock()
	defer u.m.Unlock()

	// 验证必要配置项不能为空
	if strs.IsAnyBlank(cfg.TencentCos.Bucket, cfg.TencentCos.Region, cfg.TencentCos.SecretId, cfg.TencentCos.SecretKey) {
		return fmt.Errorf("Tencent COS configuration is incomplete: Bucket, Region, SecretId, and SecretKey are required")
	}

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
