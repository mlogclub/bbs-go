package services

import (
	"math"
	"path"
	"time"

	"github.com/gorilla/feeds"
	"github.com/jinzhu/gorm"
	"github.com/mlogclub/simple"
	"github.com/sirupsen/logrus"

	"github.com/mlogclub/bbs-go/common"
	"github.com/mlogclub/bbs-go/common/config"
	"github.com/mlogclub/bbs-go/common/urls"
	"github.com/mlogclub/bbs-go/model"
	"github.com/mlogclub/bbs-go/repositories"
	"github.com/mlogclub/bbs-go/services/cache"
)

type ScanTopicCallback func(topics []model.Topic) bool

var TopicService = newTopicService()

func newTopicService() *topicService {
	return &topicService{}
}

type topicService struct{}

func (this *topicService) Get(id int64) *model.Topic {
	return repositories.TopicRepository.Get(simple.GetDB(), id)
}

func (this *topicService) Take(where ...interface{}) *model.Topic {
	return repositories.TopicRepository.Take(simple.GetDB(), where...)
}

func (this *topicService) QueryCnd(cnd *simple.QueryCnd) (list []model.Topic, err error) {
	return repositories.TopicRepository.QueryCnd(simple.GetDB(), cnd)
}

func (this *topicService) Query(queries *simple.ParamQueries) (list []model.Topic, paging *simple.Paging) {
	return repositories.TopicRepository.Query(simple.GetDB(), queries)
}

func (this *topicService) Create(t *model.Topic) error {
	return repositories.TopicRepository.Create(simple.GetDB(), t)
}

func (this *topicService) Update(t *model.Topic) error {
	return repositories.TopicRepository.Update(simple.GetDB(), t)
}

func (this *topicService) Updates(id int64, columns map[string]interface{}) error {
	return repositories.TopicRepository.Updates(simple.GetDB(), id, columns)
}

func (this *topicService) UpdateColumn(id int64, name string, value interface{}) error {
	return repositories.TopicRepository.UpdateColumn(simple.GetDB(), id, name, value)
}

func (this *topicService) Delete(id int64) error {
	err := repositories.TopicRepository.UpdateColumn(simple.GetDB(), id, "status", model.TopicStatusDeleted)
	if err == nil {
		// 删掉专栏文章
		SubjectContentService.DeleteByEntity(model.EntityTypeTopic, id)
		// 删掉标签文章
		TopicTagService.DeleteByTopicId(id)
	}
	return err
}

// 扫描
func (this *topicService) Scan(cb ScanTopicCallback) {
	var cursor int64
	for {
		list, err := repositories.TopicRepository.QueryCnd(simple.GetDB(), simple.NewQueryCnd("id > ?",
			cursor).Order("id asc").Size(100))
		if err != nil {
			break
		}
		if list == nil || len(list) == 0 {
			break
		}
		cursor = list[len(list)-1].Id
		if !cb(list) {
			break
		}
	}
}

// 倒序扫描
func (this *topicService) ScanDesc(cb ScanTopicCallback) {
	var cursor int64 = math.MaxInt64
	for {
		list, err := repositories.TopicRepository.QueryCnd(simple.GetDB(), simple.NewQueryCnd("id < ?",
			cursor).Order("id desc").Size(100))
		if err != nil {
			break
		}
		if list == nil || len(list) == 0 {
			break
		}
		cursor = list[len(list)-1].Id
		if !cb(list) {
			break
		}
	}
}

// 发表
func (this *topicService) Publish(userId int64, tags []string, title, content string, extra map[string]interface{}) (*model.Topic, *simple.CodeError) {
	if len(title) == 0 {
		return nil, simple.NewErrorMsg("标题不能为空")
	}

	if simple.RuneLen(title) > 128 {
		return nil, simple.NewErrorMsg("标题长度不能超过128")
	}

	now := simple.NowTimestamp()
	topic := &model.Topic{
		UserId:          userId,
		Title:           title,
		Content:         content,
		Status:          model.TopicStatusOk,
		LastCommentTime: now,
		CreateTime:      now,
	}
	if len(extra) > 0 {
		topic.ExtraData, _ = simple.FormatJson(extra)
	}

	err := simple.Tx(simple.GetDB(), func(tx *gorm.DB) error {
		tagIds := repositories.TagRepository.GetOrCreates(tx, tags)
		err := repositories.TopicRepository.Create(simple.GetDB(), topic)
		if err != nil {
			return err
		}

		repositories.TopicTagRepository.AddTopicTags(tx, topic.Id, tagIds)
		return nil
	})
	if err == nil {
		common.BaiduUrlPush([]string{urls.TopicUrl(topic.Id)})
	}
	return topic, simple.FromError(err)
}

// 更新
func (this *topicService) Edit(topicId int64, tags []string, title, content string) *simple.CodeError {
	if len(title) == 0 {
		return simple.NewErrorMsg("标题不能为空")
	}

	if simple.RuneLen(title) > 128 {
		return simple.NewErrorMsg("标题长度不能超过128")
	}

	err := simple.Tx(simple.GetDB(), func(tx *gorm.DB) error {
		tagIds := repositories.TagRepository.GetOrCreates(tx, tags)
		err := repositories.TopicRepository.Updates(simple.GetDB(), topicId, map[string]interface{}{
			"title":   title,
			"content": content,
		})
		if err != nil {
			return err
		}
		repositories.TopicTagRepository.RemoveTopicTags(tx, topicId)      // 先删掉所有的标签
		repositories.TopicTagRepository.AddTopicTags(tx, topicId, tagIds) // 然后重新添加标签
		return nil
	})
	return simple.FromError(err)
}

// 主题标签
func (this *topicService) GetTopicTags(topicId int64) []model.Tag {
	topicTags, err := repositories.TopicTagRepository.QueryCnd(simple.GetDB(), simple.NewQueryCnd("topic_id = ?", topicId))
	if err != nil {
		return nil
	}
	var tagIds []int64
	for _, topicTag := range topicTags {
		tagIds = append(tagIds, topicTag.TagId)
	}
	return cache.TagCache.GetList(tagIds)
}

// 指定标签下的主题列表
func (this *topicService) GetTagTopics(tagId int64, page int) (topics []model.Topic, paging *simple.Paging) {
	topicTags, paging := repositories.TopicTagRepository.Query(simple.GetDB(), simple.NewParamQueries(nil).
		Eq("tag_id", tagId).
		Eq("status", model.ArticleTagStatusOk).
		Page(page, 20).Desc("id"))
	if len(topicTags) > 0 {
		var topicIds []int64
		for _, topicTag := range topicTags {
			topicIds = append(topicIds, topicTag.TopicId)
		}
		topics = this.GetTopicInIds(topicIds)
	}
	return
}

// 根据编号批量获取主题
func (this *topicService) GetTopicInIds(topicIds []int64) []model.Topic {
	if len(topicIds) == 0 {
		return nil
	}
	var topics []model.Topic
	simple.GetDB().Where("id in (?)", topicIds).Find(&topics)
	return topics
}

// 浏览数+1
func (this *topicService) IncrViewCount(topicId int64) {
	simple.GetDB().Exec("update t_topic set view_count = view_count + 1 where id = ?", topicId)
}

// 更新最后回复时间
func (this *topicService) SetLastCommentTime(topicId, lastCommentTime int64) {
	err := this.UpdateColumn(topicId, "last_comment_time", lastCommentTime)
	if err != nil {
		logrus.Error(err)
	}
}

// rss
func (this *topicService) GenerateRss() {
	topics, err := repositories.TopicRepository.QueryCnd(simple.GetDB(),
		simple.NewQueryCnd("status = ?", model.TopicStatusOk).Order("id desc").Size(1000))
	if err != nil {
		logrus.Error(err)
		return
	}

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
			Author:      &feeds.Author{Name: user.Avatar, Email: user.Email},
			Created:     simple.TimeFromTimestamp(topic.CreateTime),
		}
		items = append(items, item)
	}
	siteTitle := cache.SysConfigCache.GetValue(model.SysConfigSiteTitle)
	siteDescription := cache.SysConfigCache.GetValue(model.SysConfigSiteDescription)
	feed := &feeds.Feed{
		Title:       siteTitle,
		Link:        &feeds.Link{Href: config.Conf.BaseUrl},
		Description: siteDescription,
		Author:      &feeds.Author{Name: siteTitle},
		Created:     time.Now(),
		Items:       items,
	}
	atom, err := feed.ToAtom()
	if err != nil {
		logrus.Error(err)
	} else {
		_ = simple.WriteString(path.Join(config.Conf.StaticPath, "topic_atom.xml"), atom, false)
	}

	rss, err := feed.ToRss()
	if err != nil {
		logrus.Error(err)
	} else {
		_ = simple.WriteString(path.Join(config.Conf.StaticPath, "topic_rss.xml"), rss, false)
	}
}
