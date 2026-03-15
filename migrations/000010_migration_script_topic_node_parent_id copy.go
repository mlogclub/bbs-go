package migrations

import (
	"github.com/mlogclub/simple/sqls"

	"bbs-go/internal/models"
)

// migrate_topic_node_parent_id 为 TopicNode 设置 parent_id=0（表结构由 AutoMigrate 添加 parent_id 列）
func migrate_topic_node_parent_id() error {
	return sqls.WithTransaction(func(txCtx *sqls.TxContext) error {
		tx := txCtx.Tx
		// 确保已有 parent_id 列的旧数据为 0（AutoMigrate 已添加列，此步保证历史数据 parent_id=0）
		return tx.Model(&models.TopicNode{}).
			Where("parent_id IS NULL OR parent_id = 0").
			Updates(map[string]interface{}{"parent_id": 0}).Error
	})
}
