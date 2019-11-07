package common

import (
	"strings"

	"github.com/sirupsen/logrus"
	"gopkg.in/resty.v1"
)

// 百度链接推送
func BaiduUrlPush(urls []string) {
	if urls == nil || len(urls) == 0 {
		return
	}
	api := "http://data.zz.baidu.com/urls?site=mlog.club&token=Y0b5JDk00GSjyMu3"
	body := strings.Join(urls, "\n")
	response, err := resty.R().SetBody(body).Post(api)
	if err != nil {
		logrus.Error(err)
	} else {
		logrus.Info("百度链接提交完成：", string(response.Body()))
	}
}
