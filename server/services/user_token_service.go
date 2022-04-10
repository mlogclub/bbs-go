package services

import (
	"bbs-go/model/constants"
	"bbs-go/pkg/common"
	"time"

	"github.com/kataras/iris/v12"
	"github.com/mlogclub/simple/common/dates"
	"github.com/mlogclub/simple/common/strs"
	"github.com/mlogclub/simple/sqls"
	"github.com/mlogclub/simple/web"
	"github.com/mlogclub/simple/web/params"

	"bbs-go/cache"
	"bbs-go/model"
	"bbs-go/repositories"
)

var UserTokenService = newUserTokenService()

func newUserTokenService() *userTokenService {
	return &userTokenService{}
}

type userTokenService struct {
}

func (s *userTokenService) Get(id int64) *model.UserToken {
	return repositories.UserTokenRepository.Get(sqls.DB(), id)
}

func (s *userTokenService) Take(where ...interface{}) *model.UserToken {
	return repositories.UserTokenRepository.Take(sqls.DB(), where...)
}

func (s *userTokenService) Find(cnd *sqls.Cnd) []model.UserToken {
	return repositories.UserTokenRepository.Find(sqls.DB(), cnd)
}

func (s *userTokenService) FindOne(cnd *sqls.Cnd) *model.UserToken {
	return repositories.UserTokenRepository.FindOne(sqls.DB(), cnd)
}

func (s *userTokenService) FindPageByParams(params *params.QueryParams) (list []model.UserToken, paging *sqls.Paging) {
	return repositories.UserTokenRepository.FindPageByParams(sqls.DB(), params)
}

func (s *userTokenService) FindPageByCnd(cnd *sqls.Cnd) (list []model.UserToken, paging *sqls.Paging) {
	return repositories.UserTokenRepository.FindPageByCnd(sqls.DB(), cnd)
}

// 获取当前登录用户的id
func (s *userTokenService) GetCurrentUserId(ctx iris.Context) int64 {
	user := s.GetCurrent(ctx)
	if user != nil {
		return user.Id
	}
	return 0
}

// 获取当前登录用户
func (s *userTokenService) GetCurrent(ctx iris.Context) *model.User {
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

// CheckLogin 检查登录状态
func (s *userTokenService) CheckLogin(ctx iris.Context) (*model.User, *web.CodeError) {
	user := s.GetCurrent(ctx)
	if user == nil {
		return nil, common.ErrorNotLogin
	}
	return user, nil
}

// 退出登录
func (s *userTokenService) Signout(ctx iris.Context) error {
	token := s.GetUserToken(ctx)
	userToken := repositories.UserTokenRepository.GetByToken(sqls.DB(), token)
	if userToken == nil {
		return nil
	}
	return repositories.UserTokenRepository.UpdateColumn(sqls.DB(), userToken.Id, "status", constants.StatusDeleted)
}

// 从请求体中获取UserToken
func (s *userTokenService) GetUserToken(ctx iris.Context) string {
	userToken := ctx.FormValue("userToken")
	if len(userToken) > 0 {
		return userToken
	}
	return ctx.GetHeader("X-User-Token")
}

// 生成
func (s *userTokenService) Generate(userId int64) (string, error) {
	token := strs.UUID()
	tokenExpireDays := SysConfigService.GetTokenExpireDays()
	expiredAt := time.Now().Add(time.Hour * 24 * time.Duration(tokenExpireDays))
	userToken := &model.UserToken{
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

// 禁用
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
