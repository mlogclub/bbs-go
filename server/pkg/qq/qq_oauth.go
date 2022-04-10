package qq

import (
	"errors"
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/mlogclub/simple/common/jsons"
	"github.com/mlogclub/simple/common/strs"
	"github.com/mlogclub/simple/common/urls"

	"github.com/go-resty/resty/v2"
	"github.com/goburrow/cache"
	"github.com/sirupsen/logrus"
	"github.com/tidwall/gjson"

	"bbs-go/pkg/common"
	"bbs-go/pkg/config"
)

type UserInfo struct {
	Ret          int    `json:"ret"`            // 返回码
	Msg          string `json:"msg"`            // 如果ret<0，会有相应的错误信息提示，返回数据全部用UTF-8编码。
	Nickname     string `json:"nickname"`       // 用户在QQ空间的昵称。
	Figureurl    string `json:"figureurl"`      // 大小为30×30像素的QQ空间头像URL。
	Figureurl1   string `json:"figureurl_1"`    // 大小为50×50像素的QQ空间头像URL。
	Figureurl2   string `json:"figureurl_2"`    // 大小为100×100像素的QQ空间头像URL。
	FigureurlQQ1 string `json:"figureurl_qq_1"` // 大小为40×40像素的QQ头像URL。
	FigureurlQQ2 string `json:"figureurl_qq_2"` // 大小为100×100像素的QQ头像URL。需要注意，不是所有的用户都拥有QQ的100x100的头像，但40x40像素则是一定会有。
	Gender       string `json:"gender"`         // 性别。 如果获取不到则默认返回"男"
	Openid       string `json:"openid"`
	Unionid      string `json:"unionid"`
}

type AccessToken struct {
	AccessToken  string `json:"access_token"`
	ExpiresIn    int    `json:"expires_in"`
	RefreshToken string `json:"refresh_token"`
}

var ctxCache = cache.New(cache.WithMaximumSize(1000), cache.WithExpireAfterAccess(10*time.Minute))

// 获取authorize url
// 文档：https://wiki.connect.qq.com/%E4%BD%BF%E7%94%A8authorization_code%E8%8E%B7%E5%8F%96access_token
// 接口：https://graph.qq.com/oauth2.0/authorize
func AuthorizeUrl(params map[string]string) string {
	// 将跳转地址写入上线文
	state := strs.UUID()
	redirectUrl := getRedirectUrl(params)
	ctxCache.Put(state, redirectUrl)

	return urls.ParseUrl("https://graph.qq.com/oauth2.0/authorize").
		AddQuery("response_type", "code").
		AddQuery("client_id", config.Instance.QQConnect.AppId).
		AddQuery("redirect_uri", redirectUrl).
		AddQuery("state", state).
		AddQuery("scope", "get_user_info").
		BuildStr()
}

// code -> accessToken
// 文档：https://wiki.connect.qq.com/%E4%BD%BF%E7%94%A8authorization_code%E8%8E%B7%E5%8F%96access_token
// 接口：https://graph.qq.com/oauth2.0/token
func AuthorizationCode(code, state string) (*AccessToken, error) {
	// 从上下文中获取跳转地址
	val, found := ctxCache.GetIfPresent(state)
	var redirectUrl string
	if found {
		redirectUrl = val.(string)
	}

	resp, err := resty.New().R().
		SetQueryParam("grant_type", "authorization_code").
		SetQueryParam("client_id", config.Instance.QQConnect.AppId).
		SetQueryParam("client_secret", config.Instance.QQConnect.AppKey).
		SetQueryParam("code", code).
		SetQueryParam("redirect_uri", redirectUrl).
		Get("https://graph.qq.com/oauth2.0/token")
	if err != nil {
		return nil, err
	}
	content := string(resp.Body())

	fmt.Println("token:" + content)

	// qq返回的数据格式如下：
	// access_token=xxx&expires_in=7776000&refresh_token=xxx
	ub := urls.ParseUrl("?" + content)
	accessToken := ub.GetQuery().Get("access_token")
	refreshToken := ub.GetQuery().Get("refresh_token")
	expiresIn, _ := strconv.Atoi(ub.GetQuery().Get("expires_in"))

	return &AccessToken{
		AccessToken:  accessToken,
		ExpiresIn:    expiresIn,
		RefreshToken: refreshToken,
	}, nil
}

// 文档：https://wiki.connect.qq.com/%E8%8E%B7%E5%8F%96%E7%94%A8%E6%88%B7openid_oauth2-0
// 接口：https://graph.qq.com/oauth2.0/me
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
	content = removeCallback(content)

	logrus.Info("me:" + content)

	return gjson.Get(content, "openid").String(), gjson.Get(content, "unionid").String(), nil
}

// 获取用户信息
// 文档：https://wiki.connect.qq.com/get_user_info
// 接口：https://graph.qq.com/user/get_user_info
func GetUserInfo(accessToken string) (*UserInfo, error) {
	openid, unionid, err := GetOpenid(accessToken)
	if err != nil {
		return nil, err
	}
	resp, err := resty.New().R().
		SetQueryParam("access_token", accessToken).
		SetQueryParam("oauth_consumer_key", config.Instance.QQConnect.AppId).
		SetQueryParam("openid", openid).
		Get("https://graph.qq.com/user/get_user_info")
	if err != nil {
		return nil, err
	}
	content := string(resp.Body())

	logrus.Info("get_user_info:" + content)

	ret := gjson.Get(content, "ret").Int()
	msg := gjson.Get(content, "msg").String()

	if ret != 0 {
		return nil, errors.New("get_user_info:ret=" + strconv.FormatInt(ret, 10) + ",msg=" + msg)
	}

	userInfo := &UserInfo{}
	err = jsons.Parse(content, userInfo)
	if err != nil {
		return nil, err
	} else {
		userInfo.Openid = openid
		userInfo.Unionid = unionid
		return userInfo, nil
	}
}

// 根据code获取用户信息
// 流程为先使用code换取accessToken，然后根据accessToken获取用户信息
func GetUserInfoByCode(code, state string) (*UserInfo, error) {
	token, err := AuthorizationCode(code, state)
	if err != nil {
		return nil, err
	}
	return GetUserInfo(token.AccessToken)
}

// 获取回调跳转地址
func getRedirectUrl(params map[string]string) string {
	redirectUrl := config.Instance.BaseUrl + "/user/qq/callback"
	if !common.IsProd() {
		redirectUrl = "http://localhost:3000/user/qq/callback"
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

// qq有些接口返回的数据带了callback，例如：callback( {"error":100020,"error_description":"code is reused error"} );
// 这里将callback去掉
func removeCallback(content string) string {
	prefix := "callback("
	suffix := ");"
	content = strings.TrimSpace(content)
	if strings.Index(content, "callback(") == 0 {
		content = content[len(prefix) : len(content)-len(suffix)]
		content = strings.TrimSpace(content)
	}
	return content
}
