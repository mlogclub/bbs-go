package search

import (
	"bbs-go/internal/cache"
	"bbs-go/internal/models"
	html2 "bbs-go/internal/pkg/html"
	"bbs-go/internal/pkg/markdown"
	"bbs-go/internal/repositories"
	"html"
	"log"
	"log/slog"

	"github.com/blevesearch/bleve/v2"
	"github.com/blevesearch/bleve/v2/mapping"
	"github.com/mitchellh/mapstructure"
	"github.com/mlogclub/simple/common/jsons"
	"github.com/mlogclub/simple/common/strs"
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

func newTextField() *mapping.FieldMapping {
	textField := bleve.NewTextFieldMapping()
	// textField.Store = true
	textField.Index = true
	textField.IncludeTermVectors = true
	textField.Analyzer = "en"
	return textField
}

func newNumField() *mapping.FieldMapping {
	numField := bleve.NewNumericFieldMapping()
	// numField.Store = true
	numField.Index = true
	numField.DocValues = true
	return numField
}

func newBoolField() *mapping.FieldMapping {
	boolField := bleve.NewBooleanFieldMapping()
	// boolField.Store = true
	boolField.Index = true
	boolField.DocValues = true
	return boolField
}

func Init(indexPath string) {
	var err error
	index, err = bleve.Open(indexPath)
	if err == bleve.ErrorIndexPathDoesNotExist {
		mapping := bleve.NewIndexMapping()
		mapping.DefaultMapping.AddFieldMappingsAt("id", newNumField())
		mapping.DefaultMapping.AddFieldMappingsAt("nodeId", newNumField())
		mapping.DefaultMapping.AddFieldMappingsAt("userId", newNumField())
		mapping.DefaultMapping.AddFieldMappingsAt("nickname", newTextField())
		mapping.DefaultMapping.AddFieldMappingsAt("title", newTextField())
		mapping.DefaultMapping.AddFieldMappingsAt("content", newTextField())
		mapping.DefaultMapping.AddFieldMappingsAt("tags", newTextField())
		mapping.DefaultMapping.AddFieldMappingsAt("recommend", newBoolField())
		mapping.DefaultMapping.AddFieldMappingsAt("status", newNumField())
		mapping.DefaultMapping.AddFieldMappingsAt("createTime", newNumField())

		// 使用 scorch 索引类型创建索引
		index, err = bleve.New(indexPath, mapping)
		if err != nil {
			log.Fatalf("创建索引失败: %v", err)
		}
	} else if err != nil {
		slog.Error(err.Error())
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

	// var result map[string]interface{}
	// if err := mapstructure.Decode(doc, &result); err != nil {
	// 	slog.Error(err.Error())
	// }
	err := index.Index(cast.ToString(topic.Id), doc)
	if err != nil {
		slog.Error(err.Error())
	}
}

// 分页查询
func SearchTopic(keyword string, nodeId int64, timeRange, page, limit int) (docs []TopicDocument, paging *sqls.Paging, err error) {
	paging = &sqls.Paging{Page: page, Limit: limit}

	query := bleve.NewBooleanQuery()
	query.AddMust(bleve.NewMatchAllQuery())

	if strs.IsNotBlank(keyword) {
		query.AddMust(bleve.NewMatchQuery(keyword))
	}

	searchRequest := bleve.NewSearchRequest(query)
	searchRequest.From = paging.Offset()
	searchRequest.Size = paging.Limit
	searchRequest.Fields = []string{"*"}
	searchRequest.Highlight = bleve.NewHighlightWithStyle("html")
	searchRequest.Highlight.AddField("title")
	searchRequest.Highlight.AddField("content")

	results, err := index.Search(searchRequest)
	if err != nil {
		slog.Error("搜索失败:", slog.Any("err", err))
	}

	for _, hit := range results.Hits {

		storedDoc := make(map[string]interface{})
		for key, field := range hit.Fields {
			storedDoc[key] = field
		}

		for field, fragments := range hit.Fragments {
			if len(fragments) > 0 {
				storedDoc[field] = fragments[0]
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
