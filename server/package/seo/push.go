package seo

import (
	"github.com/mlogclub/simple"
	"strings"

	"github.com/go-resty/resty/v2"
	"github.com/sirupsen/logrus"

	"bbs-go/package/config"
)

func Push(url string) {
	go func() {
		PushUrls([]string{url})
		PushSmUrls([]string{url})
	}()
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

// 神马链接推送
func PushSmUrls(urls []string) {
	if urls == nil || len(urls) == 0 {
		return
	}
	conf := config.Instance.SmSEO
	if simple.IsBlank(conf.Site) || simple.IsBlank(conf.UserName) || simple.IsBlank(conf.Token) {
		return
	}

	u := simple.ParseUrl("https://data.zhanzhang.sm.cn/push")
	u.AddQuery("site", conf.Site)
	u.AddQuery("user_name", conf.UserName)
	u.AddQuery("resource_name", "mip_add")
	u.AddQuery("token", conf.Token)

	body := strings.Join(urls, "\n")
	if response, err := resty.New().R().SetBody(body).Post(u.BuildStr()); err != nil {
		logrus.Error(err)
	} else {
		logrus.Info("神马搜索链接推送完成：", string(response.Body()))
	}
}
