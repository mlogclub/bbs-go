package services

import (
	"bbs-go/model/constants"
	"bbs-go/package/config"
	"bbs-go/package/es"
	"bbs-go/package/seo"
	"bbs-go/package/urls"
	"errors"
	"github.com/mlogclub/simple/date"
	"github.com/mlogclub/simple/json"
	"math"
	"net/http"
	"path"
	"time"

	"github.com/gorilla/feeds"
	"github.com/mlogclub/simple"
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"

	"bbs-go/cache"
	"bbs-go/model"
	"bbs-go/package/common"
	"bbs-go/repositories"
)

var TopicService = newTopicService()

func newTopicService() *topicService {
	return &topicService{}
}

type topicService struct{}

func (s *topicService) Get(id int64) *model.Topic {
	return repositories.TopicRepository.Get(simple.DB(), id)
}

func (s *topicService) Take(where ...interface{}) *model.Topic {
	return repositories.TopicRepository.Take(simple.DB(), where...)
}

func (s *topicService) Find(cnd *simple.SqlCnd) []model.Topic {
	return repositories.TopicRepository.Find(simple.DB(), cnd)
}

func (s *topicService) FindOne(cnd *simple.SqlCnd) *model.Topic {
	return repositories.TopicRepository.FindOne(simple.DB(), cnd)
}

func (s *topicService) FindPageByParams(params *simple.QueryParams) (list []model.Topic, paging *simple.Paging) {
	return repositories.TopicRepository.FindPageByParams(simple.DB(), params)
}

func (s *topicService) FindPageByCnd(cnd *simple.SqlCnd) (list []model.Topic, paging *simple.Paging) {
	return repositories.TopicRepository.FindPageByCnd(simple.DB(), cnd)
}

func (s *topicService) Count(cnd *simple.SqlCnd) int64 {
	return repositories.TopicRepository.Count(simple.DB(), cnd)
}

func (s *topicService) Updates(id int64, columns map[string]interface{}) error {
	if err := repositories.TopicRepository.Updates(simple.DB(), id, columns); err != nil {
		return err
	}

	// 添加索引
	es.UpdateTopicIndex(s.Get(id))

	return nil
}

func (s *topicService) UpdateColumn(id int64, name string, value interface{}) error {
	if err := repositories.TopicRepository.UpdateColumn(simple.DB(), id, name, value); err != nil {
		return err
	}

	// 添加索引
	es.UpdateTopicIndex(s.Get(id))

	return nil
}

// 删除
func (s *topicService) Delete(topicId, deleteUserId int64, r *http.Request) error {
	err := repositories.TopicRepository.UpdateColumn(simple.DB(), topicId, "status", constants.StatusDeleted)
	if err == nil {
		// 添加索引
		es.UpdateTopicIndex(s.Get(topicId))
		// 删掉标签文章
		TopicTagService.DeleteByTopicId(topicId)
		// 发送消息
		MessageService.SendTopicDeleteMsg(topicId, deleteUserId)
		// 操作日志
		OperateLogService.AddOperateLog(deleteUserId, constants.OpTypeDelete, constants.EntityTopic, topicId, "", r)
	}
	return err
}

// 取消删除
func (s *topicService) Undelete(id int64) error {
	err := repositories.TopicRepository.UpdateColumn(simple.DB(), id, "status", constants.StatusOk)
	if err == nil {
		// 删掉标签文章
		TopicTagService.UndeleteByTopicId(id)
		// 添加索引
		es.UpdateTopicIndex(s.Get(id))
	}
	return err
}

// 发表
func (s *topicService) Publish(userId int64, form model.CreateTopicForm) (*model.Topic, *simple.CodeError) {
	if form.Type == constants.TopicTypeTweet {
		if simple.IsBlank(form.Content) && len(form.ImageList) == 0 {
			return nil, simple.NewErrorMsg("内容或图片不能为空")
		}
	} else {
		if simple.IsBlank(form.Title) {
			return nil, simple.NewErrorMsg("标题不能为空")
		}

		if simple.IsBlank(form.Content) {
			return nil, simple.NewErrorMsg("内容不能为空")
		}

		if simple.RuneLen(form.Title) > 128 {
			return nil, simple.NewErrorMsg("标题长度不能超过128")
		}
	}

	if form.NodeId <= 0 {
		form.NodeId = SysConfigService.GetConfig().DefaultNodeId
		if form.NodeId <= 0 {
			return nil, simple.NewErrorMsg("请选择节点")
		}
	}
	node := repositories.TopicNodeRepository.Get(simple.DB(), form.NodeId)
	if node == nil || node.Status != constants.StatusOk {
		return nil, simple.NewErrorMsg("节点不存在")
	}

	now := date.NowTimestamp()
	topic := &model.Topic{
		Type:            form.Type,
		UserId:          userId,
		NodeId:          form.NodeId,
		Title:           form.Title,
		Content:         form.Content,
		Status:          constants.StatusOk,
		LastCommentTime: now,
		CreateTime:      now,
	}

	if len(form.ImageList) > 0 {
		imageListStr, err := json.ToStr(form.ImageList)
		if err == nil {
			topic.ImageList = imageListStr
		} else {
			logrus.Error(err)
		}
	}

	err := simple.DB().Transaction(func(tx *gorm.DB) error {
		tagIds := repositories.TagRepository.GetOrCreates(tx, form.Tags)
		err := repositories.TopicRepository.Create(tx, topic)
		if err != nil {
			return err
		}
		repositories.TopicTagRepository.AddTopicTags(tx, topic.Id, tagIds)
		return nil
	})
	if err == nil {
		// 添加索引
		es.UpdateTopicIndex(topic)
		// 用户话题计数
		UserService.IncrTopicCount(userId)
		// 获得积分
		UserService.IncrScoreForPostTopic(topic)
		// 百度链接推送
		seo.Push(urls.TopicUrl(topic.Id))
	}
	return topic, simple.FromError(err)
}

// 更新
func (s *topicService) Edit(topicId, nodeId int64, tags []string, title, content string) *simple.CodeError {
	if len(title) == 0 {
		return simple.NewErrorMsg("标题不能为空")
	}

	if simple.RuneLen(title) > 128 {
		return simple.NewErrorMsg("标题长度不能超过128")
	}

	node := repositories.TopicNodeRepository.Get(simple.DB(), nodeId)
	if node == nil || node.Status != constants.StatusOk {
		return simple.NewErrorMsg("节点不存在")
	}

	err := simple.DB().Transaction(func(tx *gorm.DB) error {
		err := repositories.TopicRepository.Updates(simple.DB(), topicId, map[string]interface{}{
			"node_id": nodeId,
			"title":   title,
			"content": content,
		})
		if err != nil {
			return err
		}

		tagIds := repositories.TagRepository.GetOrCreates(tx, tags)       // 创建帖子对应标签
		repositories.TopicTagRepository.DeleteTopicTags(tx, topicId)      // 先删掉所有的标签
		repositories.TopicTagRepository.AddTopicTags(tx, topicId, tagIds) // 然后重新添加标签
		return nil
	})
	return simple.FromError(err)
}

// 推荐
func (s *topicService) SetRecommend(topicId int64, recommend bool) error {
	topic := s.Get(topicId)
	if topic == nil || topic.Status != constants.StatusOk {
		return errors.New("帖子不存在")
	}
	if topic.Recommend == recommend { // 推荐状态没变更
		return nil
	}
	if recommend {
		if err := s.Updates(topicId, map[string]interface{}{
			"recommend":      recommend,
			"recommend_time": date.NowTimestamp(),
		}); err != nil {
			return err
		}
		MessageService.SendTopicRecommendMsg(topicId)
	} else {
		if err := s.UpdateColumn(topicId, "recommend", recommend); err != nil {
			return err
		}
	}
	return nil
}

// 话题的标签
func (s *topicService) GetTopicTags(topicId int64) []model.Tag {
	topicTags := repositories.TopicTagRepository.Find(simple.DB(), simple.NewSqlCnd().Where("topic_id = ?", topicId))

	var tagIds []int64
	for _, topicTag := range topicTags {
		tagIds = append(tagIds, topicTag.TagId)
	}
	return cache.TagCache.GetList(tagIds)
}

// 获取帖子分页列表
func (s *topicService) GetTopics(nodeId, cursor int64, recommend bool) (topics []model.Topic, nextCursor int64) {
	cnd := simple.NewSqlCnd()
	if nodeId > 0 {
		cnd.Eq("node_id", nodeId)
	}
	if cursor > 0 {
		cnd.Lt("last_comment_time", cursor)
	}
	if recommend {
		cnd.Eq("recommend", true)
	}
	cnd.Eq("status", constants.StatusOk).Desc("last_comment_time").Limit(20)
	topics = repositories.TopicRepository.Find(simple.DB(), cnd)
	if len(topics) > 0 {
		nextCursor = topics[len(topics)-1].LastCommentTime
	} else {
		nextCursor = cursor
	}
	return
}

// 指定标签下话题列表
func (s *topicService) GetTagTopics(tagId, cursor int64) (topics []model.Topic, nextCursor int64) {
	topicTags := repositories.TopicTagRepository.Find(simple.DB(), simple.NewSqlCnd().
		Eq("tag_id", tagId).
		Eq("status", constants.StatusOk).
		Desc("last_comment_time").Limit(20))
	if len(topicTags) > 0 {
		nextCursor = topicTags[len(topicTags)-1].LastCommentTime

		var topicIds []int64
		for _, topicTag := range topicTags {
			topicIds = append(topicIds, topicTag.TopicId)
		}

		topicsMap := s.GetTopicInIds(topicIds)
		if topicsMap != nil {
			for _, topicTag := range topicTags {
				if topic, found := topicsMap[topicTag.TopicId]; found {
					topics = append(topics, topic)
				}
			}
		}
	} else {
		nextCursor = cursor
	}
	return
}

// GetTopicInIds 根据编号批量获取主题
func (s *topicService) GetTopicInIds(topicIds []int64) map[int64]model.Topic {
	if len(topicIds) == 0 {
		return nil
	}
	var topics []model.Topic
	simple.DB().Where("id in (?)", topicIds).Find(&topics)

	topicsMap := make(map[int64]model.Topic, len(topics))
	for _, topic := range topics {
		topicsMap[topic.Id] = topic
	}
	return topicsMap
}

// 浏览数+1
func (s *topicService) IncrViewCount(topicId int64) {
	simple.DB().Exec("update t_topic set view_count = view_count + 1 where id = ?", topicId)
}

// 当帖子被评论的时候，更新最后回复时间、回复数量+1
func (s *topicService) OnComment(topicId int64, comment *model.Comment) {
	_ = simple.DB().Transaction(func(tx *gorm.DB) error {
		if err := repositories.TopicRepository.Updates(tx, topicId, map[string]interface{}{
			"last_comment_time":    comment.CreateTime,
			"last_comment_user_id": comment.UserId,
			"comment_count":        gorm.Expr("comment_count + 1"),
		}); err != nil {
			return err
		}
		if err := tx.Exec("update t_topic_tag set last_comment_time = ?, last_comment_user_id = ? where topic_id = ?",
			comment.CreateTime, comment.UserId, topicId).Error; err != nil {
			return err
		}
		return nil
	})
}

// rss
func (s *topicService) GenerateRss() {
	topics := repositories.TopicRepository.Find(simple.DB(),
		simple.NewSqlCnd().Where("status = ?", constants.StatusOk).Desc("id").Limit(200))

	var items []*feeds.Item
	for _, topic := range topics {
		topicUrl := urls.TopicUrl(topic.Id)
		user := cache.UserCache.Get(topic.UserId)
		if user == nil {
			continue
		}
		item := &feeds.Item{
			Title:       topic.Title,
			Link:        &feeds.Link{Href: topicUrl},
			Description: common.GetMarkdownSummary(topic.Content),
			Author:      &feeds.Author{Name: user.Avatar, Email: user.Email.String},
			Created:     date.FromTimestamp(topic.CreateTime),
		}
		items = append(items, item)
	}
	siteTitle := cache.SysConfigCache.GetValue(constants.SysConfigSiteTitle)
	siteDescription := cache.SysConfigCache.GetValue(constants.SysConfigSiteDescription)
	feed := &feeds.Feed{
		Title:       siteTitle,
		Link:        &feeds.Link{Href: config.Instance.BaseUrl},
		Description: siteDescription,
		Author:      &feeds.Author{Name: siteTitle},
		Created:     time.Now(),
		Items:       items,
	}
	atom, err := feed.ToAtom()
	if err != nil {
		logrus.Error(err)
	} else {
		_ = simple.WriteString(path.Join(config.Instance.StaticPath, "topic_atom.xml"), atom, false)
	}

	rss, err := feed.ToRss()
	if err != nil {
		logrus.Error(err)
	} else {
		_ = simple.WriteString(path.Join(config.Instance.StaticPath, "topic_rss.xml"), rss, false)
	}
}

func (s *topicService) ScanByUser(userId int64, callback func(topics []model.Topic)) {
	var cursor int64 = 0
	for {
		list := repositories.TopicRepository.Find(simple.DB(), simple.NewSqlCnd().
			Eq("user_id", userId).Gt("id", cursor).Asc("id").Limit(1000))
		if list == nil || len(list) == 0 {
			break
		}
		cursor = list[len(list)-1].Id
		callback(list)
	}
}

func (s *topicService) Scan(callback func(topics []model.Topic)) {
	var cursor int64 = 0
	for {
		list := repositories.TopicRepository.Find(simple.DB(), simple.NewSqlCnd().
			Gt("id", cursor).Asc("id").Limit(1000))
		if list == nil || len(list) == 0 {
			break
		}
		cursor = list[len(list)-1].Id
		callback(list)
	}
}

// 倒序扫描
func (s *topicService) ScanDesc(callback func(topics []model.Topic)) {
	var cursor int64 = math.MaxInt64
	for {
		list := repositories.TopicRepository.Find(simple.DB(), simple.NewSqlCnd().
			Cols("id", "status", "create_time").
			Lt("id", cursor).Desc("id").Limit(1000))
		if list == nil || len(list) == 0 {
			break
		}
		cursor = list[len(list)-1].Id
		callback(list)
	}
}

// 倒序扫描
func (s *topicService) ScanDescWithDate(dateFrom, dateTo int64, callback func(topics []model.Topic)) {
	var cursor int64 = math.MaxInt64
	for {
		list := repositories.TopicRepository.Find(simple.DB(), simple.NewSqlCnd().
			Cols("id", "status", "create_time", "update_time").
			Lt("id", cursor).Gte("create_time", dateFrom).Lt("create_time", dateTo).Desc("id").Limit(1000))
		if list == nil || len(list) == 0 {
			break
		}
		cursor = list[len(list)-1].Id
		callback(list)
	}
}

func (s *topicService) GetUserTopics(userId, cursor int64) (topics []model.Topic, nextCursor int64) {
	cnd := simple.NewSqlCnd()
	if userId > 0 {
		cnd.Eq("user_id", userId)
	}
	if cursor > 0 {
		cnd.Lt("id", cursor)
	}
	cnd.Eq("status", constants.StatusOk).Desc("id").Limit(20)
	topics = repositories.TopicRepository.Find(simple.DB(), cnd)
	if len(topics) > 0 {
		nextCursor = topics[len(topics)-1].Id
	} else {
		nextCursor = cursor
	}
	return
}
