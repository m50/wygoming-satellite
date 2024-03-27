package views

import (
	"bytes"

	"github.com/a-h/templ"
	"github.com/labstack/echo/v4"
	"github.com/m50/wygoming-satellite/views/layout"
)

func RenderView(ctx echo.Context, status int, t templ.Component) error {
    ctx.Response().Writer.WriteHeader(status)
	if ctx.Request().Header.Get("hx-request") != "true" {
		base := layout.Base()
		children := templ.WithChildren(ctx.Request().Context(), t)
		base.Render(children, ctx.Response().Writer)
	}
    return t.Render(ctx.Request().Context(), ctx.Response().Writer)
}

func ToStr(ctx echo.Context, t templ.Component) (string, error) {
	b := new(bytes.Buffer)
	if err := t.Render(ctx.Request().Context(), b); err != nil {
		return "", err
	}
	return b.String(), nil
}
