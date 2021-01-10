package services

import (
	"bbs-go/cache"
	"bbs-go/model"
	"bbs-go/model/constants"
	"bbs-go/repositories"
	"errors"
	"github.com/mlogclub/simple"
	"github.com/mlogclub/simple/date"
	"github.com/sirupsen/logrus"
	"strconv"
	"sync"
	"time"
)

var CheckInService = newCheckInService()

func newCheckInService() *checkInService {
	return &checkInService{}
}

type checkInService struct {
	m sync.Mutex
}

func (s *checkInService) Get(id int64) *model.CheckIn {
	return repositories.CheckInRepository.Get(simple.DB(), id)
}

func (s *checkInService) Take(where ...interface{}) *model.CheckIn {
	return repositories.CheckInRepository.Take(simple.DB(), where...)
}

func (s *checkInService) Find(cnd *simple.SqlCnd) []model.CheckIn {
	return repositories.CheckInRepository.Find(simple.DB(), cnd)
}

func (s *checkInService) FindOne(cnd *simple.SqlCnd) *model.CheckIn {
	return repositories.CheckInRepository.FindOne(simple.DB(), cnd)
}

func (s *checkInService) FindPageByParams(params *simple.QueryParams) (list []model.CheckIn, paging *simple.Paging) {
	return repositories.CheckInRepository.FindPageByParams(simple.DB(), params)
}

func (s *checkInService) FindPageByCnd(cnd *simple.SqlCnd) (list []model.CheckIn, paging *simple.Paging) {
	return repositories.CheckInRepository.FindPageByCnd(simple.DB(), cnd)
}

func (s *checkInService) Count(cnd *simple.SqlCnd) int64 {
	return repositories.CheckInRepository.Count(simple.DB(), cnd)
}

func (s *checkInService) Create(t *model.CheckIn) error {
	return repositories.CheckInRepository.Create(simple.DB(), t)
}

func (s *checkInService) Update(t *model.CheckIn) error {
	return repositories.CheckInRepository.Update(simple.DB(), t)
}

func (s *checkInService) Updates(id int64, columns map[string]interface{}) error {
	return repositories.CheckInRepository.Updates(simple.DB(), id, columns)
}

func (s *checkInService) UpdateColumn(id int64, name string, value interface{}) error {
	return repositories.CheckInRepository.UpdateColumn(simple.DB(), id, name, value)
}

func (s *checkInService) Delete(id int64) {
	repositories.CheckInRepository.Delete(simple.DB(), id)
}

func (s *checkInService) CheckIn(userId int64) error {
	s.m.Lock()
	defer s.m.Unlock()
	var (
		checkIn         = s.GetByUserId(userId)
		dayName         = date.GetDay(time.Now())
		yesterdayName   = date.GetDay(time.Now().Add(-time.Hour * 24))
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
			CreateTime:      date.NowTimestamp(),
			UpdateTime:      date.NowTimestamp(),
		})
	} else {
		checkIn.LatestDayName = dayName
		checkIn.ConsecutiveDays = consecutiveDays
		checkIn.UpdateTime = date.NowTimestamp()
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
	return s.FindOne(simple.NewSqlCnd().Eq("user_id", userId))
}
