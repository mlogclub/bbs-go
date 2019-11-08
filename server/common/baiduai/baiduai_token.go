package baiduai

import (
	"errors"

	"github.com/go-resty/resty/v2"
	"github.com/mlogclub/simple"
)

const (
	authUrl = "https://openapi.baidu.com/oauth/2.0/token"
)

// Authorizer 用于设置access_token
// 可以通过RESTFul api的方式从百度方获取
// 有效期为一个月，可以存至数据库中然后从数据库中获取
type Authorizer interface {
	Authorize(*Client) error
}

type Client struct {
	ClientID     string
	ClientSecret string
	AccessToken  string
	Authorizer   Authorizer
}

type AuthResponse struct {
	AccessToken      string `json:"access_token"`  // 要获取的Access Token
	ExpireIn         string `json:"expire_in"`     // Access Token的有效期(秒为单位，一般为1个月)；
	RefreshToken     string `json:"refresh_token"` // 以下参数忽略，暂时不用
	Scope            string `json:"scope"`
	SessionKey       string `json:"session_key"`
	SessionSecret    string `json:"session_secret"`
	ERROR            string `json:"error"`             // 错误码；关于错误码的详细信息请参考鉴权认证错误码(http://ai.baidu.com/docs#/Auth/top)
	ErrorDescription string `json:"error_description"` // 错误描述信息，帮助理解和解决发生的错误。
}

type DefaultAuthorizer struct{}

func (da DefaultAuthorizer) Authorize(client *Client) error {
	resp, err := resty.New().R().SetFormData(map[string]string{
		"grant_type":    "client_credentials",
		"client_id":     client.ClientID,
		"client_secret": client.ClientSecret,
	}).Post(authUrl)
	if err != nil {
		return err
	}
	authResponse := new(AuthResponse)
	if err := simple.ParseJson(string(resp.Body()), authResponse); err != nil {
		return err
	}
	if authResponse.ERROR != "" || authResponse.AccessToken == "" {
		return errors.New("授权失败:" + authResponse.ErrorDescription)
	}
	client.AccessToken = authResponse.AccessToken
	return nil
}

func (client *Client) Auth() error {
	if client.AccessToken != "" {
		// return nil
	}
	if err := client.Authorizer.Authorize(client); err != nil {
		return err
	}
	return nil
}

func (client *Client) SetAuther(auth Authorizer) {
	client.Authorizer = auth
}

func NewClient(apiKey, secretKey string) *Client {
	return &Client{
		ClientID:     apiKey,
		ClientSecret: secretKey,
		Authorizer:   DefaultAuthorizer{},
	}
}
