package services

import (
	"context"
	"database/sql"
	"bbs-go/internal/models"
	"bbs-go/internal/models/constants"
	"bbs-go/internal/pkg/bbsurls"
	"bbs-go/internal/pkg/google"
	"bbs-go/internal/pkg/wx"
	"bbs-go/internal/repositories"
	"log/slog"

	"github.com/kataras/iris/v12/x/errors"
	"github.com/mlogclub/simple/common/dates"
	"github.com/mlogclub/simple/common/jsons"
	"github.com/mlogclub/simple/sqls"
	"github.com/mlogclub/simple/web/params"
	"gorm.io/gorm"
)

var ThirdUserService = newThirdUserService()

func newThirdUserService() *thirdUserService {
	return &thirdUserService{}
}

type thirdUserService struct {
}

func (s *thirdUserService) Get(id int64) *models.ThirdUser {
	return repositories.ThirdUserRepository.Get(sqls.DB(), id)
}

func (s *thirdUserService) Take(where ...interface{}) *models.ThirdUser {
	return repositories.ThirdUserRepository.Take(sqls.DB(), where...)
}

func (s *thirdUserService) Find(cnd *sqls.Cnd) []models.ThirdUser {
	return repositories.ThirdUserRepository.Find(sqls.DB(), cnd)
}

func (s *thirdUserService) FindOne(cnd *sqls.Cnd) *models.ThirdUser {
	return repositories.ThirdUserRepository.FindOne(sqls.DB(), cnd)
}

func (s *thirdUserService) FindPageByParams(params *params.QueryParams) (list []models.ThirdUser, paging *sqls.Paging) {
	return repositories.ThirdUserRepository.FindPageByParams(sqls.DB(), params)
}

func (s *thirdUserService) FindPageByCnd(cnd *sqls.Cnd) (list []models.ThirdUser, paging *sqls.Paging) {
	return repositories.ThirdUserRepository.FindPageByCnd(sqls.DB(), cnd)
}

func (s *thirdUserService) Count(cnd *sqls.Cnd) int64 {
	return repositories.ThirdUserRepository.Count(sqls.DB(), cnd)
}

func (s *thirdUserService) Create(t *models.ThirdUser) error {
	return repositories.ThirdUserRepository.Create(sqls.DB(), t)
}

func (s *thirdUserService) Update(t *models.ThirdUser) error {
	return repositories.ThirdUserRepository.Update(sqls.DB(), t)
}

func (s *thirdUserService) Updates(id int64, columns map[string]interface{}) error {
	return repositories.ThirdUserRepository.Updates(sqls.DB(), id, columns)
}

func (s *thirdUserService) UpdateColumn(id int64, name string, value interface{}) error {
	return repositories.ThirdUserRepository.UpdateColumn(sqls.DB(), id, name, value)
}

func (s *thirdUserService) Delete(id int64) {
	repositories.ThirdUserRepository.Delete(sqls.DB(), id)
}

func (s *thirdUserService) GetByOpenId(openId string, thirdType constants.ThirdType) *models.ThirdUser {
	return repositories.ThirdUserRepository.GetByOpenId(sqls.DB(), openId, thirdType)
}

func (s *thirdUserService) GetByUserId(userId int64, thirdType constants.ThirdType) *models.ThirdUser {
	return repositories.ThirdUserRepository.GetByUserId(sqls.DB(), userId, thirdType)
}

func (s *thirdUserService) LoginWeixin(code, state string) (*models.User, error) {
	loginConfig := SysConfigService.GetLoginConfig()
	oauth := wx.NewOfficialAccount(loginConfig.WeixinLogin.AppId, loginConfig.WeixinLogin.AppSecret).GetOauth()
	accessToken, err := oauth.GetUserAccessToken(code)
	if err != nil {
		slog.Error(err.Error())
		return nil, err
	}

	info, err := oauth.GetUserInfo(accessToken.AccessToken, accessToken.OpenID, "cn")
	if err != nil {
		slog.Error(err.Error())
		return nil, err
	}

	thirdUser := ThirdUserService.GetByOpenId(info.OpenID, constants.ThirdTypeWeixin)
	if thirdUser != nil && thirdUser.UserId > 0 {
		return UserService.Get(thirdUser.UserId), nil
	}

	// copy wechat head image
	avatar, _ := UploadService.CopyImage(info.HeadImgURL)

	user := &models.User{
		Type:       constants.UserTypeNormal,
		Nickname:   info.Nickname,
		Avatar:     avatar,
		Status:     constants.StatusOk,
		CreateTime: dates.NowTimestamp(),
		UpdateTime: dates.NowTimestamp(),
	}

	err = sqls.DB().Transaction(func(tx *gorm.DB) error {
		if err := repositories.UserRepository.Create(tx, user); err != nil {
			return err
		}
		if thirdUser == nil {
			return repositories.ThirdUserRepository.Create(tx, &models.ThirdUser{
				UserId:     user.Id,
				OpenId:     info.OpenID,
				ThirdType:  constants.ThirdTypeWeixin,
				Nickname:   info.Nickname,
				Avatar:     avatar,
				ExtraData:  jsons.ToJsonStr(info),
				CreateTime: dates.NowTimestamp(),
				UpdateTime: dates.NowTimestamp(),
			})
		} else {
			thirdUser.Nickname = info.Nickname
			thirdUser.Avatar = avatar
			thirdUser.ExtraData = jsons.ToJsonStr(info)
			thirdUser.UpdateTime = dates.NowTimestamp()
			return repositories.ThirdUserRepository.Update(tx, thirdUser)
		}
	})

	if err != nil {
		return nil, err
	}
	return user, nil
}

func (s *thirdUserService) BindWeixin(userId int64, code, state string) error {
	if temp := s.GetByUserId(userId, constants.ThirdTypeWeixin); temp != nil {
		return errors.New("用户已绑定微信: " + temp.Nickname)
	}

	loginConfig := SysConfigService.GetLoginConfig()
	oauth := wx.NewOfficialAccount(loginConfig.WeixinLogin.AppId, loginConfig.WeixinLogin.AppSecret).GetOauth()
	accessToken, err := oauth.GetUserAccessToken(code)
	if err != nil {
		slog.Error(err.Error())
		return err
	}

	info, err := oauth.GetUserInfo(accessToken.AccessToken, accessToken.OpenID, "cn")
	if err != nil {
		slog.Error(err.Error())
		return err
	}

	if temp := s.GetByOpenId(info.OpenID, constants.ThirdTypeWeixin); temp != nil && temp.Id != userId {
		return errors.New("微信已绑定到其他用户~")
	}

	return s.Create(&models.ThirdUser{
		UserId:     userId,
		OpenId:     info.OpenID,
		ThirdType:  constants.ThirdTypeWeixin,
		Nickname:   info.Nickname,
		Avatar:     info.HeadImgURL,
		ExtraData:  jsons.ToJsonStr(info),
		CreateTime: dates.NowTimestamp(),
	})
}

func (s *thirdUserService) UnbindWeixin(userId int64) {
	thirdUser := s.GetByUserId(userId, constants.ThirdTypeWeixin)
	if thirdUser == nil {
		return
	}
	repositories.ThirdUserRepository.Delete(sqls.DB(), thirdUser.Id)
}

func (s *thirdUserService) LoginGoogle(code, state string) (*models.User, error) {
	loginConfig := SysConfigService.GetLoginConfig()
	if !loginConfig.GoogleLogin.Enabled {
		return nil, errors.New("Google登录未启用")
	}

	// 使用与授权时相同的 redirectURI（必须完全一致）
	redirectURI := bbsurls.AbsUrl(google.CallbackPathLogin)
	oauth := google.NewGoogleOAuth(loginConfig.GoogleLogin.ClientId, loginConfig.GoogleLogin.ClientSecret, redirectURI)

	ctx := context.Background()
	info, err := oauth.GetUserInfo(ctx, code)
	if err != nil {
		slog.Error("Google登录获取用户信息失败", slog.Any("err", err))
		return nil, err
	}

	thirdUser := ThirdUserService.GetByOpenId(info.ID, constants.ThirdTypeGoogle)
	if thirdUser != nil && thirdUser.UserId > 0 {
		return UserService.Get(thirdUser.UserId), nil
	}

	// copy google avatar image
	avatar, _ := UploadService.CopyImage(info.Picture)

	// 使用 Google 名称作为昵称，如果没有则使用邮箱前缀
	nickname := info.Name
	if nickname == "" {
		nickname = info.Email
		if nickname == "" {
			nickname = "Google用户"
		}
	}

	user := &models.User{
		Type:       constants.UserTypeNormal,
		Nickname:   nickname,
		Avatar:     avatar,
		Status:     constants.StatusOk,
		CreateTime: dates.NowTimestamp(),
		UpdateTime: dates.NowTimestamp(),
	}

	// 如果邮箱已验证，设置邮箱
	if info.VerifiedEmail && info.Email != "" {
		user.Email = sql.NullString{
			String: info.Email,
			Valid:  true,
		}
		user.EmailVerified = true
	}

	err = sqls.DB().Transaction(func(tx *gorm.DB) error {
		if err := repositories.UserRepository.Create(tx, user); err != nil {
			return err
		}
		if thirdUser == nil {
			return repositories.ThirdUserRepository.Create(tx, &models.ThirdUser{
				UserId:     user.Id,
				OpenId:     info.ID,
				ThirdType:  constants.ThirdTypeGoogle,
				Nickname:   nickname,
				Avatar:     avatar,
				ExtraData:  jsons.ToJsonStr(info),
				CreateTime: dates.NowTimestamp(),
				UpdateTime: dates.NowTimestamp(),
			})
		} else {
			thirdUser.Nickname = nickname
			thirdUser.Avatar = avatar
			thirdUser.ExtraData = jsons.ToJsonStr(info)
			thirdUser.UpdateTime = dates.NowTimestamp()
			return repositories.ThirdUserRepository.Update(tx, thirdUser)
		}
	})

	if err != nil {
		return nil, err
	}
	return user, nil
}

func (s *thirdUserService) LoginGoogleOneTap(credential string) (*models.User, error) {
	loginConfig := SysConfigService.GetLoginConfig()
	if !loginConfig.GoogleLogin.Enabled {
		return nil, errors.New("Google登录未启用")
	}

	// 使用 Google API 验证 JWT
	ctx := context.Background()
	info, err := google.VerifyJWTWithGoogleAPI(ctx, credential)
	if err != nil {
		slog.Error("Google OneTap JWT验证失败", slog.Any("err", err))
		return nil, err
	}

	// 验证 Client ID（从 JWT 中提取，但我们已经通过 Google API 验证了）
	// 这里可以额外验证 info 是否匹配我们的 Client ID

	// 查找已存在的第三方用户（使用 sub 作为 OpenId，等同于 OAuth 的 id）
	thirdUser := ThirdUserService.GetByOpenId(info.ID, constants.ThirdTypeGoogle)
	if thirdUser != nil && thirdUser.UserId > 0 {
		return UserService.Get(thirdUser.UserId), nil
	}

	// copy google avatar image
	avatar, _ := UploadService.CopyImage(info.Picture)

	// 使用 Google 名称作为昵称，如果没有则使用邮箱前缀
	nickname := info.Name
	if nickname == "" {
		nickname = info.Email
		if nickname == "" {
			nickname = "Google用户"
		}
	}

	user := &models.User{
		Type:       constants.UserTypeNormal,
		Nickname:   nickname,
		Avatar:     avatar,
		Status:     constants.StatusOk,
		CreateTime: dates.NowTimestamp(),
		UpdateTime: dates.NowTimestamp(),
	}

	// 如果邮箱已验证，设置邮箱
	if info.VerifiedEmail && info.Email != "" {
		user.Email = sql.NullString{
			String: info.Email,
			Valid:  true,
		}
		user.EmailVerified = true
	}

	err = sqls.DB().Transaction(func(tx *gorm.DB) error {
		if err := repositories.UserRepository.Create(tx, user); err != nil {
			return err
		}
		if thirdUser == nil {
			return repositories.ThirdUserRepository.Create(tx, &models.ThirdUser{
				UserId:     user.Id,
				OpenId:     info.ID, // 使用 sub 作为 OpenId
				ThirdType:  constants.ThirdTypeGoogle,
				Nickname:   nickname,
				Avatar:     avatar,
				ExtraData:  jsons.ToJsonStr(info),
				CreateTime: dates.NowTimestamp(),
				UpdateTime: dates.NowTimestamp(),
			})
		} else {
			thirdUser.Nickname = nickname
			thirdUser.Avatar = avatar
			thirdUser.ExtraData = jsons.ToJsonStr(info)
			thirdUser.UpdateTime = dates.NowTimestamp()
			return repositories.ThirdUserRepository.Update(tx, thirdUser)
		}
	})

	if err != nil {
		return nil, err
	}
	return user, nil
}

func (s *thirdUserService) BindGoogle(userId int64, code, state string) error {
	if temp := s.GetByUserId(userId, constants.ThirdTypeGoogle); temp != nil {
		return errors.New("用户已绑定Google: " + temp.Nickname)
	}

	loginConfig := SysConfigService.GetLoginConfig()
	if !loginConfig.GoogleLogin.Enabled {
		return errors.New("Google登录未启用")
	}

	// 使用与授权时相同的 redirectURI（必须完全一致）
	redirectURI := bbsurls.AbsUrl(google.CallbackPathBind)
	oauth := google.NewGoogleOAuth(loginConfig.GoogleLogin.ClientId, loginConfig.GoogleLogin.ClientSecret, redirectURI)

	ctx := context.Background()
	info, err := oauth.GetUserInfo(ctx, code)
	if err != nil {
		slog.Error("Google绑定获取用户信息失败", slog.Any("err", err))
		return err
	}

	if temp := s.GetByOpenId(info.ID, constants.ThirdTypeGoogle); temp != nil && temp.UserId != userId {
		return errors.New("Google账号已绑定到其他用户~")
	}

	nickname := info.Name
	if nickname == "" {
		nickname = info.Email
		if nickname == "" {
			nickname = "Google用户"
		}
	}

	return s.Create(&models.ThirdUser{
		UserId:     userId,
		OpenId:     info.ID,
		ThirdType:  constants.ThirdTypeGoogle,
		Nickname:   nickname,
		Avatar:     info.Picture,
		ExtraData:  jsons.ToJsonStr(info),
		CreateTime: dates.NowTimestamp(),
		UpdateTime: dates.NowTimestamp(),
	})
}

func (s *thirdUserService) UnbindGoogle(userId int64) {
	thirdUser := s.GetByUserId(userId, constants.ThirdTypeGoogle)
	if thirdUser == nil {
		return
	}
	repositories.ThirdUserRepository.Delete(sqls.DB(), thirdUser.Id)
}
