package services

import (
	"errors"
	"strings"

	"github.com/jinzhu/gorm"
	"github.com/mlogclub/simple"
	"github.com/sirupsen/logrus"
	"github.com/tidwall/gjson"

	"bbs-go/model"
	"bbs-go/repositories"
	"bbs-go/services/cache"
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
	return simple.Tx(simple.DB(), func(tx *gorm.DB) error {
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
	return simple.Tx(simple.DB(), func(tx *gorm.DB) error {
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
	sysConfig := repositories.SysConfigRepository.GetByKey(simple.DB(), key)
	if sysConfig == nil {
		sysConfig = &model.SysConfig{
			CreateTime: simple.NowTimestamp(),
		}
	}
	sysConfig.Key = key
	sysConfig.Value = value
	sysConfig.UpdateTime = simple.NowTimestamp()

	if len(name) > 0 {
		sysConfig.Name = name
	}
	if len(description) > 0 {
		sysConfig.Description = description
	}

	var err error
	if sysConfig.Id > 0 {
		err = repositories.SysConfigRepository.Update(simple.DB(), sysConfig)
	} else {
		err = repositories.SysConfigRepository.Create(simple.DB(), sysConfig)
	}
	if err != nil {
		return err
	} else {
		cache.SysConfigCache.Invalidate(key)
		return nil
	}
}

func (s *sysConfigService) GetConfigResponse() *model.ConfigResponse {
	var (
		siteTitle        = cache.SysConfigCache.GetValue(model.SysConfigSiteTitle)
		siteDescription  = cache.SysConfigCache.GetValue(model.SysConfigSiteDescription)
		siteKeywords     = cache.SysConfigCache.GetValue(model.SysConfigSiteKeywords)
		siteNavs         = cache.SysConfigCache.GetValue(model.SysConfigSiteNavs)
		siteNotification = cache.SysConfigCache.GetValue(model.SysConfigSiteNotification)
		recommendTags    = cache.SysConfigCache.GetValue(model.SysConfigRecommendTags)
		urlRedirect      = cache.SysConfigCache.GetValue(model.SysConfigUrlRedirect)
	)

	var siteKeywordsArr []string
	if err := simple.ParseJson(siteKeywords, &siteKeywordsArr); err != nil {
		logrus.Warn("站点关键词数据错误", err)
	}

	var siteNavsArr []model.SiteNav
	if err := simple.ParseJson(siteNavs, &siteNavsArr); err != nil {
		logrus.Warn("站点导航数据错误", err)
	}

	var recommendTagsArr []string
	if err := simple.ParseJson(recommendTags, &recommendTagsArr); err != nil {
		logrus.Warn("推荐标签数据错误", err)
	}

	return &model.ConfigResponse{
		SiteTitle:        siteTitle,
		SiteDescription:  siteDescription,
		SiteKeywords:     siteKeywordsArr,
		SiteNavs:         siteNavsArr,
		SiteNotification: siteNotification,
		RecommendTags:    recommendTagsArr,
		UrlRedirect:      strings.ToLower(urlRedirect) == "true",
	}
}
