package services

import (
	"bbs-go/model/constants"
	"errors"
	"github.com/mlogclub/simple/date"
	"github.com/mlogclub/simple/json"
	"github.com/mlogclub/simple/number"
	"strconv"
	"strings"

	"github.com/mlogclub/simple"
	"github.com/sirupsen/logrus"
	"github.com/tidwall/gjson"
	"gorm.io/gorm"

	"bbs-go/cache"
	"bbs-go/model"
	"bbs-go/repositories"
)

var SysConfigService = newSysConfigService()

func newSysConfigService() *sysConfigService {
	return &sysConfigService{}
}

type sysConfigService struct {
}

func (s *sysConfigService) Get(id int64) *model.SysConfig {
	return repositories.SysConfigRepository.Get(simple.DB(), id)
}

func (s *sysConfigService) Take(where ...interface{}) *model.SysConfig {
	return repositories.SysConfigRepository.Take(simple.DB(), where...)
}

func (s *sysConfigService) Find(cnd *simple.SqlCnd) []model.SysConfig {
	return repositories.SysConfigRepository.Find(simple.DB(), cnd)
}

func (s *sysConfigService) FindOne(cnd *simple.SqlCnd) *model.SysConfig {
	return repositories.SysConfigRepository.FindOne(simple.DB(), cnd)
}

func (s *sysConfigService) FindPageByParams(params *simple.QueryParams) (list []model.SysConfig, paging *simple.Paging) {
	return repositories.SysConfigRepository.FindPageByParams(simple.DB(), params)
}

func (s *sysConfigService) FindPageByCnd(cnd *simple.SqlCnd) (list []model.SysConfig, paging *simple.Paging) {
	return repositories.SysConfigRepository.FindPageByCnd(simple.DB(), cnd)
}

func (s *sysConfigService) GetAll() []model.SysConfig {
	return repositories.SysConfigRepository.Find(simple.DB(), simple.NewSqlCnd().Asc("id"))
}

func (s *sysConfigService) SetAll(configStr string) error {
	json := gjson.Parse(configStr)
	configs, ok := json.Value().(map[string]interface{})
	if !ok {
		return errors.New("配置数据格式错误")
	}
	return simple.DB().Transaction(func(tx *gorm.DB) error {
		for k, _ := range configs {
			v := json.Get(k).String()
			if err := s.setSingle(tx, k, v, "", ""); err != nil {
				return err
			}
		}
		return nil
	})
}

// 设置配置，如果配置不存在，那么创建
func (s *sysConfigService) Set(key, value, name, description string) error {
	return simple.DB().Transaction(func(tx *gorm.DB) error {
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
		sysConfig = &model.SysConfig{
			CreateTime: date.NowTimestamp(),
		}
	}
	sysConfig.Key = key
	sysConfig.Value = value
	sysConfig.UpdateTime = date.NowTimestamp()

	if simple.IsNotBlank(name) {
		sysConfig.Name = name
	}
	if simple.IsNotBlank(description) {
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

func (s *sysConfigService) GetTokenExpireDays() int {
	tokenExpireDaysStr := cache.SysConfigCache.GetValue(constants.SysConfigTokenExpireDays)
	tokenExpireDays, err := strconv.Atoi(tokenExpireDaysStr)
	if err != nil {
		tokenExpireDays = constants.DefaultTokenExpireDays
	}
	if tokenExpireDays <= 0 {
		tokenExpireDays = constants.DefaultTokenExpireDays
	}
	return tokenExpireDays
}

func (s *sysConfigService) GetLoginMethod() model.LoginMethod {
	loginMethodStr := cache.SysConfigCache.GetValue(constants.SysConfigLoginMethod)

	useDefault := true
	var loginMethod model.LoginMethod
	if simple.IsNotBlank(loginMethodStr) {
		if err := json.Parse(loginMethodStr, &loginMethod); err != nil {
			logrus.Warn("登录方式数据错误", err)
		} else {
			useDefault = false
		}
	}
	if useDefault {
		loginMethod = model.LoginMethod{
			Password: true,
			QQ:       true,
			Github:   true,
		}
	}
	return loginMethod
}

func (s *sysConfigService) GetConfig() *model.SysConfigResponse {
	var (
		siteTitle             = cache.SysConfigCache.GetValue(constants.SysConfigSiteTitle)
		siteDescription       = cache.SysConfigCache.GetValue(constants.SysConfigSiteDescription)
		siteKeywords          = cache.SysConfigCache.GetValue(constants.SysConfigSiteKeywords)
		siteNavs              = cache.SysConfigCache.GetValue(constants.SysConfigSiteNavs)
		siteNotification      = cache.SysConfigCache.GetValue(constants.SysConfigSiteNotification)
		recommendTags         = cache.SysConfigCache.GetValue(constants.SysConfigRecommendTags)
		urlRedirect           = cache.SysConfigCache.GetValue(constants.SysConfigUrlRedirect)
		scoreConfigStr        = cache.SysConfigCache.GetValue(constants.SysConfigScoreConfig)
		defaultNodeIdStr      = cache.SysConfigCache.GetValue(constants.SysConfigDefaultNodeId)
		articlePending        = cache.SysConfigCache.GetValue(constants.SysConfigArticlePending)
		topicCaptcha          = cache.SysConfigCache.GetValue(constants.SysConfigTopicCaptcha)
		userObserveSecondsStr = cache.SysConfigCache.GetValue(constants.SysConfigUserObserveSeconds)
		tokenExpireDays       = s.GetTokenExpireDays()
		loginMethod           = s.GetLoginMethod()
	)

	var siteKeywordsArr []string
	if simple.IsNotBlank(siteKeywords) {
		if err := json.Parse(siteKeywords, &siteKeywordsArr); err != nil {
			logrus.Warn("站点关键词数据错误", err)
		}
	}

	var siteNavsArr []model.ActionLink
	if simple.IsNotBlank(siteNavs) {
		if err := json.Parse(siteNavs, &siteNavsArr); err != nil {
			logrus.Warn("站点导航数据错误", err)
		}
	}

	var recommendTagsArr []string
	if simple.IsNotBlank(recommendTags) {
		if err := json.Parse(recommendTags, &recommendTagsArr); err != nil {
			logrus.Warn("推荐标签数据错误", err)
		}
	}

	var scoreConfig model.ScoreConfig
	if simple.IsNotBlank(scoreConfigStr) {
		if err := json.Parse(scoreConfigStr, &scoreConfig); err != nil {
			logrus.Warn("积分配置错误", err)
		}
	}

	var (
		defaultNodeId      = number.ToInt64(defaultNodeIdStr)
		userObserveSeconds = number.ToInt(userObserveSecondsStr)
	)

	if tokenExpireDays <= 0 {
		tokenExpireDays = 7
	}

	return &model.SysConfigResponse{
		SiteTitle:          siteTitle,
		SiteDescription:    siteDescription,
		SiteKeywords:       siteKeywordsArr,
		SiteNavs:           siteNavsArr,
		SiteNotification:   siteNotification,
		RecommendTags:      recommendTagsArr,
		UrlRedirect:        strings.ToLower(urlRedirect) == "true",
		ScoreConfig:        scoreConfig,
		DefaultNodeId:      defaultNodeId,
		ArticlePending:     strings.ToLower(articlePending) == "true",
		TopicCaptcha:       strings.ToLower(topicCaptcha) == "true",
		UserObserveSeconds: userObserveSeconds,
		TokenExpireDays:    tokenExpireDays,
		LoginMethod:        loginMethod,
	}
}

func (s *sysConfigService) GetInt(key string) int {
	value := cache.SysConfigCache.GetValue(key)
	if simple.IsBlank(value) {
		return 0
	}
	ret, _ := strconv.Atoi(value)
	return ret
}
