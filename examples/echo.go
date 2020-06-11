package examples

import (
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
}
