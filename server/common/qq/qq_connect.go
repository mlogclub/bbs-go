package qq

import (
	"github.com/mlogclub/simple"
	"golang.org/x/oauth2"

	"github.com/mlogclub/bbs-go/common"
	"github.com/mlogclub/bbs-go/common/config"
)

var oauthConfig *oauth2.Config

func GetOauthConfig(params map[string]string) *oauth2.Config {
	if oauthConfig == nil {
		oauthConfig = &oauth2.Config{
			ClientID:     config.Conf.QQConnect.AppId,
			ClientSecret: config.Conf.QQConnect.AppKey,
			RedirectURL:  getRedirectUrl(nil),
			Scopes:       []string{"public_repo", "user"},
			Endpoint: oauth2.Endpoint{
				AuthURL:  "https://graph.qq.com/oauth2.0/authorize",
				TokenURL: "https://graph.qq.com/oauth2.0/token",
			},
		}
	}
	oauthConfig.RedirectURL = getRedirectUrl(params)
	return oauthConfig
}

// 获取回调跳转地址
func getRedirectUrl(params map[string]string) string {
	redirectUrl := config.Conf.BaseUrl + "/user/qq/callback"
	if !common.IsProd() {
		redirectUrl = "http://localhost:3000/user/qq/callback"
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

func GetUserInfo() {

}
