package oauth

import (
	"github.com/kataras/iris"
	"github.com/mlogclub/mlog/services"
	"github.com/mlogclub/mlog/utils/session"
	"gopkg.in/oauth2.v3/errors"
	"gopkg.in/oauth2.v3/manage"
	"gopkg.in/oauth2.v3/server"
	"log"
	"net/http"
	"net/url"
	"strconv"
	"time"
)

var ServerInstance = NewOauthServer()

type Server struct {
	Srv *server.Server
}

func NewOauthServer() *Server {
	oauthTokenStore := NewOauthTokenStore()
	oauthClientStore := NewOauthClientStore()

	manager := manage.NewDefaultManager()
	manager.MapTokenStorage(oauthTokenStore)
	manager.MapClientStorage(oauthClientStore)
	manager.SetAuthorizeCodeTokenCfg(&manage.Config{
		AccessTokenExp:    time.Hour * 24 * 30, // 访问令牌过期时间（默认为2小时）
		RefreshTokenExp:   time.Hour * 24 * 60, // 更新令牌过期时间（默认为72小时）
		IsGenerateRefresh: true,                // 是否生成更新令牌（默认为true）
	})

	srv := server.NewServer(server.NewConfig(), manager)
	srv.SetInternalErrorHandler(func(err error) (re *errors.Response) {
		log.Println("Internal Error:", err.Error())
		return
	})
	srv.SetResponseErrorHandler(func(re *errors.Response) {
		log.Println("Response Error:", re.Error.Error())
	})
	srv.SetPasswordAuthorizationHandler(func(username, password string) (userId string, err error) {
		user, err := services.UserServiceInstance.SignIn(username, password)
		if err != nil {
			err = errors.ErrAccessDenied
		} else {
			userId = strconv.FormatInt(user.Id, 10)
		}
		return
	})
	srv.SetUserAuthorizationHandler(func(w http.ResponseWriter, r *http.Request) (userID string, err error) {
		user := session.GetCurrentUserByRequest(w, r)
		if user == nil {
			redirectToLogin(w, r)
		} else {
			userID = strconv.FormatInt(user.Id, 10)
		}
		return
	})
	return &Server{Srv: srv}
}

// 跳转到登录页面
func redirectToLogin(w http.ResponseWriter, r *http.Request) {
	query := r.URL.Query()

	clientId := query.Get("client_id")
	redirectUri := query.Get("redirect_uri")
	responseType := query.Get("response_type")
	state := query.Get("state")

	v := url.Values{}
	v.Add("client_id", clientId)
	v.Add("redirect_uri", redirectUri)
	v.Add("response_type", responseType)
	v.Add("state", state)
	redirectUrl := "/oauth/login?" + v.Encode()
	http.Redirect(w, r, redirectUrl, iris.StatusTemporaryRedirect)
}
