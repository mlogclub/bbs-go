package baiduseo

import (
	"strings"

	"github.com/sirupsen/logrus"
	"gopkg.in/resty.v1"

	"github.com/mlogclub/bbs-go/common/config"
)

func PushUrl(url string) {
	PushUrls([]string{url})
}

// 百度链接推送
func PushUrls(urls []string) {
	if urls == nil || len(urls) == 0 {
		return
	}
	if len(config.Conf.BaiduSEO.Site) == 0 || len(config.Conf.BaiduSEO.Token) == 0 {
		return
	}
	api := "http://data.zz.baidu.com/urls?site=" + config.Conf.BaiduSEO.Site + "&token=" + config.Conf.BaiduSEO.Token
	body := strings.Join(urls, "\n")
	if response, err := resty.R().SetBody(body).Post(api); err != nil {
		logrus.Error(err)
	} else {
		logrus.Info("百度链接提交完成：", string(response.Body()))
	}
}
