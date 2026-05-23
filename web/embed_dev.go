//go:build dev

package webspa

import (
	"io/fs"
	"os"
)

var SPA fs.FS = os.DirFS("web")
