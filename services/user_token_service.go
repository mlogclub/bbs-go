package services

import (
	"time"

	"github.com/kataras/iris/context"
	"github.com/mlogclub/simple"

	"github.com/mlogclub/mlog/model"
	"github.com/mlogclub/mlog/repositories"
	"github.com/mlogclub/mlog/services/cache"
)

var UserTokenService = newUserTokenService()

func newUserTokenService() *userTokenService {
	return &userTokenService{}
}

type userTokenService struct {
}

// 获取当前登录用户
func (this *userTokenService) GetCurrent(ctx context.Context) *model.User {
	token := this.getUserToken(ctx)
	userToken := cache.UserTokenCache.Get(token)
	// 没找到授权
	if userToken == nil || userToken.Status == model.UserTokenStatusDisabled {
		return nil
	}
	// 授权过期
	if userToken.ExpiredAt <= simple.NowTimestamp() {
		return nil
	}
	return cache.UserCache.Get(userToken.UserId)
}

// 从请求体中获取UserToken
func (this *userTokenService) getUserToken(ctx context.Context) string {
	userToken := ctx.FormValue("userToken")
	if len(userToken) > 0 {
		return userToken
	}
	return ctx.GetHeader("X-User-Token")
}

// 生成
func (this *userTokenService) Generate(userId int64) (string, error) {
	token := simple.Uuid()
	expiredAt := time.Now().Add(time.Hour * 24 * 7) // 7天后过期
	userToken := &model.UserToken{
		Token:      token,
		UserId:     userId,
		ExpiredAt:  simple.Timestamp(expiredAt),
		Status:     model.UserTokenStatusOk,
		CreateTime: simple.NowTimestamp(),
	}
	err := repositories.UserTokenRepository.Create(simple.GetDB(), userToken)
	if err != nil {
		return "", err
	}
	return token, nil
}

// 禁用
func (this *userTokenService) Disable(token string) error {
	t := repositories.UserTokenRepository.GetByToken(simple.GetDB(), token)
	if t == nil {
		return nil
	}
	err := repositories.UserTokenRepository.UpdateColumn(simple.GetDB(), t.Id, "status", model.UserTokenStatusDisabled)
	if err != nil {
		cache.UserTokenCache.Invalidate(token)
	}
	return err
}

func (this *userTokenService) Get(id int64) *model.UserToken {
	return repositories.UserTokenRepository.Get(simple.GetDB(), id)
}

func (this *userTokenService) Take(where ...interface{}) *model.UserToken {
	return repositories.UserTokenRepository.Take(simple.GetDB(), where...)
}

func (this *userTokenService) QueryCnd(cnd *simple.QueryCnd) (list []model.UserToken, err error) {
	return repositories.UserTokenRepository.QueryCnd(simple.GetDB(), cnd)
}

func (this *userTokenService) Query(queries *simple.ParamQueries) (list []model.UserToken, paging *simple.Paging) {
	return repositories.UserTokenRepository.Query(simple.GetDB(), queries)
}
