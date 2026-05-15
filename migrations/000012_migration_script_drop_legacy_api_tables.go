package migrations

import "github.com/mlogclub/simple/sqls"

func migrate_drop_legacy_api_tables() error {
	return sqls.WithTransaction(func(ctx *sqls.TxContext) error {
		if err := ctx.Tx.Exec("DROP TABLE IF EXISTS t_menu_api").Error; err != nil {
			return err
		}
		if err := ctx.Tx.Exec("DROP TABLE IF EXISTS t_api").Error; err != nil {
			return err
		}
		return nil
	})
}
