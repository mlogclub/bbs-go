package services

import (
	"bbs-go/internal/models/constants"
	"bbs-go/internal/models/dto"
	"errors"
	"fmt"
	"log/slog"
	"strings"

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

const (
	maxScriptInjectionCount   = 20
	maxScriptInjectionCodeLen = 20 * 1024
	maxScriptInjectionNameLen = 200
)

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

	if siteNavs := json.Get(constants.SysConfigSiteNavs); siteNavs.Exists() {
		if err := validateSiteNavs(siteNavs.String()); err != nil {
			return err
		}
	}
	if scriptInjections := json.Get(constants.SysConfigScriptInjections); scriptInjections.Exists() {
		if err := validateScriptInjections(scriptInjections.String()); err != nil {
			return err
		}
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
func (s *sysConfigService) Set(key, value string) error {
	return sqls.DB().Transaction(func(tx *gorm.DB) error {
		if err := s.setSingle(tx, key, value, "", ""); err != nil {
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

func (s *sysConfigService) GetTokenExpireDays() int {
	tokenExpireDays := cache.SysConfigCache.GetInt(constants.SysConfigTokenExpireDays)
	if tokenExpireDays <= 0 {
		tokenExpireDays = constants.DefaultTokenExpireDays
	}
	return tokenExpireDays
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

func (s *sysConfigService) IsTopicCaptcha() bool {
	return cache.SysConfigCache.GetBool(constants.SysConfigTopicCaptcha)
}

func (s *sysConfigService) GetDefaultNodeId() int64 {
	return cache.SysConfigCache.GetInt64(constants.SysConfigDefaultNodeId)
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

// GetEmailNoticeIntervalSeconds 邮件通知间隔（秒），<=0 表示不限制
func (s *sysConfigService) GetEmailNoticeIntervalSeconds() int {
	return cache.SysConfigCache.GetInt(constants.SysConfigEmailNoticeIntervalSeconds)
}

func (s *sysConfigService) IsUrlRedirect() bool {
	return cache.SysConfigCache.GetBool(constants.SysConfigUrlRedirect)
}

func (s *sysConfigService) GetLoginConfig() dto.LoginConfig {
	str := cache.SysConfigCache.GetStr(constants.SysConfigLoginConfig)
	var loginConfig dto.LoginConfig
	if err := jsons.Parse(str, &loginConfig); err != nil {
		slog.Warn("登录配置错误", slog.Any("err", err))
	}

	// 如果全部禁用，那么默认启用密码登录
	if loginConfig.IsAllDisabled() {
		loginConfig.PasswordLogin.Enabled = true
	}

	return loginConfig
}

func (s *sysConfigService) GetUploadConfig() dto.UploadConfig {
	str := cache.SysConfigCache.GetStr(constants.SysConfigUploadConfig)
	var uploadConfig dto.UploadConfig
	if err := jsons.Parse(str, &uploadConfig); err != nil {
		slog.Warn("上传配置错误", slog.Any("err", err))
	}
	return uploadConfig
}

func (s *sysConfigService) GetSmtpConfig() dto.SmtpConfig {
	str := cache.SysConfigCache.GetStr(constants.SysConfigSmtpConfig)
	var smtpConfig dto.SmtpConfig
	if strings.TrimSpace(str) == "" {
		return smtpConfig
	}
	if err := jsons.Parse(str, &smtpConfig); err != nil {
		slog.Warn("smtp配置错误", slog.Any("err", err))
	}
	return smtpConfig
}

func (s *sysConfigService) GetScriptInjections() []dto.ScriptInjection {
	str := cache.SysConfigCache.GetStr(constants.SysConfigScriptInjections)
	if strings.TrimSpace(str) == "" {
		return []dto.ScriptInjection{}
	}

	var injections []dto.ScriptInjection
	if err := jsons.Parse(str, &injections); err != nil {
		slog.Warn("脚本注入配置错误", slog.Any("err", err))
		return []dto.ScriptInjection{}
	}
	return injections
}

func (s *sysConfigService) GetBaseURL() string {
	baseURL := strings.TrimSpace(cache.SysConfigCache.GetStr(constants.SysConfigBaseURL))
	if baseURL == "" {
		return "/"
	}
	for len(baseURL) > 1 && strings.HasSuffix(baseURL, "/") {
		baseURL = strings.TrimSuffix(baseURL, "/")
	}
	return baseURL
}

func validateSiteNavs(siteNavsJson string) error {
	if strings.TrimSpace(siteNavsJson) == "" {
		return nil
	}
	var navs []dto.ActionLink
	if err := jsons.Parse(siteNavsJson, &navs); err != nil {
		return errors.New("invalid site navigation data format")
	}
	return validateActionLinks(navs, 1)
}

func validateActionLinks(navs []dto.ActionLink, depth int) error {
	if depth > 2 {
		return errors.New("site navigation supports at most two levels")
	}
	for idx, nav := range navs {
		if strings.TrimSpace(nav.Title) == "" {
			return fmt.Errorf("navigation title is required at item %d", idx+1)
		}
		if depth == 1 {
			if len(nav.Children) == 0 && strings.TrimSpace(nav.Url) == "" {
				return fmt.Errorf("primary navigation URL is required at item %d", idx+1)
			}
			if len(nav.Children) > 0 {
				if err := validateActionLinks(nav.Children, depth+1); err != nil {
					return err
				}
			}
			continue
		}
		if strings.TrimSpace(nav.Url) == "" {
			return fmt.Errorf("secondary navigation URL is required at item %d", idx+1)
		}
		if len(nav.Children) > 0 {
			return errors.New("site navigation supports at most two levels")
		}
	}
	return nil
}

func validateScriptInjections(scriptInjectionsJSON string) error {
	if strings.TrimSpace(scriptInjectionsJSON) == "" {
		return nil
	}

	var injections []dto.ScriptInjection
	if err := jsons.Parse(scriptInjectionsJSON, &injections); err != nil {
		return errors.New("invalid script injections data format")
	}

	if len(injections) > maxScriptInjectionCount {
		return fmt.Errorf("too many script injections, max %d", maxScriptInjectionCount)
	}

	for idx, injection := range injections {
		scriptName := strings.TrimSpace(injection.ScriptName)
		if scriptName == "" {
			return fmt.Errorf("script injection name is required at item %d", idx+1)
		}
		if len([]rune(scriptName)) > maxScriptInjectionNameLen {
			return fmt.Errorf("script injection name is too long at item %d", idx+1)
		}
		injectionType := strings.TrimSpace(injection.Type)
		switch injectionType {
		case "external":
			if strings.TrimSpace(injection.Src) == "" {
				return fmt.Errorf("script injection src is required at item %d", idx+1)
			}
		case "inline":
			if strings.TrimSpace(injection.Code) == "" {
				return fmt.Errorf("script injection code is required at item %d", idx+1)
			}
			if len(injection.Code) > maxScriptInjectionCodeLen {
				return fmt.Errorf("script injection code is too long at item %d", idx+1)
			}
		default:
			return fmt.Errorf("script injection type must be external or inline at item %d", idx+1)
		}
	}
	return nil
}
