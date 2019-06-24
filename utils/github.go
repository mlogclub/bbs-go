package utils

import (
	"github.com/mlogclub/simple"
	"github.com/sirupsen/logrus"
	"golang.org/x/oauth2"
	"gopkg.in/resty.v1"
)

var GithubOauthConfig *oauth2.Config

func initGithubConfig() {
	GithubOauthConfig = &oauth2.Config{
		ClientID:     Conf.GithubClientID,
		ClientSecret: Conf.GithubClientSecret,
		RedirectURL:  Conf.BaseUrl + "/user/github/callback",
		Scopes:       []string{"public_repo", "user"},
		Endpoint: oauth2.Endpoint{
			AuthURL:  "https://github.com/login/oauth/authorize",
			TokenURL: "https://github.com/login/oauth/access_token",
		},
	}
}

type GithubUserinfo struct {
	Id        int64  `json:"id"`
	Login     string `json:"login"`
	NodeId    string `json:"node_id"`
	AvatarUrl string `json:"avatar_url"`
	Url       string `json:"url"`
	HtmlUrl   string `json:"html_url"`
	Email     string `json:"email"`
	Name      string `json:"name"`
}

func GetGithubUserinfo(accessToken string) *GithubUserinfo {
	response, err := resty.R().SetQueryParam("access_token", accessToken).Get("https://api.github.com/user")
	if err != nil {
		logrus.Errorf("Get user info error %s", err)
	}
	content := string(response.Body())

	githubUserinfo := &GithubUserinfo{}
	simple.ParseJson(content, githubUserinfo)
	return githubUserinfo
}
