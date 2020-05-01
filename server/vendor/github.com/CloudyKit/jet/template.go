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

// Jet is a fast and dynamic template engine for the Go programming language, set of features
// includes very fast template execution, a dynamic and flexible language, template inheritance, low number of allocations,
// special interfaces to allow even further optimizations.

package jet

import (
	"fmt"
	"io"
	"io/ioutil"
	"path"
	"reflect"
	"strings"
	"sync"
	"text/template"
)

// Set is responsible to load,invoke parse and cache templates and relations
// every jet template is associated with one set.
// create a set with jet.NewSet(escapeeFn) returns a pointer to the Set
type Set struct {
	loader            Loader
	templates         map[string]*Template // parsed templates
	escapee           SafeWriter           // escapee to use at runtime
	globals           VarMap               // global scope for this template set
	tmx               *sync.RWMutex        // template parsing mutex
	gmx               *sync.RWMutex        // global variables map mutex
	defaultExtensions []string
	developmentMode   bool
	leftDelim         string
	rightDelim        string
}

// SetDevelopmentMode set's development mode on/off, in development mode template will be recompiled on every run
func (s *Set) SetDevelopmentMode(b bool) *Set {
	s.developmentMode = b
	return s
}

func (a *Set) LookupGlobal(key string) (val interface{}, found bool) {
	a.gmx.RLock()
	val, found = a.globals[key]
	a.gmx.RUnlock()
	return
}

// AddGlobal add or set a global variable into the Set
func (s *Set) AddGlobal(key string, i interface{}) *Set {
	s.gmx.Lock()
	if s.globals == nil {
		s.globals = make(VarMap)
	}
	s.globals[key] = reflect.ValueOf(i)
	s.gmx.Unlock()
	return s
}

func (s *Set) AddGlobalFunc(key string, fn Func) *Set {
	return s.AddGlobal(key, fn)
}

// NewSetLoader creates a new set with custom Loader
func NewSetLoader(escapee SafeWriter, loader Loader) *Set {
	return &Set{loader: loader, tmx: &sync.RWMutex{}, gmx: &sync.RWMutex{}, escapee: escapee, templates: make(map[string]*Template), defaultExtensions: append([]string{}, defaultExtensions...)}
}

// NewHTMLSetLoader creates a new set with custom Loader
func NewHTMLSetLoader(loader Loader) *Set {
	return NewSetLoader(template.HTMLEscape, loader)
}

// NewSet creates a new set, dirs is a list of directories to be searched for templates
func NewSet(escapee SafeWriter, dirs ...string) *Set {
	return NewSetLoader(escapee, &OSFileSystemLoader{dirs: dirs})
}

// NewHTMLSet creates a new set, dirs is a list of directories to be searched for templates
func NewHTMLSet(dirs ...string) *Set {
	return NewSet(template.HTMLEscape, dirs...)
}

// AddPath add path to the lookup list, when loading a template the Set will
// look into the lookup list for the file matching the provided name.
func (s *Set) AddPath(path string) {
	if loader, ok := s.loader.(hasAddPath); ok {
		loader.AddPath(path)
	} else {
		panic(fmt.Sprintf("AddPath() not supported on custom loader of type %T", s.loader))
	}
}

// AddGopathPath add path based on GOPATH env to the lookup list, when loading a template the Set will
// look into the lookup list for the file matching the provided name.
func (s *Set) AddGopathPath(path string) {
	if loader, ok := s.loader.(hasAddGopathPath); ok {
		loader.AddGopathPath(path)
	} else {
		panic(fmt.Sprintf("AddGopathPath() not supported on custom loader of type %T", s.loader))
	}
}

// Delims sets the delimiters to the specified strings. Parsed templates will
// inherit the settings. Not setting them leaves them at the default: {{ or }}.
func (s *Set) Delims(left, right string) {
	s.leftDelim = left
	s.rightDelim = right
}

// resolveName try to resolve a template name, the steps as follow
//	1. try provided path
//	2. try provided path+defaultExtensions
// ex: set.resolveName("catalog/products.list") with defaultExtensions set to []string{".html.jet",".jet"}
//	try catalog/products.list
//	try catalog/products.list.html.jet
//	try catalog/products.list.jet
func (s *Set) resolveName(name string) (newName, fileName string, foundLoaded, foundFile bool) {
	newName = name
	if _, foundLoaded = s.templates[newName]; foundLoaded {
		return
	}

	if fileName, foundFile = s.loader.Exists(name); foundFile {
		return
	}

	for _, extension := range s.defaultExtensions {
		newName = name + extension
		if _, foundLoaded = s.templates[newName]; foundLoaded {
			return
		}
		if fileName, foundFile = s.loader.Exists(newName); foundFile {
			return
		}
	}

	return
}

func (s *Set) resolveNameSibling(name, sibling string) (newName, fileName string, foundLoaded, foundFile, isRelativeName bool) {
	if sibling != "" {
		i := strings.LastIndex(sibling, "/")
		if i != -1 {
			if newName, fileName, foundLoaded, foundFile = s.resolveName(path.Join(sibling[:i+1], name)); foundFile || foundLoaded {
				isRelativeName = true
				return
			}
		}
	}
	newName, fileName, foundLoaded, foundFile = s.resolveName(name)
	return
}

// Parse parses the template, this method will link the template to the set but not the set to
func (s *Set) Parse(name, content string) (*Template, error) {
	sc := *s
	sc.developmentMode = true

	sc.tmx.RLock()
	t, err := sc.parse(name, content)
	sc.tmx.RUnlock()

	return t, err
}

func (s *Set) loadFromFile(name, fileName string) (template *Template, err error) {
	f, err := s.loader.Open(fileName)
	if err != nil {
		return nil, err
	}
	defer f.Close()
	content, err := ioutil.ReadAll(f)
	if err != nil {
		return nil, err
	}
	return s.parse(name, string(content))
}

func (s *Set) getTemplateWhileParsing(parentName, name string) (template *Template, err error) {
	name = path.Clean(name)

	if s.developmentMode {
		if newName, fileName, _, foundPath, _ := s.resolveNameSibling(name, parentName); foundPath {
			return s.loadFromFile(newName, fileName)
		} else {
			return nil, fmt.Errorf("template %s can't be loaded", name)
		}
	}

	if newName, fileName, foundLoaded, foundPath, isRelative := s.resolveNameSibling(name, parentName); foundPath {
		template, err = s.loadFromFile(newName, fileName)
		s.templates[newName] = template

		if !isRelative {
			s.templates[name] = template
		}
	} else if foundLoaded {
		template = s.templates[newName]
		if !isRelative && name != newName {
			s.templates[name] = template
		}
	} else {
		err = fmt.Errorf("template %s can't be loaded", name)
	}
	return
}

// getTemplate gets a template already loaded by name
func (s *Set) getTemplate(name, sibling string) (template *Template, err error) {
	name = path.Clean(name)

	if s.developmentMode {
		s.tmx.RLock()
		defer s.tmx.RUnlock()
		if newName, fileName, foundLoaded, foundFile, _ := s.resolveNameSibling(name, sibling); foundFile || foundLoaded {
			if foundFile {
				template, err = s.loadFromFile(newName, fileName)
			} else {
				template, _ = s.templates[newName]
			}
		} else {
			err = fmt.Errorf("template %s can't be loaded", name)
		}
		return
	}

	//fast path
	s.tmx.RLock()
	newName, fileName, foundLoaded, foundFile, isRelative := s.resolveNameSibling(name, sibling)

	if foundLoaded {
		template = s.templates[newName]
		s.tmx.RUnlock()
		if !isRelative && name != newName {
			// creates an alias
			s.tmx.Lock()
			if _, found := s.templates[name]; !found {
				s.templates[name] = template
			}
			s.tmx.Unlock()
		}
		return
	}
	s.tmx.RUnlock()

	//not found parses and cache
	s.tmx.Lock()
	defer s.tmx.Unlock()

	newName, fileName, foundLoaded, foundFile, isRelative = s.resolveNameSibling(name, sibling)
	if foundLoaded {
		template = s.templates[newName]
		if !isRelative && name != newName {
			// creates an alias
			if _, found := s.templates[name]; !found {
				s.templates[name] = template
			}
		}
	} else if foundFile {
		template, err = s.loadFromFile(newName, fileName)

		if !isRelative && name != newName {
			// creates an alias
			if _, found := s.templates[name]; !found {
				s.templates[name] = template
			}
		}

		s.templates[newName] = template
	} else {
		err = fmt.Errorf("template %s can't be loaded", name)
	}
	return
}

func (s *Set) GetTemplate(name string) (template *Template, err error) {
	template, err = s.getTemplate(name, "")
	return
}

func (s *Set) LoadTemplate(name, content string) (template *Template, err error) {
	if s.developmentMode {
		s.tmx.RLock()
		defer s.tmx.RUnlock()
		template, err = s.parse(name, content)
		return
	}

	//fast path
	var found bool
	s.tmx.RLock()
	if template, found = s.templates[name]; found {
		s.tmx.RUnlock()
		return
	}
	s.tmx.RUnlock()

	//not found parses and cache
	s.tmx.Lock()
	defer s.tmx.Unlock()

	if template, found = s.templates[name]; found {
		return
	}

	if template, err = s.parse(name, content); err == nil {
		s.templates[name] = template
	}

	return
}

func (t *Template) String() (template string) {
	if t.extends != nil {
		if len(t.Root.Nodes) > 0 && len(t.imports) == 0 {
			template += fmt.Sprintf("{{extends %q}}", t.extends.ParseName)
		} else {
			template += fmt.Sprintf("{{extends %q}}", t.extends.ParseName)
		}
	}

	for k, _import := range t.imports {
		if t.extends == nil && k == 0 {
			template += fmt.Sprintf("{{import %q}}", _import.ParseName)
		} else {
			template += fmt.Sprintf("\n{{import %q}}", _import.ParseName)
		}
	}

	if t.extends != nil || len(t.imports) > 0 {
		if len(t.Root.Nodes) > 0 {
			template += "\n" + t.Root.String()
		}
	} else {
		template += t.Root.String()
	}
	return
}

func (t *Template) addBlocks(blocks map[string]*BlockNode) {
	if len(blocks) > 0 {
		if t.processedBlocks == nil {
			t.processedBlocks = make(map[string]*BlockNode)
		}
		for key, value := range blocks {
			t.processedBlocks[key] = value
		}
	}
}

type VarMap map[string]reflect.Value

func (scope VarMap) Set(name string, v interface{}) VarMap {
	scope[name] = reflect.ValueOf(v)
	return scope
}

func (scope VarMap) SetFunc(name string, v Func) VarMap {
	scope[name] = reflect.ValueOf(v)
	return scope
}

func (scope VarMap) SetWriter(name string, v SafeWriter) VarMap {
	scope[name] = reflect.ValueOf(v)
	return scope
}

// Execute executes the template in the w Writer
func (t *Template) Execute(w io.Writer, variables VarMap, data interface{}) error {
	return t.ExecuteI18N(nil, w, variables, data)
}

type Translator interface {
	Msg(key, defaultValue string) string
	Trans(format, defaultFormat string, v ...interface{}) string
}

func (t *Template) ExecuteI18N(translator Translator, w io.Writer, variables VarMap, data interface{}) (err error) {
	st := pool_State.Get().(*Runtime)
	defer st.recover(&err)

	st.blocks = t.processedBlocks
	st.translator = translator
	st.variables = variables
	st.set = t.set
	st.Writer = w

	// resolve extended template
	for t.extends != nil {
		t = t.extends
	}

	if data != nil {
		st.context = reflect.ValueOf(data)
	}

	st.executeList(t.Root)
	return
}
