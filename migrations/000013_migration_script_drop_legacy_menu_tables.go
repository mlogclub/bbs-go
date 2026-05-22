package migrations

import "github.com/mlogclub/simple/sqls"

func migrate_drop_legacy_menu_tables() error {
	return sqls.WithTransaction(func(ctx *sqls.TxContext) error {
		if err := ctx.Tx.Exec("DROP TABLE IF EXISTS t_role_menu").Error; err != nil {
			return err
		}
		if err := ctx.Tx.Exec("DROP TABLE IF EXISTS t_menu").Error; err != nil {
			return err
		}
		return nil
	})
}
