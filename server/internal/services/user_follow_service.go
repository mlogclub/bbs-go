package services

import (
	"bbs-go/internal/cache"
	"bbs-go/internal/models"
	"bbs-go/internal/models/constants"
	"bbs-go/internal/pkg/event"
	"bbs-go/internal/repositories"

	"github.com/emirpasic/gods/sets/hashset"
	"github.com/mlogclub/simple/common/dates"
	"github.com/mlogclub/simple/sqls"
	"github.com/mlogclub/simple/web/params"
	"gorm.io/gorm"
)

var UserFollowService = newUserFollowService()

func newUserFollowService() *userFollowService {
	return &userFollowService{}
}

type userFollowService struct {
}

func (s *userFollowService) Get(id int64) *models.UserFollow {
	return repositories.UserFollowRepository.Get(sqls.DB(), id)
}

func (s *userFollowService) Take(where ...interface{}) *models.UserFollow {
	return repositories.UserFollowRepository.Take(sqls.DB(), where...)
}

func (s *userFollowService) Find(cnd *sqls.Cnd) []models.UserFollow {
	return repositories.UserFollowRepository.Find(sqls.DB(), cnd)
}

func (s *userFollowService) FindOne(cnd *sqls.Cnd) *models.UserFollow {
	return repositories.UserFollowRepository.FindOne(sqls.DB(), cnd)
}

func (s *userFollowService) FindPageByParams(params *params.QueryParams) (list []models.UserFollow, paging *sqls.Paging) {
	return repositories.UserFollowRepository.FindPageByParams(sqls.DB(), params)
}

func (s *userFollowService) FindPageByCnd(cnd *sqls.Cnd) (list []models.UserFollow, paging *sqls.Paging) {
	return repositories.UserFollowRepository.FindPageByCnd(sqls.DB(), cnd)
}

func (s *userFollowService) Count(cnd *sqls.Cnd) int64 {
	return repositories.UserFollowRepository.Count(sqls.DB(), cnd)
}

func (s *userFollowService) Create(t *models.UserFollow) error {
	return repositories.UserFollowRepository.Create(sqls.DB(), t)
}

func (s *userFollowService) Update(t *models.UserFollow) error {
	return repositories.UserFollowRepository.Update(sqls.DB(), t)
}

func (s *userFollowService) Updates(id int64, columns map[string]interface{}) error {
	return repositories.UserFollowRepository.Updates(sqls.DB(), id, columns)
}

func (s *userFollowService) UpdateColumn(id int64, name string, value interface{}) error {
	return repositories.UserFollowRepository.UpdateColumn(sqls.DB(), id, name, value)
}

func (s *userFollowService) Delete(id int64) {
	repositories.UserFollowRepository.Delete(sqls.DB(), id)
}

func (s *userFollowService) Follow(userId, otherId int64) error {
	if userId == otherId {
		// 自己关注自己，不进行处理。
		// return errors.New("自己不能关注自己")
		return nil
	}

	if s.IsFollowed(userId, otherId) {
		return nil
	}

	err := sqls.DB().Transaction(func(tx *gorm.DB) error {
		// 如果对方也关注了我，那么更新状态为互相关注
		otherFollowed := tx.Exec("update t_user_follow set status = ? where user_id = ? and other_id = ?",
			constants.FollowStatusBoth, otherId, userId).RowsAffected > 0
		status := constants.FollowStatusFollow
		if otherFollowed {
			status = constants.FollowStatusBoth
		}

		if err := repositories.UserRepository.Updates(tx, userId, map[string]interface{}{
			"follow_count": gorm.Expr("follow_count + 1"),
		}); err != nil {
			return err
		}
		cache.UserCache.Invalidate(userId)

		if err := repositories.UserRepository.Updates(tx, otherId, map[string]interface{}{
			"fans_count": gorm.Expr("fans_count + 1"),
		}); err != nil {
			return err
		}
		cache.UserCache.Invalidate(otherId)

		return repositories.UserFollowRepository.Create(tx, &models.UserFollow{
			UserId:     userId,
			OtherId:    otherId,
			Status:     status,
			CreateTime: dates.NowTimestamp(),
		})
	})
	if err != nil {
		return err
	}

	// 发送mq消息
	event.Send(event.FollowEvent{
		UserId:  userId,
		OtherId: otherId,
	})
	return nil
}

func (s *userFollowService) UnFollow(userId, otherId int64) error {
	if userId == otherId {
		// 自己关注自己，不进行处理。
		return nil
	}
	if !s.IsFollowed(userId, otherId) {
		return nil
	}
	err := sqls.DB().Transaction(func(tx *gorm.DB) error {
		success := tx.Where("user_id = ? and other_id = ?", userId, otherId).Delete(models.UserFollow{}).RowsAffected > 0
		if success {
			tx.Exec("update t_user_follow set status = ? where user_id = ? and other_id = ?",
				constants.FollowStatusFollow, otherId, userId)
		}

		if err := tx.Model(&models.User{}).Where("id = ? and follow_count > 0", userId).Updates(map[string]interface{}{
			"follow_count": gorm.Expr("follow_count - 1"),
		}).Error; err != nil {
			return err
		}
		cache.UserCache.Invalidate(userId)

		if err := tx.Model(&models.User{}).Where("id = ? and fans_count > 0", otherId).Updates(map[string]interface{}{
			"fans_count": gorm.Expr("fans_count - 1"),
		}).Error; err != nil {
			return err
		}
		cache.UserCache.Invalidate(otherId)

		return nil
	})
	if err != nil {
		return err
	}

	// 发送mq消息
	event.Send(event.UnFollowEvent{
		UserId:  userId,
		OtherId: otherId,
	})
	return nil
}

// GetFans 粉丝列表
func (s *userFollowService) GetFans(userId int64, cursor int64, limit int) (itemList []int64, nextCursor int64, hasMore bool) {
	cnd := sqls.NewCnd().Eq("other_id", userId)
	if cursor > 0 {
		cnd.Lt("id", cursor)
	}
	cnd.Desc("id").Limit(limit)
	list := repositories.UserFollowRepository.Find(sqls.DB(), cnd)

	if len(list) > 0 {
		nextCursor = list[len(list)-1].Id
		hasMore = len(list) >= limit
		for _, e := range list {
			itemList = append(itemList, e.UserId)
		}
	} else {
		nextCursor = cursor
	}
	return
}

// GetFollows 关注列表
func (s *userFollowService) GetFollows(userId int64, cursor int64, limit int) (itemList []int64, nextCursor int64, hasMore bool) {
	cnd := sqls.NewCnd().Eq("user_id", userId)
	if cursor > 0 {
		cnd.Lt("id", cursor)
	}
	cnd.Desc("id").Limit(limit)
	list := repositories.UserFollowRepository.Find(sqls.DB(), cnd)

	if len(list) > 0 {
		nextCursor = list[len(list)-1].Id
		hasMore = len(list) >= limit
		for _, e := range list {
			itemList = append(itemList, e.OtherId)
		}
	} else {
		nextCursor = cursor
	}
	return
}

// ScanFans 扫描粉丝
func (s *userFollowService) ScanFans(userId int64, handle func(fansId int64)) {
	var cursor int64 = 0
	for {
		list := s.Find(sqls.NewCnd().Eq("other_id", userId).Gt("id", cursor).Asc("id").Limit(100))
		if len(list) == 0 {
			break
		}
		cursor = list[len(list)-1].Id
		for _, item := range list {
			handle(item.UserId)
		}
	}
}

// ScanFollowed 扫描关注的用户
func (s *userFollowService) ScanFollowed(userId int64, handle func(followUserId int64)) {
	var cursor int64 = 0
	for {
		list := s.Find(sqls.NewCnd().Eq("user_id", userId).Gt("id", cursor).Asc("id").Limit(100))
		if len(list) == 0 {
			break
		}
		cursor = list[len(list)-1].Id
		for _, item := range list {
			handle(item.OtherId)
		}
	}
}

func (s *userFollowService) IsFollowed(userId, otherId int64) bool {
	if userId == otherId {
		return false
	}
	set := s.IsFollowedUsers(userId, otherId)
	return set.Contains(otherId)
}

func (s *userFollowService) IsFollowedUsers(userId int64, otherIds ...int64) hashset.Set {
	set := hashset.New()
	list := s.Find(sqls.NewCnd().Eq("user_id", userId).In("other_id", otherIds))
	for _, follow := range list {
		set.Add(follow.OtherId)
	}
	return *set
}
