package services

import (
	"strings"
	"time"

	"bbs-go/internal/cache"
	"bbs-go/internal/models"
	"bbs-go/internal/models/constants"
	"bbs-go/internal/pkg/errs"
	"bbs-go/internal/pkg/event"
	"bbs-go/internal/repositories"

	"bbs-go/internal/pkg/ginx"
	"bbs-go/internal/pkg/params"

	"github.com/gin-gonic/gin"
	"github.com/mlogclub/simple/common/dates"
	"github.com/mlogclub/simple/common/strs"
	"github.com/mlogclub/simple/sqls"
)

var UserTokenService = newUserTokenService()

func newUserTokenService() *userTokenService {
	return &userTokenService{}
}

type userTokenService struct {
}

func (s *userTokenService) GetCurrentUserId(ctx *gin.Context) int64 {
	user := s.GetCurrent(ctx)
	if user != nil {
		return user.Id
	}
	return 0
}

func (s *userTokenService) GetCurrent(ctx *gin.Context) *models.User {
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

	// 登录态访问：用于每日登录任务（带 token 打开网站即算今日登录，每用户每天仅发一次）
	trySendUserLoginEvent(ctx, user.Id)

	return user
}

// trySendUserLoginEvent 在登录态访问时发送 user.login 事件，供每日登录任务使用。每用户每天仅发一次（由布隆过滤器 TryMarkAndReturnIfNew 原子保证）。
func trySendUserLoginEvent(ctx *gin.Context, userId int64) {
	if ctx == nil || userId <= 0 {
		return
	}
	// 确保本次请求只调用一次
	ctxKeyDailyVisitSent := "daily_visit_sent"
	if _, exists := ctx.Get(ctxKeyDailyVisitSent); exists {
		return
	}
	ctx.Set(ctxKeyDailyVisitSent, true)

	// 如果今日已发送过，则不发送
	if !cache.DailyVisitCache.TryMarkAndReturnIfNew(userId) {
		return
	}
	event.Send(event.UserLoginEvent{
		UserId:     userId,
		LoginTime:  dates.NowTimestamp(),
		IsNewLogin: false,
	})
}

func (s *userTokenService) CheckLogin(ctx *gin.Context) (*models.User, error) {
	user := s.GetCurrent(ctx)
	if user == nil {
		return nil, errs.NotLogin()
	}
	return user, nil
}

func (s *userTokenService) Signout(ctx *gin.Context) error {
	token := s.GetUserToken(ctx)
	userToken := repositories.UserTokenRepository.GetByToken(sqls.DB(), token)
	if userToken == nil {
		return nil
	}
	err := repositories.UserTokenRepository.UpdateColumn(sqls.DB(), userToken.Id, "status", constants.StatusDeleted)
	if err != nil {
		return err
	}
	ginx.RemoveCookie(ctx, constants.CookieTokenKey)
	return nil
}

func (s *userTokenService) GetUserToken(ctx *gin.Context) string {
	if userToken, _ := params.Get(ctx, "userToken"); strs.IsNotBlank(userToken) {
		return userToken
	}
	if userToken := ginx.GetCookie(ctx, constants.CookieTokenKey); strs.IsNotBlank(userToken) {
		return userToken
	}
	return s.getUserTokenFromHeader(ctx)
}

func (s *userTokenService) getUserTokenFromHeader(ctx *gin.Context) string {
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
