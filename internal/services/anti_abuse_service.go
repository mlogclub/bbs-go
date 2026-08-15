package services

import (
	"bbs-go/internal/models"
	"bbs-go/internal/models/dto"
	"bbs-go/internal/pkg/locales"
	"errors"

	"github.com/mlogclub/simple/common/dates"
	"gorm.io/gorm"
)

var AntiAbuseService = newAntiAbuseService()

type publishContentType string

const (
	publishContentTopic   publishContentType = "topic"
	publishContentArticle publishContentType = "article"
	publishContentComment publishContentType = "comment"
)

type antiAbuseService struct{}

func newAntiAbuseService() *antiAbuseService { return &antiAbuseService{} }

// CheckTopic returns whether a topic should enter the existing review queue.
func (s *antiAbuseService) CheckTopic(tx *gorm.DB, userID int64, ip string) (bool, error) {
	return s.check(tx, publishContentTopic, userID, ip, true)
}

// CheckArticle returns whether an article should enter the existing review queue.
func (s *antiAbuseService) CheckArticle(tx *gorm.DB, userID int64, ip string) (bool, error) {
	return s.check(tx, publishContentArticle, userID, ip, true)
}

// CheckComment only rejects on a match. Comments do not currently have a
// corresponding administration review queue, so they must not be persisted as
// orphaned pending content.
func (s *antiAbuseService) CheckComment(tx *gorm.DB, userID int64, ip string) error {
	_, err := s.check(tx, publishContentComment, userID, ip, false)
	return err
}

func (s *antiAbuseService) check(tx *gorm.DB, contentType publishContentType, userID int64, ip string, supportsReview bool) (bool, error) {
	config := SysConfigService.GetAntiAbuseConfigTx(tx)
	review := false

	if config.User.Enabled {
		rule := rateLimitFor(config.User, contentType)
		if s.exceeded(tx, contentType, rule, "user_id", userID) {
			if config.User.Action == dto.AntiAbuseActionReview && supportsReview {
				review = true
			} else {
				return false, errors.New(locales.Get("errors.too_fast"))
			}
		}
	}

	// IP frequency control is always reject-only. IP is collected server-side by
	// the API handlers and is never read from a client supplied field.
	if config.IP.Enabled && ip != "" {
		rule := rateLimitFor(config.IP, contentType)
		if s.exceeded(tx, contentType, rule, "ip", ip) {
			return false, errors.New(locales.Get("errors.too_fast"))
		}
	}
	return review, nil
}

func rateLimitFor(config dto.PublishFrequencyConfig, contentType publishContentType) dto.PublishRateLimit {
	switch contentType {
	case publishContentTopic:
		return config.Topic
	case publishContentArticle:
		return config.Article
	default:
		return config.Comment
	}
}

func (s *antiAbuseService) exceeded(tx *gorm.DB, contentType publishContentType, rule dto.PublishRateLimit, field string, value interface{}) bool {
	if rule.DurationMinutes <= 0 || rule.MaxCount <= 0 {
		return false
	}
	var count int64
	windowStart := dates.NowTimestamp() - int64(rule.DurationMinutes)*60*1000
	query := tx.Where(field+" = ? AND create_time >= ?", value, windowStart)
	switch contentType {
	case publishContentTopic:
		query.Model(&models.Topic{}).Count(&count)
	case publishContentArticle:
		query.Model(&models.Article{}).Count(&count)
	case publishContentComment:
		query.Model(&models.Comment{}).Count(&count)
	}
	return count >= int64(rule.MaxCount)
}
