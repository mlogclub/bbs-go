package render

import (
	"bbs-go/package/urls"
	"strings"

	"github.com/PuerkitoBio/goquery"
	"github.com/mlogclub/simple"

	"bbs-go/model"
	"bbs-go/services"
)

// handleHtmlContent 处理html内容
func handleHtmlContent(htmlContent string) string {
	doc, err := goquery.NewDocumentFromReader(strings.NewReader(htmlContent))
	if err != nil {
		return htmlContent
	}

	doc.Find("a").Each(func(i int, selection *goquery.Selection) {
		href := selection.AttrOr("href", "")

		if simple.IsBlank(href) {
			return
		}

		// 不是内部链接
		if !urls.IsInternalUrl(href) {
			selection.SetAttr("target", "_blank")
			selection.SetAttr("rel", "external nofollow") // 标记站外链接，搜索引擎爬虫不传递权重值

			_config := services.SysConfigService.GetConfig()
			if _config.UrlRedirect { // 开启非内部链接跳转
				newHref := simple.ParseUrl(urls.AbsUrl("/redirect")).AddQuery("url", href).BuildStr()
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
	doc.Find("img").Each(func(i int, selection *goquery.Selection) {
		src := selection.AttrOr("src", "")
		// 处理第三方图片
		if strings.Contains(src, "qpic.cn") {
			src = simple.ParseUrl("/api/img/proxy").AddQuery("url", src).BuildStr()
			// selection.SetAttr("src", src)
		}

		// 处理图片样式
		src = HandleOssImageStyleDetail(src)

		// 处理lazyload
		selection.SetAttr("data-src", src)
		selection.RemoveAttr("src")
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
	ref - 登录来源地址，需要控制登录成功之后跳转到该地址
*/
func BuildLoginSuccess(user *model.User, ref string) *simple.JsonResult {
	token, err := services.UserTokenService.Generate(user.Id)
	if err != nil {
		return simple.JsonErrorMsg(err.Error())
	}
	return simple.NewEmptyRspBuilder().
		Put("token", token).
		Put("user", BuildUser(user)).
		Put("ref", ref).JsonResult()
}
