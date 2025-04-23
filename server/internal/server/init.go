package server

import (
	"bbs-go/internal/install"
	"bbs-go/internal/pkg/config"

	_ "github.com/golang-migrate/migrate/v4/source/file"
)

func Init() {
	install.InitConfig()
	install.InitLogger()
	if config.Instance.Installed {
		install.InitDB()
		install.InitOthers()
	}
}
