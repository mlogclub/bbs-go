package search

import (
	"bbs-go/internal/cache"
	"bbs-go/internal/models"
	"bbs-go/internal/models/constants"
	"bbs-go/internal/pkg/config"
	html2 "bbs-go/internal/pkg/html"
	"bbs-go/internal/pkg/markdown"
	"bbs-go/internal/pkg/text"
	"bbs-go/internal/repositories"
	"html"
	"log/slog"
	"math"
	"time"

	"github.com/blevesearch/bleve/v2"
	blevequery "github.com/blevesearch/bleve/v2/search/query"
	"github.com/mitchellh/mapstructure"
	"github.com/mlogclub/simple/common/dates"
	"github.com/mlogclub/simple/common/strs"
	"github.com/mlogclub/simple/sqls"
	"github.com/spf13/cast"
)

var index bleve.Index

func Init() {
	var err error
	indexPath := config.Instance.Search.IndexPath
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
		Type:       EntityTypeTopic,
		Id:         topic.Id,
		CategoryId: topic.CategoryId,
		UserId:     topic.UserId,
		Title:      html.EscapeString(topic.Title),
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
		doc.Nickname = html.EscapeString(user.Nickname)
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

func NewArticleDoc(article *models.Article) *ArticleDocument {
	if article == nil {
		return nil
	}
	doc := &ArticleDocument{
		Type:       EntityTypeArticle,
		Id:         article.Id,
		UserId:     article.UserId,
		Title:      html.EscapeString(article.Title),
		Summary:    html.EscapeString(article.Summary),
		Status:     article.Status,
		CreateTime: article.CreateTime,
	}

	content := article.Content
	if article.ContentType == constants.ContentTypeMarkdown {
		content = markdown.ToHTML(content)
	}
	content = html2.GetHtmlText(content)
	content = html.EscapeString(content)
	doc.Content = content
	if strs.IsBlank(doc.Summary) {
		doc.Summary = text.GetSummary(content, constants.SummaryLen)
	}

	user := cache.UserCache.Get(article.UserId)
	if user != nil {
		doc.Nickname = html.EscapeString(user.Nickname)
	}

	tags := getArticleTags(article.Id)
	for _, tag := range tags {
		doc.Tags = append(doc.Tags, tag.Name)
	}

	return doc
}

func NewUserDoc(user *models.User) *UserDocument {
	if user == nil {
		return nil
	}
	return &UserDocument{
		Type:         EntityTypeUser,
		Id:           user.Id,
		Username:     html.EscapeString(user.Username.String),
		Nickname:     html.EscapeString(user.Nickname),
		Avatar:       user.Avatar,
		Description:  html.EscapeString(user.Description),
		Status:       user.Status,
		TopicCount:   user.TopicCount,
		CommentCount: user.CommentCount,
		FansCount:    user.FansCount,
		FollowCount:  user.FollowCount,
		Score:        user.Score,
		Exp:          user.Exp,
		Level:        user.Level,
		CreateTime:   user.CreateTime,
	}
}

func getTopicTags(topicId int64) []models.Tag {
	topicTags := repositories.TopicTagRepository.Find(sqls.DB(), sqls.NewCnd().Where("topic_id = ?", topicId))

	var tagIds []int64
	for _, topicTag := range topicTags {
		tagIds = append(tagIds, topicTag.TagId)
	}
	return cache.TagCache.GetList(tagIds)
}

func getArticleTags(articleId int64) []models.Tag {
	tagIds := cache.ArticleTagCache.Get(articleId)
	return cache.TagCache.GetList(tagIds)
}

func UpdateTopicIndexAsync(topic *models.Topic) {
	go UpdateTopicIndex(topic)
}

// IndexData 索引数据
func UpdateTopicIndex(topic *models.Topic) {
	doc := NewTopicDoc(topic)
	if doc == nil {
		return
	}
	err := index.Index(searchDocID(EntityTypeTopic, topic.Id), doc)
	if err != nil {
		slog.Error(err.Error())
	} else {
		slog.Info("add topic search index", slog.Any("id", topic.Id))
	}
}

func DeleteTopicIndex(id int64) error {
	if err := index.Delete(searchDocID(EntityTypeTopic, id)); err != nil {
		return err
	}
	return index.Delete(cast.ToString(id))
}

func UpdateArticleIndex(article *models.Article) {
	doc := NewArticleDoc(article)
	if doc == nil {
		return
	}
	err := index.Index(searchDocID(EntityTypeArticle, article.Id), doc)
	if err != nil {
		slog.Error(err.Error())
	} else {
		slog.Info("add article search index", slog.Any("id", article.Id))
	}
}

func DeleteArticleIndex(id int64) error {
	return index.Delete(searchDocID(EntityTypeArticle, id))
}

func UpdateUserIndex(user *models.User) {
	doc := NewUserDoc(user)
	if doc == nil {
		return
	}
	err := index.Index(searchDocID(EntityTypeUser, user.Id), doc)
	if err != nil {
		slog.Error(err.Error())
	} else {
		slog.Info("add user search index", slog.Any("id", user.Id))
	}
}

func DeleteUserIndex(id int64) error {
	return index.Delete(searchDocID(EntityTypeUser, id))
}

// 分页查询
func SearchTopic(keyword string, categoryId int64, categoryIds []int64, timeRange, page, limit int) (docs []TopicDocument, paging *sqls.Paging, err error) {
	paging = &sqls.Paging{Page: page, Limit: limit}

	query := bleve.NewBooleanQuery()
	query.AddMust(bleve.NewMatchAllQuery())
	query.AddMust(typeQuery(EntityTypeTopic))

	if strs.IsNotBlank(keyword) {
		query.AddMust(keywordQuery(keyword, []string{"title", "content", "tags", "nickname"}))
	}

	if categoryId != 0 {
		if categoryId == -1 { // 推荐
			boolFieldQuery := bleve.NewBoolFieldQuery(true)
			boolFieldQuery.SetField("recommend")
			query.AddMust(boolFieldQuery)
		} else {
			nodeQuery := buildNodeQuery(categoryId, categoryIds)
			if nodeQuery != nil {
				query.AddMust(nodeQuery)
			}
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

func SearchArticle(keyword string, timeRange, page, limit int) (docs []ArticleDocument, paging *sqls.Paging, err error) {
	paging = &sqls.Paging{Page: page, Limit: limit}

	query := bleve.NewBooleanQuery()
	query.AddMust(bleve.NewMatchAllQuery())
	query.AddMust(typeQuery(EntityTypeArticle))
	if strs.IsNotBlank(keyword) {
		query.AddMust(keywordQuery(keyword, []string{"title", "summary", "content", "tags", "nickname"}))
	}
	addTimeRangeQuery(query, timeRange)

	searchRequest := bleve.NewSearchRequest(query)
	searchRequest.From = paging.Offset()
	searchRequest.Size = paging.Limit
	searchRequest.Fields = []string{"*"}
	searchRequest.Highlight = bleve.NewHighlightWithStyle("html")
	searchRequest.Highlight.AddField("title")
	searchRequest.Highlight.AddField("summary")
	searchRequest.Highlight.AddField("content")

	result, err := index.Search(searchRequest)
	if err != nil {
		slog.Error("搜索失败:", slog.Any("err", err))
		return
	}
	for _, hit := range result.Hits {
		storedDoc := hitFields(hit.Fields, hit.Fragments)
		normalizeTags(storedDoc)
		var doc ArticleDocument
		if err := mapstructure.Decode(storedDoc, &doc); err != nil {
			slog.Error(err.Error())
		}
		docs = append(docs, doc)
	}
	return
}

func SearchUser(keyword string, page, limit int) (docs []UserDocument, paging *sqls.Paging, err error) {
	paging = &sqls.Paging{Page: page, Limit: limit}

	query := bleve.NewBooleanQuery()
	query.AddMust(bleve.NewMatchAllQuery())
	query.AddMust(typeQuery(EntityTypeUser))
	if strs.IsNotBlank(keyword) {
		query.AddMust(keywordQuery(keyword, []string{"username", "nickname", "description"}))
	}

	searchRequest := bleve.NewSearchRequest(query)
	searchRequest.From = paging.Offset()
	searchRequest.Size = paging.Limit
	searchRequest.Fields = []string{"*"}
	searchRequest.Highlight = bleve.NewHighlightWithStyle("html")
	searchRequest.Highlight.AddField("nickname")
	searchRequest.Highlight.AddField("username")
	searchRequest.Highlight.AddField("description")

	result, err := index.Search(searchRequest)
	if err != nil {
		slog.Error("搜索失败:", slog.Any("err", err))
		return
	}
	for _, hit := range result.Hits {
		storedDoc := hitFields(hit.Fields, hit.Fragments)
		var doc UserDocument
		if err := mapstructure.Decode(storedDoc, &doc); err != nil {
			slog.Error(err.Error())
		}
		docs = append(docs, doc)
	}
	return
}

func SearchAll(keyword string, limit int) (AllResult, error) {
	if limit <= 0 {
		limit = 5
	}
	topics, _, err := SearchTopic(keyword, 0, nil, 0, 1, limit)
	if err != nil {
		return AllResult{}, err
	}
	articles, _, err := SearchArticle(keyword, 0, 1, limit)
	if err != nil {
		return AllResult{}, err
	}
	users, _, err := SearchUser(keyword, 1, limit)
	if err != nil {
		return AllResult{}, err
	}
	return AllResult{Topics: topics, Articles: articles, Users: users}, nil
}

func buildNodeQuery(categoryId int64, categoryIds []int64) blevequery.Query {
	if len(categoryIds) == 0 {
		return buildExactNodeQuery(categoryId)
	}
	if len(categoryIds) == 1 {
		return buildExactNodeQuery(categoryIds[0])
	}
	queries := make([]blevequery.Query, 0, len(categoryIds))
	for _, id := range categoryIds {
		queries = append(queries, buildExactNodeQuery(id))
	}
	return bleve.NewDisjunctionQuery(queries...)
}

func buildExactNodeQuery(categoryId int64) blevequery.Query {
	f := float64(categoryId)
	b := true
	categoryIdQuery := bleve.NewNumericRangeInclusiveQuery(&f, &f, &b, &b)
	categoryIdQuery.SetField("categoryId")
	return categoryIdQuery
}

func typeQuery(entityType string) blevequery.Query {
	q := bleve.NewTermQuery(entityType)
	q.SetField("type")
	return q
}

func keywordQuery(keyword string, fields []string) blevequery.Query {
	queries := make([]blevequery.Query, 0, len(fields))
	for _, field := range fields {
		q := bleve.NewMatchQuery(keyword)
		q.SetField(field)
		queries = append(queries, q)
	}
	return bleve.NewDisjunctionQuery(queries...)
}

func addTimeRangeQuery(query *blevequery.BooleanQuery, timeRange int) {
	if timeRange == 0 {
		return
	}
	var beginTime int64
	if timeRange == 1 {
		beginTime = dates.Timestamp(time.Now().Add(-24 * time.Hour))
	} else if timeRange == 2 {
		beginTime = dates.Timestamp(time.Now().Add(-7 * 24 * time.Hour))
	} else if timeRange == 3 {
		beginTime = dates.Timestamp(time.Now().AddDate(0, -1, 0))
	} else if timeRange == 4 {
		beginTime = dates.Timestamp(time.Now().AddDate(-1, 0, 0))
	}
	if beginTime == 0 {
		return
	}
	min := float64(beginTime)
	max := float64(math.MaxInt64)
	createTimeQuery := bleve.NewNumericRangeQuery(&min, &max)
	createTimeQuery.SetField("createTime")
	query.AddMust(createTimeQuery)
}

func hitFields(fields map[string]interface{}, fragments map[string][]string) map[string]interface{} {
	storedDoc := make(map[string]interface{})
	for key, field := range fields {
		storedDoc[key] = field
	}
	for field, values := range fragments {
		if len(values) > 0 {
			storedDoc[field] = values[0]
		}
	}
	return storedDoc
}

func normalizeTags(storedDoc map[string]interface{}) {
	tagField, ok := storedDoc["tags"]
	if !ok {
		return
	}
	switch v := tagField.(type) {
	case string:
		storedDoc["tags"] = []string{v}
	case []interface{}:
		var tags []string
		for _, tag := range v {
			if name, ok := tag.(string); ok {
				tags = append(tags, name)
			}
		}
		storedDoc["tags"] = tags
	}
}
