package models

import (
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"gorm.io/gorm"
)

// OIDCIdentityHash indexes the full, case-sensitive identity without putting
// two 512-character utf8mb4 columns into a MySQL unique index.
func OIDCIdentityHash(issuer, subject string) string {
	encoded, _ := json.Marshal([2]string{issuer, subject})
	digest := sha256.Sum256(encoded)
	return hex.EncodeToString(digest[:])
}

func (identity *OIDCIdentity) BeforeCreate(_ *gorm.DB) error {
	identity.IdentityHash = OIDCIdentityHash(identity.Issuer, identity.Subject)
	return nil
}
