package api

import (
	"bbs-go/internal/models/constants"
	"bbs-go/internal/pkg/common"
	"bbs-go/internal/pkg/errs"
	"bbs-go/internal/pkg/markdown"
	"bbs-go/internal/spam"
	"math/rand"
	"strconv"
	"strings"

	"github.com/kataras/iris/v12"
	"github.com/mlogclub/simple/common/strs"
	"github.com/mlogclub/simple/sqls"
	"github.com/mlogclub/simple/web"
	"github.com/mlogclub/simple/web/params"

	"bbs-go/internal/cache"
	"bbs-go/internal/controllers/render"
	"bbs-go/internal/models"
	"bbs-go/internal/services"
)

type TopicController struct {
	Ctx iris.Context
}

func (c *TopicController) GetNode_navs() *web.JsonResult {
	nodes := []models.NodeResponse{
		{
			Id:   0,
			Name: "最新",
		},
		{
			Id:   -1,
			Name: "推荐",
		},
		{
			Id:   -2,
			Name: "关注",
		},
	}
	realNodes := render.BuildNodes(services.TopicNodeService.GetNodes())
	nodes = append(nodes, realNodes...)
	return web.JsonData(nodes)
}

// 节点
func (c *TopicController) GetNodes() *web.JsonResult {
	nodes := render.BuildNodes(services.TopicNodeService.GetNodes())
	return web.JsonData(nodes)
}

// 节点信息
func (c *TopicController) GetNode() *web.JsonResult {
	nodeId := params.FormValueInt64Default(c.Ctx, "nodeId", 0)
	node := services.TopicNodeService.Get(nodeId)
	return web.JsonData(render.BuildNode(node))
}

// 发表帖子
func (c *TopicController) PostCreate() *web.JsonResult {
	user := services.UserTokenService.GetCurrent(c.Ctx)
	if err := services.UserService.CheckPostStatus(user); err != nil {
		return web.JsonError(err)
	}
	form := models.GetCreateTopicForm(c.Ctx)

	if err := spam.CheckTopic(user, form); err != nil {
		return web.JsonError(err)
	}

	topic, err := services.TopicPublishService.Publish(user.Id, form)
	if err != nil {
		return web.JsonError(err)
	}
	return web.JsonData(render.BuildSimpleTopic(topic))
}

// 编辑时获取详情
func (c *TopicController) GetEditBy(topicId int64) *web.JsonResult {
	user := services.UserTokenService.GetCurrent(c.Ctx)
	if err := services.UserService.CheckPostStatus(user); err != nil {
		return web.JsonError(err)
	}

	topic := services.TopicService.Get(topicId)
	if topic == nil || topic.Status != constants.StatusOk {
		return web.JsonErrorMsg("话题不存在或已被删除")
	}
	if topic.Type != constants.TopicTypeTopic {
		return web.JsonErrorMsg("当前类型帖子不支持修改")
	}

	// 非作者、且非管理员
	if topic.UserId != user.Id && !user.HasAnyRole(constants.RoleAdmin, constants.RoleOwner) {
		return web.JsonErrorMsg("无权限")
	}

	tags := services.TopicService.GetTopicTags(topicId)
	var tagNames []string
	if len(tags) > 0 {
		for _, tag := range tags {
			tagNames = append(tagNames, tag.Name)
		}
	}

	return web.NewEmptyRspBuilder().
		Put("id", topic.Id).
		Put("nodeId", topic.NodeId).
		Put("title", topic.Title).
		Put("content", topic.Content).
		Put("hideContent", topic.HideContent).
		Put("tags", tagNames).
		JsonResult()
}

// 编辑帖子
func (c *TopicController) PostEditBy(topicId int64) *web.JsonResult {
	user := services.UserTokenService.GetCurrent(c.Ctx)
	if err := services.UserService.CheckPostStatus(user); err != nil {
		return web.JsonError(err)
	}

	topic := services.TopicService.Get(topicId)
	if topic == nil || topic.Status != constants.StatusOk {
		return web.JsonErrorMsg("话题不存在或已被删除")
	}

	// 非作者、且非管理员
	if topic.UserId != user.Id && !user.HasAnyRole(constants.RoleAdmin, constants.RoleOwner) {
		return web.JsonErrorMsg("无权限")
	}

	var (
		nodeId      = params.FormValueInt64Default(c.Ctx, "nodeId", 0)
		title       = strings.TrimSpace(params.FormValue(c.Ctx, "title"))
		content     = strings.TrimSpace(params.FormValue(c.Ctx, "content"))
		hideContent = strings.TrimSpace(params.FormValue(c.Ctx, "hideContent"))
		tags        = params.FormValueStringArray(c.Ctx, "tags")
	)

	err := services.TopicService.Edit(topicId, nodeId, tags, title, content, hideContent)
	if err != nil {
		return web.JsonError(err)
	}
	// 操作日志
	services.OperateLogService.AddOperateLog(user.Id, constants.OpTypeUpdate, constants.EntityTopic, topicId,
		"", c.Ctx.Request())
	return web.JsonData(render.BuildSimpleTopic(topic))
}

// 删除帖子
func (c *TopicController) PostDeleteBy(topicId int64) *web.JsonResult {
	user := services.UserTokenService.GetCurrent(c.Ctx)
	if err := services.UserService.CheckPostStatus(user); err != nil {
		return web.JsonError(err)
	}

	topic := services.TopicService.Get(topicId)
	if topic == nil || topic.Status != constants.StatusOk {
		return web.JsonSuccess()
	}

	// 非作者、且非管理员
	if topic.UserId != user.Id && !user.HasAnyRole(constants.RoleAdmin, constants.RoleOwner) {
		return web.JsonErrorMsg("无权限")
	}

	if err := services.TopicService.Delete(topicId, user.Id, c.Ctx.Request()); err != nil {
		return web.JsonError(err)
	}
	return web.JsonSuccess()
}

// PostRecommendBy 设为推荐
func (c *TopicController) PostRecommendBy(topicId int64) *web.JsonResult {
	recommend, err := params.FormValueBool(c.Ctx, "recommend")
	if err != nil {
		return web.JsonError(err)
	}
	user := services.UserTokenService.GetCurrent(c.Ctx)
	if user == nil {
		return web.JsonError(errs.NotLogin)
	}
	if !user.HasAnyRole(constants.RoleOwner, constants.RoleAdmin) {
		return web.JsonErrorMsg("无权限")
	}

	err = services.TopicService.SetRecommend(topicId, recommend)
	if err != nil {
		return web.JsonError(err)
	}
	return web.JsonSuccess()
}

// 帖子详情
func (c *TopicController) GetBy(topicId int64) *web.JsonResult {
	topic := services.TopicService.Get(topicId)
	if topic == nil || topic.Status == constants.StatusDeleted {
		return web.JsonErrorMsg("帖子不存在")
	}

	// 审核中文章控制展示
	user := services.UserTokenService.GetCurrent(c.Ctx)
	if topic.Status == constants.StatusReview {
		if user != nil {
			if topic.UserId != user.Id && !user.IsOwnerOrAdmin() {
				return web.JsonErrorCode(403, "文章审核中")
			}
		} else {
			return web.JsonErrorCode(403, "文章审核中")
		}
	}

	services.TopicService.IncrViewCount(topicId) // 增加浏览量
	return web.JsonData(render.BuildTopic(topic, user))
}

// 点赞用户
func (c *TopicController) GetRecentlikesBy(topicId int64) *web.JsonResult {
	likes := services.UserLikeService.Recent(constants.EntityTopic, topicId, 5)
	var users []models.UserInfo
	for _, like := range likes {
		userInfo := render.BuildUserInfoDefaultIfNull(like.UserId)
		if userInfo != nil {
			users = append(users, *userInfo)
		}
	}
	return web.JsonData(users)
}

// 最新帖子
func (c *TopicController) GetRecent() *web.JsonResult {
	user := services.UserTokenService.GetCurrent(c.Ctx)
	topics := services.TopicService.Find(sqls.NewCnd().Where("status = ?", constants.StatusOk).Desc("id").Limit(10))
	return web.JsonData(render.BuildSimpleTopics(topics, user))
}

// 用户帖子列表
func (c *TopicController) GetUserTopics() *web.JsonResult {
	userId, err := params.FormValueInt64(c.Ctx, "userId")
	if err != nil {
		return web.JsonError(err)
	}
	cursor := params.FormValueInt64Default(c.Ctx, "cursor", 0)
	user := services.UserTokenService.GetCurrent(c.Ctx)
	topics, cursor, hasMore := services.TopicService.GetUserTopics(userId, cursor)
	return web.JsonCursorData(render.BuildSimpleTopics(topics, user), strconv.FormatInt(cursor, 10), hasMore)
}

// 帖子列表
func (c *TopicController) GetTopics() *web.JsonResult {
	var (
		cursor = params.FormValueInt64Default(c.Ctx, "cursor", 0)
		nodeId = params.FormValueInt64Default(c.Ctx, "nodeId", 0)
		user   = services.UserTokenService.GetCurrent(c.Ctx)
	)
	if nodeId == constants.NodeIdFollow && user == nil {
		return web.JsonError(errs.NotLogin)
	}

	var temp []models.Topic
	if cursor <= 0 {
		stickyTopics := services.TopicService.GetStickyTopics(nodeId, 3)
		temp = append(temp, stickyTopics...)
	}
	topics, cursor, hasMore := services.TopicService.GetTopics(user, nodeId, cursor)
	for _, topic := range topics {
		topic.Sticky = false // 正常列表不要渲染置顶
		temp = append(temp, topic)
	}
	list := common.Distinct(temp, func(t models.Topic) any {
		return t.Id
	})
	return web.JsonCursorData(render.BuildSimpleTopics(list, user), strconv.FormatInt(cursor, 10), hasMore)
}

// 标签帖子列表
func (c *TopicController) GetTagTopics() *web.JsonResult {
	var (
		cursor     = params.FormValueInt64Default(c.Ctx, "cursor", 0)
		tagId, err = params.FormValueInt64(c.Ctx, "tagId")
		user       = services.UserTokenService.GetCurrent(c.Ctx)
	)
	if err != nil {
		return web.JsonError(err)
	}
	topics, cursor, hasMore := services.TopicService.GetTagTopics(tagId, cursor)
	return web.JsonCursorData(render.BuildSimpleTopics(topics, user), strconv.FormatInt(cursor, 10), hasMore)
}

// 收藏
func (c *TopicController) GetFavoriteBy(topicId int64) *web.JsonResult {
	user := services.UserTokenService.GetCurrent(c.Ctx)
	if user == nil {
		return web.JsonError(errs.NotLogin)
	}
	err := services.FavoriteService.AddTopicFavorite(user.Id, topicId)
	if err != nil {
		return web.JsonError(err)
	}
	return web.JsonSuccess()
}

// 推荐话题列表（目前逻辑为取最近50条数据随机展示）
func (c *TopicController) GetRecommend() *web.JsonResult {
	topics := cache.TopicCache.GetRecommendTopics()
	if len(topics) == 0 {
		return web.JsonSuccess()
	} else {
		dest := make([]models.Topic, len(topics))
		perm := rand.Perm(len(topics))
		for i, v := range perm {
			dest[v] = topics[i]
		}
		end := 10
		if end > len(topics) {
			end = len(topics)
		}
		ret := dest[0:end]
		return web.JsonData(render.BuildSimpleTopics(ret, nil))
	}
}

// 最新话题
func (c *TopicController) GetNewest() *web.JsonResult {
	topics := services.TopicService.Find(sqls.NewCnd().Eq("status", constants.StatusOk).Desc("id").Limit(6))
	return web.JsonData(render.BuildSimpleTopics(topics, nil))
}

// 设置置顶
func (c *TopicController) PostStickyBy(topicId int64) *web.JsonResult {
	user := services.UserTokenService.GetCurrent(c.Ctx)
	if user == nil {
		return web.JsonError(errs.NotLogin)
	}
	if !user.HasAnyRole(constants.RoleOwner, constants.RoleAdmin) {
		return web.JsonErrorMsg("无权限")
	}

	var (
		sticky = params.FormValueBoolDefault(c.Ctx, "sticky", false) // 是否指定
	)
	if err := services.TopicService.SetSticky(topicId, sticky); err != nil {
		return web.JsonError(err)
	}
	return web.JsonSuccess()
}

func (c *TopicController) GetHide_content() *web.JsonResult {
	topicId := params.FormValueInt64Default(c.Ctx, "topicId", 0)
	var (
		exists      = false // 是否有隐藏内容
		show        = false // 是否显示隐藏内容
		hideContent = ""    // 隐藏内容
	)
	topic := services.TopicService.Get(topicId)
	if topic != nil && topic.Status == constants.StatusOk && strs.IsNotBlank(topic.HideContent) {
		exists = true
		if user := services.UserTokenService.GetCurrent(c.Ctx); user != nil {
			if user.Id == topic.UserId || services.CommentService.IsCommented(user.Id, constants.EntityTopic, topic.Id) {
				show = true
				hideContent = markdown.ToHTML(topic.HideContent)
			}
		}
	}
	return web.JsonData(map[string]interface{}{
		"exists":  exists,
		"show":    show,
		"content": hideContent,
	})
}
