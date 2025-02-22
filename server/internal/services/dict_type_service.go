package services

import (
	"bbs-go/internal/models"
	"bbs-go/internal/repositories"

	"github.com/mlogclub/simple/sqls"
	"github.com/mlogclub/simple/web/params"
)

var DictTypeService = newDictTypeService()

func newDictTypeService() *dictTypeService {
	return &dictTypeService{}
}

type dictTypeService struct {
}

func (s *dictTypeService) Get(id int64) *models.DictType {
	return repositories.DictTypeRepository.Get(sqls.DB(), id)
}

func (s *dictTypeService) Take(where ...interface{}) *models.DictType {
	return repositories.DictTypeRepository.Take(sqls.DB(), where...)
}

func (s *dictTypeService) Find(cnd *sqls.Cnd) []models.DictType {
	return repositories.DictTypeRepository.Find(sqls.DB(), cnd)
}

func (s *dictTypeService) FindOne(cnd *sqls.Cnd) *models.DictType {
	return repositories.DictTypeRepository.FindOne(sqls.DB(), cnd)
}

func (s *dictTypeService) FindPageByParams(params *params.QueryParams) (list []models.DictType, paging *sqls.Paging) {
	return repositories.DictTypeRepository.FindPageByParams(sqls.DB(), params)
}

func (s *dictTypeService) FindPageByCnd(cnd *sqls.Cnd) (list []models.DictType, paging *sqls.Paging) {
	return repositories.DictTypeRepository.FindPageByCnd(sqls.DB(), cnd)
}

func (s *dictTypeService) Count(cnd *sqls.Cnd) int64 {
	return repositories.DictTypeRepository.Count(sqls.DB(), cnd)
}

func (s *dictTypeService) Create(t *models.DictType) error {
	return repositories.DictTypeRepository.Create(sqls.DB(), t)
}

func (s *dictTypeService) Update(t *models.DictType) error {
	return repositories.DictTypeRepository.Update(sqls.DB(), t)
}

func (s *dictTypeService) Updates(id int64, columns map[string]interface{}) error {
	return repositories.DictTypeRepository.Updates(sqls.DB(), id, columns)
}

func (s *dictTypeService) UpdateColumn(id int64, name string, value interface{}) error {
	return repositories.DictTypeRepository.UpdateColumn(sqls.DB(), id, name, value)
}

func (s *dictTypeService) Delete(id int64) {
	repositories.DictTypeRepository.Delete(sqls.DB(), id)
}

func (s *dictTypeService) GetByCode(code string) *models.DictType {
	return s.FindOne(sqls.NewCnd().Eq("code", code))
}
