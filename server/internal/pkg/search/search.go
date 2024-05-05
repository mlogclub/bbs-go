package search

import (
	"bbs-go/internal/cache"
	"bbs-go/internal/models"
	html2 "bbs-go/internal/pkg/html"
	"bbs-go/internal/pkg/markdown"
	"bbs-go/internal/repositories"
	"fmt"
	"html"
	"log"
	"log/slog"

	"github.com/blevesearch/bleve/v2"
	"github.com/blevesearch/bleve/v2/index/scorch"
	"github.com/fatih/structs"
	"github.com/mlogclub/simple/common/jsons"
	"github.com/mlogclub/simple/sqls"
	"github.com/spf13/cast"
)

var index bleve.Index

type TopicDocument struct {
	Id         int64    `json:"id"`
	NodeId     int64    `json:"nodeId"`
	UserId     int64    `json:"userId"`
	Nickname   string   `json:"nickname"`
	Title      string   `json:"title"`
	Content    string   `json:"content"`
	Tags       []string `json:"tags"`
	Recommend  bool     `json:"recommend"`
	Status     int      `json:"status"`
	CreateTime int64    `json:"createTime"`
}

func (t *TopicDocument) ToStr() string {
	str, err := jsons.ToStr(t)
	if err != nil {
		slog.Error(err.Error(), slog.Any("err", err))
	}
	return str
}

func Init(indexPath string) {
	var err error
	index, err = bleve.Open(indexPath)
	if err != nil {
		textField := bleve.NewTextFieldMapping()
		textField.Store = true
		textField.Index = true
		textField.IncludeTermVectors = true
		textField.Analyzer = "en"

		numField := bleve.NewNumericFieldMapping()
		numField.DocValues = true
		numField.Store = true
		numField.Index = true

		boolField := bleve.NewBooleanFieldMapping()
		boolField.DocValues = true
		boolField.Store = true
		boolField.Index = true

		indexMapping := bleve.NewIndexMapping()
		indexMapping.DefaultMapping.AddFieldMappingsAt("id", numField)
		indexMapping.DefaultMapping.AddFieldMappingsAt("nodeId", numField)
		indexMapping.DefaultMapping.AddFieldMappingsAt("userId", numField)
		indexMapping.DefaultMapping.AddFieldMappingsAt("nickname", textField)
		indexMapping.DefaultMapping.AddFieldMappingsAt("title", textField)
		indexMapping.DefaultMapping.AddFieldMappingsAt("content", textField)
		// TODO tags
		indexMapping.DefaultMapping.AddFieldMappingsAt("recommend", boolField)
		indexMapping.DefaultMapping.AddFieldMappingsAt("status", numField)
		indexMapping.DefaultMapping.AddFieldMappingsAt("createTime", numField)

		// 使用 scorch 索引类型创建索引
		index, err = bleve.NewUsing(indexPath, indexMapping, scorch.Name, scorch.Name, nil)
		if err != nil {
			log.Fatalf("创建索引失败: %v", err)
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

	err := index.Index(cast.ToString(topic.Id), structs.Map(topic))
	if err != nil {
		slog.Error(err.Error())
	}
}

// 分页查询
func SearchTopic(keyword string, nodeId int64, timeRange, page, limit int) (docs []TopicDocument, paging *sqls.Paging, err error) {
	paging = &sqls.Paging{Page: page, Limit: limit}

	boolQuery := bleve.NewBooleanQuery()
	searchRequest := bleve.NewSearchRequest(boolQuery)
	searchRequest.From = paging.Offset()
	searchRequest.Size = paging.Limit

	results, err := index.Search(searchRequest)
	if err != nil {
		slog.Error("搜索失败:", slog.Any("err", err))
	}

	for _, hit := range results.Hits {
		fmt.Println(hit)

		// var doc TopicDocument

		// doc.Type = strings.Split(hit.ID, "-")[0]

		// if title, ok := hit.Fields["title"].(string); ok {
		// 	doc.Title = title
		// }
		// if content, ok := hit.Fields["content"].(string); ok {
		// 	doc.Content = content
		// }

		// if userId, ok := hit.Fields["userId"].(float64); ok {
		// 	doc.UserId = userId
		// }

		// if did, ok := hit.Fields["did"].(float64); ok {
		// 	doc.Id = did
		// }

		// if createTime, ok := hit.Fields["createTime"].(float64); ok {
		// 	doc.CreateTime = createTime
		// }

		// docs = append(docs, doc)
	}

	return
}
