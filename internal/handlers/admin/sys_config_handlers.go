package admin

import (
	"io"
	"strconv"

	"github.com/gin-gonic/gin"

	"bbs-go/internal/pkg/ginx"
	"bbs-go/internal/pkg/params"

	"github.com/mlogclub/simple/common/strs"
	"github.com/mlogclub/simple/web"

	"bbs-go/internal/cache"
	"bbs-go/internal/models/constants"
	"bbs-go/internal/models/dto"
	"bbs-go/internal/services"
)

func SysConfigDetail(ctx *gin.Context) {
	id, err := strconv.ParseInt(ctx.Param("id"), 10, 64)
	if err != nil {
		ginx.WriteJSON(ctx, err)
		return
	}

	t := services.SysConfigService.Get(id)
	if t == nil {
		ginx.WriteJSON(ctx, ginx.ErrorMessage("Not found, id="+strconv.FormatInt(id, 10)))
		return
	}
	ginx.WriteJSON(ctx, t)

}

func SysConfigList(ctx *gin.Context) {
	list, paging := services.SysConfigService.FindPageByParams(params.NewQueryParams(ctx).PageByReq().Desc("id"))
	ginx.WriteJSON(ctx, &web.PageResult{Results: list, Page: paging})

}

func SysConfigConfigs(ctx *gin.Context) {

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
		DefaultCategoryId:          services.SysConfigService.GetDefaultCategoryId(),
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
	ginx.WriteJSON(ctx, resp)

}

func SysConfigSave(ctx *gin.Context) {
	body, err := io.ReadAll(ctx.Request.Body)
	if err != nil {
		ginx.WriteJSON(ctx, err)
		return
	}
	if err := services.SysConfigService.SetAll(string(body)); err != nil {
		ginx.WriteJSON(ctx, err)
		return
	}
	ginx.WriteJSON(ctx, nil)

}
