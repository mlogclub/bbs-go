package services

import (
	"bbs-go/cache"
	"bbs-go/model"
	"bbs-go/model/constants"
	"bbs-go/repositories"
	"errors"
	"strconv"
	"sync"
	"time"

	"github.com/mlogclub/simple/common/dates"
	"github.com/mlogclub/simple/sqls"
	"github.com/mlogclub/simple/web/params"
	"github.com/sirupsen/logrus"
)

var CheckInService = newCheckInService()

func newCheckInService() *checkInService {
	return &checkInService{}
}

type checkInService struct {
	m sync.Mutex
}

func (s *checkInService) Get(id int64) *model.CheckIn {
	return repositories.CheckInRepository.Get(sqls.DB(), id)
}

func (s *checkInService) Take(where ...interface{}) *model.CheckIn {
	return repositories.CheckInRepository.Take(sqls.DB(), where...)
}

func (s *checkInService) Find(cnd *sqls.Cnd) []model.CheckIn {
	return repositories.CheckInRepository.Find(sqls.DB(), cnd)
}

func (s *checkInService) FindOne(cnd *sqls.Cnd) *model.CheckIn {
	return repositories.CheckInRepository.FindOne(sqls.DB(), cnd)
}

func (s *checkInService) FindPageByParams(params *params.QueryParams) (list []model.CheckIn, paging *sqls.Paging) {
	return repositories.CheckInRepository.FindPageByParams(sqls.DB(), params)
}

func (s *checkInService) FindPageByCnd(cnd *sqls.Cnd) (list []model.CheckIn, paging *sqls.Paging) {
	return repositories.CheckInRepository.FindPageByCnd(sqls.DB(), cnd)
}

func (s *checkInService) Count(cnd *sqls.Cnd) int64 {
	return repositories.CheckInRepository.Count(sqls.DB(), cnd)
}

func (s *checkInService) Create(t *model.CheckIn) error {
	return repositories.CheckInRepository.Create(sqls.DB(), t)
}

func (s *checkInService) Update(t *model.CheckIn) error {
	return repositories.CheckInRepository.Update(sqls.DB(), t)
}

func (s *checkInService) Updates(id int64, columns map[string]interface{}) error {
	return repositories.CheckInRepository.Updates(sqls.DB(), id, columns)
}

func (s *checkInService) UpdateColumn(id int64, name string, value interface{}) error {
	return repositories.CheckInRepository.UpdateColumn(sqls.DB(), id, name, value)
}

func (s *checkInService) Delete(id int64) {
	repositories.CheckInRepository.Delete(sqls.DB(), id)
}

func (s *checkInService) CheckIn(userId int64) error {
	s.m.Lock()
	defer s.m.Unlock()
	var (
		checkIn         = s.GetByUserId(userId)
		dayName         = dates.GetDay(time.Now())
		yesterdayName   = dates.GetDay(time.Now().Add(-time.Hour * 24))
		consecutiveDays = 1
		err             error
	)

	if checkIn != nil && checkIn.LatestDayName == dayName {
		return errors.New("你已签到")
	}

	if checkIn != nil && checkIn.LatestDayName == yesterdayName {
		consecutiveDays = checkIn.ConsecutiveDays + 1
	}

	if checkIn == nil {
		err = s.Create(&model.CheckIn{
			Model:           model.Model{},
			UserId:          userId,
			LatestDayName:   dayName,
			ConsecutiveDays: consecutiveDays,
			CreateTime:      dates.NowTimestamp(),
			UpdateTime:      dates.NowTimestamp(),
		})
	} else {
		checkIn.LatestDayName = dayName
		checkIn.ConsecutiveDays = consecutiveDays
		checkIn.UpdateTime = dates.NowTimestamp()
		err = s.Update(checkIn)
	}
	if err == nil {
		// 清理签到排行榜缓存
		cache.UserCache.RefreshCheckInRank()
		// 处理签到积分
		config := SysConfigService.GetConfig()
		if config.ScoreConfig.CheckInScore > 0 {
			_ = UserService.IncrScore(userId, config.ScoreConfig.CheckInScore, constants.EntityCheckIn,
				strconv.FormatInt(userId, 10), "签到"+strconv.Itoa(dayName))
		} else {
			logrus.Warn("签到积分未配置...")
		}
	}
	return err
}

func (s *checkInService) GetByUserId(userId int64) *model.CheckIn {
	return s.FindOne(sqls.NewCnd().Eq("user_id", userId))
}
