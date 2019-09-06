package github

import (
	"github.com/go-resty/resty/v2"
	"github.com/mlogclub/simple"
	"github.com/sirupsen/logrus"
	"golang.org/x/oauth2"

	"github.com/mlogclub/bbs-go/common"
	"github.com/mlogclub/bbs-go/common/config"
)

var oauthConfig *oauth2.Config

// params callback携带的参数
func GetOauthConfig(params map[string]string) *oauth2.Config {
	if oauthConfig == nil {
		oauthConfig = &oauth2.Config{
			ClientID:     config.Conf.Github.ClientID,
			ClientSecret: config.Conf.Github.ClientSecret,
			RedirectURL:  getRedirectUrl(nil),
			Scopes:       []string{"public_repo", "user"},
			Endpoint: oauth2.Endpoint{
				AuthURL:  "https://github.com/login/oauth/authorize",
				TokenURL: "https://github.com/login/oauth/access_token",
			},
		}
	}
	oauthConfig.RedirectURL = getRedirectUrl(params)
	return oauthConfig
}

// 获取回调跳转地址
func getRedirectUrl(params map[string]string) string {
	redirectUrl := config.Conf.BaseUrl + "/user/github/callback"
	if !common.IsProd() {
		redirectUrl = "http://localhost:3000/user/github/callback"
	}
	if len(params) > 0 {
		ub := simple.ParseUrl(redirectUrl)
		for k, v := range params {
			ub.AddQuery(k, v)
		}
		redirectUrl = ub.BuildStr()
	}
	return redirectUrl
}

type UserInfo struct {
	Id        int64  `json:"id"`
	Login     string `json:"login"`
	NodeId    string `json:"node_id"`
	AvatarUrl string `json:"avatar_url"`
	Url       string `json:"url"`
	HtmlUrl   string `json:"html_url"`
	Email     string `json:"email"`
	Name      string `json:"name"`
	Bio       string `json:"bio"`
	Company   string `json:"company"`
	Blog      string `json:"blog"`
	Location  string `json:"location"`
}

func GetUserInfo(accessToken string) (*UserInfo, error) {
	response, err := resty.New().R().SetQueryParam("access_token", accessToken).Get("https://api.github.com/user")
	if err != nil {
		logrus.Errorf("Get user info error %s", err)
		return nil, err
	}
	content := string(response.Body())

	userInfo := &UserInfo{}
	err = simple.ParseJson(content, userInfo)
	if err != nil {
		return nil, err
	}
	return userInfo, nil
}
