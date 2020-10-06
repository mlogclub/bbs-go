package uploader

import (
	"sync"

	"github.com/mlogclub/simple"

	"bbs-go/config"
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
	enable := config.Instance.Uploader.Enable
	if simple.EqualsIgnoreCase(enable, "aliyun") || simple.EqualsIgnoreCase(enable, "oss") ||
		simple.EqualsIgnoreCase(enable, "aliyunOss") {
		return aliyun
	} else {
		return local
	}
}
