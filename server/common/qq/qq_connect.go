package qq

import (
	"context"
	"errors"
	"strconv"

	"github.com/mlogclub/simple"
	"github.com/sirupsen/logrus"
	"github.com/tidwall/gjson"
	"golang.org/x/oauth2"
	"gopkg.in/resty.v1"

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

type UserInfo struct {
	Ret          string `json:"ret"`            // 返回码
	Msg          string `json:"msg"`            // 如果ret<0，会有相应的错误信息提示，返回数据全部用UTF-8编码。
	Nickname     string `json:"nickname"`       // 用户在QQ空间的昵称。
	Figureurl    string `json:"figureurl"`      // 大小为30×30像素的QQ空间头像URL。
	Figureurl1   string `json:"figureurl_1"`    // 大小为50×50像素的QQ空间头像URL。
	Figureurl2   string `json:"figureurl_2"`    // 大小为100×100像素的QQ空间头像URL。
	FigureurlQQ1 string `json:"figureurl_qq_1"` // 大小为40×40像素的QQ头像URL。
	FigureurlQQ2 string `json:"figureurl_qq_2"` // 大小为100×100像素的QQ头像URL。需要注意，不是所有的用户都拥有QQ的100x100的头像，但40x40像素则是一定会有。
	Gender       string `json:"gender"`         // 性别。 如果获取不到则默认返回"男"
	OpenId       string `json:"openId"`
	UnionId      string `json:"unionId"`
}

// 根据code获取用户信息
// 流程为先使用code换取accessToken，然后根据accessToken获取用户信息
func GetUserInfoByCode(code string) (*UserInfo, error) {
	token, err := GetOauthConfig(nil).Exchange(context.TODO(), code)
	if err != nil {
		return nil, err
	}
	return GetUserInfo(token.AccessToken)
}

// 获取用户信息
func GetUserInfo(accessToken string) (*UserInfo, error) {
	openid, unionid, err := GetOpenid(accessToken)
	if err != nil {
		return nil, err
	}
	resp, err := resty.New().R().
		SetQueryParam("access_token", accessToken).
		SetQueryParam("oauth_consumer_key", config.Conf.QQConnect.AppId).
		SetQueryParam("openid", openid).
		Get("https://graph.qq.com/user/get_user_info")
	content := string(resp.Body())

	logrus.Info("get_user_info:" + content)

	ret := gjson.Get(content, "ret").Int()
	msg := gjson.Get(content, "msg").String()

	if ret != 0 {
		return nil, errors.New("get_user_info:ret=" + strconv.FormatInt(ret, 10) + ",msg=" + msg)
	}

	userInfo := &UserInfo{}
	err = simple.ParseJson(content, userInfo)
	if err != nil {
		return nil, err
	} else {
		userInfo.OpenId = openid
		userInfo.UnionId = unionid
		return userInfo, nil
	}
}

// 获取openId
func GetOpenid(accessToken string) (string, string, error) {
	resp, err := resty.New().R().
		SetQueryParam("access_token", accessToken).
		SetQueryParam("unionid", "1"). // 申请unionId，0：不申请，1：申请
		Get("https://graph.qq.com/oauth2.0/me")
	if err != nil {
		logrus.Errorf("QQ: Get openid error", err)
		return "", "", err
	}
	content := string(resp.Body())

	logrus.Info("me:" + content)

	return gjson.Get(content, "openid").String(), gjson.Get(content, "unionid").String(), nil
}
