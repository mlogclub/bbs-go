package github

import (
	"github.com/mlogclub/mlog/utils/config"
	"github.com/mlogclub/simple"
	"github.com/sirupsen/logrus"
	"golang.org/x/oauth2"
	"gopkg.in/resty.v1"
)

var OauthConfig *oauth2.Config

func InitConfig() {
	OauthConfig = &oauth2.Config{
		ClientID:     config.Conf.Github.ClientID,
		ClientSecret: config.Conf.Github.ClientSecret,
		RedirectURL:  config.Conf.BaseUrl + "/user/github/callback",
		Scopes:       []string{"public_repo", "user"},
		Endpoint: oauth2.Endpoint{
			AuthURL:  "https://github.com/login/oauth/authorize",
			TokenURL: "https://github.com/login/oauth/access_token",
		},
	}
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
}

func GetUserInfo(accessToken string) (*UserInfo, error) {
	response, err := resty.R().SetQueryParam("access_token", accessToken).Get("https://api.github.com/user")
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
