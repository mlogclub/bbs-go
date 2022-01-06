package model

import (
	"bbs-go/model/constants"
	"bbs-go/pkg/common"
	"strings"

	"github.com/kataras/iris/v12"
	"github.com/mlogclub/simple"
	"github.com/tidwall/gjson"
)

type CreateTopicForm struct {
	Type        constants.TopicType
	CaptchaId   string
	CaptchaCode string
	NodeId      int64
	Title       string
	Content     string
	Tags        []string
	ImageList   []ImageDTO
	UserAgent   string
	Ip          string
}

type CreateArticleForm struct {
	Title       string
	Summary     string
	Content     string
	ContentType string
	Tags        []string
	SourceUrl   string
}

// CreateCommentForm 发表评论
type CreateCommentForm struct {
	EntityType  string     `form:"entityType"`
	EntityId    int64      `form:"entityId"`
	Content     string     `form:"content"`
	ImageList   []ImageDTO `form:"imageList"`
	QuoteId     int64      `form:"quoteId"`
	ContentType string     `form:"contentType"`
	UserAgent   string     `form:"userAgent"`
	Ip          string     `form:"ip"`
}

type ImageDTO struct {
	Url string `json:"url"`
}

func GetCreateTopicForm(ctx iris.Context) CreateTopicForm {
	var (
		topicType = simple.FormValueIntDefault(ctx, "type", int(constants.TopicTypeTopic))
	)
	return CreateTopicForm{
		Type:        constants.TopicType(topicType),
		CaptchaId:   simple.FormValue(ctx, "captchaId"),
		CaptchaCode: simple.FormValue(ctx, "captchaCode"),
		NodeId:      simple.FormValueInt64Default(ctx, "nodeId", 0),
		Title:       strings.TrimSpace(simple.FormValue(ctx, "title")),
		Content:     strings.TrimSpace(simple.FormValue(ctx, "content")),
		Tags:        simple.FormValueStringArray(ctx, "tags"),
		ImageList:   GetImageList(ctx, "imageList"),
		UserAgent:   common.GetUserAgent(ctx.Request()),
		Ip:          common.GetRequestIP(ctx.Request()),
	}
}

func GetCreateCommentForm(ctx iris.Context) CreateCommentForm {
	form := CreateCommentForm{
		EntityType:  simple.FormValue(ctx, "entityType"),
		EntityId:    simple.FormValueInt64Default(ctx, "entityId", 0),
		Content:     strings.TrimSpace(simple.FormValue(ctx, "content")),
		ImageList:   GetImageList(ctx, "imageList"),
		QuoteId:     simple.FormValueInt64Default(ctx, "quoteId", 0),
		ContentType: simple.FormValueDefault(ctx, "contentType", constants.ContentTypeText),
		UserAgent:   common.GetUserAgent(ctx.Request()),
		Ip:          common.GetRequestIP(ctx.Request()),
	}
	return form
}

func GetImageList(ctx iris.Context, paramName string) []ImageDTO {
	imageListStr := simple.FormValue(ctx, paramName)
	var imageList []ImageDTO
	if simple.IsNotBlank(imageListStr) {
		ret := gjson.Parse(imageListStr)
		if ret.IsArray() {
			for _, item := range ret.Array() {
				url := item.Get("url").String()
				imageList = append(imageList, ImageDTO{
					Url: url,
				})
			}
		}
	}
	return imageList
}
