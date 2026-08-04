package uploader

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"log/slog"
	"net/url"
	"strings"
	"sync"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/mlogclub/simple/common/strs"

	"bbs-go/internal/models/dto"
	"bbs-go/internal/pkg/bbsurls"
)

// S3Uploader 通用 S3 兼容存储上传器，兼容 Cloudflare R2 / MinIO / RustFS 等。
// 复用 aws-sdk-go-v2 实现，通过自定义 Endpoint 与 PathStyle 适配各家 S3 兼容服务。
type S3Uploader struct {
	m          sync.Mutex
	client     *s3.Client
	currentCfg dto.UploadConfig
}

func (u *S3Uploader) PutObject(cfg dto.UploadConfig, key string, body io.Reader, opts *PutOptions) (string, error) {
	if err := u.initClient(cfg); err != nil {
		return "", err
	}
	putInput := &s3.PutObjectInput{
		Bucket: aws.String(cfg.S3.Bucket),
		Key:    aws.String(key),
		Body:   body,
	}
	if opts != nil {
		if opts.ContentLength > 0 {
			putInput.ContentLength = aws.Int64(opts.ContentLength)
		}
		if opts.ContentType != "" {
			putInput.ContentType = aws.String(opts.ContentType)
		}
		if opts.ContentDisposition != "" {
			putInput.ContentDisposition = aws.String(opts.ContentDisposition)
		}
	}
	if _, err := u.client.PutObject(context.Background(), putInput); err != nil {
		slog.Error("S3 PutObject failed", slog.Any("err", err), slog.String("bucket", cfg.S3.Bucket), slog.String("key", key))
		return "", fmt.Errorf("failed to upload object to S3: %w", err)
	}
	return S3ObjectURL(cfg, key), nil
}

func (u *S3Uploader) CopyImage(cfg dto.UploadConfig, originUrl string) (string, error) {
	data, ct, err := download(originUrl)
	if err != nil {
		return "", err
	}
	ct = NormalizeImageContentType(ct)
	key := GenerateImageKey(data, ct)
	opts := &PutOptions{ContentType: ct, ContentLength: int64(len(data))}
	return u.PutObject(cfg, key, bytes.NewReader(data), opts)
}

func (u *S3Uploader) initClient(cfg dto.UploadConfig) error {
	if !u.isCfgChange(cfg) {
		return nil
	}

	u.m.Lock()
	defer u.m.Unlock()

	if strs.IsAnyBlank(cfg.S3.Endpoint, cfg.S3.Bucket, cfg.S3.AccessKeyId, cfg.S3.AccessKeySecret) {
		return fmt.Errorf("S3 configuration is incomplete: Endpoint, Bucket, AccessKeyId, and AccessKeySecret are required")
	}

	region := strings.TrimSpace(cfg.S3.Region)
	if region == "" {
		// Cloudflare R2 使用 auto；MinIO/RustFS 也接受任意非空值
		region = "auto"
	}

	awsCfg, err := config.LoadDefaultConfig(context.TODO(),
		config.WithRegion(region),
		config.WithCredentialsProvider(credentials.NewStaticCredentialsProvider(
			cfg.S3.AccessKeyId,
			cfg.S3.AccessKeySecret,
			"",
		)),
	)
	if err != nil {
		slog.Error("Failed to load S3 config", slog.Any("err", err))
		return fmt.Errorf("failed to load S3 config: %w", err)
	}

	u.client = s3.NewFromConfig(awsCfg, func(o *s3.Options) {
		o.BaseEndpoint = aws.String(strings.TrimRight(cfg.S3.Endpoint, "/"))
		o.UsePathStyle = cfg.S3.PathStyle
	})

	u.currentCfg = cfg
	return nil
}

func (u *S3Uploader) isCfgChange(cfg dto.UploadConfig) bool {
	if u.client == nil {
		return true
	}
	if u.currentCfg.S3.Endpoint != cfg.S3.Endpoint ||
		u.currentCfg.S3.Region != cfg.S3.Region ||
		u.currentCfg.S3.Bucket != cfg.S3.Bucket ||
		u.currentCfg.S3.AccessKeyId != cfg.S3.AccessKeyId ||
		u.currentCfg.S3.AccessKeySecret != cfg.S3.AccessKeySecret ||
		u.currentCfg.S3.PathStyle != cfg.S3.PathStyle {
		return true
	}
	return false
}

// S3ObjectURL 生成 S3 兼容存储对象的访问 URL。
// 优先使用配置的 Host 作为前缀（自定义域名 / CDN / R2 公开域名）。
// 未配置 Host 时按 Endpoint 推导：
//   - PathStyle=true:  {endpoint}/{bucket}/{key}
//   - PathStyle=false: https://{bucket}.{endpoint-host}/{key}
func S3ObjectURL(cfg dto.UploadConfig, key string) string {
	if host := strings.TrimSpace(cfg.S3.Host); host != "" {
		return bbsurls.UrlJoin(host, key)
	}

	endpoint := strings.TrimRight(strings.TrimSpace(cfg.S3.Endpoint), "/")
	if endpoint == "" {
		return key
	}

	if cfg.S3.PathStyle {
		return bbsurls.UrlJoin(endpoint, cfg.S3.Bucket, key)
	}

	parsed, err := url.Parse(endpoint)
	if err != nil || parsed.Host == "" {
		// 退化到 path-style，保证可用
		return bbsurls.UrlJoin(endpoint, cfg.S3.Bucket, key)
	}
	scheme := parsed.Scheme
	if scheme == "" {
		scheme = "https"
	}
	return fmt.Sprintf("%s://%s.%s/%s", scheme, cfg.S3.Bucket, parsed.Host, key)
}
