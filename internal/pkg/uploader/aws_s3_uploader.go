package uploader

import (
	"bbs-go/internal/models/dto"
	"bytes"
	"context"
	"fmt"
	"log/slog"
	"sync"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/mlogclub/simple/common/strs"
)

type AwsS3Uploader struct {
	m          sync.Mutex
	client     *s3.Client
	currentCfg dto.UploadConfig
}

func (u *AwsS3Uploader) PutImage(cfg dto.UploadConfig, data []byte, contentType string) (string, error) {
	if strs.IsBlank(contentType) {
		contentType = "image/jpeg"
	}
	key := generateImageKey(data, contentType)
	return u.PutObject(cfg, key, data, contentType)
}

func (u *AwsS3Uploader) PutObject(cfg dto.UploadConfig, key string, data []byte, contentType string) (string, error) {
	if err := u.initClient(cfg); err != nil {
		return "", err
	}

	putInput := &s3.PutObjectInput{
		Bucket:      aws.String(cfg.AwsS3.Bucket),
		Key:         aws.String(key),
		Body:        bytes.NewReader(data),
		ContentType: aws.String(contentType),
	}

	_, err := u.client.PutObject(context.Background(), putInput)
	if err != nil {
		slog.Error("AWS S3 PutObject failed", slog.Any("err", err), slog.String("bucket", cfg.AwsS3.Bucket), slog.String("key", key))
		return "", fmt.Errorf("failed to upload object to S3: %w", err)
	}

	// 生成固定格式的 URL: https://{bucket}.s3.{region}.amazonaws.com/{key}
	url := fmt.Sprintf("https://%s.s3.%s.amazonaws.com/%s", cfg.AwsS3.Bucket, cfg.AwsS3.Region, key)
	return url, nil
}

func (u *AwsS3Uploader) CopyImage(cfg dto.UploadConfig, originUrl string) (string, error) {
	data, contentType, err := download(originUrl)
	if err != nil {
		return "", err
	}
	return u.PutImage(cfg, data, contentType)
}

func (u *AwsS3Uploader) initClient(cfg dto.UploadConfig) error {
	if !u.isCfgChange(cfg) {
		return nil
	}

	u.m.Lock()
	defer u.m.Unlock()

	// 验证必要配置项不能为空
	if strs.IsAnyBlank(cfg.AwsS3.Region, cfg.AwsS3.Bucket, cfg.AwsS3.AccessKeyId, cfg.AwsS3.AccessKeySecret) {
		return fmt.Errorf("AWS S3 configuration is incomplete: Region, Bucket, AccessKeyId, and AccessKeySecret are required")
	}

	// 创建 AWS 配置（使用标准 AWS S3）
	awsCfg, err := config.LoadDefaultConfig(context.TODO(),
		config.WithRegion(cfg.AwsS3.Region),
		config.WithCredentialsProvider(credentials.NewStaticCredentialsProvider(
			cfg.AwsS3.AccessKeyId,
			cfg.AwsS3.AccessKeySecret,
			"",
		)),
	)
	if err != nil {
		slog.Error("Failed to load AWS config", slog.Any("err", err))
		return fmt.Errorf("failed to load AWS config: %w", err)
	}

	// 创建 S3 客户端
	u.client = s3.NewFromConfig(awsCfg)

	u.currentCfg = cfg
	return nil
}

func (u *AwsS3Uploader) isCfgChange(cfg dto.UploadConfig) bool {
	if u.client == nil {
		return true
	}

	if u.currentCfg.AwsS3.Region != cfg.AwsS3.Region ||
		u.currentCfg.AwsS3.Bucket != cfg.AwsS3.Bucket ||
		u.currentCfg.AwsS3.AccessKeyId != cfg.AwsS3.AccessKeyId ||
		u.currentCfg.AwsS3.AccessKeySecret != cfg.AwsS3.AccessKeySecret {
		return true
	}

	return false
}
