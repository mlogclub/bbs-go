package collect

import (
	"bbs-go/model/constants"
	"bbs-go/package/uploader"
	"errors"
	"github.com/mlogclub/simple/date"
	"strings"

	"bbs-go/services"

	"github.com/mlogclub/simple"

	"bbs-go/model"
)

type WxbotApi struct {
}

func NewWxbotApi() *WxbotApi {
	return &WxbotApi{}
}

func (api *WxbotApi) Publish(article *WxArticle) (*model.Article, error) {
	if len(article.Title) == 0 || len(article.HtmlContent) == 0 {
		return nil, errors.New("内容为空")
	}

	userId, _ := api.initUser(article)
	tags := api.initTags(article)
	summary := article.Summary
	if simple.RuneLen(summary) == 0 {
		summary = simple.GetSummary(article.TextContent, constants.SummaryLen)
	}

	return services.ArticleService.Publish(userId, article.Title, summary, article.HtmlContent, constants.ContentTypeHtml,
		tags, article.Url)
}

func (api *WxbotApi) initUser(article *WxArticle) (int64, error) {
	user := services.UserService.GetByUsername(article.AppID)
	if user != nil {
		user.Nickname = article.AppName
		user.Description = article.WxIntro
		_ = services.UserService.Update(user)
		return user.Id, nil
	} else {
		avatar, err := uploader.CopyImage(article.OriHead)
		if err != nil {
			return 0, err
		}
		user := &model.User{
			Username:    simple.SqlNullString(article.AppID),
			Nickname:    article.AppName,
			Description: article.WxIntro,
			Avatar:      avatar,
			Status:      constants.StatusOk,
			Type:        constants.UserTypeGzh,
			CreateTime:  date.NowTimestamp(),
			UpdateTime:  date.NowTimestamp(),
		}
		err = services.UserService.Create(user)
		if err != nil {
			return 0, err
		}
		return user.Id, nil
	}
}

func (api *WxbotApi) initTags(wxArticle *WxArticle) []string {
	var tagNames []string

	if len(wxArticle.Categories) > 0 {
		ss := strings.Split(wxArticle.Categories, ",")
		if ss != nil && len(ss) > 0 {
			for _, s := range ss {
				s = strings.TrimSpace(s)
				if len(s) > 0 {
					tagNames = append(tagNames, s)
				}
			}
		}
	}

	if len(wxArticle.Tags) > 0 {
		ss := strings.Split(wxArticle.Tags, ",")
		if ss != nil && len(ss) > 0 {
			for _, s := range ss {
				s = strings.TrimSpace(s)
				if len(s) > 0 {
					tagNames = append(tagNames, s)
				}
			}
		}
	}
	return tagNames
}

// WxArticle 微信文章
type WxArticle struct {
	Id          int64  `json:"id"`          // 编号
	Title       string `json:"title"`       // 标题
	Author      string `json:"author"`      // 作者
	AppName     string `json:"appName"`     // 公众号名称
	AppID       string `json:"appId"`       // 公众号ID
	Cover       string `json:"cover"`       // 文章封面
	Intro       string `json:"intro"`       // 描述
	HtmlContent string `json:"htmlContent"` // 公众号文章html内容
	MdContent   string `json:"mdContent"`   // 公众号文章md内容
	TextContent string `json:"textContent"` // 文本内容
	Summary     string `json:"summary"`     // 摘要
	PubAt       string `json:"pubAt"`       // 发布时间
	UrlMd5      string `json:"urlMd5"`      // 链接地址的md5
	RoundHead   string `json:"roundHead"`   // 圆头像
	OriHead     string `json:"oriHead"`     // 原头像
	Url         string `json:"url"`         // 微信文章链接地址
	SourceURL   string `json:"sourceUrl"`   // 公众号原文地址
	ArticleId   int64  `json:"articleId"`   // 发布线上返回的id
	Tags        string `json:"tags"`        // 标签字符串
	Category    string `json:"category"`    // 一级分类
	Categories  string `json:"categories"`  // 二级分类
	Copyright   string `json:"copyright"`   // 已经 0,1,2   微小宝那 1 标识为原创
	Video       string `json:"video"`       // 视频地址
	Audio       string `json:"audio"`       // 音频地址
	WxID        string `json:"wxId"`        // 微信公众号ID
	WxIntro     string `json:"wxIntro"`     // 微信公众号介绍
	Images      string `json:"images"`      // 图片
	PublishTime int64  `json:"publishTime"` // 采集器发布时间
	CreateTime  int64  `json:"createTime"`  // 创建时间
	UpdatedTime int64  `json:"updatedTime"` // 更新时间
}
