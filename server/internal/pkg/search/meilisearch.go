package search

import (
	"context"
	"encoding/json"
	"fmt"
	"log/slog"
	"time"

	"bbs-go/internal/models"
	"bbs-go/internal/models/constants"
	"bbs-go/internal/pkg/config"
	"bbs-go/internal/repositories"

	"github.com/meilisearch/meilisearch-go"
	"github.com/spf13/cast"

	"bbs-go/internal/pkg/simple/common/dates"
	"bbs-go/internal/pkg/simple/common/strs"
	"bbs-go/internal/pkg/simple/sqls"
)

const (
	TopicIndexName = "topics"
)

type MeiliSearchClient struct {
	client meilisearch.ServiceManager
	ctx    context.Context
}

var meiliClient *MeiliSearchClient

func InitMeiliSearch() {
	if !config.Instance.MeiliSearch.Enabled {
		slog.Info("MeiliSearch is disabled, skipping initialization")
		return
	}

	cfg := config.Instance.MeiliSearch
	if cfg.Host == "" {
		cfg.Host = "localhost"
	}
	if cfg.Port == 0 {
		cfg.Port = 7700
	}
	if cfg.Index == "" {
		cfg.Index = TopicIndexName
	}

	serverURL := fmt.Sprintf("http://%s:%d", cfg.Host, cfg.Port)

	client := meilisearch.New(serverURL, meilisearch.WithAPIKey(cfg.APIKey))

	meiliClient = &MeiliSearchClient{
		client: client,
		ctx:    context.Background(),
	}

	err := meiliClient.setupIndex()
	if err != nil {
		slog.Error("Failed to setup MeiliSearch index", slog.Any("error", err))
	} else {
		slog.Info("MeiliSearch initialized successfully", slog.String("url", serverURL), slog.String("index", cfg.Index))
	}
}

func (m *MeiliSearchClient) setupIndex() error {
	index := m.client.Index(config.Instance.MeiliSearch.Index)

	settings := &meilisearch.Settings{
		SearchableAttributes: []string{
			"title",
			"content",
			"nickname",
			"tags",
		},
		FilterableAttributes: []string{
			"nodeId",
			"userId",
			"status",
			"recommend",
			"createTime",
		},
		SortableAttributes: []string{
			"createTime",
			"id",
		},
		DisplayedAttributes: []string{"*"},
	}

	_, err := index.UpdateSettings(settings)
	if err != nil {
		return fmt.Errorf("failed to update index settings: %w", err)
	}

	return nil
}

func UpdateTopicIndexMeili(topic *models.Topic) {
	if meiliClient == nil {
		slog.Warn("MeiliSearch client not initialized")
		return
	}

	doc := NewTopicDoc(topic)
	if doc == nil {
		return
	}

	index := meiliClient.client.Index(config.Instance.MeiliSearch.Index)

	// Convert to the correct format for MeiliSearch
	docMap := map[string]interface{}{
		"id":         doc.Id,
		"nodeId":     doc.NodeId,
		"userId":     doc.UserId,
		"nickname":   doc.Nickname,
		"title":      doc.Title,
		"content":    doc.Content,
		"tags":       doc.Tags,
		"recommend":  doc.Recommend,
		"status":     doc.Status,
		"createTime": doc.CreateTime,
	}

	_, err := index.AddDocuments([]map[string]interface{}{docMap}, nil)
	if err != nil {
		slog.Error("Failed to add document to MeiliSearch", slog.Any("error", err), slog.Any("id", topic.Id))
	} else {
		slog.Info("Added topic to MeiliSearch index", slog.Any("id", topic.Id))
	}
}

func DeleteTopicIndexMeili(id int64) error {
	if meiliClient == nil {
		return fmt.Errorf("MeiliSearch client not initialized")
	}

	index := meiliClient.client.Index(config.Instance.MeiliSearch.Index)

	_, err := index.DeleteDocument(cast.ToString(id))
	if err != nil {
		return fmt.Errorf("failed to delete document from MeiliSearch: %w", err)
	}

	return nil
}

func SearchTopicMeili(keyword string, nodeId int64, timeRange, page, limit int) (docs []TopicDocument, paging *sqls.Paging, err error) {
	if meiliClient == nil {
		return nil, nil, fmt.Errorf("MeiliSearch client not initialized")
	}

	paging = &sqls.Paging{Page: page, Limit: limit}

	index := meiliClient.client.Index(config.Instance.MeiliSearch.Index)

	var filters []string

	if nodeId != 0 {
		if nodeId == -1 {
			filters = append(filters, "recommend = true")
		} else {
			filters = append(filters, fmt.Sprintf("nodeId = %d", nodeId))
		}
	}

	if timeRange != 0 {
		var beginTime int64
		switch timeRange {
		case 1:
			beginTime = dates.Timestamp(time.Now().Add(-24 * time.Hour))
		case 2:
			beginTime = dates.Timestamp(time.Now().Add(-7 * 24 * time.Hour))
		case 3:
			beginTime = dates.Timestamp(time.Now().AddDate(0, -1, 0))
		case 4:
			beginTime = dates.Timestamp(time.Now().AddDate(-1, 0, 0))
		}

		if beginTime > 0 {
			filters = append(filters, fmt.Sprintf("createTime >= %d", beginTime))
		}
	}

	searchParams := &meilisearch.SearchRequest{
		Limit:                 int64(limit),
		Offset:                int64(paging.Offset()),
		AttributesToHighlight: []string{"title", "content"},
		HighlightPreTag:       "<em>",
		HighlightPostTag:      "</em>",
	}

	if len(filters) > 0 {
		searchParams.Filter = filters
	}

	query := ""
	if strs.IsNotBlank(keyword) {
		query = keyword
	}

	searchResult, err := index.Search(query, searchParams)
	if err != nil {
		slog.Error("MeiliSearch search failed", slog.Any("error", err))
		return nil, paging, err
	}

	for _, hit := range searchResult.Hits {
		var doc TopicDocument

		// Helper function to unmarshal json.RawMessage to a specific type
		unmarshallField := func(fieldName string, target interface{}) bool {
			if rawData, exists := hit[fieldName]; exists {
				if err := json.Unmarshal(rawData, target); err == nil {
					return true
				}
			}
			return false
		}

		var idFloat float64
		if unmarshallField("id", &idFloat) {
			doc.Id = int64(idFloat)
		}

		var nodeIdFloat float64
		if unmarshallField("nodeId", &nodeIdFloat) {
			doc.NodeId = int64(nodeIdFloat)
		}

		var userIdFloat float64
		if unmarshallField("userId", &userIdFloat) {
			doc.UserId = int64(userIdFloat)
		}

		unmarshallField("nickname", &doc.Nickname)
		unmarshallField("title", &doc.Title)
		unmarshallField("content", &doc.Content)
		unmarshallField("recommend", &doc.Recommend)

		var statusFloat float64
		if unmarshallField("status", &statusFloat) {
			doc.Status = int(statusFloat)
		}

		var createTimeFloat float64
		if unmarshallField("createTime", &createTimeFloat) {
			doc.CreateTime = int64(createTimeFloat)
		}

		unmarshallField("tags", &doc.Tags)

		// Handle formatted/highlighted results
		var formatted map[string]interface{}
		if unmarshallField("_formatted", &formatted) {
			if highlightedTitle, ok := formatted["title"].(string); ok {
				doc.Title = highlightedTitle
			}
			if highlightedContent, ok := formatted["content"].(string); ok {
				doc.Content = highlightedContent
			}
		}

		docs = append(docs, doc)
	}

	return docs, paging, nil
}

func ReindexAllTopicsMeili() error {
	if meiliClient == nil {
		return fmt.Errorf("MeiliSearch client not initialized")
	}

	slog.Info("Starting MeiliSearch reindexing")

	index := meiliClient.client.Index(config.Instance.MeiliSearch.Index)

	_, err := index.DeleteAllDocuments()
	if err != nil {
		slog.Error("Failed to clear MeiliSearch index", slog.Any("error", err))
		return err
	}

	var totalIndexed int
	batchSize := 100
	var batch []TopicDocument
	offset := 0

	for {
		cnd := sqls.NewCnd().Where("status != ?", constants.StatusDeleted).
			Limit(batchSize).
			Desc("id")
		cnd.Paging.Page = (offset / batchSize) + 1
		cnd.Paging.Limit = batchSize

		topics := repositories.TopicRepository.Find(sqls.DB(), cnd)

		if len(topics) == 0 {
			break
		}

		for _, topic := range topics {
			doc := NewTopicDoc(&topic)
			if doc != nil {
				batch = append(batch, *doc)
			}
		}

		if len(batch) > 0 {
			if err := indexBatchMeili(index, batch); err != nil {
				slog.Error("Failed to index batch", slog.Any("error", err))
			} else {
				totalIndexed += len(batch)
				slog.Info("Indexed batch", slog.Int("count", len(batch)), slog.Int("total", totalIndexed))
			}
			batch = nil
		}

		if len(topics) < batchSize {
			break
		}

		offset += batchSize
	}

	slog.Info("MeiliSearch reindexing completed", slog.Int("totalIndexed", totalIndexed))
	return nil
}

func indexBatchMeili(index meilisearch.IndexManager, docs []TopicDocument) error {
	// Convert documents to the correct format
	docMaps := make([]map[string]interface{}, len(docs))
	for i, doc := range docs {
		docMaps[i] = map[string]interface{}{
			"id":         doc.Id,
			"nodeId":     doc.NodeId,
			"userId":     doc.UserId,
			"nickname":   doc.Nickname,
			"title":      doc.Title,
			"content":    doc.Content,
			"tags":       doc.Tags,
			"recommend":  doc.Recommend,
			"status":     doc.Status,
			"createTime": doc.CreateTime,
		}
	}

	_, err := index.AddDocuments(docMaps, nil)
	return err
}
