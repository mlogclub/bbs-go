package api

import (
	"github.com/kataras/iris/v12"
	"github.com/mlogclub/simple/web"

	"bbs-go/internal/cache"
	"bbs-go/internal/models/constants"
	"bbs-go/internal/models/dto"
	"bbs-go/internal/pkg/config"
	"bbs-go/internal/services"
)

type ConfigController struct {
	Ctx iris.Context
}

func (c *ConfigController) GetConfigs() *web.JsonResult {
	cfg := config.Instance

	var b *web.RspBuilder
	if cfg.Installed {
		loginConfig := services.SysConfigService.GetLoginConfig()
		sysConfig := &dto.SysConfigOpenResponse{
			SiteTitle:                  cache.SysConfigCache.GetStr(constants.SysConfigSiteTitle),
			SiteDescription:            cache.SysConfigCache.GetStr(constants.SysConfigSiteDescription),
			SiteKeywords:               cache.SysConfigCache.GetStrArr(constants.SysConfigSiteKeywords),
			SiteLogo:                   cache.SysConfigCache.GetStr(constants.SysConfigSiteLogo),
			SiteNavs:                   services.SysConfigService.GetSiteNavs(),
			SiteNotification:           cache.SysConfigCache.GetStr(constants.SysConfigSiteNotification),
			RecommendTags:              cache.SysConfigCache.GetStrArr(constants.SysConfigRecommendTags),
			UrlRedirect:                services.SysConfigService.IsUrlRedirect(),
			DefaultNodeId:              services.SysConfigService.GetDefaultNodeId(),
			ArticlePending:             services.SysConfigService.IsArticlePending(),
			TopicCaptcha:               services.SysConfigService.IsTopicCaptcha(),
			UserObserveSeconds:         cache.SysConfigCache.GetInt(constants.SysConfigUserObserveSeconds),
			TokenExpireDays:            services.SysConfigService.GetTokenExpireDays(),
			CreateTopicEmailVerified:   services.SysConfigService.IsCreateTopicEmailVerified(),
			CreateArticleEmailVerified: services.SysConfigService.IsCreateArticleEmailVerified(),
			CreateCommentEmailVerified: services.SysConfigService.IsCreateCommentEmailVerified(),
			EnableHideContent:          services.SysConfigService.IsEnableHideContent(),
			Modules:                    services.SysConfigService.GetModules(),
			EmailNoticeIntervalSeconds: services.SysConfigService.GetEmailNoticeIntervalSeconds(),
			LoginConfig: dto.OpenLoginConfig{
				PasswordLogin: loginConfig.PasswordLogin,
				WeixinLogin:   dto.EnabledConfig{Enabled: loginConfig.WeixinLogin.Enabled},
				SmsLogin:      dto.EnabledConfig{Enabled: loginConfig.SmsLogin.Enabled},
				GoogleLogin:   dto.EnabledConfig{Enabled: loginConfig.GoogleLogin.Enabled},
			},
		}
		b = web.NewRspBuilder(sysConfig)
	} else {
		b = web.NewEmptyRspBuilder()
	}
	b.Put("installed", cfg.Installed)
	b.Put("language", cfg.Language)
	return b.JsonResult()
}
