package main

import (
	"fmt"

	"clardy.eu/wygoming-satellite/handlers/ws"
	"clardy.eu/wygoming-satellite/views"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/labstack/gommon/log"
)


func main() {
	count := 0
	wsHandler := ws.GetWS()

	e := echo.New()
	e.Logger.SetLevel(log.DEBUG)
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())
	e.Static("/assets", "assets")
	e.GET("/", func(c echo.Context) error {
		wsHandler.Feed <- ws.WebSocketMessage{Message: "Welcome!"}
		return views.RenderView(c, 200, views.Index(count))
	})
	e.POST("/count", func(c echo.Context) error {
		count++
		wsHandler.Feed <- ws.WebSocketMessage{Message: fmt.Sprintf("New count: %d", count)}
		return views.RenderView(c, 200, views.Count(count))
	})
	e.GET("/ws", wsHandler.Handle)

	e.Logger.Fatal(e.Start(":8765"))
}
