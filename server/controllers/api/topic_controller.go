package api

import (
	"bbs-go/model/constants"
	"bbs-go/package/common"
	"bbs-go/package/es"
	"math/rand"
	"strconv"
	"strings"

	"github.com/dchest/captcha"

	"github.com/kataras/iris/v12"
	"github.com/mlogclub/simple"

	"bbs-go/cache"
	"bbs-go/controllers/render"
	"bbs-go/model"
	"bbs-go/services"
)

type TopicController struct {
	Ctx iris.Context
}

func (c *TopicController) GetEs() *simple.JsonResult {
	go func() {
		services.TopicService.Scan(func(topics []model.Topic) {
			for _, topic := range topics {
				es.UpdateTopicIndex(&topic)
			}
		})
	}()
	return simple.JsonSuccess()
}

// 节点
func (c *TopicController) GetNodes() *simple.JsonResult {
	nodes := services.TopicNodeService.GetNodes()
	return simple.JsonData(render.BuildNodes(nodes))
}

// 节点信息
func (c *TopicController) GetNode() *simple.JsonResult {
	nodeId := simple.FormValueInt64Default(c.Ctx, "nodeId", 0)
	node := services.TopicNodeService.Get(nodeId)
	return simple.JsonData(render.BuildNode(node))
}

// 发表帖子
func (c *TopicController) PostCreate() *simple.JsonResult {
	user := services.UserTokenService.GetCurrent(c.Ctx)
	if err := services.UserService.CheckPostStatus(user); err != nil {
		return simple.JsonError(err)
	}

	form := model.GetCreateTopicForm(c.Ctx)

	if services.SysConfigService.GetConfig().TopicCaptcha && !captcha.VerifyString(form.CaptchaId, form.CaptchaCode) {
		return simple.JsonError(common.CaptchaError)
	}

	topic, err := services.TopicService.Publish(user.Id, form)
	if err != nil {
		return simple.JsonError(err)
	}
	return simple.JsonData(render.BuildSimpleTopic(topic))
}

// 编辑时获取详情
func (c *TopicController) GetEditBy(topicId int64) *simple.JsonResult {
	user := services.UserTokenService.GetCurrent(c.Ctx)
	if err := services.UserService.CheckPostStatus(user); err != nil {
		return simple.JsonError(err)
	}

	topic := services.TopicService.Get(topicId)
	if topic == nil || topic.Status != constants.StatusOk {
		return simple.JsonErrorMsg("话题不存在或已被删除")
	}
	if topic.Type != constants.TopicTypeTopic {
		return simple.JsonErrorMsg("当前类型帖子不支持修改")
	}

	// 非作者、且非管理员
	if topic.UserId != user.Id && !user.HasAnyRole(constants.RoleAdmin, constants.RoleOwner) {
		return simple.JsonErrorMsg("无权限")
	}

	tags := services.TopicService.GetTopicTags(topicId)
	var tagNames []string
	if len(tags) > 0 {
		for _, tag := range tags {
			tagNames = append(tagNames, tag.Name)
		}
	}

	return simple.NewEmptyRspBuilder().
		Put("topicId", topic.Id).
		Put("nodeId", topic.NodeId).
		Put("title", topic.Title).
		Put("content", topic.Content).
		Put("tags", tagNames).
		JsonResult()
}

// 编辑帖子
func (c *TopicController) PostEditBy(topicId int64) *simple.JsonResult {
	user := services.UserTokenService.GetCurrent(c.Ctx)
	if err := services.UserService.CheckPostStatus(user); err != nil {
		return simple.JsonError(err)
	}

	topic := services.TopicService.Get(topicId)
	if topic == nil || topic.Status != constants.StatusOk {
		return simple.JsonErrorMsg("话题不存在或已被删除")
	}

	// 非作者、且非管理员
	if topic.UserId != user.Id && !user.HasAnyRole(constants.RoleAdmin, constants.RoleOwner) {
		return simple.JsonErrorMsg("无权限")
	}

	nodeId := simple.FormValueInt64Default(c.Ctx, "nodeId", 0)
	title := strings.TrimSpace(simple.FormValue(c.Ctx, "title"))
	content := strings.TrimSpace(simple.FormValue(c.Ctx, "content"))
	tags := simple.FormValueStringArray(c.Ctx, "tags")

	err := services.TopicService.Edit(topicId, nodeId, tags, title, content)
	if err != nil {
		return simple.JsonError(err)
	}
	// 操作日志
	services.OperateLogService.AddOperateLog(user.Id, constants.OpTypeUpdate, constants.EntityTopic, topicId,
		"", c.Ctx.Request())
	return simple.JsonData(render.BuildSimpleTopic(topic))
}

// 删除帖子
func (c *TopicController) PostDeleteBy(topicId int64) *simple.JsonResult {
	user := services.UserTokenService.GetCurrent(c.Ctx)
	if err := services.UserService.CheckPostStatus(user); err != nil {
		return simple.JsonError(err)
	}

	topic := services.TopicService.Get(topicId)
	if topic == nil || topic.Status != constants.StatusOk {
		return simple.JsonSuccess()
	}

	// 非作者、且非管理员
	if topic.UserId != user.Id && !user.HasAnyRole(constants.RoleAdmin, constants.RoleOwner) {
		return simple.JsonErrorMsg("无权限")
	}

	if err := services.TopicService.Delete(topicId, user.Id, c.Ctx.Request()); err != nil {
		return simple.JsonErrorMsg(err.Error())
	}
	return simple.JsonSuccess()
}

// PostRecommendBy 设为推荐
func (c *TopicController) PostRecommendBy(topicId int64) *simple.JsonResult {
	recommend, err := simple.FormValueBool(c.Ctx, "recommend")
	if err != nil {
		return simple.JsonErrorMsg(err.Error())
	}
	user := services.UserTokenService.GetCurrent(c.Ctx)
	if user == nil {
		return simple.JsonError(simple.ErrorNotLogin)
	}
	if !user.HasAnyRole(constants.RoleOwner, constants.RoleAdmin) {
		return simple.JsonErrorMsg("无权限")
	}

	err = services.TopicService.SetRecommend(topicId, recommend)
	if err != nil {
		return simple.JsonErrorMsg(err.Error())
	}
	return simple.JsonSuccess()
}

// 帖子详情
func (c *TopicController) GetBy(topicId int64) *simple.JsonResult {
	topic := services.TopicService.Get(topicId)
	if topic == nil || topic.Status != constants.StatusOk {
		return simple.JsonErrorMsg("主题不存在")
	}
	services.TopicService.IncrViewCount(topicId) // 增加浏览量
	return simple.JsonData(render.BuildTopic(topic))
}

// 点赞
func (c *TopicController) PostLikeBy(topicId int64) *simple.JsonResult {
	user := services.UserTokenService.GetCurrent(c.Ctx)
	if user == nil {
		return simple.JsonError(simple.ErrorNotLogin)
	}
	err := services.UserLikeService.TopicLike(user.Id, topicId)
	if err != nil {
		return simple.JsonErrorMsg(err.Error())
	}
	return simple.JsonSuccess()
}

// 点赞用户
func (c *TopicController) GetRecentlikesBy(topicId int64) *simple.JsonResult {
	likes := services.UserLikeService.Recent(constants.EntityTopic, topicId, 5)
	var users []model.UserInfo
	for _, like := range likes {
		userInfo := render.BuildUserById(like.UserId)
		if userInfo != nil {
			users = append(users, *userInfo)
		}
	}
	return simple.JsonData(users)
}

// 最新帖子
func (c *TopicController) GetRecent() *simple.JsonResult {
	user := services.UserTokenService.GetCurrent(c.Ctx)
	topics := services.TopicService.Find(simple.NewSqlCnd().Where("status = ?", constants.StatusOk).Desc("id").Limit(10))
	return simple.JsonData(render.BuildSimpleTopics(topics, user))
}

// 用户帖子列表
func (c *TopicController) GetUserTopics() *simple.JsonResult {
	userId, err := simple.FormValueInt64(c.Ctx, "userId")
	if err != nil {
		return simple.JsonErrorMsg(err.Error())
	}
	cursor := simple.FormValueInt64Default(c.Ctx, "cursor", 0)
	user := services.UserTokenService.GetCurrent(c.Ctx)
	topics, cursor := services.TopicService.GetUserTopics(userId, cursor)
	return simple.JsonCursorData(render.BuildSimpleTopics(topics, user), strconv.FormatInt(cursor, 10))
}

// 帖子列表
func (c *TopicController) GetTopics() *simple.JsonResult {
	var (
		cursor       = simple.FormValueInt64Default(c.Ctx, "cursor", 0)
		nodeId       = simple.FormValueInt64Default(c.Ctx, "nodeId", 0)
		recommend, _ = simple.FormValueBool(c.Ctx, "recommend")
		user         = services.UserTokenService.GetCurrent(c.Ctx)
	)
	topics, cursor := services.TopicService.GetTopics(nodeId, cursor, recommend)
	return simple.JsonCursorData(render.BuildSimpleTopics(topics, user), strconv.FormatInt(cursor, 10))
}

// 标签帖子列表
func (c *TopicController) GetTagTopics() *simple.JsonResult {
	var (
		cursor     = simple.FormValueInt64Default(c.Ctx, "cursor", 0)
		tagId, err = simple.FormValueInt64(c.Ctx, "tagId")
		user       = services.UserTokenService.GetCurrent(c.Ctx)
	)
	if err != nil {
		return simple.JsonErrorMsg(err.Error())
	}
	topics, cursor := services.TopicService.GetTagTopics(tagId, cursor)
	return simple.JsonCursorData(render.BuildSimpleTopics(topics, user), strconv.FormatInt(cursor, 10))
}

// 收藏
func (c *TopicController) GetFavoriteBy(topicId int64) *simple.JsonResult {
	user := services.UserTokenService.GetCurrent(c.Ctx)
	if user == nil {
		return simple.JsonError(simple.ErrorNotLogin)
	}
	err := services.FavoriteService.AddTopicFavorite(user.Id, topicId)
	if err != nil {
		return simple.JsonErrorMsg(err.Error())
	}
	return simple.JsonSuccess()
}

// 推荐话题列表（目前逻辑为取最近50条数据随机展示）
func (c *TopicController) GetRecommend() *simple.JsonResult {
	topics := cache.TopicCache.GetRecommendTopics()
	if topics == nil || len(topics) == 0 {
		return simple.JsonSuccess()
	} else {
		dest := make([]model.Topic, len(topics))
		perm := rand.Perm(len(topics))
		for i, v := range perm {
			dest[v] = topics[i]
		}
		end := 10
		if end > len(topics) {
			end = len(topics)
		}
		ret := dest[0:end]
		return simple.JsonData(render.BuildSimpleTopics(ret, nil))
	}
}

// 最新话题
func (c *TopicController) GetNewest() *simple.JsonResult {
	topics := services.TopicService.Find(simple.NewSqlCnd().Eq("status", constants.StatusOk).Desc("id").Limit(6))
	return simple.JsonData(render.BuildSimpleTopics(topics, nil))
}
