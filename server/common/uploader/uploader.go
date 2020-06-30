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

func CopyImage(originUrl string) (string, error) {
	return getUploader().CopyImage(originUrl)
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
