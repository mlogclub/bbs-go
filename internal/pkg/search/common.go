package search

import (
	"log/slog"

	"github.com/blevesearch/bleve/v2"
	"github.com/blevesearch/bleve/v2/mapping"
	"github.com/mlogclub/simple/common/jsons"
)

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

func newIndex(indexPath string) bleve.Index {
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
