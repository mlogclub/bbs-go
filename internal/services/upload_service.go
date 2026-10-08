package services

import (
	"bytes"
	"fmt"
	"io"
	"net/http"
	"strings"
	"sync"

	"github.com/mlogclub/simple/common/strs"
	"github.com/mlogclub/simple/common/urls"

	"bbs-go/internal/models/dto"
	"bbs-go/internal/pkg/bbsurls"
	"bbs-go/internal/pkg/respath"
	"bbs-go/internal/pkg/uploader"
)

// allowedImageContentTypes 允许上传的真实图片 MIME 类型，基于文件内容嗅探而非客户端
// 声明的 Content-Type，防止以图片名义上传可在浏览器中被当作其他类型执行的内容（如 HTML）。
// 不包含 image/svg+xml：SVG 可嵌入脚本，不属于安全的光栅图片格式。
var allowedImageContentTypes = map[string]bool{
	"image/jpeg":   true,
	"image/png":    true,
	"image/gif":    true,
	"image/webp":   true,
	"image/bmp":    true,
	"image/x-icon": true,
}

// sniffImageContentType 读取数据前缀以嗅探真实内容类型并校验是否为允许的图片格式，
// 而不是信任调用方传入的 contentType（可被客户端任意伪造）。返回嗅探到的 Content-Type
// 与包含完整数据的可读流（供后续写入存储，不丢失已读取的前缀字节）。
func sniffImageContentType(body io.Reader) (string, io.Reader, error) {
	sniffBuf := make([]byte, 512)
	n, err := io.ReadFull(body, sniffBuf)
	if err != nil && err != io.ErrUnexpectedEOF && err != io.EOF {
		return "", nil, err
	}
	sniffBuf = sniffBuf[:n]
	detected := http.DetectContentType(sniffBuf)
	if idx := strings.Index(detected, ";"); idx >= 0 {
		detected = detected[:idx]
	}
	if !allowedImageContentTypes[detected] {
		return "", nil, fmt.Errorf("unsupported file type: %s", detected)
	}
	return detected, io.MultiReader(bytes.NewReader(sniffBuf), body), nil
}

var UploadService = newUploadService()

type uploadService struct {
	uploaderMap map[dto.UploadMethod]uploader.Uploader
	once        sync.Once
}

func newUploadService() *uploadService {
	return &uploadService{
		uploaderMap: make(map[dto.UploadMethod]uploader.Uploader),
	}
}

func (s *uploadService) putObject(key string, body io.Reader, opts *uploader.PutOptions) (string, error) {
	u, err := s.getUploader()
	if err != nil {
		return "", err
	}
	cfg := SysConfigService.GetUploadConfig()
	return u.PutObject(cfg, key, body, opts)
}

// PutObject 按 key 流式上传；opts 可设置 ContentType、ContentDisposition、ContentLength。
func (s *uploadService) PutObject(key string, body io.Reader, opts *uploader.PutOptions) (string, error) {
	return s.putObject(key, body, opts)
}

func (s *uploadService) ObjectURL(key string) string {
	cfg := SysConfigService.GetUploadConfig()
	if strs.IsBlank(string(cfg.EnableUploadMethod)) {
		cfg.EnableUploadMethod = dto.Local
	}

	switch cfg.EnableUploadMethod {
	case dto.AliyunOss:
		return bbsurls.UrlJoin(cfg.AliyunOss.Host, key)
	case dto.TencentCos:
		return fmt.Sprintf("https://%s.cos.%s.myqcloud.com/%s", cfg.TencentCos.Bucket, cfg.TencentCos.Region, key)
	case dto.AwsS3:
		return fmt.Sprintf("https://%s.s3.%s.amazonaws.com/%s", cfg.AwsS3.Bucket, cfg.AwsS3.Region, key)
	case dto.S3:
		return uploader.S3ObjectURL(cfg, key)
	default:
		return respath.UploadsURLPrefix + key
	}
}

// PutImage 上传图片（已有完整字节）；key 使用内容 MD5，供 CopyImage 等场景。
// contentType 入参已不再使用调用方声明的值，实际类型通过内容嗅探得出。
func (s *uploadService) PutImage(data []byte, _ string) (string, error) {
	detectedType, _, err := sniffImageContentType(bytes.NewReader(data))
	if err != nil {
		return "", err
	}
	key := uploader.GenerateImageKey(data, detectedType)
	opts := &uploader.PutOptions{ContentType: detectedType, ContentLength: int64(len(data))}
	return s.putObject(key, bytes.NewReader(data), opts)
}

// PutImageStream 流式上传图片；key 使用 UUID，无需先读完整 body。
// contentType 入参已不再使用调用方声明的值，实际类型通过内容嗅探得出。
func (s *uploadService) PutImageStream(body io.Reader, contentLength int64, _ string) (string, error) {
	detectedType, fullBody, err := sniffImageContentType(body)
	if err != nil {
		return "", err
	}
	key := uploader.GenerateImageKeyByContentType(detectedType)
	opts := &uploader.PutOptions{ContentType: detectedType, ContentLength: contentLength}
	return s.putObject(key, fullBody, opts)
}

func (s *uploadService) CopyImage(url string) (string, error) {
	u, err := s.getUploader()
	if err != nil {
		return "", err
	}
	u1 := urls.ParseUrl(url).GetURL()
	u2 := urls.ParseUrl(SysConfigService.GetBaseURL()).GetURL()
	if u1.Host == u2.Host {
		return url, nil
	}
	cfg := SysConfigService.GetUploadConfig()
	return u.CopyImage(cfg, url)
}

func (s *uploadService) getUploader() (uploader.Uploader, error) {
	s.once.Do(func() {
		s.uploaderMap[dto.Local] = &uploader.LocalUploader{}
		s.uploaderMap[dto.AliyunOss] = &uploader.AliyunOssUploader{}
		s.uploaderMap[dto.TencentCos] = &uploader.TencentCosUploader{}
		s.uploaderMap[dto.AwsS3] = &uploader.AwsS3Uploader{}
		s.uploaderMap[dto.S3] = &uploader.S3Uploader{}
	})
	cfg := SysConfigService.GetUploadConfig()

	if strs.IsBlank(string(cfg.EnableUploadMethod)) {
		cfg.EnableUploadMethod = dto.Local
	}

	u, ok := s.uploaderMap[cfg.EnableUploadMethod]
	if !ok {
		return nil, fmt.Errorf("error: Upload method: %s not found", cfg.EnableUploadMethod)
	}
	return u, nil
}
