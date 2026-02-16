package req

import (
	"bbs-go/internal/models/constants"
	"bbs-go/internal/pkg/common"
	"log/slog"
	"strings"

	"github.com/kataras/iris/v12"
	"github.com/mlogclub/simple/common/jsons"
	"github.com/mlogclub/simple/common/strs"
	"github.com/mlogclub/simple/web"
	"github.com/mlogclub/simple/web/params"
	"github.com/tidwall/gjson"
)

type CreateTopicForm struct {
	Type        constants.TopicType   `json:"type"`
	NodeId      int64                 `json:"nodeId"`
	Title       string                `json:"title"`
	Content     string                `json:"content"`
	ContentType constants.ContentType `json:"contentType"`
	HideContent string                `json:"hideContent"`
	Tags        []string              `json:"tags"`
	ImageList   []ImageDTO            `json:"imageList"`
	Vote        *CreateVoteForm       `json:"vote"`
	UserAgent   string                `json:"userAgent"`
	Ip          string                `json:"ip"`

	CaptchaId       string `json:"captchaId"`
	CaptchaCode     string `json:"captchaCode"`
	CaptchaProtocol int    `json:"captchaProtocol"`
}

type CreateVoteForm struct {
	Type      constants.VoteType     `json:"type"`
	Title     string                 `json:"title"`
	ExpiredAt int64                  `json:"expiredAt"`
	VoteNum   int                    `json:"voteNum"`
	Options   []CreateVoteOptionForm `json:"options"`
}

type CreateVoteOptionForm struct {
	Content string `json:"content"`
}

type VoteCastForm struct {
	VoteId    int64   `json:"voteId"`
	OptionIds []int64 `json:"optionIds"`
}

type EditTopicForm struct {
	NodeId      int64    `json:"nodeId"`
	Title       string   `json:"title"`
	Content     string   `json:"content"`
	HideContent string   `json:"hideContent"`
	Tags        []string `json:"tags"`
}

type CreateArticleForm struct {
	Title       string
	Summary     string
	Content     string
	ContentType constants.ContentType
	Cover       *ImageDTO
	Tags        []string
	SourceUrl   string
}

// CreateCommentForm 发表评论
type CreateCommentForm struct {
	EntityType string     `form:"entityType"`
	EntityId   int64      `form:"entityId"`
	Content    string     `form:"content"`
	ImageList  []ImageDTO `form:"imageList"`
	QuoteId    int64      `form:"quoteId"`
	UserAgent  string     `form:"userAgent"`
	Ip         string     `form:"ip"`
}

type ImageDTO struct {
	Url string `json:"url"`
}

func GetCreateTopicForm(ctx iris.Context) CreateTopicForm {
	var form *CreateTopicForm
	if ctx.GetHeader("Content-Type") == "application/json" {
		if err := ctx.ReadJSON(&form); err != nil {
			slog.Error(err.Error(), slog.Any("err", err))
		}
	} else {
		form = &CreateTopicForm{
			Type:            constants.TopicType(params.FormValueIntDefault(ctx, "type", int(constants.TopicTypeTopic))),
			NodeId:          params.FormValueInt64Default(ctx, "nodeId", 0),
			Title:           strings.TrimSpace(params.FormValue(ctx, "title")),
			Content:         strings.TrimSpace(params.FormValue(ctx, "content")),
			ContentType:     constants.ContentType(params.FormValue(ctx, "contentType")),
			HideContent:     strings.TrimSpace(params.FormValue(ctx, "hideContent")),
			Tags:            params.FormValueStringArray(ctx, "tags"),
			ImageList:       GetImageList(ctx, "imageList"),
			CaptchaId:       params.FormValue(ctx, "captchaId"),
			CaptchaCode:     params.FormValue(ctx, "captchaCode"),
			CaptchaProtocol: params.FormValueIntDefault(ctx, "captchaProtocol", 0),
		}
	}

	if form.Type == constants.TopicTypeTweet {
		form.ContentType = constants.ContentTypeText
	}

	form.Ip = web.GetRequestIP(ctx.Request())
	form.UserAgent = web.GetUserAgent(ctx.Request())
	return *form
}

func GetCreateCommentForm(ctx iris.Context) CreateCommentForm {
	form := CreateCommentForm{
		EntityType: params.FormValue(ctx, "entityType"),
		EntityId:   common.GetID(ctx, "entityId"),
		Content:    strings.TrimSpace(params.FormValue(ctx, "content")),
		ImageList:  GetImageList(ctx, "imageList"),
		QuoteId:    params.FormValueInt64Default(ctx, "quoteId", 0),
		UserAgent:  web.GetUserAgent(ctx.Request()),
		Ip:         web.GetRequestIP(ctx.Request()),
	}
	return form
}

func GetCreateArticleForm(ctx iris.Context) CreateArticleForm {
	var (
		title   = ctx.PostValue("title")
		summary = ctx.PostValue("summary")
		content = ctx.PostValue("content")
		tags    = params.FormValueStringArray(ctx, "tags")
		cover   = GetImageDTO(ctx, "cover")
	)
	return CreateArticleForm{
		Title:       title,
		Summary:     summary,
		Content:     content,
		ContentType: constants.ContentTypeMarkdown,
		Cover:       cover,
		Tags:        tags,
	}
}

func GetImageList(ctx iris.Context, paramName string) []ImageDTO {
	imageListStr := params.FormValue(ctx, paramName)
	var imageList []ImageDTO
	if strs.IsNotBlank(imageListStr) {
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

func GetImageDTO(ctx iris.Context, paramName string) (img *ImageDTO) {
	str := params.FormValue(ctx, paramName)
	if strs.IsBlank(str) {
		return
	}
	if err := jsons.Parse(str, &img); err != nil {
		slog.Error(err.Error(), slog.Any("err", err))
	}
	return
}
