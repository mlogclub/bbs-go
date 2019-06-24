package oauth

import (
	"github.com/mlogclub/mlog/controllers/render"
	"github.com/mlogclub/mlog/model"
	"gopkg.in/oauth2.v3"
	"net/http"
	"strconv"

	"github.com/sirupsen/logrus"
)

func GetUserInfoByRequest(r *http.Request) *model.UserResponse {
	token, err := ServerInstance.Srv.ValidationBearerToken(r)
	if err != nil {
		logrus.Errorln(err)
		return nil
	}
	return GetUserInfoByToken(token)
}

func GetUserInfo(accessToken string) *model.UserResponse {
	if len(accessToken) == 0 {
		return nil
	}
	token, err := ServerInstance.Srv.Manager.LoadAccessToken(accessToken)
	if err != nil {
		logrus.Errorln(err)
		return nil
	}
	return GetUserInfoByToken(token)
}

func GetUserInfoByToken(token oauth2.TokenInfo) *model.UserResponse {
	userId, err := strconv.ParseInt(token.GetUserID(), 10, 64)
	if err != nil {
		logrus.Errorln(err)
		return nil
	}
	if userId <= 0 {
		return nil
	}
	return render.BuildUserById(userId)
}
