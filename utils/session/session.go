package session

import (
	context2 "context"
	"github.com/go-session/redis"
	"github.com/go-session/session"
	"github.com/kataras/iris/context"
	"github.com/mlogclub/mlog/utils/config"
	"github.com/mlogclub/simple"
	"github.com/sirupsen/logrus"
	"net/http"

	"github.com/mlogclub/mlog/model"
)

const (
	SessionCurrentUser = "CurrentUser"
)

func InitSessionManager() {
	session.InitManager(
		session.SetStore(redis.NewRedisStore(&redis.Options{
			Addr: config.Conf.RedisAddr,
		})),
		session.SetCookieName("mlog_session_id"),
		session.SetExpired(86400),
		session.SetEnableSIDInURLQuery(false),
		session.SetEnableSIDInHTTPHeader(false),
	)
}

func Start(ctx context.Context) session.Store {
	return StartByRequest(ctx.ResponseWriter(), ctx.Request())
}

func StartByRequest(w http.ResponseWriter, r *http.Request) session.Store {
	store, err := session.Start(context2.Background(), w, r)
	if err != nil {
		logrus.Error(err)
	}
	return store
}

func SetCurrentUser(ctx context.Context, user *model.User) {
	store := Start(ctx)
	store.Set(SessionCurrentUser, user)
	err := store.Save()
	if err != nil {
		logrus.Error(err)
	}
}

func GetCurrentUser(ctx context.Context) *model.User {
	return GetCurrentUserByRequest(ctx.ResponseWriter(), ctx.Request())
}

func GetCurrentUserByRequest(w http.ResponseWriter, r *http.Request) *model.User {
	val, exists := StartByRequest(w, r).Get(SessionCurrentUser)
	if exists {
		json, err := simple.FormatJson(val)
		if err != nil {
			logrus.Error(err)
			return nil
		}
		user := &model.User{}
		err = simple.ParseJson(json, user)
		if err != nil {
			logrus.Error(err)
		}
		return user
	}
	return nil
}

func DelCurrentUser(ctx context.Context) {
	store := Start(ctx)
	store.Delete(SessionCurrentUser)
	err := store.Save()
	if err != nil {
		logrus.Error(err)
	}
}
