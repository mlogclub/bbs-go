package services

import (
	"bbs-go/internal/models"
	"bbs-go/internal/models/constants"
	"bbs-go/internal/pkg/search"
	"log/slog"
	"sync"

	"github.com/mlogclub/simple/common/dates"
	"github.com/mlogclub/simple/sqls"
)

var SearchReindexService = newSearchReindexService()

type SearchReindexStatus struct {
	Running          bool   `json:"running"`
	Processed        int64  `json:"processed"`
	Total            int64  `json:"total"`
	TopicProcessed   int64  `json:"topicProcessed"`
	TopicTotal       int64  `json:"topicTotal"`
	ArticleProcessed int64  `json:"articleProcessed"`
	ArticleTotal     int64  `json:"articleTotal"`
	UserProcessed    int64  `json:"userProcessed"`
	UserTotal        int64  `json:"userTotal"`
	StartedAt        int64  `json:"startedAt"`
	FinishedAt       int64  `json:"finishedAt"`
	Error            string `json:"error"`
}

type searchReindexService struct {
	mu     sync.Mutex
	status SearchReindexStatus
}

func newSearchReindexService() *searchReindexService {
	return &searchReindexService{}
}

func (s *searchReindexService) Status() SearchReindexStatus {
	s.mu.Lock()
	defer s.mu.Unlock()
	return s.status
}

func (s *searchReindexService) Start() (SearchReindexStatus, bool) {
	s.mu.Lock()
	if s.status.Running {
		status := s.status
		s.mu.Unlock()
		return status, false
	}

	s.status = SearchReindexStatus{
		Running:   true,
		StartedAt: dates.NowTimestamp(),
	}
	status := s.status
	s.mu.Unlock()

	go s.run()
	return status, true
}

func (s *searchReindexService) run() {
	defer func() {
		if r := recover(); r != nil {
			slog.Error("search reindex failed", slog.Any("err", r))
			s.finishWithError("search reindex failed")
		}
	}()

	topicTotal := TopicService.Count(sqls.NewCnd().Where("status <> ?", constants.StatusDeleted))
	var articleTotal int64
	sqls.DB().Model(&models.Article{}).Where("status <> ?", constants.StatusDeleted).Count(&articleTotal)
	var userTotal int64
	sqls.DB().Model(&models.User{}).Where("status <> ?", constants.StatusDeleted).Count(&userTotal)
	s.setTotals(topicTotal, articleTotal, userTotal)

	TopicService.ScanDesc(func(topics []models.Topic) {
		for _, topic := range topics {
			if topic.Status == constants.StatusDeleted {
				continue
			}
			search.UpdateTopicIndex(&topic)
			s.incrementTopicProcessed()
		}
	})
	ArticleService.ScanDesc(func(articles []models.Article) {
		for _, article := range articles {
			if article.Status == constants.StatusDeleted {
				continue
			}
			search.UpdateArticleIndex(&article)
			s.incrementArticleProcessed()
		}
	})
	UserService.Scan(func(users []models.User) {
		for _, user := range users {
			if user.Status == constants.StatusDeleted {
				continue
			}
			search.UpdateUserIndex(&user)
			s.incrementUserProcessed()
		}
	})
	s.finishWithError("")
}

func (s *searchReindexService) setTotals(topicTotal, articleTotal, userTotal int64) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.status.TopicTotal = topicTotal
	s.status.ArticleTotal = articleTotal
	s.status.UserTotal = userTotal
	s.status.Total = topicTotal + articleTotal + userTotal
}

func (s *searchReindexService) incrementTopicProcessed() {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.status.Processed++
	s.status.TopicProcessed++
}

func (s *searchReindexService) incrementArticleProcessed() {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.status.Processed++
	s.status.ArticleProcessed++
}

func (s *searchReindexService) incrementUserProcessed() {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.status.Processed++
	s.status.UserProcessed++
}

func (s *searchReindexService) finishWithError(message string) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.status.Running = false
	s.status.FinishedAt = dates.NowTimestamp()
	s.status.Error = message
}
