//go:build !dev

package webspa

import "embed"

// `all:` is required because React Router emits route chunks whose filenames
// begin with an underscore. Without it, Go's embed package silently skips
// those files and a packaged binary serves index.html for their asset URLs.
// The browser then reloads indefinitely while trying to load the route module.
//go:embed all:build/spa
var SPA embed.FS
