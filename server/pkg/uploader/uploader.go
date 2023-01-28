package uploader

import (
	"strings"
	"sync"

	"github.com/mlogclub/simple/common/strs"
	"github.com/mlogclub/simple/common/urls"

	"server/pkg/config"
)

type uploader interface {
	PutImage(data []byte, contentType string) (string, error)
	PutObject(key string, data []byte, contentType string) (string, error)
	CopyImage(originUrl string) (string, error)
}

var (
	aliyun = &aliyunOssUploader{
		once:   sync.Once{},
		bucket: nil,
	}
	txy = &txyCosUploader{
		once:   sync.Once{},
		object: nil,
	}
	server = &serverUploader{}
	local  = &localUploader{}
)

func PutImage(data []byte, contentType string) (string, error) {
	return getUploader().PutImage(data, contentType)
}

func PutObject(key string, data []byte, contentType string) (string, error) {
	return getUploader().PutObject(key, data, contentType)
}

func CopyImage(url string) (string, error) {
	u1 := urls.ParseUrl(url).GetURL()
	u2 := urls.ParseUrl(config.Instance.BaseUrl).GetURL()
	// 本站host，不下载
	if u1.Host == u2.Host {
		return url, nil
	}
	return getUploader().CopyImage(url)
}

func getUploader() uploader {
	enable := config.Instance.Uploader.Enable
	switch enable {
	case "aly":
		return aliyun
	case "txy":
		return txy
	case "server":
		return server
	default:
		return local
	}
}

// IsEnabledOss 是否启用阿里云oss
func IsEnabledOss() bool {
	enable := config.Instance.Uploader.Enable
	return strs.EqualsIgnoreCase(enable, "aliyun") || strs.EqualsIgnoreCase(enable, "oss") ||
		strs.EqualsIgnoreCase(enable, "aliyunOss")
}

// IsOssImageUrl 是否是存放在阿里云oss中的图片
func IsOssImageUrl(url string) bool {
	host := urls.ParseUrl(config.Instance.Uploader.AliyunOss.Host).GetURL().Host
	return strings.Contains(url, host)
}
