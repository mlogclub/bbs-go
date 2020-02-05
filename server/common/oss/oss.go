package oss

import (
	"os"
	"strconv"
	"time"

	"github.com/go-resty/resty/v2"
	"github.com/minio/minio-go/v6"
	"github.com/mlogclub/simple"
	"github.com/sirupsen/logrus"

	"bbs-go/common/config"
)

// 上传图片
func PutImage(data []byte) (string, error) {
	md5 := simple.MD5Bytes(data)
	key := "images/" + simple.TimeFormat(time.Now(), "2006/01/02/") + md5 + ".jpg"
	return PutObject(key, data)
}

// 上传
func PutObject(key string, data []byte) (string, error) {
	client, err := getClient()
	if err != nil {
		return "", err
	}
	tepPath := save(data)
	defer os.Remove(tepPath)
	n, err := client.FPutObject("bbs", key, tepPath, minio.PutObjectOptions{})
	logrus.Info(n)
	if err != nil {
		return "", err
	}
	return config.Conf.Minio.Host + key, nil
}

//将图片copy到oss
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

func save(data []byte) string {
	tempFile := strconv.Itoa(time.Now().Nanosecond())
	file, err := os.OpenFile(tempFile, os.O_RDWR|os.O_APPEND|os.O_CREATE, 066) // For read access.
	//2、关闭文件

	// 3、写入数据
	count, err := file.Write(data)

	if err != nil {
		logrus.Error(simple.JsonErrorMsg(err.Error()))
	} else {
		logrus.Info(count)
	}
	defer file.Close()
	return tempFile
}

// 图片url签名
// func SignUrl(url string) string {
// 	// 非oss资源不进行签名
// 	host := config.Conf.Minio.Host
// 	if strings.Index(url, host) == -1 {
// 		return url
// 	}

// 	bucket, err := getBucket()
// 	if err != nil {
// 		return url
// 	}
// 	key := getObjectKey(url)
// 	ret, err := bucket.SignURL(key, oss.HTTPGet, 60*3) // 签名，有效期3分钟
// 	if err != nil {
// 		logrus.Error(err)
// 		return key
// 	}
// 	urlBuilder := simple.ParseUrl(url)
// 	params := simple.ParseUrl(ret).GetQuery()
// 	for k := range params {
// 		v := params.Get(k)
// 		urlBuilder.AddQuery(k, v)
// 	}
// 	return urlBuilder.BuildStr()
// }

// 根据URL获取ObjectKey
func getObjectKey(u string) string {
	urlBuilder := simple.ParseUrl(u)
	objectKey := urlBuilder.GetURL().Path
	objectKey = objectKey[1:]
	return objectKey
}

func getClient() (*minio.Client, error) {
	//client,err := minio.New(endpoint, accessKeyID, secretAccessKey, useSSL)
	client, err := minio.New(config.Conf.Minio.Endpoint, config.Conf.Minio.AccessId, config.Conf.Minio.AccessSecret, false)
	if err != nil {
		return nil, err
	}
	return client, err
}
