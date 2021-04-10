package uploader

import (
	"strings"
	"sync"

	"github.com/mlogclub/simple"

	"bbs-go/package/config"
)

type uploader interface {
	PutImage(data []byte) (string, error)
	PutObject(key string, data []byte) (string, error)
	CopyImage(originUrl string) (string, error)
}

var (
	aliyun = &aliyunOssUploader{
		once:   sync.Once{},
		bucket: nil,
	}
	local = &localUploader{}
)

func PutImage(data []byte) (string, error) {
	return getUploader().PutImage(data)
}

func PutObject(key string, data []byte) (string, error) {
	return getUploader().PutObject(key, data)
}

func CopyImage(url string) (string, error) {
	u1 := simple.ParseUrl(url).GetURL()
	u2 := simple.ParseUrl(config.Instance.BaseUrl).GetURL()
	// 本站host，不下载
	if u1.Host == u2.Host {
		return url, nil
	}
	return getUploader().CopyImage(url)
}

func getUploader() uploader {
	if IsEnabledOss() {
		return aliyun
	} else {
		return local
	}
}

// IsEnabledOss 是否启用阿里云oss
func IsEnabledOss() bool {
	enable := config.Instance.Uploader.Enable
	return simple.EqualsIgnoreCase(enable, "aliyun") || simple.EqualsIgnoreCase(enable, "oss") ||
		simple.EqualsIgnoreCase(enable, "aliyunOss")
}

// IsOssImageUrl 是否是存放在阿里云oss中的图片
func IsOssImageUrl(url string) bool {
	host := simple.ParseUrl(config.Instance.Uploader.AliyunOss.Host).GetURL().Host
	return strings.Contains(url, host)
}
