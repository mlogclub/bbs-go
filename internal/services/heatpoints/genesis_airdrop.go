package heatpoints

import (
	"bbs-go/internal/cache"
	"bbs-go/internal/models"
	"bbs-go/internal/models/constants"
	"log/slog"

	"github.com/mlogclub/simple/common/dates"
	"github.com/mlogclub/simple/sqls"
	"gorm.io/gorm"
)

// GenesisAirdropService 创世空投服务
type GenesisAirdropService struct{}

var Airdrop = &GenesisAirdropService{}

// Execute 批量执行空投：给所有尚未收到空投的正常用户发放
func (s *GenesisAirdropService) Execute() error {
	// 找出所有已收到空投的用户 ID
	var airdroppedUserIds []int64
	sqls.DB().Model(&models.UserHeatLog{}).
		Where("change_type = ?", constants.HeatMintTypeGenesisAirdrop).
		Pluck("user_id", &airdroppedUserIds)

	// 找出尚未收到空投的正常用户
	query := sqls.DB().Where("status = ?", constants.StatusOk)
	if len(airdroppedUserIds) > 0 {
		query = query.Where("id NOT IN ?", airdroppedUserIds)
	}

	var users []models.User
	query.Find(&users)

	if len(users) == 0 {
		return nil // 所有用户都已空投
	}

	slog.Info("创世空投开始", "new_user_count", len(users))

	return sqls.DB().Transaction(func(tx *gorm.DB) error {
		for _, user := range users {
			if err := s.airdropUser(tx, user.Id); err != nil {
				return err
			}
		}

		mintLog := models.SystemMintLog{
			MintType:    constants.HeatMintTypeGenesisAirdrop,
			Amount:      len(users) * constants.HeatPointsGenesisAirdrop,
			RecipientId: 0,
			Remark:      "创世空投（增量）",
			CreateTime:  dates.NowTimestamp(),
		}
		if err := tx.Create(&mintLog).Error; err != nil {
			return err
		}

		slog.Info("创世空投完成", "total_minted", len(users)*constants.HeatPointsGenesisAirdrop)
		return nil
	})
}

// AirdropSingleUser 给单个用户发放空投（注册时调用）
func (s *GenesisAirdropService) AirdropSingleUser(userId int64) error {
	// 检查该用户是否已收到
	var count int64
	sqls.DB().Model(&models.UserHeatLog{}).
		Where("user_id = ? AND change_type = ?", userId, constants.HeatMintTypeGenesisAirdrop).
		Count(&count)
	if count > 0 {
		return nil // 已空投
	}

	return sqls.DB().Transaction(func(tx *gorm.DB) error {
		return s.airdropUser(tx, userId)
	})
}

// airdropUser 在事务内给单个用户发放空投
func (s *GenesisAirdropService) airdropUser(tx *gorm.DB, userId int64) error {
	// 发放热度点
	if err := tx.Model(&models.User{}).Where("id = ?", userId).
		Update("heat_points", constants.HeatPointsGenesisAirdrop).Error; err != nil {
		return err
	}

	cache.UserCache.Invalidate(userId)

	// 初始化统计
	heatStats := models.UserHeatStats{
		UserId:      userId,
		TotalPoints: constants.HeatPointsGenesisAirdrop,
		UpdateTime:  dates.NowTimestamp(),
	}
	if err := tx.Where("user_id = ?", userId).FirstOrCreate(&heatStats).Error; err != nil {
		tx.Model(&models.UserHeatStats{}).Where("user_id = ?", userId).
			Update("total_points", constants.HeatPointsGenesisAirdrop)
	}

	// 记录流水
	heatLog := models.UserHeatLog{
		UserId:     userId,
		ChangeType: constants.HeatMintTypeGenesisAirdrop,
		Amount:     constants.HeatPointsGenesisAirdrop,
		Balance:    constants.HeatPointsGenesisAirdrop,
		RefId:      "genesis",
		Remark:     "创世空投",
		CreateTime: dates.NowTimestamp(),
	}
	return tx.Create(&heatLog).Error
}
