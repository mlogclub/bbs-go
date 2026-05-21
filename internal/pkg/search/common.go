package search

import (
	"log/slog"
	"strconv"

	"github.com/blevesearch/bleve/v2"
	"github.com/blevesearch/bleve/v2/mapping"
	"github.com/mlogclub/simple/common/jsons"
)

const (
	EntityTypeTopic   = "topic"
	EntityTypeArticle = "article"
	EntityTypeUser    = "user"
)

type TopicDocument struct {
	Type       string   `json:"type"`
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

type ArticleDocument struct {
	Type       string   `json:"type"`
	Id         int64    `json:"id"`
	UserId     int64    `json:"userId"`
	Nickname   string   `json:"nickname"`
	Title      string   `json:"title"`
	Summary    string   `json:"summary"`
	Content    string   `json:"content"`
	Tags       []string `json:"tags"`
	Status     int      `json:"status"`
	CreateTime int64    `json:"createTime"`
}

type UserDocument struct {
	Type         string `json:"type"`
	Id           int64  `json:"id"`
	Username     string `json:"username"`
	Nickname     string `json:"nickname"`
	Avatar       string `json:"avatar"`
	Description  string `json:"description"`
	Status       int    `json:"status"`
	TopicCount   int    `json:"topicCount"`
	CommentCount int    `json:"commentCount"`
	FansCount    int    `json:"fansCount"`
	FollowCount  int    `json:"followCount"`
	Score        int    `json:"score"`
	Exp          int    `json:"exp"`
	Level        int    `json:"level"`
	CreateTime   int64  `json:"createTime"`
}

type AllResult struct {
	Topics   []TopicDocument   `json:"topics"`
	Articles []ArticleDocument `json:"articles"`
	Users    []UserDocument    `json:"users"`
}

func (t *TopicDocument) ToStr() string {
	str, err := jsons.ToStr(t)
	if err != nil {
		slog.Error(err.Error(), slog.Any("err", err))
	}
	return str
}

func searchDocID(entityType string, id int64) string {
	return entityType + ":" + strconv.FormatInt(id, 10)
}

func newIndex(indexPath string) bleve.Index {
	mapping := bleve.NewIndexMapping()
	mapping.DefaultMapping.AddFieldMappingsAt("type", newTextField())
	mapping.DefaultMapping.AddFieldMappingsAt("id", newNumField())
	mapping.DefaultMapping.AddFieldMappingsAt("nodeId", newNumField())
	mapping.DefaultMapping.AddFieldMappingsAt("userId", newNumField())
	mapping.DefaultMapping.AddFieldMappingsAt("username", newTextField())
	mapping.DefaultMapping.AddFieldMappingsAt("nickname", newTextField())
	mapping.DefaultMapping.AddFieldMappingsAt("avatar", newTextField())
	mapping.DefaultMapping.AddFieldMappingsAt("title", newTextField())
	mapping.DefaultMapping.AddFieldMappingsAt("summary", newTextField())
	mapping.DefaultMapping.AddFieldMappingsAt("content", newTextField())
	mapping.DefaultMapping.AddFieldMappingsAt("description", newTextField())
	mapping.DefaultMapping.AddFieldMappingsAt("tags", newTextField())
	mapping.DefaultMapping.AddFieldMappingsAt("recommend", newBoolField())
	mapping.DefaultMapping.AddFieldMappingsAt("status", newNumField())
	mapping.DefaultMapping.AddFieldMappingsAt("topicCount", newNumField())
	mapping.DefaultMapping.AddFieldMappingsAt("commentCount", newNumField())
	mapping.DefaultMapping.AddFieldMappingsAt("fansCount", newNumField())
	mapping.DefaultMapping.AddFieldMappingsAt("followCount", newNumField())
	mapping.DefaultMapping.AddFieldMappingsAt("score", newNumField())
	mapping.DefaultMapping.AddFieldMappingsAt("exp", newNumField())
	mapping.DefaultMapping.AddFieldMappingsAt("level", newNumField())
	mapping.DefaultMapping.AddFieldMappingsAt("createTime", newNumField())

	index, err := bleve.New(indexPath, mapping)
	if err != nil {
		slog.Info("创建索引失败", slog.Any("err", err))
	}
	return index
}

func newTextField() *mapping.FieldMapping {
	textField := bleve.NewTextFieldMapping()
	// textField.Store = true
	textField.Index = true
	textField.IncludeTermVectors = true
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
