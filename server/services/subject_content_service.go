package services

import (
	"github.com/jinzhu/gorm"
	"github.com/mlogclub/simple"
	"github.com/sirupsen/logrus"

	"github.com/mlogclub/bbs-go/model"
	"github.com/mlogclub/bbs-go/repositories"
)

var SubjectContentService = newSubjectContentService()

func newSubjectContentService() *subjectContentService {
	return &subjectContentService{}
}

type subjectContentService struct {
}

func (this *subjectContentService) Get(id int64) *model.SubjectContent {
	return repositories.SubjectContentRepository.Get(simple.DB(), id)
}

func (this *subjectContentService) Take(where ...interface{}) *model.SubjectContent {
	return repositories.SubjectContentRepository.Take(simple.DB(), where...)
}

func (this *subjectContentService) Find(cnd *simple.SqlCnd) []model.SubjectContent {
	return repositories.SubjectContentRepository.Find(simple.DB(), cnd)
}

func (this *subjectContentService) FindOne(db *gorm.DB, cnd *simple.SqlCnd) (ret *model.SubjectContent) {
	cnd.FindOne(db, &ret)
	return
}

func (this *subjectContentService) FindPageByParams(params *simple.QueryParams) (list []model.SubjectContent, paging *simple.Paging) {
	return repositories.SubjectContentRepository.FindPageByParams(simple.DB(), params)
}

func (this *subjectContentService) FindPageByCnd(cnd *simple.SqlCnd) (list []model.SubjectContent, paging *simple.Paging) {
	return repositories.SubjectContentRepository.FindPageByCnd(simple.DB(), cnd)
}

func (this *subjectContentService) DeleteByEntity(entityType string, entityId int64) {
	t := this.GetByEntity(entityType, entityId)
	if t != nil {
		this.Delete(t.Id)
	}
}

func (this *subjectContentService) Delete(id int64) {
	err := repositories.SubjectContentRepository.UpdateColumn(simple.DB(), id, "deleted", true)
	if err != nil {
		logrus.Error(err)
	}
}

func (this *subjectContentService) GetByEntity(entityType string, entityId int64) *model.SubjectContent {
	return repositories.SubjectContentRepository.Take(simple.DB(), "entity_type = ? and entity_id = ?", entityType, entityId)
}

func (this *subjectContentService) GetSubjectContents(subjectId, cursor int64) (contents []model.SubjectContent, nextCursor int64) {
	cnd := simple.NewSqlCnd().Desc("id").Limit(20)
	if subjectId > 0 {
		cnd.Eq("subject_id", subjectId)
	}
	if cursor > 0 {
		cnd.Lt("id", cursor)
	}
	contents = repositories.SubjectContentRepository.Find(simple.DB(), cnd)
	if len(contents) > 0 {
		nextCursor = contents[len(contents)-1].Id
	} else {
		nextCursor = cursor
	}
	return
}

// 发布
func (this *subjectContentService) Publish(subjectId int64, entityType string, entityId int64, title, summary string) (c *model.SubjectContent, err error) {
	c = this.GetByEntity(entityType, entityId)
	if c != nil {
		c.SubjectId = subjectId
		c.EntityType = entityType
		c.EntityId = entityId
		c.Title = title
		c.Summary = summary
		c.Deleted = false
		c.CreateTime = simple.NowTimestamp()
		err = repositories.SubjectContentRepository.Update(simple.DB(), c)
	} else {
		c := &model.SubjectContent{
			SubjectId:  subjectId,
			EntityType: entityType,
			EntityId:   entityId,
			Title:      title,
			Summary:    summary,
			Deleted:    false,
			CreateTime: simple.NowTimestamp(),
		}
		err = repositories.SubjectContentRepository.Create(simple.DB(), c)
	}
	return
}
