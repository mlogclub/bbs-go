package search

import (
	"bbs-go/internal/cache"
	"bbs-go/internal/models"
	html2 "bbs-go/internal/pkg/html"
	"bbs-go/internal/pkg/markdown"
	"bbs-go/internal/repositories"
	"html"
	"log/slog"
	"math"
	"time"

	"github.com/blevesearch/bleve/v2"
	"github.com/mitchellh/mapstructure"
	"github.com/mlogclub/simple/common/dates"
	"github.com/mlogclub/simple/common/strs"
	"github.com/mlogclub/simple/sqls"
	"github.com/spf13/cast"
)

var index bleve.Index

func Init(indexPath string) {
	var err error
	if index, err = bleve.Open(indexPath); err != nil {
		if err == bleve.ErrorIndexPathDoesNotExist {
			index = newIndex(indexPath)
		} else {
			slog.Error(err.Error())
		}
	}
}

func NewTopicDoc(topic *models.Topic) *TopicDocument {
	if topic == nil {
		return nil
	}
	doc := &TopicDocument{
		Id:         topic.Id,
		NodeId:     topic.NodeId,
		UserId:     topic.UserId,
		Title:      topic.Title,
		Status:     topic.Status,
		Recommend:  topic.Recommend,
		CreateTime: topic.CreateTime,
	}

	// 处理内容
	content := markdown.ToHTML(topic.Content)
	content = html2.GetHtmlText(content)
	content = html.EscapeString(content)

	doc.Content = content

	// 处理用户
	user := cache.UserCache.Get(topic.UserId)
	if user != nil {
		doc.Nickname = user.Nickname
	}

	// 处理标签
	tags := getTopicTags(topic.Id)
	var tagsArr []string
	for _, tag := range tags {
		tagsArr = append(tagsArr, tag.Name)
	}
	tagsArr = append(tagsArr, "hello")
	doc.Tags = tagsArr

	return doc
}

func getTopicTags(topicId int64) []models.Tag {
	topicTags := repositories.TopicTagRepository.Find(sqls.DB(), sqls.NewCnd().Where("topic_id = ?", topicId))

	var tagIds []int64
	for _, topicTag := range topicTags {
		tagIds = append(tagIds, topicTag.TagId)
	}
	return cache.TagCache.GetList(tagIds)
}

// IndexData 索引数据
func UpdateTopicIndex(topic *models.Topic) {
	doc := NewTopicDoc(topic)
	if doc == nil {
		return
	}
	err := index.Index(cast.ToString(topic.Id), doc)
	if err != nil {
		slog.Error(err.Error())
	} else {
		slog.Info("add topic search index", slog.Any("id", topic.Id))
	}
}

func DeleteTopicIndex(id int64) error {
	return index.Delete(cast.ToString(id))
}

// 分页查询
func SearchTopic(keyword string, nodeId int64, timeRange, page, limit int) (docs []TopicDocument, paging *sqls.Paging, err error) {
	paging = &sqls.Paging{Page: page, Limit: limit}

	query := bleve.NewBooleanQuery()
	query.AddMust(bleve.NewMatchAllQuery())

	if strs.IsNotBlank(keyword) {
		query.AddMust(bleve.NewMatchQuery(keyword))
	}

	if nodeId != 0 {
		if nodeId == -1 { // 推荐
			boolFieldQuery := bleve.NewBoolFieldQuery(true)
			boolFieldQuery.SetField("recommend")
			query.AddMust(boolFieldQuery)
		} else {
			f := float64(nodeId)
			b := true
			nodeIdQuery := bleve.NewNumericRangeInclusiveQuery(&f, &f, &b, &b)
			nodeIdQuery.SetField("nodeId")
			query.AddMust(nodeIdQuery)
		}
	}
	if timeRange != 0 {
		var beginTime int64
		if timeRange == 1 { // 一天内
			beginTime = dates.Timestamp(time.Now().Add(-24 * time.Hour))
		} else if timeRange == 2 { // 一周内
			beginTime = dates.Timestamp(time.Now().Add(-7 * 24 * time.Hour))
		} else if timeRange == 3 { // 一月内
			beginTime = dates.Timestamp(time.Now().AddDate(0, -1, 0))
		} else if timeRange == 4 { // 一年内
			beginTime = dates.Timestamp(time.Now().AddDate(-1, 0, 0))
		}

		min := float64(beginTime)
		max := float64(math.MaxInt64)
		createTimeQuery := bleve.NewNumericRangeQuery(&min, &max)
		createTimeQuery.SetField("createTime")
		query.AddMust(createTimeQuery)
	}

	searchRequest := bleve.NewSearchRequest(query)
	searchRequest.From = paging.Offset()
	searchRequest.Size = paging.Limit
	searchRequest.Fields = []string{"*"}
	searchRequest.Highlight = bleve.NewHighlightWithStyle("html")
	searchRequest.Highlight.AddField("title")
	searchRequest.Highlight.AddField("content")

	result, err := index.Search(searchRequest)
	if err != nil {
		slog.Error("搜索失败:", slog.Any("err", err))
	}

	for _, hit := range result.Hits {

		storedDoc := make(map[string]interface{})
		for key, field := range hit.Fields {
			storedDoc[key] = field
		}

		for field, fragments := range hit.Fragments {
			if len(fragments) > 0 {
				storedDoc[field] = fragments[0]
			}
		}

		if tagField, ok := storedDoc["tags"]; ok {
			switch v := tagField.(type) {
			case string:
				storedDoc["tags"] = []string{v}
			case []interface{}:
				var tags []string
				for _, tag := range v {
					tags = append(tags, tag.(string))
				}
				storedDoc["tags"] = tags
			}
		}

		var doc TopicDocument
		if err := mapstructure.Decode(storedDoc, &doc); err != nil {
			slog.Error(err.Error())
		}
		docs = append(docs, doc)
	}

	return
}
