package webspa

import (
	"io/fs"
	"strings"
	"testing"
)

func TestEmbeddedSPAIncludesUnderscoreRouteChunks(t *testing.T) {
	found := false
	err := fs.WalkDir(SPA, "build/spa/assets", func(path string, entry fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if !entry.IsDir() && strings.HasPrefix(entry.Name(), "_") {
			found = true
		}
		return nil
	})
	if err != nil {
		t.Fatal(err)
	}
	if !found {
		t.Fatal("embedded SPA is missing underscore-prefixed route chunks")
	}
}
