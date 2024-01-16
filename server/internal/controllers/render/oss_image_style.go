package render

import (
	"bbs-go/internal/pkg/config"
	"bbs-go/internal/pkg/uploader"
	"strings"

	"github.com/mlogclub/simple/common/strs"
)

func HandleOssImageStyleAvatar(url string) string {
	if !uploader.IsEnabledOss() {
		return url
	}
	return HandleOssImageStyle(url, config.Instance.Uploader.AliyunOss.StyleAvatar)
}

func HandleOssImageStyleDetail(url string) string {
	if !uploader.IsEnabledOss() {
		return url
	}
	return HandleOssImageStyle(url, config.Instance.Uploader.AliyunOss.StyleDetail)
}

func HandleOssImageStyleSmall(url string) string {
	if !uploader.IsEnabledOss() {
		return url
	}
	return HandleOssImageStyle(url, config.Instance.Uploader.AliyunOss.StyleSmall)
}

func HandleOssImageStylePreview(url string) string {
	if !uploader.IsEnabledOss() {
		return url
	}
	return HandleOssImageStyle(url, config.Instance.Uploader.AliyunOss.StylePreview)
}

func HandleOssImageStyle(url, style string) string {
	if strs.IsBlank(style) || strs.IsBlank(url) {
		return url
	}
	if !uploader.IsOssImageUrl(url) {
		return url
	}
	if strings.HasSuffix(strings.ToLower(url), ".gif") {
		return url
	}
	sep := config.Instance.Uploader.AliyunOss.StyleSplitter
	if strs.IsBlank(sep) {
		return url
	}
	return strings.Join([]string{url, style}, sep)
}
