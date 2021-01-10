package model

import (
	"bbs-go/model/constants"
	"github.com/kataras/iris/v12"
	"github.com/mlogclub/simple"
	"github.com/tidwall/gjson"
	"strings"
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
}

func GetCreateTopicForm(ctx iris.Context) CreateTopicForm {
	var (
		topicType    = simple.FormValueIntDefault(ctx, "type", int(constants.TopicTypeTopic))
		imageListStr = simple.FormValue(ctx, "imageList")
	)
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
	return CreateTopicForm{
		Type:        constants.TopicType(topicType),
		CaptchaId:   simple.FormValue(ctx, "captchaId"),
		CaptchaCode: simple.FormValue(ctx, "captchaCode"),
		NodeId:      simple.FormValueInt64Default(ctx, "nodeId", 0),
		Title:       strings.TrimSpace(simple.FormValue(ctx, "title")),
		Content:     strings.TrimSpace(simple.FormValue(ctx, "content")),
		Tags:        simple.FormValueStringArray(ctx, "tags"),
		ImageList:   imageList,
	}
}

// 发表评论
type CreateCommentForm struct {
	EntityType  string `form:"entityType"`
	EntityId    int64  `form:"entityId"`
	Content     string `form:"content"`
	QuoteId     int64  `form:"quoteId"`
	ContentType string `form:"contentType"`
}

type ImageDTO struct {
	Url string `json:"url"`
}
