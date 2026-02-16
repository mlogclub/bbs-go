package services

import (
	"bbs-go/internal/cache"
	"bbs-go/internal/models"
	"bbs-go/internal/repositories"

	"github.com/mlogclub/simple/common/dates"
	"github.com/mlogclub/simple/sqls"
	"github.com/mlogclub/simple/web/params"
)

var UserBadgeService = newUserBadgeService()

func newUserBadgeService() *userBadgeService {
	return &userBadgeService{}
}

type userBadgeService struct {
}

func (s *userBadgeService) Get(id int64) *models.UserBadge {
	return repositories.UserBadgeRepository.Get(sqls.DB(), id)
}

func (s *userBadgeService) Take(where ...interface{}) *models.UserBadge {
	return repositories.UserBadgeRepository.Take(sqls.DB(), where...)
}

func (s *userBadgeService) Find(cnd *sqls.Cnd) []models.UserBadge {
	return repositories.UserBadgeRepository.Find(sqls.DB(), cnd)
}

func (s *userBadgeService) FindOne(cnd *sqls.Cnd) *models.UserBadge {
	return repositories.UserBadgeRepository.FindOne(sqls.DB(), cnd)
}

func (s *userBadgeService) FindPageByParams(params *params.QueryParams) (list []models.UserBadge, paging *sqls.Paging) {
	return repositories.UserBadgeRepository.FindPageByParams(sqls.DB(), params)
}

func (s *userBadgeService) FindPageByCnd(cnd *sqls.Cnd) (list []models.UserBadge, paging *sqls.Paging) {
	return repositories.UserBadgeRepository.FindPageByCnd(sqls.DB(), cnd)
}

func (s *userBadgeService) Count(cnd *sqls.Cnd) int64 {
	return repositories.UserBadgeRepository.Count(sqls.DB(), cnd)
}

func (s *userBadgeService) Give(ctx *sqls.TxContext, userId int64, badgeId int64, sourceType string, sourceId string) error {
	if repositories.UserBadgeRepository.GetBy(ctx.Tx, userId, badgeId) != nil {
		return nil
	}
	err := repositories.UserBadgeRepository.Create(ctx.Tx, &models.UserBadge{
		UserId:     userId,
		BadgeId:    badgeId,
		SourceType: sourceType,
		SourceId:   sourceId,
		IsWorn:     false,
		SortNo:     0,
		CreateTime: dates.NowTimestamp(),
		UpdateTime: dates.NowTimestamp(),
	})
	if err != nil {
		return err
	}
	cache.UserBadgeCache.Invalidate(userId)
	return nil
}
