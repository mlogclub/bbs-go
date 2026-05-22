package server

import (
	"fmt"
	"io"
	"os"

	"bbs-go/internal/pkg/config"
)

func printBanner() {
	printBannerTo(os.Stdout, config.Instance, config.GetEnv())
}

func printBannerTo(w io.Writer, cfg *config.Config, env string) {
	_, _ = fmt.Fprint(w, renderBanner(cfg, env))
}

func renderBanner(cfg *config.Config, env string) string {
	if cfg == nil {
		cfg = &config.Config{}
	}
	if env == "" {
		env = config.EnvDev
	}

	return fmt.Sprintf(`
 ____  ____  ____         ____  ___
| __ )| __ )/ ___|       / ___|/ _ \
|  _ \|  _ \\___ \_____ | |  _| | | |
| |_) | |_) |___) |_____| |_| | |_| |
|____/|____/|____/       \____|\___/

:: BBS-GO ::  https://bbs-go.com

Environment : %s
Port        : %d
Language    : %s
Installed   : %t
Address     : http://127.0.0.1:%d

`, env, cfg.Port, cfg.Language, cfg.Installed, cfg.Port)
}
