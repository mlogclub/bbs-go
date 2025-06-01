package uploader

import (
	"bbs-go/internal/models/dto"
	"bbs-go/internal/pkg/bbsurls"
	"bytes"
	"log/slog"
	"sync"

	"github.com/mlogclub/simple/common/strs"

	"github.com/aliyun/aliyun-oss-go-sdk/oss"
)

type AliyunOssUploader struct {
	m      sync.Mutex
	bucket *oss.Bucket
}

func (u *AliyunOssUploader) PutImage(cfg dto.UploadConfig, data []byte, contentType string) (string, error) {
	if strs.IsBlank(contentType) {
		contentType = "image/jpeg"
	}
	key := generateImageKey(data, contentType)
	return u.PutObject(cfg, key, data, contentType)
}

func (u *AliyunOssUploader) PutObject(cfg dto.UploadConfig, key string, data []byte, contentType string) (string, error) {
	if err := u.initBucket(cfg); err != nil {
		return "", err
	}
	var options []oss.Option
	if strs.IsNotBlank(contentType) {
		options = append(options, oss.ContentType(contentType))
	}
	if err := u.bucket.PutObject(key, bytes.NewReader(data), options...); err != nil {
		return "", err
	}
	return bbsurls.UrlJoin(cfg.AliyunOss.Host, key), nil
}

func (u *AliyunOssUploader) CopyImage(cfg dto.UploadConfig, originUrl string) (string, error) {
	data, contentType, err := download(originUrl)
	if err != nil {
		return "", err
	}
	return u.PutImage(cfg, data, contentType)
}

func (u *AliyunOssUploader) initBucket(cfg dto.UploadConfig) error {
	if !u.isCfgChange(cfg) {
		return nil
	}

	u.m.Lock()
	defer u.m.Unlock()

	client, err := oss.New(cfg.AliyunOss.Endpoint, cfg.AliyunOss.AccessKeyId, cfg.AliyunOss.AccessKeySecret)
	if err != nil {
		slog.Error(err.Error(), slog.Any("err", err))
		return err
	}

	bucket, err := client.Bucket(cfg.AliyunOss.Bucket)
	if err != nil {
		slog.Error(err.Error(), slog.Any("err", err))
		return err
	}

	u.bucket = bucket
	return nil
}

func (u *AliyunOssUploader) isCfgChange(cfg dto.UploadConfig) bool {
	if u.bucket == nil {
		return true
	}

	if u.bucket.Client.Config.Endpoint != cfg.AliyunOss.Endpoint ||
		u.bucket.Client.Config.AccessKeyID != cfg.AliyunOss.AccessKeyId ||
		u.bucket.Client.Config.AccessKeySecret != cfg.AliyunOss.AccessKeySecret ||
		u.bucket.BucketName != cfg.AliyunOss.Bucket {
		return true
	}

	return false
}
