package services

import (
	"bbs-go/internal/models/dto"
	"bbs-go/internal/pkg/uploader"
	"fmt"
	"sync"

	"github.com/mlogclub/simple/common/strs"
	"github.com/mlogclub/simple/common/urls"
)

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

func (s *uploadService) PutImage(data []byte, contentType string) (string, error) {
	u, err := s.getUploader()
	if err != nil {
		return "", err
	}
	cfg := SysConfigService.GetUploadConfig()
	return u.PutImage(cfg, data, contentType)
}

func (s *uploadService) PutObject(key string, data []byte, contentType string) (string, error) {
	u, err := s.getUploader()
	if err != nil {
		return "", err
	}
	cfg := SysConfigService.GetUploadConfig()
	return u.PutObject(cfg, key, data, contentType)
}

func (s *uploadService) CopyImage(url string) (string, error) {
	u, err := s.getUploader()
	if err != nil {
		return "", err
	}
	u1 := urls.ParseUrl(url).GetURL()
	u2 := urls.ParseUrl(SysConfigService.GetBaseURL()).GetURL()
	// 本站host，不下载
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
