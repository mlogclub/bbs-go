package server

import (
	"strings"
	"testing"

	"bbs-go/internal/pkg/config"
)

func TestRenderBannerIncludesBrandAndRuntimeInfo(t *testing.T) {
	banner := renderBanner(&config.Config{
		Language:  config.LanguageZhCN,
		Port:      9090,
		Installed: true,
	}, config.EnvProd)

	for _, want := range []string{
		":: BBS-GO ::  https://bbs-go.com",
		"Environment : prod",
		"Port        : 9090",
		"Language    : zh-CN",
		"Installed   : true",
	} {
		if !strings.Contains(banner, want) {
			t.Fatalf("banner missing %q:\n%s", want, banner)
		}
	}
}

func TestRenderBannerHandlesNilConfig(t *testing.T) {
	banner := renderBanner(nil, "")

	for _, want := range []string{
		":: BBS-GO ::  https://bbs-go.com",
		"Environment : dev",
		"Port        : 0",
		"Language    :",
		"Installed   : false",
	} {
		if !strings.Contains(banner, want) {
			t.Fatalf("banner missing %q:\n%s", want, banner)
		}
	}
}
