// Copyright © 2020 Alexandre Kovac <contact@kovacou.com>.
//
// Use of this source code is governed by an MIT-style
// license that can be found in the LICENSE file.

package examples

import (
	"github.com/gin-gonic/gin"

	"github.com/kovacou/go-templates"
)

func main() {
	e := gin.Default()
	e.HTMLRender = templates.Gin(templates.Config{
		Layout:    "layout",
		Extension: "html",
		Debug:     false,
		Directory: "resources/templates",
		Funcs:     templates.FuncMap,
	})
}
