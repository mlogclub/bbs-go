// Copyright 2016 José Santos <henrique_1609@me.com>
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
// http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package jet

import (
	"errors"
	"io"
	"os"
	"path"
	"path/filepath"
)

// Loader is a minimal interface required for loading templates.
type Loader interface {
	// Open opens the underlying reader with template content.
	Open(name string) (io.ReadCloser, error)
	// Exists checks for template existence and returns full path.
	Exists(name string) (string, bool)
}

// hasAddPath is an optional Loader interface. Most probably useful for OS file system only, thus unexported.
type hasAddPath interface {
	AddPath(path string)
}

// hasAddGopathPath is an optional Loader interface. Most probably useful for OS file system only, thus unexported.
type hasAddGopathPath interface {
	AddGopathPath(path string)
}

// OSFileSystemLoader implements Loader interface using OS file system (os.File).
type OSFileSystemLoader struct {
	dirs []string
}

// NewOSFileSystemLoader returns an initialized OSFileSystemLoader.
func NewOSFileSystemLoader(paths ...string) *OSFileSystemLoader {
	return &OSFileSystemLoader{dirs: paths}
}

// Open opens a file from OS file system.
func (l *OSFileSystemLoader) Open(name string) (io.ReadCloser, error) {
	return os.Open(name)
}

// Exists checks if the template name exists by walking the list of template paths
// returns string with the full path of the template and bool true if the template file was found
func (l *OSFileSystemLoader) Exists(name string) (string, bool) {
	for i := 0; i < len(l.dirs); i++ {
		fileName := path.Join(l.dirs[i], name)
		if _, err := os.Stat(fileName); err == nil {
			return fileName, true
		}
	}
	return "", false
}

// AddPath adds the path to the internal list of paths searched when loading templates.
func (l *OSFileSystemLoader) AddPath(path string) {
	l.dirs = append(l.dirs, path)
}

// AddGopathPath adds a path located in the GOPATH.
// Example: l.AddGopathPath("github.com/CloudyKit/jet/example/views")
func (l *OSFileSystemLoader) AddGopathPath(path string) {
	paths := filepath.SplitList(os.Getenv("GOPATH"))
	for i := 0; i < len(paths); i++ {
		var err error
		path, err = filepath.Abs(filepath.Join(paths[i], "src", path))
		if err != nil {
			panic(errors.New("Can't add this path err: " + err.Error()))
		}

		if fstats, err := os.Stat(path); os.IsNotExist(err) == false && fstats.IsDir() {
			l.AddPath(path)
			return
		}
	}

	if fstats, err := os.Stat(path); os.IsNotExist(err) == false && fstats.IsDir() {
		l.AddPath(path)
	}
}
