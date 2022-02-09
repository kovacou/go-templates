// Copyright © 2020 Alexandre Kovac <contact@kovacou.com>.
//
// Use of this source code is governed by an MIT-style
// license that can be found in the LICENSE file.

package templates

import (
	"errors"
	"fmt"
	"html/template"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"regexp"
	"strings"

	"github.com/kovacou/go-types"
)

var (
	ErrTemplate = errors.New("template can't be nil")
	ErrName     = errors.New("template must have a name")

	includeRegexp = regexp.MustCompile(`{{(\s*)template(\s*)"includes/([a-zA-Z0-9/-]+)"(.*)}}`)
)

type Config struct {
	Layout    string
	Extension string
	Debug     bool // Debug mode will enable files watcher.
	Directory string
	Funcs     types.Map
}

func (c Config) getLayout() string {
	return c.Layout + "." + c.Extension
}

// Create
func Create(config Config) Engine {
	if len(config.Funcs) == 0 {
		config.Funcs = FuncMap.Copy()
	}

	e := &engine{
		funcs:  config.Funcs,
		tpl:    map[string]*template.Template{},
		files:  map[string]types.Strings{},
		config: config,
	}

	files := types.Strings{}
	filepath.Walk(config.Directory, func(path string, i os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		if !i.IsDir() {
			files = append(files, path)
		}
		return nil
	})

	for _, f := range files {
		if e.isLayout(f) || strings.Contains(f, "includes/") {
			continue
		}

		layout := e.layoutOf(f)
		isChild := files.Contains(layout)
		tree := types.Strings{}

		if isChild {
			tree = append(tree, layout, f)
		} else {
			tree = append(tree, f)
		}

		if sub := e.computeInclude(tree, isChild); !sub.Empty() {
			tree = append(tree, sub...)
		}

		e.addFromFiles(e.getTemplateName(f), tree)
	}

	return e
}

type Engine interface {
	// Parse(tpl string, v interface{})
	Load(string) *template.Template
}

type engine struct {
	tpl    map[string]*template.Template
	files  map[string]types.Strings
	funcs  map[string]interface{}
	config Config
	logger *log.Logger
}

func (e *engine) Load(tpl string) *template.Template {
	if e.config.Debug {
		return e.load(tpl)
	}
	return e.tpl[tpl]
}

func (e *engine) load(tpl string) *template.Template {
	return e.loadFromFiles(e.files[tpl])
}

func (e *engine) loadFromFiles(files []string) *template.Template {
	basename := filepath.Base(files[0])
	out, err := template.New(basename).Funcs(e.funcs).ParseFiles(files...)
	if err != nil {
		panic(err)
	}
	return out
}

func (e *engine) add(name string, tpl *template.Template) {
	if tpl == nil {
		panic(ErrTemplate)
	}

	if len(name) == 0 {
		panic(ErrName)
	}

	e.tpl[name] = tpl
}

func (e *engine) addFromFiles(name string, files []string) {
	if e.config.Debug {
		e.files[name] = files
	} else {
		e.add(name, e.loadFromFiles(files))
	}
}

func (e *engine) computeInclude(files []string, layout bool) (inc types.Strings) {
	for _, f := range files {
		body, _ := ioutil.ReadFile(f)
		for _, m := range includeRegexp.FindAllStringSubmatch(string(body), -1) {
			inc = append(
				inc,
				filepath.Clean(
					fmt.Sprintf("%s/includes/%s.%s",
						e.config.Directory,
						m[3],
						e.config.Extension,
					),
				))
		}
	}

	return
}

func (e *engine) getTemplateName(path string) string {
	dir, file := filepath.Split(path)
	dir = strings.Replace(dir, e.config.Directory, "", 1)
	file = strings.TrimSuffix(file, "."+e.config.Extension)
	return strings.Trim(filepath.Clean(dir+file), string(os.PathSeparator))
}

func (e *engine) isLayout(tpl string) bool {
	return filepath.Base(tpl) == e.config.getLayout()
}

func (e *engine) layoutOf(tpl string) string {
	return filepath.Dir(tpl) + "/" + e.config.getLayout()
}
