package repositories

import (
	"errors"
	"github.com/jinzhu/gorm"
	"github.com/mlogclub/simple"
	"strings"

	"github.com/mlogclub/mlog/model"
)

type TagRepository struct {
}

func NewTagRepository() *TagRepository {
	return &TagRepository{}
}

func (this *TagRepository) Get(db *gorm.DB, id int64) *model.Tag {
	ret := &model.Tag{}
	if err := db.First(ret, "id = ?", id).Error; err != nil {
		return nil
	}
	return ret
}

func (this *TagRepository) Take(db *gorm.DB, where ...interface{}) *model.Tag {
	ret := &model.Tag{}
	if err := db.Take(ret, where...).Error; err != nil {
		return nil
	}
	return ret
}

func (this *TagRepository) QueryCnd(db *gorm.DB, cnd *simple.QueryCnd) (list []model.Tag, err error) {
	err = cnd.DoQuery(db).Find(&list).Error
	return
}

func (this *TagRepository) Query(db *gorm.DB, queries *simple.ParamQueries) (list []model.Tag, paging *simple.Paging) {
	queries.StartQuery(db).Find(&list)
	queries.StartCount(db).Model(&model.Tag{}).Count(&queries.Paging.Total)
	paging = queries.Paging
	return
}

func (this *TagRepository) Create(db *gorm.DB, t *model.Tag) (err error) {
	err = db.Create(t).Error
	return
}

func (this *TagRepository) Update(db *gorm.DB, t *model.Tag) (err error) {
	err = db.Save(t).Error
	return
}

func (this *TagRepository) Updates(db *gorm.DB, id int64, columns map[string]interface{}) (err error) {
	err = db.Model(&model.Tag{}).Where("id = ?", id).Updates(columns).Error
	return
}

func (this *TagRepository) UpdateColumn(db *gorm.DB, id int64, name string, value interface{}) (err error) {
	err = db.Model(&model.Tag{}).Where("id = ?", id).UpdateColumn(name, value).Error
	return
}

func (this *TagRepository) Delete(db *gorm.DB, id int64) {
	db.Model(&model.Tag{}).Delete("id", id)
}

func (this *TagRepository) GetTagInIds(tagIds []int64) []model.Tag {
	if len(tagIds) == 0 {
		return nil
	}
	var tags []model.Tag
	simple.GetDB().Where("id in (?)", tagIds).Find(&tags)
	return tags
}

func (this *TagRepository) FindByName(name string) *model.Tag {
	if len(name) == 0 {
		return nil
	}
	return this.Take(simple.GetDB(), "name = ?", name)
}

func (this *TagRepository) GetOrCreate(db *gorm.DB, name string) (*model.Tag, error) {
	if len(name) == 0 {
		return nil, errors.New("标签为空")
	}
	tag := this.FindByName(name)
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

func (this *TagRepository) GetOrCreates(db *gorm.DB, tags []string) (tagIds []int64) {
	for _, tagName := range tags {
		tagName = strings.TrimSpace(tagName)
		tag, err := this.GetOrCreate(db, tagName)
		if err == nil {
			tagIds = append(tagIds, tag.Id)
		}
	}
	return
}
