package heatpoints

import (
	"bbs-go/internal/models"
	"bbs-go/internal/models/constants"
	"log/slog"
	"math/rand"
	"time"

	"github.com/mlogclub/simple/common/dates"
	"github.com/mlogclub/simple/sqls"
)

// SnapshotService 快照服务
type SnapshotService struct{}

var HeatSnapshot = &SnapshotService{}

// TakeAllSnapshots 执行所有快照
func (s *SnapshotService) TakeAllSnapshots() error {
	slog.Info("开始执行热度点快照")

	// 1. 生成交互快照（含质押总量）
	s.takeTopicSnapshots()

	// 2. 生成流通快照
	s.takeCirculationSnapshot()

	// 3. 生成火焰偏移
	s.generateFlameOffset()

	slog.Info("热度点快照完成")
	return nil
}

// takeTopicSnapshots 为所有正常帖子生成互动快照（分批处理）
func (s *SnapshotService) takeTopicSnapshots() {
	today := time.Now().Format("20060102")

	// 批量查询每个帖子的质押总量（一条 SQL 替代 N+1）
	type topicStakeTotal struct {
		TopicId    int64
		StakeTotal int
	}
	var stakeTotals []topicStakeTotal
	sqls.DB().Model(&models.TopicStake{}).
		Where("status = ?", constants.StakeStatusActive).
		Group("topic_id").
		Select("topic_id, COALESCE(SUM(heat_points), 0) as stake_total").
		Scan(&stakeTotals)

	// 转为 map 方便查找
	stakeMap := make(map[int64]int)
	for _, st := range stakeTotals {
		stakeMap[st.TopicId] = st.StakeTotal
	}

	// 批量查询所有帖子的今日点赞数据（去重用户数 + 总点赞数）
	type likeRow struct {
		EntityId    int64
		UniqueUsers int
		TotalCount  int
	}
	var likeRows []likeRow
	sqls.DB().Model(&models.UserLike{}).
		Where("entity_type = ?", "topic").
		Group("entity_id").
		Select("entity_id, COUNT(DISTINCT user_id) as unique_users, COUNT(*) as total_count").
		Scan(&likeRows)

	likeMap := make(map[int64]likeRow)
	for _, r := range likeRows {
		likeMap[r.EntityId] = r
	}

	// 批量查询所有帖子的评论数据（去重用户数 + 总评论数 + 有效评论数）
	type commentRow struct {
		EntityId         int64
		UniqueCommenters int
		TotalComments    int
		ValidComments    int
	}
	var commentRows []commentRow
	sqls.DB().Model(&models.Comment{}).
		Where("entity_type = ?", "topic").
		Group("entity_id").
		Select("entity_id, COUNT(DISTINCT user_id) as unique_commenters, COUNT(*) as total_comments, SUM(CASE WHEN LENGTH(content) >= 5 THEN 1 ELSE 0 END) as valid_comments").
		Scan(&commentRows)

	commentMap := make(map[int64]commentRow)
	for _, r := range commentRows {
		commentMap[r.EntityId] = r
	}

	// 分批处理帖子（每批 500 条，游标式扫描）
	var lastId int64 = 0
	for {
		var topics []models.Topic
		sqls.DB().Where("status = ? AND id > ?", constants.StatusOk, lastId).
			Order("id asc").Limit(500).Find(&topics)

		if len(topics) == 0 {
			break
		}

		for _, topic := range topics {
			lastId = topic.Id

			stakeTotal := stakeMap[topic.Id]
			likeData := likeMap[topic.Id]
			commentData := commentMap[topic.Id]

			snapshot := models.TopicInteractionSnapshot{
				TopicId:          topic.Id,
				SnapshotDate:     today,
				UniqueLikers:     likeData.UniqueUsers,
				UniqueCommenters: commentData.UniqueCommenters,
				ValidComments:    commentData.ValidComments,
				TotalLikes:       likeData.TotalCount,
				TotalComments:    commentData.TotalComments,
				StakeTotal:       stakeTotal,
				CreateTime:       dates.NowTimestamp(),
			}

			// FirstOrCreate 保证幂等（重复运行不报错）
			sqls.DB().Where("topic_id = ? AND snapshot_date = ?", topic.Id, today).
				FirstOrCreate(&snapshot)
		}

		if len(topics) < 500 {
			break
		}
	}
}

// takeCirculationSnapshot 生成活跃流通快照（修复 N+1 查询）
func (s *SnapshotService) takeCirculationSnapshot() {
	today := time.Now().Format("20060102")

	// 一次性查询所有用户的质押总量（替代逐用户查询）
	type userStakeRow struct {
		UserId int64
		Total  int
	}
	var userStakes []userStakeRow
	sqls.DB().Model(&models.TopicStake{}).
		Where("status = ?", constants.StakeStatusActive).
		Group("user_id").
		Select("user_id, COALESCE(SUM(heat_points), 0) as total").
		Scan(&userStakes)

	stakeMap := make(map[int64]int)
	for _, row := range userStakes {
		stakeMap[row.UserId] = row.Total
	}

	// 查询所有用户统计
	var stats []models.UserHeatStats
	sqls.DB().Find(&stats)

	totalSupply := 0
	activeCirculation := 0
	stakedTotal := 0
	activeUserCount := 0

	for _, stat := range stats {
		userStaked := stakeMap[stat.UserId]
		totalSupply += stat.TotalPoints
		stakedTotal += userStaked

		// 活跃贡献 = max(当前质押中, 近7天累计质押量)
		activeContribution := userStaked
		if stat.StakedInWindow > activeContribution {
			activeContribution = stat.StakedInWindow
		}
		activeCirculation += activeContribution

		if stat.StakedInWindow > 0 || userStaked > 0 {
			activeUserCount++
		}
	}

	snapshot := models.HeatCirculationSnapshot{
		SnapshotDate:      today,
		TotalSupply:       totalSupply,
		ActiveCirculation: activeCirculation,
		StakedTotal:       stakedTotal,
		ActiveUserCount:   activeUserCount,
		CreateTime:        dates.NowTimestamp(),
	}
	sqls.DB().Where("snapshot_date = ?", today).FirstOrCreate(&snapshot)
}

// generateFlameOffset 生成火焰偏移量（每个边界独立生成 + 周级平滑）
func (s *SnapshotService) generateFlameOffset() {
	today := time.Now().Format("20060102")

	var existing models.DailyFlameOffset
	if err := sqls.DB().Where("date = ?", today).First(&existing).Error; err == nil {
		return // 已生成
	}

	// 周级种子：同周内基准偏移相近，跨周独立
	year, week := time.Now().ISOWeek()
	weekSeed := int64(year*100 + week)
	dayOfWeek := int(time.Now().Weekday())

	// 每个边界独立生成偏移量
	generateOffset := func(boundaryIdx int64) float64 {
		rng := rand.New(rand.NewSource(weekSeed + boundaryIdx*31))
		base := 0.85 + rng.Float64()*0.30 // [0.85, 1.15]

		// 日内微调 ±5%，平滑周内跳变
		jitter := (float64(dayOfWeek) - 3.0) / 60.0 // [-0.05, +0.05]
		result := base + jitter

		// 钳位到 [0.80, 1.20]
		if result < 0.80 {
			result = 0.80
		}
		if result > 1.20 {
			result = 1.20
		}
		return result
	}

	offset := models.DailyFlameOffset{
		Date:          today,
		Phase12Offset: generateOffset(1),
		Phase23Offset: generateOffset(2),
		Flame2Offset:  generateOffset(3),
		Flame3Offset:  generateOffset(4),
		Flame4Offset:  generateOffset(5),
		Flame5Offset:  generateOffset(6),
		CreateTime:    dates.NowTimestamp(),
	}
	sqls.DB().Create(&offset)
}

// GetActiveCirculation 获取活跃流通量
func (s *SnapshotService) GetActiveCirculation() int {
	yesterday := time.Now().AddDate(0, 0, -1).Format("20060102")
	var snapshot models.HeatCirculationSnapshot
	if err := sqls.DB().Where("snapshot_date = ?", yesterday).First(&snapshot).Error; err != nil {
		return constants.HeatPointsDefaultCirculation
	}
	if snapshot.ActiveCirculation <= 0 {
		return constants.HeatPointsDefaultCirculation
	}
	return snapshot.ActiveCirculation
}
