// Copyright © 2020 Alexandre Kovac <contact@kovacou.com>.
//
// Use of this source code is governed by an MIT-style
// license that can be found in the LICENSE file.

package templates

import (
	"net/http"

	"go.mongodb.org/mongo-driver/bson/primitive"

	"github.com/kovacou/go-types"
)

// contextID is an unique ID for the current runtime.
var contextID = primitive.NewObjectID().Hex()

// RenderFunc
type RenderFunc func(code int, name string, data interface{}) error

// New creates a new Renderer with the given callback.
func New(rf RenderFunc) Renderer {
	return &renderer{
		data:       types.SyncMap(),
		renderFunc: rf,
	}
}

type Renderer interface {
	Set(key string, v interface{})
	Parse404() error
	Parse500() error
	Parse(code int, name string) error
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

func (e *renderer) Parse(code int, name string) error {
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
