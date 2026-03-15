package uploader

import (
	"bytes"
	"fmt"
	"io"
	"log/slog"
	"sync"

	"github.com/aliyun/aliyun-oss-go-sdk/oss"
	"github.com/mlogclub/simple/common/strs"

	"bbs-go/internal/models/dto"
	"bbs-go/internal/pkg/bbsurls"
)

type AliyunOssUploader struct {
	m      sync.Mutex
	bucket *oss.Bucket
}

func (u *AliyunOssUploader) PutObject(cfg dto.UploadConfig, key string, body io.Reader, opts *PutOptions) (string, error) {
	if err := u.initBucket(cfg); err != nil {
		return "", err
	}
	var options []oss.Option
	if opts != nil {
		if opts.ContentType != "" {
			options = append(options, oss.ContentType(opts.ContentType))
		}
		if opts.ContentDisposition != "" {
			options = append(options, oss.ContentDisposition(opts.ContentDisposition))
		}
	}
	if err := u.bucket.PutObject(key, body, options...); err != nil {
		return "", err
	}
	return bbsurls.UrlJoin(cfg.AliyunOss.Host, key), nil
}

func (u *AliyunOssUploader) CopyImage(cfg dto.UploadConfig, originUrl string) (string, error) {
	data, ct, err := download(originUrl)
	if err != nil {
		return "", err
	}
	ct = NormalizeImageContentType(ct)
	key := GenerateImageKey(data, ct)
	opts := &PutOptions{ContentType: ct, ContentLength: int64(len(data))}
	return u.PutObject(cfg, key, bytes.NewReader(data), opts)
}

func (u *AliyunOssUploader) initBucket(cfg dto.UploadConfig) error {
	if !u.isCfgChange(cfg) {
		return nil
	}

	u.m.Lock()
	defer u.m.Unlock()

	// 验证必要配置项不能为空
	if strs.IsAnyBlank(cfg.AliyunOss.Endpoint, cfg.AliyunOss.AccessKeyId, cfg.AliyunOss.AccessKeySecret, cfg.AliyunOss.Bucket, cfg.AliyunOss.Host) {
		return fmt.Errorf("Aliyun OSS configuration is incomplete: Endpoint, AccessKeyId, AccessKeySecret, Bucket, and Host are required")
	}

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
