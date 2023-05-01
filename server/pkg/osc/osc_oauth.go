package osc

import (
	"time"

	"github.com/mlogclub/simple/common/jsons"
	"github.com/mlogclub/simple/common/strs"
	"github.com/mlogclub/simple/common/urls"
	"github.com/tidwall/gjson"

	"bbs-go/pkg/common"
	"bbs-go/pkg/config"

	"github.com/go-resty/resty/v2"
	"github.com/goburrow/cache"
	"github.com/sirupsen/logrus"
)

var ctxCache = cache.New(cache.WithMaximumSize(1000), cache.WithExpireAfterAccess(10*time.Minute))

type UserInfo struct {
	Id       int64  `json:"id"`
	Email    string `json:"email"`
	Name     string `json:"name"`
	Gender   string `json:"gender"`
	Avatar   string `json:"avatar"`
	Location string `json:"location"`
	Url      string `json:"url"`
}

func AuthCodeURL(params map[string]string) string {
	var (
		state       = strs.UUID()
		redirectUrl = getRedirectUrl(params)
	)
	ctxCache.Put(state, redirectUrl) // 将跳转地址写入上线文
	return urls.ParseUrl("https://www.oschina.net/action/oauth2/authorize").AddQueries(map[string]string{
		"client_id":     config.Instance.OSChina.ClientID,
		"response_type": "code",
		"redirect_uri":  redirectUrl,
		"state":         state,
	}).BuildStr()
}

// GetUserInfoByCode 根据code获取用户信息
// 流程为先使用code换取accessToken，然后根据accessToken获取用户信息
func GetUserInfoByCode(code, state string) (*UserInfo, error) {
	// 从上下文中获取跳转地址
	val, found := ctxCache.GetIfPresent(state)
	var redirectUrl string
	if found {
		redirectUrl = val.(string)
	}

	resp, err := resty.New().R().SetQueryParams(map[string]string{
		"client_id":     config.Instance.OSChina.ClientID,
		"client_secret": config.Instance.OSChina.ClientSecret,
		"grant_type":    "authorization_code",
		"redirect_uri":  redirectUrl,
		"code":          code,
		"dataType":      "json",
	}).Get("https://www.oschina.net/action/openapi/token")
	if err != nil {
		return nil, err
	}
	body := string(resp.Body())
	accessToken := gjson.Get(body, "access_token")
	return GetUserInfo(accessToken.String())
}

// GetUserInfo 根据accessToken获取用户信息
func GetUserInfo(accessToken string) (*UserInfo, error) {
	response, err := resty.New().R().SetQueryParam("access_token", accessToken).
		Get("https://www.oschina.net/action/openapi/user")
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
	redirectUrl := config.Instance.BaseUrl + "/user/osc/callback"
	if !common.IsProd() {
		redirectUrl = "http://localhost:3000/user/osc/callback"
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
