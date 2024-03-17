package web

import (
	"fmt"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/labstack/gommon/log"
	"github.com/m50/wygoming-satellite/services/config"
	"github.com/m50/wygoming-satellite/services/web/handlers/ws"
	"github.com/m50/wygoming-satellite/views"
)

type WebServer struct {
	conf   *config.Config
	server *echo.Echo
}

func NewWebServer(conf *config.Config, logger *log.Logger) WebServer {
	server := echo.New()
	server.Logger = logger
	server.HideBanner = true
	server.Use(middleware.LoggerWithConfig(middleware.LoggerConfig{
		Format: "${time_rfc3339} \x1b[34mRQST\x1b[0m ${remote_ip} -> ${method} http://${host}${uri} : ${status} ${error}\n",
	}))
	server.Use(middleware.Recover())
	server.Static("/assets", "assets")

	wsHandler := ws.GetWS()
	server.GET("/", func(c echo.Context) error {
		return views.RenderView(c, 200, views.Index())
	})
	server.GET("/chat", func(c echo.Context) error {
		return views.RenderView(c, 200, views.Chat())
	})
	server.GET("/config", func(c echo.Context) error {
		return views.RenderView(c, 200, views.Config(conf))
	})
	server.GET("/ws", wsHandler.Handle)

	return WebServer{
		conf,
		server,
	}
}

func (w *WebServer) Run() {
	w.server.Logger.Fatal(w.server.Start(fmt.Sprintf(":%d", w.conf.WebUIPort)))
}
