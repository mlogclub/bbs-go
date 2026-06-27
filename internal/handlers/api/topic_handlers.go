package api

import (
	"bbs-go/internal/models/constants"
	"bbs-go/internal/models/req"
	"bbs-go/internal/models/resp"
	"bbs-go/internal/permissions"
	"bbs-go/internal/pkg/common"
	"bbs-go/internal/pkg/errs"
	"bbs-go/internal/pkg/idcodec"
	"bbs-go/internal/pkg/locales"
	"bbs-go/internal/pkg/markdown"
	"bbs-go/internal/spam"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"

	"bbs-go/internal/pkg/ginx"
	"bbs-go/internal/pkg/params"

	"github.com/mlogclub/simple/common/strs"
	"github.com/mlogclub/simple/sqls"
	"github.com/mlogclub/simple/web"

	"bbs-go/internal/handlers/render"
	"bbs-go/internal/models"
	"bbs-go/internal/services"
)

func topicGetBuiltInCategories() []resp.CategoryResponse {
	return []resp.CategoryResponse{
		{
			Id:   0,
			Name: locales.Get("topic.category.latest"),
			Logo: "/res/images/category_latest.png",
		},
		{
			Id:   -1,
			Name: locales.Get("topic.category.recommend"),
			Logo: "/res/images/category_recommend.png",
		},
		{
			Id:   -2,
			Name: locales.Get("topic.category.follow"),
			Logo: "/res/images/category_follow.png",
		},
	}
}

// 节点
// 节点信息
// 发表帖子
// 编辑时获取详情
// 编辑帖子
// 删除帖子
// PostRecommendBy 设为推荐
// 帖子详情
// 点赞用户
// 最新帖子
// 用户帖子列表
// 帖子列表
// 采纳答案
// 取消采纳答案
// 标签帖子列表
// 收藏
// 设置置顶
func CategoryNavs(ctx *gin.Context) {

	categories := render.BuildCategoryResponses(services.CategoryService.GetTopLevelCategories())
	ginx.WriteJSON(ctx, categories)

}

func Categories(ctx *gin.Context) {
	topicType := constants.TopicType(params.FormValueIntDefault(ctx, "type", -1))
	var categoryList []models.Category
	if topicType >= 0 {
		categoryList = services.CategoryService.GetCategoriesByTopicType(topicType)
	} else {
		categoryList = services.CategoryService.GetCategories()
	}
	categories := render.BuildCategoryResponseTree(0, categoryList)
	ginx.WriteJSON(ctx, categories)

}

func Category(ctx *gin.Context) {
	categoryId, _ := params.GetInt64(ctx, "categoryId")
	if categoryId <= 0 {
		for _, category := range topicGetBuiltInCategories() {
			if category.Id == categoryId {
				ginx.WriteJSON(ctx, category)
				return
			}
		}
	}
	category := services.CategoryService.Get(categoryId)
	if category == nil {
		ginx.WriteJSON(ctx, ginx.ErrorMessage(locales.Get("common.not_found")))
		return
	}
	ginx.WriteJSON(ctx, render.BuildCategoryWithChildren(category))

}

func TopicCreate(ctx *gin.Context) {
	user := common.GetCurrentUser(ctx)
	if err := services.UserService.CheckPostStatus(user); err != nil {
		ginx.WriteJSON(ctx, err)
		return
	}

	var form req.CreateTopicReq
	if err := ginx.BindJSON(ctx, &form); err != nil {
		ginx.WriteJSON(ctx, err)
		return
	}
	form.Title = strings.TrimSpace(form.Title)
	form.Content = strings.TrimSpace(form.Content)
	form.HideContent = strings.TrimSpace(form.HideContent)
	if constants.IsTweetTopicType(form.Type) {
		form.ContentType = constants.ContentTypeText
	}
	form.Ip = web.GetRequestIP(ctx.Request)
	form.UserAgent = web.GetUserAgent(ctx.Request)

	if err := spam.CheckTopic(user, form); err != nil {
		ginx.WriteJSON(ctx, err)
		return
	}

	topic, err := services.TopicPublishService.Publish(user.Id, form)
	if err != nil {
		ginx.WriteJSON(ctx, err)
		return
	}
	ginx.WriteJSON(ctx, render.BuildSimpleTopic(topic))

}

func TopicEditForm(ctx *gin.Context) {
	topicIdStr := ctx.Param("id")

	topicId := idcodec.Decode(topicIdStr)
	user := common.GetCurrentUser(ctx)
	if err := services.UserService.CheckPostStatus(user); err != nil {
		ginx.WriteJSON(ctx, err)
		return
	}

	topic := services.TopicService.Get(topicId)
	if topic == nil || topic.Status != constants.StatusOk {
		ginx.WriteJSON(ctx, ginx.ErrorMessage(locales.Get("common.not_found")))
		return
	}
	if constants.IsTweetTopicType(topic.Type) {
		ginx.WriteJSON(ctx, ginx.ErrorMessage(locales.Get("topic.type_not_supported")))
		return
	}

	// 非作者、且非站长
	if topic.UserId != user.Id && !user.HasRole(constants.RoleOwner) {
		ginx.WriteJSON(ctx, ginx.ErrorMessage(locales.Get("topic.no_permission")))
		return
	}

	tags := services.TopicService.GetTopicTags(topicId)
	var tagNames []string
	if len(tags) > 0 {
		for _, tag := range tags {
			tagNames = append(tagNames, tag.Name)
		}
	}

	attachments := render.BuildAttachmentResponses(services.AttachmentService.ListByTopicId(topicId), nil)

	ginx.WriteJSON(ctx, map[string]any{
		"id":          idcodec.Encode(topic.Id),
		"type":        topic.Type,
		"categoryId":  topic.CategoryId,
		"title":       topic.Title,
		"content":     topic.Content,
		"contentType": topic.ContentType,
		"hideContent": topic.HideContent,
		"tags":        tagNames,
		"attachments": attachments,
	})

}

func TopicEdit(ctx *gin.Context) {
	topicIdStr := ctx.Param("id")

	topicId := idcodec.Decode(topicIdStr)
	user := common.GetCurrentUser(ctx)
	if err := services.UserService.CheckPostStatus(user); err != nil {
		ginx.WriteJSON(ctx, err)
		return
	}

	topic := services.TopicService.Get(topicId)
	if topic == nil || topic.Status != constants.StatusOk {
		ginx.WriteJSON(ctx, ginx.ErrorMessage(locales.Get("common.not_found")))
		return
	}

	// 非作者、且非站长
	if topic.UserId != user.Id && !user.HasRole(constants.RoleOwner) {
		ginx.WriteJSON(ctx, ginx.ErrorMessage(locales.Get("topic.no_permission")))
		return
	}

	var form req.EditTopicReq
	if err := ginx.BindJSON(ctx, &form); err != nil {
		ginx.WriteJSON(ctx, err)
		return
	}
	form.Title = strings.TrimSpace(form.Title)
	form.Content = strings.TrimSpace(form.Content)
	form.HideContent = strings.TrimSpace(form.HideContent)

	err := services.TopicService.Edit(user.Id, topicId, form)
	if err != nil {
		ginx.WriteJSON(ctx, err)
		return
	}
	ginx.WriteJSON(ctx, render.BuildSimpleTopic(topic))

}

func TopicRemove(ctx *gin.Context) {
	topicIdStr := ctx.Param("id")

	topicId := idcodec.Decode(topicIdStr)
	user := common.GetCurrentUser(ctx)
	if err := services.UserService.CheckPostStatus(user); err != nil {
		ginx.WriteJSON(ctx, err)
		return
	}

	topic := services.TopicService.Get(topicId)
	if topic == nil || topic.Status != constants.StatusOk {
		ginx.WriteJSON(ctx, nil)
		return
	}

	// 非作者、且非站长
	if topic.UserId != user.Id && !user.HasRole(constants.RoleOwner) {
		ginx.WriteJSON(ctx, ginx.ErrorMessage(locales.Get("topic.no_permission")))
		return
	}

	if err := services.TopicService.Delete(topicId, user.Id, ctx.Request); err != nil {
		ginx.WriteJSON(ctx, err)
		return
	}
	ginx.WriteJSON(ctx, nil)

}

func TopicRecommend(ctx *gin.Context) {
	topicIdStr := ctx.Param("id")

	topicId := idcodec.Decode(topicIdStr)
	recommend, err := params.FormValueBool(ctx, "recommend")
	if err != nil {
		ginx.WriteJSON(ctx, err)
		return
	}
	user := common.GetCurrentUser(ctx)
	if user == nil {
		ginx.WriteJSON(ctx, errs.NotLogin())
		return
	}
	if !services.PermissionService.HasPermission(user, permissions.PermissionTopicRecommend.Code) {
		ginx.WriteJSON(ctx, ginx.ErrorMessage(locales.Get("topic.no_permission")))
		return
	}

	err = services.TopicService.SetRecommend(topicId, recommend)
	if err != nil {
		ginx.WriteJSON(ctx, err)
		return
	}
	ginx.WriteJSON(ctx, nil)

}

func TopicDetail(ctx *gin.Context) {
	topicIdStr := ctx.Param("id")

	topicId := idcodec.Decode(topicIdStr)
	topic := services.TopicService.Get(topicId)
	if topic == nil || topic.Status == constants.StatusDeleted {
		ginx.WriteJSON(ctx, ginx.ErrorMessage(locales.Get("common.not_found")))
		return
	}

	// 审核中文章控制展示
	user := common.GetCurrentUser(ctx)
	if topic.Status == constants.StatusReview {
		if user != nil {
			if topic.UserId != user.Id && !user.IsOwner() {
				ginx.WriteJSON(ctx, ginx.ErrorCode(403, locales.Get("topic.under_review")))
				return
			}
		} else {
			ginx.WriteJSON(ctx, ginx.ErrorCode(403, locales.Get("topic.under_review")))
			return
		}
	}

	services.TopicService.IncrViewCount(topicId) // 增加浏览量
	ginx.WriteJSON(ctx, render.BuildTopic(ctx, topic))

}

func TopicRecentlikes(ctx *gin.Context) {
	topicIdStr := ctx.Param("id")

	topicId := idcodec.Decode(topicIdStr)
	likes := services.UserLikeService.Recent(constants.EntityTopic, topicId, 5)
	var users []resp.UserInfo
	for _, like := range likes {
		userInfo := render.BuildUserInfoDefaultIfNull(like.UserId)
		if userInfo != nil {
			users = append(users, *userInfo)
		}
	}
	ginx.WriteJSON(ctx, users)

}

func TopicRecent(ctx *gin.Context) {
	topics := services.TopicService.Find(sqls.NewCnd().Where("status = ?", constants.StatusOk).Desc("id").Limit(10))
	ginx.WriteJSON(ctx, render.BuildSimpleTopics(ctx, topics))

}

func TopicUserTopics(ctx *gin.Context) {
	userId := common.GetID(ctx, "userId")
	if userId <= 0 {
		ginx.WriteJSON(ctx, ginx.ErrorMessage("param: userId required"))
		return
	}
	cursor := params.FormValueInt64Default(ctx, "cursor", 0)
	topics, cursor, hasMore := services.TopicService.GetUserTopics(userId, cursor)
	ginx.WriteJSON(ctx, ginx.CursorData(render.BuildSimpleTopics(ctx, topics), strconv.FormatInt(cursor, 10), hasMore))

}

func TopicTopics(ctx *gin.Context) {
	var (
		cursor     = params.FormValueInt64Default(ctx, "cursor", 0)
		categoryId = params.FormValueInt64Default(ctx, "categoryId", 0)
		qaStatus   = strings.TrimSpace(params.FormValue(ctx, "qaStatus"))
		sort       = strings.TrimSpace(params.FormValue(ctx, "sort"))
		user       = common.GetCurrentUser(ctx)
	)
	if categoryId == constants.CategoryIdFollow && user == nil {
		ginx.WriteJSON(ctx, errs.NotLogin())
		return
	}

	var temp []models.Topic
	if cursor <= 0 {
		stickyTopics := services.TopicService.GetStickyTopics(categoryId, 3, qaStatus)
		temp = append(temp, stickyTopics...)
	}
	topics, cursor, hasMore := services.TopicService.GetTopics(user, categoryId, cursor, qaStatus, sort)
	for _, topic := range topics {
		topic.Sticky = false // 正常列表不要渲染置顶
		temp = append(temp, topic)
	}
	list := common.Distinct(temp, func(t models.Topic) any {
		return t.Id
	})
	ginx.WriteJSON(ctx, ginx.CursorData(render.BuildSimpleTopics(ctx, list), strconv.FormatInt(cursor, 10), hasMore))

}

func TopicAcceptAnswer(ctx *gin.Context) {
	topicIdStr := ctx.Param("id")

	topicId := idcodec.Decode(topicIdStr)
	commentId := params.FormValueInt64Default(ctx, "commentId", 0)
	if commentId <= 0 {
		ginx.WriteJSON(ctx, ginx.ErrorMessage("commentId is required"))
		return
	}
	user := common.GetCurrentUser(ctx)
	if user == nil {
		ginx.WriteJSON(ctx, errs.NotLogin())
		return
	}
	if err := services.TopicService.AcceptAnswer(topicId, commentId, user.Id, user.IsOwner()); err != nil {
		ginx.WriteJSON(ctx, err)
		return
	}
	ginx.WriteJSON(ctx, nil)

}

func TopicUnacceptAnswer(ctx *gin.Context) {
	topicIdStr := ctx.Param("id")

	topicId := idcodec.Decode(topicIdStr)
	user := common.GetCurrentUser(ctx)
	if user == nil {
		ginx.WriteJSON(ctx, errs.NotLogin())
		return
	}
	if err := services.TopicService.UnacceptAnswer(topicId, user.Id, user.IsOwner()); err != nil {
		ginx.WriteJSON(ctx, err)
		return
	}
	ginx.WriteJSON(ctx, nil)

}

func TopicTagTopics(ctx *gin.Context) {
	var (
		cursor     = params.FormValueInt64Default(ctx, "cursor", 0)
		tagId, err = params.FormValueInt64(ctx, "tagId")
	)
	if err != nil {
		ginx.WriteJSON(ctx, err)
		return
	}
	topics, cursor, hasMore := services.TopicService.GetTagTopics(tagId, cursor)
	ginx.WriteJSON(ctx, ginx.CursorData(render.BuildSimpleTopics(ctx, topics), strconv.FormatInt(cursor, 10), hasMore))

}

func TopicFavorite(ctx *gin.Context) {
	topicIdStr := ctx.Param("id")

	topicId := idcodec.Decode(topicIdStr)
	user := common.GetCurrentUser(ctx)
	if user == nil {
		ginx.WriteJSON(ctx, errs.NotLogin())
		return
	}
	err := services.FavoriteService.AddTopicFavorite(user.Id, topicId)
	if err != nil {
		ginx.WriteJSON(ctx, err)
		return
	}
	ginx.WriteJSON(ctx, nil)

}

func TopicSticky(ctx *gin.Context) {
	topicIdStr := ctx.Param("id")

	topicId := idcodec.Decode(topicIdStr)
	user := common.GetCurrentUser(ctx)
	if user == nil {
		ginx.WriteJSON(ctx, errs.NotLogin())
		return
	}
	if !services.PermissionService.HasPermission(user, permissions.PermissionTopicSticky.Code) {
		ginx.WriteJSON(ctx, ginx.ErrorMessage(locales.Get("topic.no_permission")))
		return
	}

	var (
		sticky = params.FormValueBoolDefault(ctx, "sticky", false) // 是否指定
	)
	if err := services.TopicService.SetSticky(topicId, sticky); err != nil {
		ginx.WriteJSON(ctx, err)
		return
	}
	ginx.WriteJSON(ctx, nil)

}

func TopicHideContent(ctx *gin.Context) {
	topicId := common.GetID(ctx, "topicId")
	var (
		exists      = false // 是否有隐藏内容
		show        = false // 是否显示隐藏内容
		hideContent = ""    // 隐藏内容
	)
	topic := services.TopicService.Get(topicId)
	if topic != nil && topic.Status == constants.StatusOk && strs.IsNotBlank(topic.HideContent) {
		exists = true
		if user := common.GetCurrentUser(ctx); user != nil {
			if user.Id == topic.UserId || services.CommentService.IsCommented(user.Id, constants.EntityTopic, topic.Id) {
				show = true
				hideContent = markdown.ToHTML(topic.HideContent)
			}
		}
	}
	ginx.WriteJSON(ctx, map[string]interface{}{
		"exists":  exists,
		"show":    show,
		"content": hideContent,
	})

}
