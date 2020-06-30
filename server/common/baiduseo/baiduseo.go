package baiduseo

import (
	"strings"

	"github.com/go-resty/resty/v2"
	"github.com/sirupsen/logrus"

	"bbs-go/config"
)

func PushUrl(url string) {
	PushUrls([]string{url})
}

// 百度链接推送
func PushUrls(urls []string) {
	if urls == nil || len(urls) == 0 {
		return
	}
	if len(config.Instance.BaiduSEO.Site) == 0 || len(config.Instance.BaiduSEO.Token) == 0 {
		return
	}
	api := "http://data.zz.baidu.com/urls?site=" + config.Instance.BaiduSEO.Site + "&token=" +
		config.Instance.BaiduSEO.Token
	body := strings.Join(urls, "\n")
	if response, err := resty.New().R().SetBody(body).Post(api); err != nil {
		logrus.Error(err)
	} else {
		logrus.Info("百度链接提交完成：", string(response.Body()))
	}
}
