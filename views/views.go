package views

import (
	"bytes"

	"github.com/a-h/templ"
	"github.com/labstack/echo/v4"
)

func RenderView(ctx echo.Context, status int, t templ.Component) error {
    ctx.Response().Writer.WriteHeader(status)
    return t.Render(ctx.Request().Context(), ctx.Response().Writer)
}

func ToStr(ctx echo.Context, t templ.Component) (string, error) {
	b := new(bytes.Buffer)
	if err := t.Render(ctx.Request().Context(), b); err != nil {
		return "", err
	}
	return b.String(), nil
}
