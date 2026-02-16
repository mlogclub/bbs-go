package services

import (
	"bbs-go/internal/cache"
	"bbs-go/internal/models"
	"bbs-go/internal/repositories"

	"github.com/mlogclub/simple/sqls"
	"github.com/mlogclub/simple/web/params"
)

var BadgeService = newBadgeService()

func newBadgeService() *badgeService {
	return &badgeService{}
}

type badgeService struct {
}

func (s *badgeService) Get(id int64) *models.Badge {
	return repositories.BadgeRepository.Get(sqls.DB(), id)
}

func (s *badgeService) Take(where ...interface{}) *models.Badge {
	return repositories.BadgeRepository.Take(sqls.DB(), where...)
}

func (s *badgeService) Find(cnd *sqls.Cnd) []models.Badge {
	return repositories.BadgeRepository.Find(sqls.DB(), cnd)
}

func (s *badgeService) FindOne(cnd *sqls.Cnd) *models.Badge {
	return repositories.BadgeRepository.FindOne(sqls.DB(), cnd)
}

func (s *badgeService) FindPageByParams(params *params.QueryParams) (list []models.Badge, paging *sqls.Paging) {
	return repositories.BadgeRepository.FindPageByParams(sqls.DB(), params)
}

func (s *badgeService) FindPageByCnd(cnd *sqls.Cnd) (list []models.Badge, paging *sqls.Paging) {
	return repositories.BadgeRepository.FindPageByCnd(sqls.DB(), cnd)
}

func (s *badgeService) Count(cnd *sqls.Cnd) int64 {
	return repositories.BadgeRepository.Count(sqls.DB(), cnd)
}

func (s *badgeService) Create(t *models.Badge) error {
	if err := repositories.BadgeRepository.Create(sqls.DB(), t); err != nil {
		return err
	}
	cache.BadgeCache.Reload()
	return nil
}

func (s *badgeService) Update(t *models.Badge) error {
	if err := repositories.BadgeRepository.Update(sqls.DB(), t); err != nil {
		return err
	}
	cache.BadgeCache.Reload()
	return nil
}

func (s *badgeService) Updates(id int64, columns map[string]interface{}) error {
	if err := repositories.BadgeRepository.Updates(sqls.DB(), id, columns); err != nil {
		return err
	}
	cache.BadgeCache.Reload()
	return nil
}
