package oss

import (
	"bytes"
	"strings"
	"time"

	"github.com/aliyun/aliyun-oss-go-sdk/oss"
	"github.com/go-resty/resty/v2"
	"github.com/mlogclub/simple"
	"github.com/sirupsen/logrus"

	"github.com/mlogclub/bbs-go/common/config"
)

// 上传图片
func PutImage(data []byte) (string, error) {
	md5 := simple.MD5Bytes(data)
	key := "images/" + simple.TimeFormat(time.Now(), "2006/01/02/") + md5 + ".jpg"
	return PutObject(key, data)
}

// 上传
func PutObject(key string, data []byte) (string, error) {
	bucket, err := getBucket()
	if err != nil {
		return "", err
	}
	err = bucket.PutObject(key, bytes.NewReader(data))
	if err != nil {
		return "", err
	}
	return config.Conf.AliyunOss.Host + key, nil
}

// 将图片copy到oss
func CopyImage(inputUrl string) (string, error) {
	data, err := download(inputUrl)
	if err != nil {
		return "", err
	}
	return PutImage(data)
}

// 下载
func download(url string) ([]byte, error) {
	rsp, err := resty.New().R().Get(url)
	if err != nil {
		return nil, err
	}
	return rsp.Body(), nil
}

// 图片url签名
func SignUrl(url string) string {
	// 非oss资源不进行签名
	host := config.Conf.AliyunOss.Host
	if strings.Index(url, host) == -1 {
		return url
	}

	bucket, err := getBucket()
	if err != nil {
		return url
	}
	key := getObjectKey(url)
	ret, err := bucket.SignURL(key, oss.HTTPGet, 60*3) // 签名，有效期3分钟
	if err != nil {
		logrus.Error(err)
		return key
	}
	urlBuilder := simple.ParseUrl(url)
	params := simple.ParseUrl(ret).GetQuery()
	for k := range params {
		v := params.Get(k)
		urlBuilder.AddQuery(k, v)
	}
	return urlBuilder.BuildStr()
}

// 根据URL获取ObjectKey
func getObjectKey(u string) string {
	urlBuilder := simple.ParseUrl(u)
	objectKey := urlBuilder.GetURL().Path
	objectKey = objectKey[1:]
	return objectKey
}

func getBucket() (*oss.Bucket, error) {
	client, err := oss.New(config.Conf.AliyunOss.Endpoint, config.Conf.AliyunOss.AccessId, config.Conf.AliyunOss.AccessSecret)
	if err != nil {
		return nil, err
	}
	return client.Bucket(config.Conf.AliyunOss.Bucket)
}
