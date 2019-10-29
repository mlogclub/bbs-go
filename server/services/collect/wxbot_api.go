package collect

import (
	"strings"

	"github.com/mlogclub/bbs-go/common/oss"
	"github.com/mlogclub/bbs-go/services"

	"github.com/jinzhu/gorm"
	"github.com/kataras/iris/core/errors"
	"github.com/mlogclub/simple"

	"github.com/mlogclub/bbs-go/model"
)

type WxbotApi struct {
}

func NewWxbotApi() *WxbotApi {
	return &WxbotApi{}
}

func (this *WxbotApi) Publish(wxArticle *WxArticle) (*model.Article, error) {
	if len(wxArticle.Title) == 0 || len(wxArticle.HtmlContent) == 0 {
		return nil, errors.New("内容为空")
	}

	userId, _ := this.initUser(wxArticle)
	categoryId := this.initCategory(simple.GetDB(), wxArticle)
	tags := this.initTags(simple.GetDB(), wxArticle)
	summary := wxArticle.Summary
	if simple.RuneLen(summary) == 0 {
		summary = simple.GetSummary(wxArticle.TextContent, 256)
	}

	return services.ArticleService.Publish(userId, wxArticle.Title, summary,
		wxArticle.HtmlContent, model.ContentTypeHtml, categoryId, tags, wxArticle.Url, true)
}

func (this *WxbotApi) initUser(article *WxArticle) (int64, error) {
	user := services.UserService.GetByUsername(article.AppID)
	if user != nil {
		user.Nickname = article.AppName
		user.Description = article.WxIntro
		_ = services.UserService.Update(user)
		return user.Id, nil
	} else {
		avatar, err := oss.CopyImage(article.OriHead)
		if err != nil {
			return 0, err
		}
		user := &model.User{
			Username:    simple.SqlNullString(article.AppID),
			Nickname:    article.AppName,
			Description: article.WxIntro,
			Avatar:      avatar,
			Status:      model.UserStatusOk,
			Type:        model.UserTypeGzh,
			CreateTime:  simple.NowTimestamp(),
			UpdateTime:  simple.NowTimestamp(),
		}
		err = services.UserService.Create(user)
		if err != nil {
			return 0, err
		}
		return user.Id, nil
	}
}

func (this *WxbotApi) initCategory(db *gorm.DB, wxArticle *WxArticle) int64 {
	if len(wxArticle.Category) == 0 {
		return 0
	}
	cat, _ := services.CategoryService.GetOrCreate(wxArticle.Category)
	if cat != nil {
		return cat.Id
	}
	return 0
}

func (this *WxbotApi) initTags(db *gorm.DB, wxArticle *WxArticle) []string {
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
