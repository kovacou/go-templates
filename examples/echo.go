// Copyright © 2020 Alexandre Kovac <contact@kovacou.com>.
//
// Use of this source code is governed by an MIT-style
// license that can be found in the LICENSE file.

package examples

import (
	"log"
	"net/http"

	"github.com/labstack/echo"

	"github.com/kovacou/go-templates"
)

func main() {
	e := echo.New()
	e.Renderer = templates.Echo(templates.Config{
		Layout:    "layout",
		Extension: "html",
		Debug:     false,
		Directory: "resources/templates",
		Funcs:     templates.FuncMap,
	})

	e.Use(func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(ctx echo.Context) (err error) {
			ctx.Set("render", templates.New(ctx.Render))
			return next(ctx)
		}
	})

	e.GET("/welcome", func(ctx echo.Context) error {
		return ctx.Get("render").(templates.Renderer).Parse(http.StatusOK, "welcome", map[string]interface{}{
			"who": "world",
		})
	})

	e.GET("/articles", func(ctx echo.Context) error {
		return ctx.Get("render").(templates.Renderer).Parse(http.StatusOK, "articles/listing")
	})

	e.GET("/articles/show", func(ctx echo.Context) error {
		return ctx.Get("render").(templates.Renderer).Parse(http.StatusOK, "articles/article")
	})

	log.Panic(e.Start(":8080"))
}
