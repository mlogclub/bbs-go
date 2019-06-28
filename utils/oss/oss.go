package oss

import (
	"bytes"
	"github.com/aliyun/aliyun-oss-go-sdk/oss"
	"github.com/mlogclub/mlog/utils/config"
	"github.com/mlogclub/simple"
	"gopkg.in/resty.v1"
	"time"
)

func UploadImage(data []byte) (string, error) {
	md5 := simple.MD5Bytes(data)
	key := "images/" + simple.TimeFormat(time.Now(), "2006/01/02/") + md5 + ".jpg"
	return Upload(key, data)
}

func CopyImage(inputUrl string) (string, error) {
	data, err := download(inputUrl)
	if err != nil {
		return "", err
	}
	return UploadImage(data)
}

func Upload(key string, data []byte) (string, error) {
	// 创建OSSClient实例。
	client, err := oss.New(config.Conf.AliyunOss.Endpoint, config.Conf.AliyunOss.AccessId, config.Conf.AliyunOss.AccessSecret)
	if err != nil {
		return "", err
	}

	// 获取存储空间。
	bucket, err := client.Bucket(config.Conf.AliyunOss.Bucket)
	if err != nil {
		return "", err
	}

	// 上传字符串。
	err = bucket.PutObject(key, bytes.NewReader(data))
	if err != nil {
		return "", err
	}
	return config.Conf.AliyunOss.Host + key, nil
}

func download(url string) ([]byte, error) {
	rsp, err := resty.R().Get(url)
	if err != nil {
		return nil, err
	}
	return rsp.Body(), nil
}
