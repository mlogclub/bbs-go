package migrations

import (
	"bbs-go/internal/models"
	"fmt"
	"strings"
	"testing"
	"time"

	"github.com/glebarez/sqlite"
	"github.com/mlogclub/simple/sqls"
	"gorm.io/gorm"
	"gorm.io/gorm/schema"
)

func TestOIDCIdentityHashMigration(t *testing.T) {
	db, err := gorm.Open(sqlite.Open(fmt.Sprintf("file:oidc_hash_%d?mode=memory&cache=shared", time.Now().UnixNano())), &gorm.Config{NamingStrategy: schema.NamingStrategy{TablePrefix: "t_", SingularTable: true}})
	if err != nil {
		t.Fatal(err)
	}
	conn, _ := db.DB()
	t.Cleanup(func() { _ = conn.Close() })
	sqls.SetDB(db)
	// Model the deployed schema before the hash column existed.
	if err := db.Exec(`CREATE TABLE t_o_id_c_identity (id INTEGER PRIMARY KEY, provider_key TEXT, issuer TEXT, subject TEXT, user_id INTEGER)`).Error; err != nil {
		t.Fatal(err)
	}
	issuer := "https://example.com/" + strings.Repeat("a", 480)
	prefix := strings.Repeat("s", 500)
	for i, subject := range []string{prefix + "A", prefix + "a"} {
		if err := db.Exec("INSERT INTO t_o_id_c_identity (id,provider_key,issuer,subject,user_id) VALUES (?,?,?,?,?)", i+1, "provider", issuer, subject, i+1).Error; err != nil {
			t.Fatal(err)
		}
	}
	if err := db.AutoMigrate(&models.OIDCIdentity{}); err != nil {
		t.Fatal(err)
	}
	if err := migrate_oidc_identity_hash(); err != nil {
		t.Fatal(err)
	}
	if err := migrate_oidc_identity_hash(); err != nil {
		t.Fatalf("migration is not idempotent: %v", err)
	}
	var rows []models.OIDCIdentity
	if err := db.Order("id").Find(&rows).Error; err != nil {
		t.Fatal(err)
	}
	if len(rows) != 2 || rows[0].IdentityHash == rows[1].IdentityHash {
		t.Fatal("long case-sensitive identities were not preserved")
	}
	for _, row := range rows {
		if row.IdentityHash != models.OIDCIdentityHash(issuer, row.Subject) {
			t.Fatal("incorrect backfill")
		}
	}
	duplicate := rows[0]
	duplicate.Id = 0
	duplicate.UserId = 99
	duplicate.IdentityHash = ""
	if err := db.Create(&duplicate).Error; err == nil {
		t.Fatal("duplicate identity was accepted")
	}
	fresh := models.OIDCIdentity{ProviderKey: "provider", Issuer: issuer, Subject: prefix + "B", UserId: 3}
	if err := db.Create(&fresh).Error; err != nil {
		t.Fatal(err)
	}
	if len(fresh.IdentityHash) != 64 {
		t.Fatal("create hook did not compute the hash")
	}
	// The new composite index must fit MySQL's 3072-byte utf8mb4 limit.
	stmt := &gorm.Statement{DB: db}
	if err := stmt.Parse(&models.OIDCIdentity{}); err != nil {
		t.Fatal(err)
	}
	bytes := 0
	for _, index := range stmt.Schema.ParseIndexes() {
		if index.Name != "uk_oidc_identity_hash" {
			continue
		}
		for _, field := range index.Fields {
			bytes += field.Field.Size * 4
		}
	}
	if bytes == 0 || bytes > 3072 {
		t.Fatalf("invalid MySQL index size: %d", bytes)
	}
}
