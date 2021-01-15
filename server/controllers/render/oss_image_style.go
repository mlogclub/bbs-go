package render

import (
	"bbs-go/common/uploader"
	"bbs-go/config"
	"github.com/mlogclub/simple"
	"strings"
)

func HandleOssImageStyleAvatar(url string) string {
	if uploader.EnabledObjectStorage() == -1 {
		return url
	}
	return HandleOssImageStyle(url, config.Instance.Uploader.ObjectStorage.StyleAvatar)
}

func HandleOssImageStyleDetail(url string) string {
	if uploader.EnabledObjectStorage() == -1 {
		return url
	}
	return HandleOssImageStyle(url, config.Instance.Uploader.ObjectStorage.StyleDetail)
}

func HandleOssImageStyleSmall(url string) string {
	if uploader.EnabledObjectStorage() == -1 {
		return url
	}
	return HandleOssImageStyle(url, config.Instance.Uploader.ObjectStorage.StyleSmall)
}

func HandleOssImageStylePreview(url string) string {
	if uploader.EnabledObjectStorage() == -1 {
		return url
	}
	return HandleOssImageStyle(url, config.Instance.Uploader.ObjectStorage.StylePreview)
}

func HandleOssImageStyle(url, style string) string {
	if simple.IsBlank(style) || simple.IsBlank(url) {
		return url
	}
	if uploader.EnabledObjectStorage() == -1 {
		return url
	}
	sep := config.Instance.Uploader.ObjectStorage.StyleSplitter
	if simple.IsBlank(sep) {
		return url
	}
	return strings.Join([]string{url, style}, sep)
}
