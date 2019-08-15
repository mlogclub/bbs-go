package oss

// import (
// 	"bytes"
// 	"time"
//
// 	"github.com/aliyun/aliyun-oss-go-sdk/oss"
// 	"github.com/mlogclub/simple"
// 	"github.com/sirupsen/logrus"
// 	"gopkg.in/resty.v1"
//
// 	"github.com/mlogclub/mlog/utils/config"
// )
//
// var AliyunOss *aliyunOss
//
// func init() {
// 	AliyunOss, _ = New(config.Conf.AliyunOss.Endpoint, config.Conf.AliyunOss.AccessId,
// 		config.Conf.AliyunOss.AccessSecret, config.Conf.AliyunOss.Bucket, config.Conf.AliyunOss.Host)
// }
//
// type aliyunOss struct {
// 	host   string
// 	client *oss.Client
// 	bucket *oss.Bucket
// }
//
// func New(endpoint, accessId, accessKey, bucketName, host string) (*aliyunOss, error) {
// 	client, err := oss.New(endpoint, accessId, accessKey, oss.UseCname(true))
// 	if err != nil {
// 		return nil, err
// 	}
// 	bucket, err := client.Bucket(bucketName)
// 	if err != nil {
// 		return nil, err
// 	}
// 	return &aliyunOss{host: host, client: client, bucket: bucket}, nil
// }
//
// // 上传图片
// func (this *aliyunOss) UploadImage(data []byte) (string, error) {
// 	md5 := simple.MD5Bytes(data)
// 	key := "images/" + simple.TimeFormat(time.Now(), "2006/01/02/") + md5 + ".jpg"
// 	return this.Upload(key, data)
// }
//
// // 将图片copy到oss
// func (this *aliyunOss) CopyImage(inputUrl string) (string, error) {
// 	data, err := this.download(inputUrl)
// 	if err != nil {
// 		return "", err
// 	}
// 	return this.UploadImage(data)
// }
//
// // 上传
// func (this *aliyunOss) Upload(key string, data []byte) (string, error) {
// 	err := this.bucket.PutObject(key, bytes.NewReader(data))
// 	if err != nil {
// 		return "", err
// 	}
// 	return this.host + key, nil
// }
//
// // 图片url签名
// func (this *aliyunOss) SignUrl(objectKey string) string {
// 	ret, err := this.bucket.SignURL(objectKey, oss.HTTPGet, 3600)
// 	if err != nil {
// 		logrus.Error(err)
// 		return objectKey
// 	}
// 	return ret
// }
//
// // 下载
// func (this *aliyunOss) download(url string) ([]byte, error) {
// 	rsp, err := resty.R().Get(url)
// 	if err != nil {
// 		return nil, err
// 	}
// 	return rsp.Body(), nil
// }
