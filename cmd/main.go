package main

import (
	"github.com/labstack/gommon/log"
	"github.com/m50/wygoming-satellite/services/config"
	"github.com/m50/wygoming-satellite/services/homeassistant"
	"github.com/m50/wygoming-satellite/services/web"
)

func newLogger(conf *config.Config) *log.Logger {
	logger := log.New("wygoming-satellite")
	logger.SetHeader("${time_rfc3339} ${level} ${short_file}:${line}")
	conf.SetLogLevel(logger)

	return logger
}

func main() {
	conf, err := config.ReadConfig("./config.json")
	if err != nil {
		panic(err)
	}

	logger := newLogger(conf)
	ha := homeassistant.NewHomeAssistantConnection(conf, logger)
	webServer := web.NewWebServer(conf, logger)
	go ha.Run()
	go webServer.Run()

	for {
	} // Keep alive
}
