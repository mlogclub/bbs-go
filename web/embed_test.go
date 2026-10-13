package webspa_test

import (
	"io/fs"
	"testing"

	webspa "bbs-go/web"
)

func TestSPA_ImplementsFS(t *testing.T) {
	var _ fs.FS = webspa.SPA
	if webspa.SPA == nil {
		t.Fatal("expected SPA filesystem to not be nil")
	}
}
