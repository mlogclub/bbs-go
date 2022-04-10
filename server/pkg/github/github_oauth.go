package github

import (
	"context"
	"time"

	"github.com/mlogclub/simple/common/jsons"
	"github.com/mlogclub/simple/common/strs"
	"github.com/mlogclub/simple/common/urls"

	"github.com/go-resty/resty/v2"
	"github.com/goburrow/cache"
	"github.com/sirupsen/logrus"
	"golang.org/x/oauth2"

	"bbs-go/pkg/common"
	"bbs-go/pkg/config"
)

var ctxCache = cache.New(cache.WithMaximumSize(1000), cache.WithExpireAfterAccess(10*time.Minute))

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

// params callback携带的参数
func newOauthConfig(redirectUrl string) *oauth2.Config {
	return &oauth2.Config{
		ClientID:     config.Instance.Github.ClientID,
		ClientSecret: config.Instance.Github.ClientSecret,
		RedirectURL:  redirectUrl,
		Scopes:       []string{},
		Endpoint: oauth2.Endpoint{
			AuthURL:  "https://github.com/login/oauth/authorize",
			TokenURL: "https://github.com/login/oauth/access_token",
		},
	}
}

func AuthCodeURL(params map[string]string) string {
	// 将跳转地址写入上线文
	state := strs.UUID()
	redirectUrl := getRedirectUrl(params)
	ctxCache.Put(state, redirectUrl)

	return newOauthConfig(redirectUrl).AuthCodeURL(state)
}

// 根据code获取用户信息
// 流程为先使用code换取accessToken，然后根据accessToken获取用户信息
func GetUserInfoByCode(code, state string) (*UserInfo, error) {
	// 从上下文中获取跳转地址
	val, found := ctxCache.GetIfPresent(state)
	var redirectUrl string
	if found {
		redirectUrl = val.(string)
	}

	token, err := newOauthConfig(redirectUrl).Exchange(context.TODO(), code)
	if err != nil {
		return nil, err
	}
	return GetUserInfo(token.AccessToken)
}

// 根据accessToken获取用户信息
func GetUserInfo(accessToken string) (*UserInfo, error) {
	response, err := resty.New().R().SetHeader("Authorization", "token "+accessToken).Get("https://api.github.com/user")
	if err != nil {
		logrus.Errorf("Get user info error %s", err)
		return nil, err
	}
	content := string(response.Body())

	userInfo := &UserInfo{}
	err = jsons.Parse(content, userInfo)
	if err != nil {
		return nil, err
	}
	return userInfo, nil
}

// 获取回调跳转地址
func getRedirectUrl(params map[string]string) string {
	redirectUrl := config.Instance.BaseUrl + "/user/github/callback"
	if !common.IsProd() {
		redirectUrl = "http://localhost:3000/user/github/callback"
	}
	if len(params) > 0 {
		ub := urls.ParseUrl(redirectUrl)
		for k, v := range params {
			ub.AddQuery(k, v)
		}
		redirectUrl = ub.BuildStr()
	}
	return redirectUrl
}
