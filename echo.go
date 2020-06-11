package templates

import (
	"io"

	"github.com/labstack/echo"
)

type EchoRenderer interface {
	Render(io.Writer, string, interface{}, echo.Context) error
}

type echoRenderer struct {
	Engine
}

func (r *echoRenderer) Render(w io.Writer, name string, data interface{}, ctx echo.Context) error {
	return r.Engine.Load(name).Execute(w, data)
}

// Echo creates a new renderer for echo framework.
func Echo(cfg Config) EchoRenderer {
	return &echoRenderer{Create(cfg)}
}
