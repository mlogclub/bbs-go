package migrations

import (
	"github.com/mlogclub/simple/sqls"
)

// migrate_sync_user_topic_comment_counts 重建用户发帖/评论计数，使其与「已发布/正常」数据一致。
// 修复历史遗留的 topic_count / comment_count 漂移（删除/待审核未正确扣减或计入）。
func migrate_sync_user_topic_comment_counts() error {
	return sqls.WithTransaction(func(ctx *sqls.TxContext) error {
		if err := ctx.Tx.Exec(`
			UPDATE t_user SET topic_count = (
				SELECT COUNT(*) FROM t_topic
				WHERE t_topic.user_id = t_user.id AND t_topic.status = 0
			), comment_count = (
				SELECT COUNT(*) FROM t_comment
				WHERE t_comment.user_id = t_user.id AND t_comment.status = 0
			)
		`).Error; err != nil {
			return err
		}
		return nil
	})
}
