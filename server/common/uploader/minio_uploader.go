package uploader

import (
	"bbs-go/common/config"
	"os"
	"strconv"

	"time"

	"github.com/go-resty/resty/v2"
	"github.com/minio/minio-go/v6"
	"github.com/mlogclub/simple"
	"github.com/sirupsen/logrus"
)

type minioUploader struct {
	Bucket string `bucket`
}

// 上传图片
func (*minioUploader) PutImage(data []byte) (string, error) {
	md5 := simple.MD5Bytes(data)
	key := "images/" + simple.TimeFormat(time.Now(), "2006/01/02/") + md5 + ".jpg"
	return PutObject(key, data)
}

// 上传
func (miu *minioUploader) PutObject(key string, d []byte) (string, error) {
	client, err := miu.getClient()
	if err != nil {
		return "", err
	}
	tepPath := Save(d)
	defer os.Remove(tepPath)
	n, err := client.FPutObject(miu.Bucket, key, tepPath, minio.PutObjectOptions{})
	logrus.Info(n)
	if err != nil {
		return "", err
	}
	return config.Conf.Uploader.Minio.Host + key + "?token=", nil
}

//将图片copy到oss
func (*minioUploader) CopyImage(inputUrl string) (string, error) {
	data, err := download(inputUrl)
	if err != nil {
		return "", err
	}
	return PutImage(data)
}

// 下载
func (*minioUploader) download(url string) ([]byte, error) {
	rsp, err := resty.New().R().Get(url)
	if err != nil {
		return nil, err
	}
	return rsp.Body(), nil
}

func Save(d []byte) string {
	tempFile := strconv.Itoa(time.Now().Nanosecond())
	file, err := os.OpenFile(tempFile, os.O_RDWR|os.O_APPEND|os.O_CREATE, 066) // For read access.
	//2、关闭文件

	// 3、写入数据
	count, err := file.Write(d)

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

func (mi *minioUploader) getClient() (*minio.Client, error) {
	//client,err := minio.New(endpoint, accessKeyID, secretAccessKey, useSSL)
	client, err := minio.New(config.Conf.Uploader.Minio.Endpoint, config.Conf.Uploader.Minio.AccessId, config.Conf.Uploader.Minio.AccessSecret, false)
	if err != nil {
		print(err)
		return nil, err
	}
	mi.Bucket=config.Conf.Uploader.Minio.Bucket
	return client, err
}
