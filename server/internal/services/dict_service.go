package services

import (
	"bbs-go/internal/models"
	"bbs-go/internal/models/constants"
	"bbs-go/internal/repositories"

	"github.com/mlogclub/simple/sqls"
	"github.com/mlogclub/simple/web/params"
	"gorm.io/gorm"
)

var DictService = newDictService()

func newDictService() *dictService {
	return &dictService{}
}

type dictService struct {
}

func (s *dictService) Get(id int64) *models.Dict {
	return repositories.DictRepository.Get(sqls.DB(), id)
}

func (s *dictService) Take(where ...interface{}) *models.Dict {
	return repositories.DictRepository.Take(sqls.DB(), where...)
}

func (s *dictService) Find(cnd *sqls.Cnd) []models.Dict {
	return repositories.DictRepository.Find(sqls.DB(), cnd)
}

func (s *dictService) FindOne(cnd *sqls.Cnd) *models.Dict {
	return repositories.DictRepository.FindOne(sqls.DB(), cnd)
}

func (s *dictService) FindPageByParams(params *params.QueryParams) (list []models.Dict, paging *sqls.Paging) {
	return repositories.DictRepository.FindPageByParams(sqls.DB(), params)
}

func (s *dictService) FindPageByCnd(cnd *sqls.Cnd) (list []models.Dict, paging *sqls.Paging) {
	return repositories.DictRepository.FindPageByCnd(sqls.DB(), cnd)
}

func (s *dictService) Count(cnd *sqls.Cnd) int64 {
	return repositories.DictRepository.Count(sqls.DB(), cnd)
}

func (s *dictService) Create(t *models.Dict) error {
	return repositories.DictRepository.Create(sqls.DB(), t)
}

func (s *dictService) Update(t *models.Dict) error {
	return repositories.DictRepository.Update(sqls.DB(), t)
}

func (s *dictService) Updates(id int64, columns map[string]interface{}) error {
	return repositories.DictRepository.Updates(sqls.DB(), id, columns)
}

func (s *dictService) UpdateColumn(id int64, name string, value interface{}) error {
	return repositories.DictRepository.UpdateColumn(sqls.DB(), id, name, value)
}

func (s *dictService) Delete(id int64) {
	repositories.DictRepository.Delete(sqls.DB(), id)
}

func (s *dictService) GetNextSortNo() int {
	if max := s.FindOne(sqls.NewCnd().Desc("sort_no")); max != nil {
		return max.SortNo + 1
	}
	return 0
}

func (s *dictService) UpdateSort(ids []int64) error {
	return sqls.DB().Transaction(func(tx *gorm.DB) error {
		for i, id := range ids {
			if err := repositories.DictRepository.UpdateColumn(tx, id, "sort_no", i); err != nil {
				return err
			}
		}
		return nil
	})
}

func (s *dictService) FindByTypeId(typeId int64) []models.Dict {
	return s.Find(sqls.NewCnd().Eq("type_id", typeId).Eq("status", constants.StatusOk).Asc("sort_no").Desc("id"))
}

func (s *dictService) GetBy(typeId int64, name string) *models.Dict {
	return s.FindOne(sqls.NewCnd().Where("type_id = ? and name = ?", typeId, name))
}
