package admin

import (
	"strconv"

	"github.com/kataras/iris/v12"
	"github.com/mlogclub/simple/common/strs"
	"github.com/mlogclub/simple/web"
	"github.com/mlogclub/simple/web/params"

	"bbs-go/internal/cache"
	"bbs-go/internal/models/constants"
	"bbs-go/internal/models/dto"
	"bbs-go/internal/services"
)

type SysConfigController struct {
	Ctx iris.Context
}

func (c *SysConfigController) GetBy(id int64) *web.JsonResult {
	t := services.SysConfigService.Get(id)
	if t == nil {
		return web.JsonErrorMsg("Not found, id=" + strconv.FormatInt(id, 10))
	}
	return web.JsonData(t)
}

func (c *SysConfigController) AnyList() *web.JsonResult {
	list, paging := services.SysConfigService.FindPageByParams(params.NewQueryParams(c.Ctx).PageByReq().Desc("id"))
	return web.JsonData(&web.PageResult{Results: list, Page: paging})
}

func (c *SysConfigController) GetConfigs() *web.JsonResult {
	resp := &dto.SysConfigAdminResponse{
		SiteTitle:                  cache.SysConfigCache.GetStr(constants.SysConfigSiteTitle),
		SiteDescription:            cache.SysConfigCache.GetStr(constants.SysConfigSiteDescription),
		BaseURL:                    services.SysConfigService.GetBaseURL(),
		SiteKeywords:               cache.SysConfigCache.GetStrArr(constants.SysConfigSiteKeywords),
		SiteLogo:                   cache.SysConfigCache.GetStr(constants.SysConfigSiteLogo),
		SiteNavs:                   services.SysConfigService.GetSiteNavs(),
		SiteNotification:           cache.SysConfigCache.GetStr(constants.SysConfigSiteNotification),
		AboutPageConfig:            services.SysConfigService.GetAboutPageConfig(),
		FooterLinks:                services.SysConfigService.GetFooterLinks(),
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
		EnableQaBounty:             services.SysConfigService.IsEnableQaBounty(),
		QaBountyMin:                services.SysConfigService.GetQaBountyMin(),
		QaBountyMax:                services.SysConfigService.GetQaBountyMax(),
		QaBountyRequired:           services.SysConfigService.IsQaBountyRequired(),
		Modules:                    services.SysConfigService.GetModules(),
		EmailWhitelist:             services.SysConfigService.GetEmailWhitelist(),
		EmailNoticeIntervalSeconds: services.SysConfigService.GetEmailNoticeIntervalSeconds(),
		NotificationTypes:          services.SysConfigService.GetNotificationTypes(),
		LoginConfig:                services.SysConfigService.GetLoginConfig(),
		SmtpConfig:                 services.SysConfigService.GetSmtpConfig(),
		UploadConfig:               services.SysConfigService.GetUploadConfig(),
		AttachmentConfig:           services.SysConfigService.GetAttachmentConfig(),
		ScriptInjections:           services.SysConfigService.GetScriptInjections(),
	}
	if strs.IsBlank(resp.SiteLogo) {
		resp.SiteLogo = "/res/images/logo.png"
	}
	return web.JsonData(resp)
}

func (c *SysConfigController) PostSave() *web.JsonResult {
	body, err := c.Ctx.GetBody()
	if err != nil {
		return web.JsonError(err)
	}
	if err := services.SysConfigService.SetAll(string(body)); err != nil {
		return web.JsonError(err)
	}
	return web.JsonSuccess()
}
