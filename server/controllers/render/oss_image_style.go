package render

import (
	"bbs-go/package/config"
	uploader2 "bbs-go/package/uploader"
	"github.com/mlogclub/simple"
	"strings"
)

func HandleOssImageStyleAvatar(url string) string {
	if !uploader2.IsEnabledOss() {
		return url
	}
	return HandleOssImageStyle(url, config.Instance.Uploader.AliyunOss.StyleAvatar)
}

func HandleOssImageStyleDetail(url string) string {
	if !uploader2.IsEnabledOss() {
		return url
	}
	return HandleOssImageStyle(url, config.Instance.Uploader.AliyunOss.StyleDetail)
}

func HandleOssImageStyleSmall(url string) string {
	if !uploader2.IsEnabledOss() {
		return url
	}
	return HandleOssImageStyle(url, config.Instance.Uploader.AliyunOss.StyleSmall)
}

func HandleOssImageStylePreview(url string) string {
	if !uploader2.IsEnabledOss() {
		return url
	}
	return HandleOssImageStyle(url, config.Instance.Uploader.AliyunOss.StylePreview)
}

func HandleOssImageStyle(url, style string) string {
	if simple.IsBlank(style) || simple.IsBlank(url) {
		return url
	}
	if !uploader2.IsOssImageUrl(url) {
		return url
	}
	sep := config.Instance.Uploader.AliyunOss.StyleSplitter
	if simple.IsBlank(sep) {
		return url
	}
	return strings.Join([]string{url, style}, sep)
}
