package services

import (
	"bbs-go/model/constants"
	"github.com/mlogclub/simple/date"
	"time"

	"github.com/kataras/iris/v12"
	"github.com/mlogclub/simple"

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
	return repositories.UserTokenRepository.Get(simple.DB(), id)
}

func (s *userTokenService) Take(where ...interface{}) *model.UserToken {
	return repositories.UserTokenRepository.Take(simple.DB(), where...)
}

func (s *userTokenService) Find(cnd *simple.SqlCnd) []model.UserToken {
	return repositories.UserTokenRepository.Find(simple.DB(), cnd)
}

func (s *userTokenService) FindOne(cnd *simple.SqlCnd) *model.UserToken {
	return repositories.UserTokenRepository.FindOne(simple.DB(), cnd)
}

func (s *userTokenService) FindPageByParams(params *simple.QueryParams) (list []model.UserToken, paging *simple.Paging) {
	return repositories.UserTokenRepository.FindPageByParams(simple.DB(), params)
}

func (s *userTokenService) FindPageByCnd(cnd *simple.SqlCnd) (list []model.UserToken, paging *simple.Paging) {
	return repositories.UserTokenRepository.FindPageByCnd(simple.DB(), cnd)
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
	if userToken.ExpiredAt <= date.NowTimestamp() {
		return nil
	}
	user := cache.UserCache.Get(userToken.UserId)
	if user == nil || user.Status != constants.StatusOk {
		return nil
	}
	return user
}

// CheckLogin 检查登录状态
func (s *userTokenService) CheckLogin(ctx iris.Context) (*model.User, *simple.CodeError) {
	user := s.GetCurrent(ctx)
	if user == nil {
		return nil, simple.ErrorNotLogin
	}
	return user, nil
}

// 退出登录
func (s *userTokenService) Signout(ctx iris.Context) error {
	token := s.GetUserToken(ctx)
	userToken := repositories.UserTokenRepository.GetByToken(simple.DB(), token)
	if userToken == nil {
		return nil
	}
	return repositories.UserTokenRepository.UpdateColumn(simple.DB(), userToken.Id, "status", constants.StatusDeleted)
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
	token := simple.UUID()
	tokenExpireDays := SysConfigService.GetTokenExpireDays()
	expiredAt := time.Now().Add(time.Hour * 24 * time.Duration(tokenExpireDays))
	userToken := &model.UserToken{
		Token:      token,
		UserId:     userId,
		ExpiredAt:  date.Timestamp(expiredAt),
		Status:     constants.StatusOk,
		CreateTime: date.NowTimestamp(),
	}
	err := repositories.UserTokenRepository.Create(simple.DB(), userToken)
	if err != nil {
		return "", err
	}
	return token, nil
}

// 禁用
func (s *userTokenService) Disable(token string) error {
	t := repositories.UserTokenRepository.GetByToken(simple.DB(), token)
	if t == nil {
		return nil
	}
	err := repositories.UserTokenRepository.UpdateColumn(simple.DB(), t.Id, "status", constants.StatusDeleted)
	if err != nil {
		cache.UserTokenCache.Invalidate(token)
	}
	return err
}
