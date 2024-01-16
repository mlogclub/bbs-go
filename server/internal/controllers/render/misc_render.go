package render

import (
	"bbs-go/internal/models/constants"
	"bbs-go/internal/pkg/bbsurls"
	"strings"
	"time"

	"github.com/kataras/iris/v12"
	"github.com/kataras/iris/v12/context"
	"github.com/microcosm-cc/bluemonday"

	"github.com/PuerkitoBio/goquery"
	"github.com/mlogclub/simple/common/strs"
	"github.com/mlogclub/simple/common/urls"
	"github.com/mlogclub/simple/web"

	"bbs-go/internal/models"
	"bbs-go/internal/services"
)

func xssProtection(htmlContent string) string {
	ugcProtection := bluemonday.UGCPolicy() // 用户生成内容模式
	ugcProtection.AllowAttrs("class").OnElements("code")
	ugcProtection.AllowAttrs("start").OnElements("ol", "ul", "li")
	return ugcProtection.Sanitize(htmlContent)
}

// handleHtmlContent 处理html内容
func handleHtmlContent(htmlContent string) string {
	htmlContent = xssProtection(htmlContent)
	doc, err := goquery.NewDocumentFromReader(strings.NewReader(htmlContent))
	if err != nil {
		return htmlContent
	}

	doc.Find("a").Each(func(_ int, selection *goquery.Selection) {
		href := selection.AttrOr("href", "")

		if strs.IsBlank(href) {
			return
		}

		// 不是内部链接
		if !bbsurls.IsInternalUrl(href) {
			selection.SetAttr("target", "_blank")
			selection.SetAttr("rel", "external nofollow") // 标记站外链接，搜索引擎爬虫不传递权重值

			_config := services.SysConfigService.GetConfig()
			if _config.UrlRedirect { // 开启非内部链接跳转
				newHref := urls.ParseUrl(bbsurls.AbsUrl("/redirect")).AddQuery("url", href).BuildStr()
				selection.SetAttr("href", newHref)
			}
		}

		// 如果a标签没有title，那么设置title
		title := selection.AttrOr("title", "")
		if len(title) == 0 {
			selection.SetAttr("title", selection.Text())
		}
	})

	// 处理图片
	doc.Find("img").Each(func(_ int, selection *goquery.Selection) {
		src := selection.AttrOr("src", "")

		// 处理第三方图片
		if strings.Contains(src, "qpic.cn") {
			src = urls.ParseUrl("/api/img/proxy").AddQuery("url", src).BuildStr()
		}

		// 处理图片样式
		src = HandleOssImageStyleDetail(src)

		// // 处理lazyload
		// selection.SetAttr("data-src", src)
		// selection.RemoveAttr("src")

		selection.SetAttr("src", src)
	})

	if htmlStr, err := doc.Find("body").Html(); err == nil {
		return htmlStr
	}
	return htmlContent
}

/*
BuildLoginSuccess 处理登录成功后的返回数据

Parameter:

	user - login user
	redirect - 登录来源地址，需要控制登录成功之后跳转到该地址
*/
func BuildLoginSuccess(ctx iris.Context, user *models.User, redirect string) *web.JsonResult {
	token, err := services.UserTokenService.Generate(user.Id)
	if err != nil {
		return web.JsonError(err)
	}
	ctx.SetCookieKV(constants.CookieTokenKey, token, context.CookieHTTPOnly(true), context.CookieExpires(365*24*time.Hour))
	return web.NewEmptyRspBuilder().
		Put("token", token).
		Put("user", BuildUserProfile(user)).
		Put("redirect", redirect).JsonResult()
}
