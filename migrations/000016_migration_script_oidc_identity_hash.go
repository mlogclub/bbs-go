package migrations

import (
	"bbs-go/internal/models"
	"github.com/mlogclub/simple/sqls"
	"gorm.io/gorm"
)

// The nullable hash column lets AutoMigrate add the new index before existing
// identities are backfilled. Old issuer/subject values are retained in full.
func migrate_oidc_identity_hash() error {
	return sqls.DB().Transaction(func(tx *gorm.DB) error {
		var identities []models.OIDCIdentity
		return tx.Where("identity_hash IS NULL OR identity_hash = ?", "").FindInBatches(&identities, 100, func(batch *gorm.DB, _ int) error {
			for _, identity := range identities {
				if err := batch.Model(&models.OIDCIdentity{}).Where("id = ?", identity.Id).UpdateColumn("identity_hash", models.OIDCIdentityHash(identity.Issuer, identity.Subject)).Error; err != nil {
					return err
				}
			}
			return nil
		}).Error
	})
}
