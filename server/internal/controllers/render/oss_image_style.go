package render

import (
	"bbs-go/internal/models/dto"
	"bbs-go/internal/services"
	"strings"

	"github.com/mlogclub/simple/common/strs"
)

func HandleOssImageStyleAvatar(url string) string {
	cfg := services.SysConfigService.GetUploadConfig()
	if cfg.EnableUploadMethod != dto.AliyunOss {
		return url
	}
	return HandleOssImageStyle(url, cfg.AliyunOss.StyleSplitter, cfg.AliyunOss.StyleAvatar)
}

func HandleOssImageStyleDetail(url string) string {
	cfg := services.SysConfigService.GetUploadConfig()
	if cfg.EnableUploadMethod != dto.AliyunOss {
		return url
	}
	return HandleOssImageStyle(url, cfg.AliyunOss.StyleSplitter, cfg.AliyunOss.StyleDetail)
}

func HandleOssImageStyleSmall(url string) string {
	cfg := services.SysConfigService.GetUploadConfig()
	if cfg.EnableUploadMethod != dto.AliyunOss {
		return url
	}
	return HandleOssImageStyle(url, cfg.AliyunOss.StyleSplitter, cfg.AliyunOss.StyleSmall)
}

func HandleOssImageStylePreview(url string) string {
	cfg := services.SysConfigService.GetUploadConfig()
	if cfg.EnableUploadMethod != dto.AliyunOss {
		return url
	}
	return HandleOssImageStyle(url, cfg.AliyunOss.StyleSplitter, cfg.AliyunOss.StylePreview)
}

func HandleOssImageStyle(url, splitter, style string) string {
	if strs.IsBlank(style) || strs.IsBlank(url) {
		return url
	}
	if strings.HasSuffix(strings.ToLower(url), ".gif") {
		return url
	}
	if strs.IsBlank(splitter) {
		return url
	}
	return strings.Join([]string{url, style}, splitter)
}
