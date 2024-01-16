package services

import (
	"bbs-go/internal/models/constants"
	"bbs-go/internal/pkg/errs"
	"strings"
	"time"

	"bbs-go/internal/cache"
	"bbs-go/internal/models"
	"bbs-go/internal/repositories"

	"github.com/kataras/iris/v12"
	"github.com/mlogclub/simple/common/dates"
	"github.com/mlogclub/simple/common/strs"
	"github.com/mlogclub/simple/sqls"
	"github.com/mlogclub/simple/web"
)

var UserTokenService = newUserTokenService()

func newUserTokenService() *userTokenService {
	return &userTokenService{}
}

type userTokenService struct {
}

func (s *userTokenService) GetCurrentUserId(ctx iris.Context) int64 {
	user := s.GetCurrent(ctx)
	if user != nil {
		return user.Id
	}
	return 0
}

func (s *userTokenService) GetCurrent(ctx iris.Context) *models.User {
	token := s.GetUserToken(ctx)
	userToken := cache.UserTokenCache.Get(token)
	// 没找到授权
	if userToken == nil || userToken.Status == constants.StatusDeleted {
		return nil
	}
	// 授权过期
	if userToken.ExpiredAt <= dates.NowTimestamp() {
		return nil
	}
	user := cache.UserCache.Get(userToken.UserId)
	if user == nil || user.Status != constants.StatusOk {
		return nil
	}
	return user
}

func (s *userTokenService) CheckLogin(ctx iris.Context) (*models.User, *web.CodeError) {
	user := s.GetCurrent(ctx)
	if user == nil {
		return nil, errs.NotLogin
	}
	return user, nil
}

func (s *userTokenService) Signout(ctx iris.Context) error {
	token := s.GetUserToken(ctx)
	userToken := repositories.UserTokenRepository.GetByToken(sqls.DB(), token)
	if userToken == nil {
		return nil
	}
	err := repositories.UserTokenRepository.UpdateColumn(sqls.DB(), userToken.Id, "status", constants.StatusDeleted)
	if err != nil {
		return err
	}
	ctx.RemoveCookie(constants.CookieTokenKey)
	return nil
}

func (s *userTokenService) GetUserToken(ctx iris.Context) string {
	if userToken := ctx.GetCookie(constants.CookieTokenKey); strs.IsNotBlank(userToken) {
		return userToken
	}
	if userToken := s.getUserTokenFromHeader(ctx); strs.IsNotBlank(userToken) {
		return userToken
	}
	return ctx.FormValue("userToken")
}

func (s *userTokenService) getUserTokenFromHeader(ctx iris.Context) string {
	if authorization := ctx.GetHeader("Authorization"); strs.IsNotBlank(authorization) {
		userToken, _ := strings.CutPrefix(authorization, "Bearer ")
		return userToken
	}
	return ctx.GetHeader("X-User-Token")
}

func (s *userTokenService) Generate(userId int64) (string, error) {
	token := strs.UUID()
	tokenExpireDays := SysConfigService.GetTokenExpireDays()
	expiredAt := time.Now().Add(time.Hour * 24 * time.Duration(tokenExpireDays))
	userToken := &models.UserToken{
		Token:      token,
		UserId:     userId,
		ExpiredAt:  dates.Timestamp(expiredAt),
		Status:     constants.StatusOk,
		CreateTime: dates.NowTimestamp(),
	}
	err := repositories.UserTokenRepository.Create(sqls.DB(), userToken)
	if err != nil {
		return "", err
	}
	return token, nil
}

func (s *userTokenService) Disable(token string) error {
	t := repositories.UserTokenRepository.GetByToken(sqls.DB(), token)
	if t == nil {
		return nil
	}
	err := repositories.UserTokenRepository.UpdateColumn(sqls.DB(), t.Id, "status", constants.StatusDeleted)
	if err != nil {
		cache.UserTokenCache.Invalidate(token)
	}
	return err
}
