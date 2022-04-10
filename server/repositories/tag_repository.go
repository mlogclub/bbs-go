package repositories

import (
	"bbs-go/model/constants"
	"errors"
	"strings"

	"github.com/mlogclub/simple/common/dates"
	"github.com/mlogclub/simple/sqls"
	"github.com/mlogclub/simple/web/params"
	"gorm.io/gorm"

	"bbs-go/model"
)

var TagRepository = newTagRepository()

func newTagRepository() *tagRepository {
	return &tagRepository{}
}

type tagRepository struct {
}

func (r *tagRepository) Get(db *gorm.DB, id int64) *model.Tag {
	ret := &model.Tag{}
	if err := db.First(ret, "id = ?", id).Error; err != nil {
		return nil
	}
	return ret
}

func (r *tagRepository) Take(db *gorm.DB, where ...interface{}) *model.Tag {
	ret := &model.Tag{}
	if err := db.Take(ret, where...).Error; err != nil {
		return nil
	}
	return ret
}

func (r *tagRepository) Find(db *gorm.DB, cnd *sqls.Cnd) (list []model.Tag) {
	cnd.Find(db, &list)
	return
}

func (r *tagRepository) FindOne(db *gorm.DB, cnd *sqls.Cnd) *model.Tag {
	ret := &model.Tag{}
	if err := cnd.FindOne(db, &ret); err != nil {
		return nil
	}
	return ret
}

func (r *tagRepository) FindPageByParams(db *gorm.DB, params *params.QueryParams) (list []model.Tag, paging *sqls.Paging) {
	return r.FindPageByCnd(db, &params.Cnd)
}

func (r *tagRepository) FindPageByCnd(db *gorm.DB, cnd *sqls.Cnd) (list []model.Tag, paging *sqls.Paging) {
	cnd.Find(db, &list)
	count := cnd.Count(db, &model.Tag{})

	paging = &sqls.Paging{
		Page:  cnd.Paging.Page,
		Limit: cnd.Paging.Limit,
		Total: count,
	}
	return
}

func (r *tagRepository) Create(db *gorm.DB, t *model.Tag) (err error) {
	err = db.Create(t).Error
	return
}

func (r *tagRepository) Update(db *gorm.DB, t *model.Tag) (err error) {
	err = db.Save(t).Error
	return
}

func (r *tagRepository) Updates(db *gorm.DB, id int64, columns map[string]interface{}) (err error) {
	err = db.Model(&model.Tag{}).Where("id = ?", id).Updates(columns).Error
	return
}

func (r *tagRepository) UpdateColumn(db *gorm.DB, id int64, name string, value interface{}) (err error) {
	err = db.Model(&model.Tag{}).Where("id = ?", id).UpdateColumn(name, value).Error
	return
}

func (r *tagRepository) Delete(db *gorm.DB, id int64) {
	db.Delete(&model.Tag{}, "id = ?", id)
}

func (r *tagRepository) GetTagInIds(tagIds []int64) []model.Tag {
	if len(tagIds) == 0 {
		return nil
	}
	var tags []model.Tag
	sqls.DB().Where("id in (?)", tagIds).Find(&tags)
	return tags
}

func (r *tagRepository) GetByName(name string) *model.Tag {
	if len(name) == 0 {
		return nil
	}
	return r.Take(sqls.DB(), "name = ?", name)
}

func (r *tagRepository) GetOrCreate(db *gorm.DB, name string) (*model.Tag, error) {
	if len(name) == 0 {
		return nil, errors.New("标签为空")
	}
	tag := r.GetByName(name)
	if tag != nil {
		return tag, nil
	} else {
		tag = &model.Tag{
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

func (r *tagRepository) GetOrCreates(db *gorm.DB, tags []string) (tagIds []int64) {
	for _, tagName := range tags {
		tagName = strings.TrimSpace(tagName)
		tag, err := r.GetOrCreate(db, tagName)
		if err == nil {
			tagIds = append(tagIds, tag.Id)
		}
	}
	return
}
