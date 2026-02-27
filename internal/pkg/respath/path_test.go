package respath

import (
	"path/filepath"
	"testing"
)

func TestUploadsPath(t *testing.T) {
	got := UploadsPath("a", "b.txt")
	want := filepath.Join(".", "res", "uploads", "a", "b.txt")
	if got != want {
		t.Fatalf("UploadsPath()=%q want=%q", got, want)
	}
}
