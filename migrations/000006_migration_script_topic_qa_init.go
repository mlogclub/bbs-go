package migrations

import (
	"bbs-go/internal/models"
	"bbs-go/internal/models/constants"

	"github.com/mlogclub/simple/sqls"
)

func migrate_topic_qa_init() error {
	return sqls.WithTransaction(func(txCtx *sqls.TxContext) error {
		tx := txCtx.Tx

		if err := tx.Model(&models.Topic{}).
			Where("qa_status = '' OR qa_status IS NULL").
			Updates(map[string]interface{}{
				"qa_status":           constants.QaStatusUnsolved,
				"accepted_comment_id": 0,
				"solved_at":           0,
			}).Error; err != nil {
			return err
		}

		if err := tx.Model(&models.TopicNode{}).
			Where("type = '' OR type IS NULL").
			Update("type", constants.TopicNodeTypeNormal).Error; err != nil {
			return err
		}

		return nil
	})
}
