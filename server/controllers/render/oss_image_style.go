package render

import (
	"bbs-go/pkg/config"
	"bbs-go/pkg/uploader"
	"strings"

	"github.com/mlogclub/simple"
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
	if simple.IsBlank(style) || simple.IsBlank(url) {
		return url
	}
	if !uploader.IsOssImageUrl(url) {
		return url
	}
	if strings.HasSuffix(strings.ToLower(url), ".gif") {
		return url
	}
	sep := config.Instance.Uploader.AliyunOss.StyleSplitter
	if simple.IsBlank(sep) {
		return url
	}
	return strings.Join([]string{url, style}, sep)
}
