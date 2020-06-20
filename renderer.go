// Copyright © 2020 Alexandre Kovac <contact@kovacou.com>.
//
// Use of this source code is governed by an MIT-style
// license that can be found in the LICENSE file.

package templates

import (
	"net/http"

	"github.com/rs/xid"

	"github.com/kovacou/go-types"
)

// contextID is an unique ID for the current runtime.
var contextID = xid.New().String()

// RenderFunc
type RenderFunc func(code int, name string, data interface{}) error

// New creates a new Renderer with the given callback.
func New(rf RenderFunc) Renderer {
	return &renderer{
		data:       types.SyncMap(),
		renderFunc: rf,
	}
}

// Renderer
type Renderer interface {
	Set(key string, v interface{})
	Parse404() error
	Parse500() error
	Parse(code int, name string, values ...map[string]interface{}) error
	AddPostParsing(...func(data types.TSafeMap) error)
}

type renderer struct {
	engine     Engine
	data       types.TSafeMap
	renderFunc RenderFunc
	pc         []func(types.TSafeMap) error
}

func (e *renderer) Set(key string, v interface{}) {
	e.data.Set(key, v)
}

func (e *renderer) AddPostParsing(c ...func(data types.TSafeMap) error) {
	e.pc = append(e.pc, c...)
}

func (e *renderer) Parse(code int, name string, values ...map[string]interface{}) error {
	if len(values) > 0 {
		for k, v := range values[0] {
			e.data.Set(k, v)
		}
	}

	e.data.Set("context", contextID)
	e.data.Set("activeTemplate", name)

	for _, c := range e.pc {
		if err := c(e.data); err != nil {
			return err
		}
	}

	return e.renderFunc(code, name, e.data.Map())
}

func (e *renderer) Parse404() (err error) {
	return e.Parse(http.StatusNotFound, "errors/404")
}

func (e *renderer) Parse500() (err error) {
	return e.Parse(http.StatusInternalServerError, "errors/500")
}
