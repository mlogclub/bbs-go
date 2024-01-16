package repositories

import (
	"bbs-go/internal/models/constants"
	"errors"
	"strings"

	"github.com/mlogclub/simple/common/dates"
	"github.com/mlogclub/simple/sqls"
	"github.com/mlogclub/simple/web/params"
	"gorm.io/gorm"

	"bbs-go/internal/models"
)

var TagRepository = newTagRepository()

func newTagRepository() *tagRepository {
	return &tagRepository{}
}

type tagRepository struct {
}

func (r *tagRepository) Get(db *gorm.DB, id int64) *models.Tag {
	ret := &models.Tag{}
	if err := db.First(ret, "id = ?", id).Error; err != nil {
		return nil
	}
	return ret
}

func (r *tagRepository) Take(db *gorm.DB, where ...interface{}) *models.Tag {
	ret := &models.Tag{}
	if err := db.Take(ret, where...).Error; err != nil {
		return nil
	}
	return ret
}

func (r *tagRepository) Find(db *gorm.DB, cnd *sqls.Cnd) (list []models.Tag) {
	cnd.Find(db, &list)
	return
}

func (r *tagRepository) FindOne(db *gorm.DB, cnd *sqls.Cnd) *models.Tag {
	ret := &models.Tag{}
	if err := cnd.FindOne(db, &ret); err != nil {
		return nil
	}
	return ret
}

func (r *tagRepository) FindPageByParams(db *gorm.DB, params *params.QueryParams) (list []models.Tag, paging *sqls.Paging) {
	return r.FindPageByCnd(db, &params.Cnd)
}

func (r *tagRepository) FindPageByCnd(db *gorm.DB, cnd *sqls.Cnd) (list []models.Tag, paging *sqls.Paging) {
	cnd.Find(db, &list)
	count := cnd.Count(db, &models.Tag{})

	paging = &sqls.Paging{
		Page:  cnd.Paging.Page,
		Limit: cnd.Paging.Limit,
		Total: count,
	}
	return
}

func (r *tagRepository) Create(db *gorm.DB, t *models.Tag) (err error) {
	err = db.Create(t).Error
	return
}

func (r *tagRepository) Update(db *gorm.DB, t *models.Tag) (err error) {
	err = db.Save(t).Error
	return
}

func (r *tagRepository) Updates(db *gorm.DB, id int64, columns map[string]interface{}) (err error) {
	err = db.Model(&models.Tag{}).Where("id = ?", id).Updates(columns).Error
	return
}

func (r *tagRepository) UpdateColumn(db *gorm.DB, id int64, name string, value interface{}) (err error) {
	err = db.Model(&models.Tag{}).Where("id = ?", id).UpdateColumn(name, value).Error
	return
}

func (r *tagRepository) Delete(db *gorm.DB, id int64) {
	db.Delete(&models.Tag{}, "id = ?", id)
}

func (r *tagRepository) GetTagInIds(tagIds []int64) []models.Tag {
	if len(tagIds) == 0 {
		return nil
	}
	var tags []models.Tag
	sqls.DB().Where("id in (?)", tagIds).Find(&tags)
	return tags
}

func (r *tagRepository) GetByName(name string) *models.Tag {
	if len(name) == 0 {
		return nil
	}
	return r.Take(sqls.DB(), "name = ?", name)
}

func (r *tagRepository) GetOrCreate(db *gorm.DB, name string) (*models.Tag, error) {
	if len(name) == 0 {
		return nil, errors.New("标签为空")
	}
	tag := r.GetByName(name)
	if tag != nil {
		return tag, nil
	} else {
		tag = &models.Tag{
			Name:       name,
			Status:     constants.StatusOk,
			CreateTime: dates.NowTimestamp(),
			UpdateTime: dates.NowTimestamp(),
		}
		err := r.Create(db, tag)
		if err != nil {
			return nil, err
		}
		return tag, nil
	}
}

func (r *tagRepository) GetOrCreates(db *gorm.DB, tags []string) (tagIds []int64, err error) {
	for _, tagName := range tags {
		var tag *models.Tag
		tag, err = r.GetOrCreate(db, strings.TrimSpace(tagName))
		if err != nil {
			return
		}
		tagIds = append(tagIds, tag.Id)
	}
	return
}
