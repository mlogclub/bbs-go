package services

import (
	"bbs-go/internal/models/constants"
	"bbs-go/internal/models/dto"
	"errors"
	"log/slog"

	"github.com/mlogclub/simple/common/dates"
	"github.com/mlogclub/simple/common/jsons"
	"github.com/mlogclub/simple/common/strs"
	"github.com/mlogclub/simple/sqls"
	"github.com/mlogclub/simple/web/params"
	"github.com/tidwall/gjson"

	"gorm.io/gorm"

	"bbs-go/internal/cache"
	"bbs-go/internal/models"
	"bbs-go/internal/repositories"
)

var SysConfigService = newSysConfigService()

func newSysConfigService() *sysConfigService {
	return &sysConfigService{}
}

type sysConfigService struct {
}

func (s *sysConfigService) Get(id int64) *models.SysConfig {
	return repositories.SysConfigRepository.Get(sqls.DB(), id)
}

func (s *sysConfigService) Take(where ...interface{}) *models.SysConfig {
	return repositories.SysConfigRepository.Take(sqls.DB(), where...)
}

func (s *sysConfigService) Find(cnd *sqls.Cnd) []models.SysConfig {
	return repositories.SysConfigRepository.Find(sqls.DB(), cnd)
}

func (s *sysConfigService) FindOne(cnd *sqls.Cnd) *models.SysConfig {
	return repositories.SysConfigRepository.FindOne(sqls.DB(), cnd)
}

func (s *sysConfigService) FindPageByParams(params *params.QueryParams) (list []models.SysConfig, paging *sqls.Paging) {
	return repositories.SysConfigRepository.FindPageByParams(sqls.DB(), params)
}

func (s *sysConfigService) FindPageByCnd(cnd *sqls.Cnd) (list []models.SysConfig, paging *sqls.Paging) {
	return repositories.SysConfigRepository.FindPageByCnd(sqls.DB(), cnd)
}

func (s *sysConfigService) GetAll() []models.SysConfig {
	return repositories.SysConfigRepository.Find(sqls.DB(), sqls.NewCnd().Asc("id"))
}

func (s *sysConfigService) SetAll(configStr string) error {
	json := gjson.Parse(configStr)
	configs, ok := json.Value().(map[string]interface{})
	if !ok {
		return errors.New("配置数据格式错误")
	}
	return sqls.DB().Transaction(func(tx *gorm.DB) error {
		for k := range configs {
			v := json.Get(k).String()
			if err := s.setSingle(tx, k, v, "", ""); err != nil {
				return err
			}
		}
		return nil
	})
}

// Set 设置配置，如果配置不存在，那么创建
func (s *sysConfigService) Set(key, value, name, description string) error {
	return sqls.DB().Transaction(func(tx *gorm.DB) error {
		if err := s.setSingle(tx, key, value, name, description); err != nil {
			return err
		}
		return nil
	})
}

func (s *sysConfigService) setSingle(db *gorm.DB, key, value, name, description string) error {
	if len(key) == 0 {
		return errors.New("sys config key is null")
	}
	sysConfig := repositories.SysConfigRepository.GetByKey(db, key)
	if sysConfig == nil {
		sysConfig = &models.SysConfig{
			CreateTime: dates.NowTimestamp(),
		}
	}
	sysConfig.Key = key
	sysConfig.Value = value
	sysConfig.UpdateTime = dates.NowTimestamp()

	if strs.IsNotBlank(name) {
		sysConfig.Name = name
	}
	if strs.IsNotBlank(description) {
		sysConfig.Description = description
	}

	var err error
	if sysConfig.Id > 0 {
		err = repositories.SysConfigRepository.Update(db, sysConfig)
	} else {
		err = repositories.SysConfigRepository.Create(db, sysConfig)
	}
	if err != nil {
		return err
	} else {
		cache.SysConfigCache.Invalidate(key)
		return nil
	}
}

func (s *sysConfigService) GetConfig() *dto.SysConfigResponse {
	return &dto.SysConfigResponse{
		SiteTitle:                  cache.SysConfigCache.GetStr(constants.SysConfigSiteTitle),
		SiteDescription:            cache.SysConfigCache.GetStr(constants.SysConfigSiteDescription),
		SiteKeywords:               cache.SysConfigCache.GetStrArr(constants.SysConfigSiteKeywords),
		SiteLogo:                   cache.SysConfigCache.GetStr(constants.SysConfigSiteLogo),
		SiteNavs:                   s.GetSiteNavs(),
		SiteNotification:           cache.SysConfigCache.GetStr(constants.SysConfigSiteNotification),
		RecommendTags:              cache.SysConfigCache.GetStrArr(constants.SysConfigRecommendTags),
		UrlRedirect:                cache.SysConfigCache.GetBool(constants.SysConfigUrlRedirect),
		ScoreConfig:                s.GetScoreConfig(),
		DefaultNodeId:              cache.SysConfigCache.GetInt64(constants.SysConfigDefaultNodeId),
		ArticlePending:             s.IsArticlePending(),
		TopicCaptcha:               cache.SysConfigCache.GetBool(constants.SysConfigTopicCaptcha),
		UserObserveSeconds:         cache.SysConfigCache.GetInt(constants.SysConfigUserObserveSeconds),
		TokenExpireDays:            s.GetTokenExpireDays(),
		LoginMethod:                s.GetLoginMethod(),
		CreateTopicEmailVerified:   s.IsCreateTopicEmailVerified(),
		CreateArticleEmailVerified: s.IsCreateArticleEmailVerified(),
		CreateCommentEmailVerified: s.IsCreateCommentEmailVerified(),
		EnableHideContent:          s.IsEnableHideContent(),
		Modules:                    s.GetModules(),
		EmailWhitelist:             s.GetEmailWhitelist(),
	}
}

func (s *sysConfigService) GetTokenExpireDays() int {
	tokenExpireDays := cache.SysConfigCache.GetInt(constants.SysConfigTokenExpireDays)
	if tokenExpireDays <= 0 {
		tokenExpireDays = constants.DefaultTokenExpireDays
	}
	return tokenExpireDays
}

func (s *sysConfigService) GetLoginMethod() dto.LoginMethod {
	loginMethodStr := cache.SysConfigCache.GetStr(constants.SysConfigLoginMethod)

	useDefault := true
	var loginMethod dto.LoginMethod
	if strs.IsNotBlank(loginMethodStr) {
		if err := jsons.Parse(loginMethodStr, &loginMethod); err != nil {
			slog.Warn("登录方式数据错误", slog.Any("err", err))
		} else {
			useDefault = false
		}
	}
	if useDefault {
		loginMethod = dto.LoginMethod{
			Password: true,
			QQ:       true,
			Github:   true,
		}
	}
	return loginMethod
}

func (s *sysConfigService) IsCreateTopicEmailVerified() bool {
	return cache.SysConfigCache.GetBool(constants.SysConfigCreateTopicEmailVerified)
}

func (s *sysConfigService) IsCreateArticleEmailVerified() bool {
	return cache.SysConfigCache.GetBool(constants.SysConfigCreateArticleEmailVerified)
}

func (s *sysConfigService) IsCreateCommentEmailVerified() bool {
	return cache.SysConfigCache.GetBool(constants.SysConfigCreateCommentEmailVerified)
}

func (s *sysConfigService) IsEnableHideContent() bool {
	return cache.SysConfigCache.GetBool(constants.SysConfigEnableHideContent)
}

func (s *sysConfigService) IsArticlePending() bool {
	return cache.SysConfigCache.GetBool(constants.SysConfigArticlePending)
}

func (s *sysConfigService) GetSiteNavs() []dto.ActionLink {
	siteNavs := cache.SysConfigCache.GetStr(constants.SysConfigSiteNavs)
	var siteNavsArr []dto.ActionLink
	if strs.IsNotBlank(siteNavs) {
		if err := jsons.Parse(siteNavs, &siteNavsArr); err != nil {
			slog.Warn("站点导航数据错误", slog.Any("err", err))
		}
	}
	return siteNavsArr
}

func (s *sysConfigService) GetModules() dto.ModulesConfig {
	str := cache.SysConfigCache.GetStr(constants.SysConfigModules)

	useDefault := true
	var modulesConfig dto.ModulesConfig
	if strs.IsNotBlank(str) {
		if err := jsons.Parse(str, &modulesConfig); err != nil {
			slog.Warn("启用模块配置错误", slog.Any("err", err))
		} else {
			useDefault = false
		}
	}
	if useDefault {
		modulesConfig = dto.ModulesConfig{
			Tweet:   true,
			Topic:   true,
			Article: true,
		}
	}
	return modulesConfig
}

// GetEmailWhitelist 邮箱白名单
func (s *sysConfigService) GetEmailWhitelist() []string {
	str := cache.SysConfigCache.GetStr(constants.SysConfigEmailWhitelist)
	var emailWhitelist []string
	if strs.IsNotBlank(str) {
		_ = jsons.Parse(str, &emailWhitelist)
	}
	return emailWhitelist
}

func (s *sysConfigService) GetScoreConfig() dto.ScoreConfig {
	str := cache.SysConfigCache.GetStr(constants.SysConfigScoreConfig)
	var scoreConfig dto.ScoreConfig
	if err := jsons.Parse(str, &scoreConfig); err != nil {
		slog.Warn("积分配置错误", slog.Any("err", err))
	}
	return scoreConfig
}
