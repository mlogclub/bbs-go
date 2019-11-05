package repositories

import (
	"errors"
	"strings"

	"github.com/jinzhu/gorm"
	"github.com/mlogclub/simple"

	"github.com/mlogclub/bbs-go/model"
)

var TagRepository = newTagRepository()

func newTagRepository() *tagRepository {
	return &tagRepository{}
}

type tagRepository struct {
}

func (this *tagRepository) Get(db *gorm.DB, id int64) *model.Tag {
	ret := &model.Tag{}
	if err := db.First(ret, "id = ?", id).Error; err != nil {
		return nil
	}
	return ret
}

func (this *tagRepository) Take(db *gorm.DB, where ...interface{}) *model.Tag {
	ret := &model.Tag{}
	if err := db.Take(ret, where...).Error; err != nil {
		return nil
	}
	return ret
}

func (this *tagRepository) Find(db *gorm.DB, cnd *simple.SqlCnd) (list []model.Tag) {
	cnd.Find(db, &list)
	return
}

func (this *tagRepository) FindOne(db *gorm.DB, cnd *simple.SqlCnd) (ret *model.Tag) {
	cnd.FindOne(db, &ret)
	return
}

func (this *tagRepository) FindPageByParams(db *gorm.DB, params *simple.QueryParams) (list []model.Tag, paging *simple.Paging) {
	return this.FindPageByCnd(db, &params.SqlCnd)
}

func (this *tagRepository) FindPageByCnd(db *gorm.DB, cnd *simple.SqlCnd) (list []model.Tag, paging *simple.Paging) {
	cnd.Find(db, &list)
	count := cnd.Count(db, &model.Tag{})

	paging = &simple.Paging{
		Page:  cnd.Paging.Page,
		Limit: cnd.Paging.Limit,
		Total: count,
	}
	return
}

func (this *tagRepository) Create(db *gorm.DB, t *model.Tag) (err error) {
	err = db.Create(t).Error
	return
}

func (this *tagRepository) Update(db *gorm.DB, t *model.Tag) (err error) {
	err = db.Save(t).Error
	return
}

func (this *tagRepository) Updates(db *gorm.DB, id int64, columns map[string]interface{}) (err error) {
	err = db.Model(&model.Tag{}).Where("id = ?", id).Updates(columns).Error
	return
}

func (this *tagRepository) UpdateColumn(db *gorm.DB, id int64, name string, value interface{}) (err error) {
	err = db.Model(&model.Tag{}).Where("id = ?", id).UpdateColumn(name, value).Error
	return
}

func (this *tagRepository) Delete(db *gorm.DB, id int64) {
	db.Delete(&model.Tag{}, "id = ?", id)
}

func (this *tagRepository) GetTagInIds(tagIds []int64) []model.Tag {
	if len(tagIds) == 0 {
		return nil
	}
	var tags []model.Tag
	simple.DB().Where("id in (?)", tagIds).Find(&tags)
	return tags
}

func (this *tagRepository) GetByName(name string) *model.Tag {
	if len(name) == 0 {
		return nil
	}
	return this.Take(simple.DB(), "name = ?", name)
}

func (this *tagRepository) GetOrCreate(db *gorm.DB, name string) (*model.Tag, error) {
	if len(name) == 0 {
		return nil, errors.New("标签为空")
	}
	tag := this.GetByName(name)
	if tag != nil {
		return tag, nil
	} else {
		tag = &model.Tag{
			Name:       name,
			Status:     model.TagStatusOk,
			CreateTime: simple.NowTimestamp(),
			UpdateTime: simple.NowTimestamp(),
		}
		err := this.Create(db, tag)
		if err != nil {
			return nil, err
		}
		return tag, nil
	}
}

func (this *tagRepository) GetOrCreates(db *gorm.DB, tags []string) (tagIds []int64) {
	for _, tagName := range tags {
		tagName = strings.TrimSpace(tagName)
		tag, err := this.GetOrCreate(db, tagName)
		if err == nil {
			tagIds = append(tagIds, tag.Id)
		}
	}
	return
}
