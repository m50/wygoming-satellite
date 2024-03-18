package handlers

import (
	"encoding/json"
	"os"
	"strconv"

	"github.com/labstack/echo/v4"
	"github.com/m50/wygoming-satellite/services/config"
	"github.com/m50/wygoming-satellite/views"
)

type ConfigHandler struct {
	conf *config.Config
}

func NewConfigHandler(conf *config.Config) ConfigHandler {
	return ConfigHandler{
		conf,
	}
}

func (h *ConfigHandler) HandleGet(c echo.Context) error {
	if c.Request().Header.Get("hx-request") != "true" {
		return views.RenderView(c, 200, views.RootConfig(h.conf.Values))
	}
	return views.RenderView(c, 200, views.Config(h.conf.Values))
}

func (h *ConfigHandler) HandleUpdate(c echo.Context) error {
	p, err := strconv.Atoi(c.FormValue("WebUIPort"))
	if err != nil {
		return err
	}
	h.conf.Values.WebUIPort = uint16(p)
	h.conf.Values.LogLevel = c.FormValue("LogLevel")
	h.conf.Values.HomeAssistant.Address = c.FormValue("HomeAssistantAddress")
	h.conf.Values.HomeAssistant.AccessToken = c.FormValue("HomeAssistantAccessToken")
	h.conf.Values.MQTTBroker = c.FormValue("MQTTBroker")
	cnf, err := json.MarshalIndent(h.conf.Values, "", "  ")
	if err != nil {
		return err
	}
	os.WriteFile(h.conf.ConfigPath, cnf, 664)
	return views.RenderView(c, 200, views.Config(h.conf.Values))
}
