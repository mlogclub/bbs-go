package services

import (
	. "bbs-go/base"
	"bbs-go/model"
	"bbs-go/model/constants"
	"bbs-go/pkg/event"
	"bbs-go/repositories"

	"github.com/mlogclub/simple/sqls"
	"github.com/mlogclub/simple/web/params"
	"gorm.io/gorm"
)

var PhService = newPhService()

func newPhService() *phService {
	return &phService{}
}

type phService struct {
}

func (s *phService) Get(id int64) *model.PurchaseHistory {
	return repositories.PhRepository.Get(sqls.DB(), id)
}

func (s *phService) Take(where ...interface{}) *model.PurchaseHistory {
	return repositories.PhRepository.Take(sqls.DB(), where...)
}

func (s *phService) Find(cnd *sqls.Cnd) []model.PurchaseHistory {
	return repositories.PhRepository.Find(sqls.DB(), cnd)
}

func (s *phService) FindOne(cnd *sqls.Cnd) *model.PurchaseHistory {
	return repositories.PhRepository.FindOne(sqls.DB(), cnd)
}

func (s *phService) FindPageByParams(params *params.QueryParams) (list []model.PurchaseHistory, paging *sqls.Paging) {
	return repositories.PhRepository.FindPageByParams(sqls.DB(), params)
}

func (s *phService) FindPageByCnd(cnd *sqls.Cnd) (list []model.PurchaseHistory, paging *sqls.Paging) {
	return repositories.PhRepository.FindPageByCnd(sqls.DB(), cnd)
}

func (s *phService) BuyTopic(user *model.User, topic *model.Topic) error {

	ph := model.PurchaseHistory{
		BuyId:  topic.Id,
		UserId: user.Id,
		Score:  topic.Score,
	}

	if err := sqls.DB().Transaction(func(tx *gorm.DB) error {

		// 购买记录
		if err := repositories.PhRepository.Create(tx, &ph); err != nil {
			return err
		}

		// 更新积分
		if err := repositories.UserRepository.Update(tx, user); err != nil {
			return err
		}

		// 更新发帖人积分
		err := UserService.IncrScore(topic.UserId, GetScore(topic.Score), "pay", user.Nickname, "隐藏内容售出")
		if err != nil {
			return err
		}
		return nil
	}); err != nil {
		return err
	}

	event.Send(event.TopicBuyEvent{
		UserId:       user.Id,
		ToUserId:     topic.UserId,
		QuoteContent: topic.Title,
	})

	return nil
}

func (s *phService) Update(t *model.PurchaseHistory) error {
	err := repositories.PhRepository.Update(sqls.DB(), t)
	return err
}

func (s *phService) Updates(id int64, columns map[string]interface{}) error {
	err := repositories.PhRepository.Updates(sqls.DB(), id, columns)
	return err
}

func (s *phService) UpdateColumn(id int64, name string, value interface{}) error {
	err := repositories.PhRepository.UpdateColumn(sqls.DB(), id, name, value)
	return err
}

func (s *phService) Delete(id int64) error {
	err := repositories.PhRepository.UpdateColumn(sqls.DB(), id, "status", constants.StatusDeleted)
	if err == nil {
		// 删掉标签文章
		ArticleTagService.DeleteByArticleId(id)
	}
	return err
}

func (s *phService) IsBuy(userId int64, buyId int64) bool {
	return s.FindOne(sqls.NewCnd().Where("user_id = ? and buy_id = ?", userId, buyId)) != nil
}
