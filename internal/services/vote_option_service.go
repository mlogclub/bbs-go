package services

import (
	"bbs-go/internal/models"
	"bbs-go/internal/repositories"

	"github.com/mlogclub/simple/sqls"
	"github.com/mlogclub/simple/web/params"
)

var VoteOptionService = newVoteOptionService()

func newVoteOptionService() *voteOptionService {
	return &voteOptionService{}
}

type voteOptionService struct {
}

func (s *voteOptionService) Get(id int64) *models.VoteOption {
	return repositories.VoteOptionRepository.Get(sqls.DB(), id)
}

func (s *voteOptionService) Take(where ...interface{}) *models.VoteOption {
	return repositories.VoteOptionRepository.Take(sqls.DB(), where...)
}

func (s *voteOptionService) Find(cnd *sqls.Cnd) []models.VoteOption {
	return repositories.VoteOptionRepository.Find(sqls.DB(), cnd)
}

func (s *voteOptionService) FindOne(cnd *sqls.Cnd) *models.VoteOption {
	return repositories.VoteOptionRepository.FindOne(sqls.DB(), cnd)
}

func (s *voteOptionService) FindPageByParams(params *params.QueryParams) (list []models.VoteOption, paging *sqls.Paging) {
	return repositories.VoteOptionRepository.FindPageByParams(sqls.DB(), params)
}

func (s *voteOptionService) FindPageByCnd(cnd *sqls.Cnd) (list []models.VoteOption, paging *sqls.Paging) {
	return repositories.VoteOptionRepository.FindPageByCnd(sqls.DB(), cnd)
}

func (s *voteOptionService) Count(cnd *sqls.Cnd) int64 {
	return repositories.VoteOptionRepository.Count(sqls.DB(), cnd)
}

func (s *voteOptionService) Create(t *models.VoteOption) error {
	return repositories.VoteOptionRepository.Create(sqls.DB(), t)
}

func (s *voteOptionService) Update(t *models.VoteOption) error {
	return repositories.VoteOptionRepository.Update(sqls.DB(), t)
}

func (s *voteOptionService) Updates(id int64, columns map[string]interface{}) error {
	return repositories.VoteOptionRepository.Updates(sqls.DB(), id, columns)
}

func (s *voteOptionService) UpdateColumn(id int64, name string, value interface{}) error {
	return repositories.VoteOptionRepository.UpdateColumn(sqls.DB(), id, name, value)
}

func (s *voteOptionService) Delete(id int64) {
	repositories.VoteOptionRepository.Delete(sqls.DB(), id)
}

func (s *voteOptionService) FindByVoteId(voteId int64) []models.VoteOption {
	return repositories.VoteOptionRepository.Find(sqls.DB(), sqls.NewCnd().Eq("vote_id", voteId).Asc("sort_no").Asc("id"))
}
