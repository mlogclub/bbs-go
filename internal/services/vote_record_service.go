package services

import (
	"bbs-go/internal/models"
	"bbs-go/internal/repositories"

	"github.com/mlogclub/simple/sqls"
	"github.com/mlogclub/simple/web/params"
)

var VoteRecordService = newVoteRecordService()

func newVoteRecordService() *voteRecordService {
	return &voteRecordService{}
}

type voteRecordService struct {
}

func (s *voteRecordService) Get(id int64) *models.VoteRecord {
	return repositories.VoteRecordRepository.Get(sqls.DB(), id)
}

func (s *voteRecordService) Take(where ...interface{}) *models.VoteRecord {
	return repositories.VoteRecordRepository.Take(sqls.DB(), where...)
}

func (s *voteRecordService) Find(cnd *sqls.Cnd) []models.VoteRecord {
	return repositories.VoteRecordRepository.Find(sqls.DB(), cnd)
}

func (s *voteRecordService) FindOne(cnd *sqls.Cnd) *models.VoteRecord {
	return repositories.VoteRecordRepository.FindOne(sqls.DB(), cnd)
}

func (s *voteRecordService) FindPageByParams(params *params.QueryParams) (list []models.VoteRecord, paging *sqls.Paging) {
	return repositories.VoteRecordRepository.FindPageByParams(sqls.DB(), params)
}

func (s *voteRecordService) FindPageByCnd(cnd *sqls.Cnd) (list []models.VoteRecord, paging *sqls.Paging) {
	return repositories.VoteRecordRepository.FindPageByCnd(sqls.DB(), cnd)
}

func (s *voteRecordService) Count(cnd *sqls.Cnd) int64 {
	return repositories.VoteRecordRepository.Count(sqls.DB(), cnd)
}

func (s *voteRecordService) Create(t *models.VoteRecord) error {
	return repositories.VoteRecordRepository.Create(sqls.DB(), t)
}

func (s *voteRecordService) Update(t *models.VoteRecord) error {
	return repositories.VoteRecordRepository.Update(sqls.DB(), t)
}

func (s *voteRecordService) Updates(id int64, columns map[string]interface{}) error {
	return repositories.VoteRecordRepository.Updates(sqls.DB(), id, columns)
}

func (s *voteRecordService) UpdateColumn(id int64, name string, value interface{}) error {
	return repositories.VoteRecordRepository.UpdateColumn(sqls.DB(), id, name, value)
}

func (s *voteRecordService) Delete(id int64) {
	repositories.VoteRecordRepository.Delete(sqls.DB(), id)
}

func (s *voteRecordService) GetBy(userId, voteId int64) *models.VoteRecord {
	return s.FindOne(sqls.NewCnd().Where("user_id = ? and vote_id = ?", userId, voteId))
}
