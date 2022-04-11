package services

import (
	"bbs-go/model/constants"
	"errors"
	"strconv"
	"strings"

	"github.com/mlogclub/simple/common/dates"
	"github.com/mlogclub/simple/common/jsons"
	"github.com/mlogclub/simple/common/numbers"
	"github.com/mlogclub/simple/common/strs"
	"github.com/mlogclub/simple/sqls"
	"github.com/mlogclub/simple/web/params"
	"github.com/tidwall/gjson"

	"github.com/sirupsen/logrus"
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
	return repositories.SysConfigRepository.Get(sqls.DB(), id)
}

func (s *sysConfigService) Take(where ...interface{}) *model.SysConfig {
	return repositories.SysConfigRepository.Take(sqls.DB(), where...)
}

func (s *sysConfigService) Find(cnd *sqls.Cnd) []model.SysConfig {
	return repositories.SysConfigRepository.Find(sqls.DB(), cnd)
}

func (s *sysConfigService) FindOne(cnd *sqls.Cnd) *model.SysConfig {
	return repositories.SysConfigRepository.FindOne(sqls.DB(), cnd)
}

func (s *sysConfigService) FindPageByParams(params *params.QueryParams) (list []model.SysConfig, paging *sqls.Paging) {
	return repositories.SysConfigRepository.FindPageByParams(sqls.DB(), params)
}

func (s *sysConfigService) FindPageByCnd(cnd *sqls.Cnd) (list []model.SysConfig, paging *sqls.Paging) {
	return repositories.SysConfigRepository.FindPageByCnd(sqls.DB(), cnd)
}

func (s *sysConfigService) GetAll() []model.SysConfig {
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
		sysConfig = &model.SysConfig{
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
	if strs.IsNotBlank(loginMethodStr) {
		if err := jsons.Parse(loginMethodStr, &loginMethod); err != nil {
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

func (s *sysConfigService) IsCreateTopicEmailVerified() bool {
	value := cache.SysConfigCache.GetValue(constants.SysConfigCreateTopicEmailVerified)
	return strs.EqualsIgnoreCase(value, "true") || strs.EqualsIgnoreCase(value, "1")
}

func (s *sysConfigService) IsCreateArticleEmailVerified() bool {
	value := cache.SysConfigCache.GetValue(constants.SysConfigCreateArticleEmailVerified)
	return strs.EqualsIgnoreCase(value, "true") || strs.EqualsIgnoreCase(value, "1")
}

func (s *sysConfigService) IsCreateCommentEmailVerified() bool {
	value := cache.SysConfigCache.GetValue(constants.SysConfigCreateCommentEmailVerified)
	return strs.EqualsIgnoreCase(value, "true") || strs.EqualsIgnoreCase(value, "1")
}

func (s *sysConfigService) IsEnableHideContent() bool {
	value := cache.SysConfigCache.GetValue(constants.SysConfigEnableHideContent)
	return strs.EqualsIgnoreCase(value, "true") || strs.EqualsIgnoreCase(value, "1")
}

func (s *sysConfigService) GetSiteNavs() []model.ActionLink {
	siteNavs := cache.SysConfigCache.GetValue(constants.SysConfigSiteNavs)
	var siteNavsArr []model.ActionLink
	if strs.IsNotBlank(siteNavs) {
		if err := jsons.Parse(siteNavs, &siteNavsArr); err != nil {
			logrus.Warn("站点导航数据错误", err)
		}
	}
	return siteNavsArr
}

func (s *sysConfigService) GetConfig() *model.SysConfigResponse {
	var (
		siteTitle                  = cache.SysConfigCache.GetValue(constants.SysConfigSiteTitle)
		siteDescription            = cache.SysConfigCache.GetValue(constants.SysConfigSiteDescription)
		siteKeywords               = cache.SysConfigCache.GetValue(constants.SysConfigSiteKeywords)
		siteNotification           = cache.SysConfigCache.GetValue(constants.SysConfigSiteNotification)
		recommendTags              = cache.SysConfigCache.GetValue(constants.SysConfigRecommendTags)
		urlRedirect                = cache.SysConfigCache.GetValue(constants.SysConfigUrlRedirect)
		scoreConfigStr             = cache.SysConfigCache.GetValue(constants.SysConfigScoreConfig)
		defaultNodeIdStr           = cache.SysConfigCache.GetValue(constants.SysConfigDefaultNodeId)
		articlePending             = cache.SysConfigCache.GetValue(constants.SysConfigArticlePending)
		topicCaptcha               = cache.SysConfigCache.GetValue(constants.SysConfigTopicCaptcha)
		userObserveSecondsStr      = cache.SysConfigCache.GetValue(constants.SysConfigUserObserveSeconds)
		siteNavs                   = s.GetSiteNavs()
		tokenExpireDays            = s.GetTokenExpireDays()
		loginMethod                = s.GetLoginMethod()
		createTopicEmailVerified   = s.IsCreateTopicEmailVerified()
		createArticleEmailVerified = s.IsCreateArticleEmailVerified()
		createCommentEmailVerified = s.IsCreateCommentEmailVerified()
		enableHideContent          = s.IsEnableHideContent()
	)

	var siteKeywordsArr []string
	if strs.IsNotBlank(siteKeywords) {
		if err := jsons.Parse(siteKeywords, &siteKeywordsArr); err != nil {
			logrus.Warn("站点关键词数据错误", err)
		}
	}

	var recommendTagsArr []string
	if strs.IsNotBlank(recommendTags) {
		if err := jsons.Parse(recommendTags, &recommendTagsArr); err != nil {
			logrus.Warn("推荐标签数据错误", err)
		}
	}

	var scoreConfig model.ScoreConfig
	if strs.IsNotBlank(scoreConfigStr) {
		if err := jsons.Parse(scoreConfigStr, &scoreConfig); err != nil {
			logrus.Warn("积分配置错误", err)
		}
	}

	var (
		defaultNodeId      = numbers.ToInt64(defaultNodeIdStr)
		userObserveSeconds = numbers.ToInt(userObserveSecondsStr)
	)

	if tokenExpireDays <= 0 {
		tokenExpireDays = 7
	}

	return &model.SysConfigResponse{
		SiteTitle:                  siteTitle,
		SiteDescription:            siteDescription,
		SiteKeywords:               siteKeywordsArr,
		SiteNavs:                   siteNavs,
		SiteNotification:           siteNotification,
		RecommendTags:              recommendTagsArr,
		UrlRedirect:                strings.ToLower(urlRedirect) == "true",
		ScoreConfig:                scoreConfig,
		DefaultNodeId:              defaultNodeId,
		ArticlePending:             strings.ToLower(articlePending) == "true",
		TopicCaptcha:               strings.ToLower(topicCaptcha) == "true",
		UserObserveSeconds:         userObserveSeconds,
		TokenExpireDays:            tokenExpireDays,
		LoginMethod:                loginMethod,
		CreateTopicEmailVerified:   createTopicEmailVerified,
		CreateArticleEmailVerified: createArticleEmailVerified,
		CreateCommentEmailVerified: createCommentEmailVerified,
		EnableHideContent:          enableHideContent,
	}
}

func (s *sysConfigService) GetInt(key string) int {
	value := cache.SysConfigCache.GetValue(key)
	if strs.IsBlank(value) {
		return 0
	}
	ret, _ := strconv.Atoi(value)
	return ret
}
