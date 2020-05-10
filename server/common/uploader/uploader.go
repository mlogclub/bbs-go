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
	miuploader= &minioUploader{}
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


	enable := config.Conf.Uploader.Enable
	if (enable=="Aliyun") {
		return aliyun
	}else if (enable=="Minio"){
println("get uploader minio")
		return  miuploader
	}	else {
		return local
	}
}
