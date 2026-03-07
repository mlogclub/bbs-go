package migrations

import (
	"bbs-go/internal/models"
	"bbs-go/internal/models/constants"
	"bbs-go/internal/repositories"

	"github.com/mlogclub/simple/common/dates"
	"github.com/mlogclub/simple/sqls"
)

func migrate_qa_bounty_config() error {
	return sqls.WithTransaction(func(txCtx *sqls.TxContext) error {
		tx := txCtx.Tx
		now := dates.NowTimestamp()

		seeds := []struct {
			Key         string
			Value       string
			Name        string
			Description string
		}{
			{constants.SysConfigEnableQaBounty, "true", "开启问答悬赏", "开启后发表问答帖时可设置悬赏积分"},
			{constants.SysConfigQaBountyMin, "0", "悬赏积分下限", "问答悬赏积分下限，0 表示不校验"},
			{constants.SysConfigQaBountyMax, "0", "悬赏积分上限", "问答悬赏积分上限，0 表示不校验"},
			{constants.SysConfigQaBountyRequired, "false", "问答必填悬赏", "是否要求问答帖必须设置悬赏"},
		}

		for _, s := range seeds {
			existing := repositories.SysConfigRepository.GetByKey(tx, s.Key)
			if existing != nil {
				continue
			}
			cfg := &models.SysConfig{
				Key:         s.Key,
				Value:       s.Value,
				Name:        s.Name,
				Description: s.Description,
				CreateTime:  now,
				UpdateTime:  now,
			}
			if err := repositories.SysConfigRepository.Create(tx, cfg); err != nil {
				return err
			}
		}
		return nil
	})
}
