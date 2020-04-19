package uploader

import (
	"sync"

	"bbs-go/common/config"
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
	if config.Conf.Uploader.Enable == "aliyun" || config.Conf.Uploader.Enable == "Oss" {
		return aliyun
	} else {
		return local
	}
}
