package services

import (
	"bbs-go/internal/models"
	"bbs-go/internal/models/constants"
	"bbs-go/internal/repositories"
	"fmt"
	"strconv"
	"testing"
	"time"

	"github.com/mlogclub/simple/sqls"

	// "gorm.io/driver/sqlite" // Sqlite driver based on CGO
	"github.com/glebarez/sqlite" // Pure go SQLite driver, checkout https://github.com/glebarez/sqlite for details

	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"gorm.io/gorm/schema"
)

func setupTestDB(t *testing.T) *gorm.DB {
	t.Helper()

	dsn := fmt.Sprintf("file:task_engine_test_%d?mode=memory&cache=shared&_fk=1", time.Now().UnixNano())
	db, err := gorm.Open(sqlite.Open(dsn), &gorm.Config{
		NamingStrategy: schema.NamingStrategy{
			TablePrefix:   "t_",
			SingularTable: true,
		},
		Logger: logger.Default.LogMode(logger.Silent),
	})
	if err != nil {
		t.Fatalf("open sqlite db: %v", err)
	}

	sqlDB, err := db.DB()
	if err != nil {
		t.Fatalf("get sql db: %v", err)
	}
	sqlDB.SetMaxOpenConns(1)
	sqlDB.SetMaxIdleConns(1)
	t.Cleanup(func() {
		_ = sqlDB.Close()
	})

	sqls.SetDB(db)

	if err := db.AutoMigrate(
		&models.User{},
		&models.TaskConfig{},
		&models.UserTaskEvent{},
		&models.UserTaskLog{},
		&models.UserScoreLog{},
		&models.UserExpLog{},
		&models.UserBadge{},
		&models.LevelConfig{},
	); err != nil {
		t.Fatalf("auto migrate: %v", err)
	}
	return db
}

func mustCreateUser(t *testing.T, now int64) *models.User {
	t.Helper()

	u := &models.User{
		Nickname:   "u",
		Status:     constants.StatusOk,
		CreateTime: now,
		UpdateTime: now,
	}
	if err := repositories.UserRepository.Create(sqls.DB(), u); err != nil {
		t.Fatalf("create user: %v", err)
	}
	return u
}

func mustCreateTaskConfig(t *testing.T, cfg *models.TaskConfig) *models.TaskConfig {
	t.Helper()

	if err := repositories.TaskConfigRepository.Create(sqls.DB(), cfg); err != nil {
		t.Fatalf("create task config: %v", err)
	}
	return cfg
}

func TestTaskEngineService_HandleUserEvent_MultiTaskConfigs(t *testing.T) {
	db := setupTestDB(t)
	now := time.Date(2025, 1, 2, 9, 0, 0, 0, time.UTC).UnixMilli()
	user := mustCreateUser(t, now)

	mustCreateTaskConfig(t, &models.TaskConfig{
		EventType:      constants.TaskEventTypeCheckIn,
		Title:          "t1",
		Description:    "t1",
		Period:         0,
		MaxFinishCount: 1,
		EventCount:     1,
		SortNo:         10,
		Status:         constants.StatusOk,
		CreateTime:     now,
		UpdateTime:     now,
	})
	mustCreateTaskConfig(t, &models.TaskConfig{
		EventType:      constants.TaskEventTypeCheckIn,
		Title:          "t2",
		Description:    "t2",
		Period:         0,
		MaxFinishCount: 1,
		EventCount:     1,
		SortNo:         20,
		Status:         constants.StatusOk,
		CreateTime:     now,
		UpdateTime:     now,
	})

	TaskEngineService.HandleUserEvent(user.Id, constants.TaskEventTypeCheckIn, now)

	var logs []models.UserTaskLog
	if err := db.Find(&logs).Error; err != nil {
		t.Fatalf("find logs: %v", err)
	}
	if len(logs) != 2 {
		t.Fatalf("expected 2 logs, got %d", len(logs))
	}

	var events []models.UserTaskEvent
	if err := db.Find(&events).Error; err != nil {
		t.Fatalf("find events: %v", err)
	}
	if len(events) != 2 {
		t.Fatalf("expected 2 events, got %d", len(events))
	}
}

func TestTaskEngineService_HandleUserEvent_EventCountAndMaxFinish(t *testing.T) {
	db := setupTestDB(t)
	now := time.Date(2025, 1, 2, 9, 0, 0, 0, time.UTC).UnixMilli()
	user := mustCreateUser(t, now)

	cfg := mustCreateTaskConfig(t, &models.TaskConfig{
		EventType:      constants.TaskEventTypeCommentCreate,
		Title:          "acc_3",
		Description:    "acc_3",
		Period:         0,
		MaxFinishCount: 1,
		EventCount:     3,
		SortNo:         10,
		Status:         constants.StatusOk,
		CreateTime:     now,
		UpdateTime:     now,
	})

	for i := 0; i < 2; i++ {
		TaskEngineService.HandleUserEvent(user.Id, constants.TaskEventTypeCommentCreate, now+int64(i))
	}

	var logCount int64
	if err := db.Model(&models.UserTaskLog{}).Count(&logCount).Error; err != nil {
		t.Fatalf("count logs: %v", err)
	}
	if logCount != 0 {
		t.Fatalf("expected 0 logs, got %d", logCount)
	}

	ev := repositories.UserTaskEventRepository.Take(db, "user_id = ? AND period_key = ? AND task_id = ?", user.Id, 0, cfg.Id)
	if ev == nil {
		t.Fatalf("expected user task event")
	}
	if ev.EventCount != 2 || ev.TaskFinishCount != 0 {
		t.Fatalf("expected eventCount=2 finishCount=0, got eventCount=%d finishCount=%d", ev.EventCount, ev.TaskFinishCount)
	}

	TaskEngineService.HandleUserEvent(user.Id, constants.TaskEventTypeCommentCreate, now+10)

	if err := db.Model(&models.UserTaskLog{}).Count(&logCount).Error; err != nil {
		t.Fatalf("count logs: %v", err)
	}
	if logCount != 1 {
		t.Fatalf("expected 1 log, got %d", logCount)
	}

	ev = repositories.UserTaskEventRepository.Take(db, "user_id = ? AND period_key = ? AND task_id = ?", user.Id, 0, cfg.Id)
	if ev.EventCount != 0 || ev.TaskFinishCount != 1 {
		t.Fatalf("expected eventCount=0 finishCount=1, got eventCount=%d finishCount=%d", ev.EventCount, ev.TaskFinishCount)
	}

	TaskEngineService.HandleUserEvent(user.Id, constants.TaskEventTypeCommentCreate, now+11)
	if err := db.Model(&models.UserTaskLog{}).Count(&logCount).Error; err != nil {
		t.Fatalf("count logs: %v", err)
	}
	if logCount != 1 {
		t.Fatalf("expected still 1 log, got %d", logCount)
	}
}

func TestTaskEngineService_HandleUserEvent_DailyPeriodKey(t *testing.T) {
	db := setupTestDB(t)

	day1 := time.Date(2025, 1, 2, 9, 0, 0, 0, time.UTC).UnixMilli()
	day2 := time.Date(2025, 1, 3, 9, 0, 0, 0, time.UTC).UnixMilli()
	user := mustCreateUser(t, day1)

	cfg := mustCreateTaskConfig(t, &models.TaskConfig{
		EventType:      constants.TaskEventTypeLikeCreate,
		Title:          "daily_like",
		Description:    "daily_like",
		Period:         1,
		MaxFinishCount: 1,
		EventCount:     1,
		SortNo:         10,
		Status:         constants.StatusOk,
		CreateTime:     day1,
		UpdateTime:     day1,
	})

	TaskEngineService.HandleUserEvent(user.Id, constants.TaskEventTypeLikeCreate, day1)
	TaskEngineService.HandleUserEvent(user.Id, constants.TaskEventTypeLikeCreate, day1+1)
	TaskEngineService.HandleUserEvent(user.Id, constants.TaskEventTypeLikeCreate, day2)

	var logCount int64
	if err := db.Model(&models.UserTaskLog{}).Count(&logCount).Error; err != nil {
		t.Fatalf("count logs: %v", err)
	}
	if logCount != 2 {
		t.Fatalf("expected 2 logs, got %d", logCount)
	}

	var eventCount int64
	if err := db.Model(&models.UserTaskEvent{}).Count(&eventCount).Error; err != nil {
		t.Fatalf("count events: %v", err)
	}
	if eventCount != 2 {
		t.Fatalf("expected 2 events, got %d", eventCount)
	}

	pk1, err := strconv.Atoi(TaskEngineService.GetPeriodKey(constants.TaskPeriodDaily, day1))
	if err != nil {
		t.Fatalf("parse pk1: %v", err)
	}
	pk2, err := strconv.Atoi(TaskEngineService.GetPeriodKey(constants.TaskPeriodDaily, day2))
	if err != nil {
		t.Fatalf("parse pk2: %v", err)
	}
	if pk1 == pk2 {
		t.Fatalf("expected different period keys, got pk1=%d pk2=%d", pk1, pk2)
	}

	if repositories.UserTaskEventRepository.Take(db, "user_id = ? AND period_key = ? AND task_id = ?", user.Id, pk1, cfg.Id) == nil {
		t.Fatalf("expected event row for day1")
	}
	if repositories.UserTaskEventRepository.Take(db, "user_id = ? AND period_key = ? AND task_id = ?", user.Id, pk2, cfg.Id) == nil {
		t.Fatalf("expected event row for day2")
	}
}

func TestTaskEngineService_HandleUserEvent_TaskActiveGuards(t *testing.T) {
	db := setupTestDB(t)
	now := time.Date(2025, 1, 2, 9, 0, 0, 0, time.UTC).UnixMilli()
	user := mustCreateUser(t, now)

	mustCreateTaskConfig(t, &models.TaskConfig{
		EventType:      "evt.inactive.status",
		Title:          "inactive_status",
		Description:    "inactive_status",
		Period:         0,
		MaxFinishCount: 1,
		EventCount:     1,
		SortNo:         10,
		Status:         constants.StatusDeleted,
		CreateTime:     now,
		UpdateTime:     now,
	})
	mustCreateTaskConfig(t, &models.TaskConfig{
		EventType:      "evt.inactive.future",
		Title:          "inactive_future",
		Description:    "inactive_future",
		Period:         0,
		MaxFinishCount: 1,
		EventCount:     1,
		SortNo:         10,
		Status:         constants.StatusOk,
		StartTime:      now + 3600_000,
		CreateTime:     now,
		UpdateTime:     now,
	})
	mustCreateTaskConfig(t, &models.TaskConfig{
		EventType:      "evt.inactive.past",
		Title:          "inactive_past",
		Description:    "inactive_past",
		Period:         0,
		MaxFinishCount: 1,
		EventCount:     1,
		SortNo:         10,
		Status:         constants.StatusOk,
		EndTime:        now - 1,
		CreateTime:     now,
		UpdateTime:     now,
	})

	TaskEngineService.HandleUserEvent(user.Id, "evt.inactive.status", now)
	TaskEngineService.HandleUserEvent(user.Id, "evt.inactive.future", now)
	TaskEngineService.HandleUserEvent(user.Id, "evt.inactive.past", now)

	var logCount int64
	if err := db.Model(&models.UserTaskLog{}).Count(&logCount).Error; err != nil {
		t.Fatalf("count logs: %v", err)
	}
	if logCount != 0 {
		t.Fatalf("expected 0 logs, got %d", logCount)
	}
}

func TestTaskEngineService_HandleUserEvent_GrantReward(t *testing.T) {
	db := setupTestDB(t)
	now := time.Date(2025, 1, 2, 9, 0, 0, 0, time.UTC).UnixMilli()
	user := mustCreateUser(t, now)

	mustCreateTaskConfig(t, &models.TaskConfig{
		EventType:      constants.TaskEventTypeTopicCreate,
		Title:          "reward",
		Description:    "reward",
		Score:          5,
		Exp:            10,
		BadgeId:        3,
		Period:         0,
		MaxFinishCount: 1,
		EventCount:     1,
		SortNo:         10,
		Status:         constants.StatusOk,
		CreateTime:     now,
		UpdateTime:     now,
	})

	TaskEngineService.HandleUserEvent(user.Id, constants.TaskEventTypeTopicCreate, now)

	var u models.User
	if err := db.First(&u, "id = ?", user.Id).Error; err != nil {
		t.Fatalf("get user: %v", err)
	}
	if u.Score != 5 {
		t.Fatalf("expected score=5, got %d", u.Score)
	}
	if u.Exp != 10 {
		t.Fatalf("expected exp=10, got %d", u.Exp)
	}

	var scoreLogs int64
	if err := db.Model(&models.UserScoreLog{}).Count(&scoreLogs).Error; err != nil {
		t.Fatalf("count user score logs: %v", err)
	}
	if scoreLogs != 1 {
		t.Fatalf("expected 1 user score log, got %d", scoreLogs)
	}

	var expLogs int64
	if err := db.Model(&models.UserExpLog{}).Count(&expLogs).Error; err != nil {
		t.Fatalf("count user exp logs: %v", err)
	}
	if expLogs != 1 {
		t.Fatalf("expected 1 user exp log, got %d", expLogs)
	}

	var badges int64
	if err := db.Model(&models.UserBadge{}).Count(&badges).Error; err != nil {
		t.Fatalf("count user badges: %v", err)
	}
	if badges != 1 {
		t.Fatalf("expected 1 user badge, got %d", badges)
	}
}

func TestTaskEngineService_HandleUserEvent_DuplicateLogIgnored(t *testing.T) {
	db := setupTestDB(t)
	now := time.Date(2025, 1, 2, 9, 0, 0, 0, time.UTC).UnixMilli()
	user := mustCreateUser(t, now)

	cfg := mustCreateTaskConfig(t, &models.TaskConfig{
		EventType:      "evt.dup",
		Title:          "dup",
		Description:    "dup",
		Period:         0,
		MaxFinishCount: 2,
		EventCount:     1,
		SortNo:         10,
		Status:         constants.StatusOk,
		CreateTime:     now,
		UpdateTime:     now,
	})

	// Create a mismatched state that can happen under concurrency:
	// log exists (finishNo=1) but UserTaskEvent has not recorded the completion.
	if err := repositories.UserTaskEventRepository.Create(db, &models.UserTaskEvent{
		UserId:          user.Id,
		PeriodKey:       0,
		TaskId:          cfg.Id,
		EventCount:      0,
		TaskFinishCount: 0,
		CreateTime:      now,
		UpdateTime:      now,
	}); err != nil {
		t.Fatalf("create user task event: %v", err)
	}
	if err := repositories.UserTaskLogRepository.Create(db, &models.UserTaskLog{
		UserId:     user.Id,
		PeriodKey:  0,
		TaskId:     cfg.Id,
		FinishNo:   1,
		CreateTime: now,
		UpdateTime: now,
	}); err != nil {
		t.Fatalf("create user task log: %v", err)
	}

	TaskEngineService.HandleUserEvent(user.Id, "evt.dup", now)

	var logCount int64
	if err := db.Model(&models.UserTaskLog{}).Count(&logCount).Error; err != nil {
		t.Fatalf("count logs: %v", err)
	}
	if logCount != 1 {
		t.Fatalf("expected 1 log after duplicate insert ignored, got %d", logCount)
	}

	TaskEngineService.HandleUserEvent(user.Id, "evt.dup", now+1)
	if err := db.Model(&models.UserTaskLog{}).Count(&logCount).Error; err != nil {
		t.Fatalf("count logs: %v", err)
	}
	if logCount != 2 {
		t.Fatalf("expected 2 logs after next completion, got %d", logCount)
	}
}
