package config

import (
	"encoding/json"
	"errors"
	"os"
	"strings"

	"github.com/labstack/echo/v4"
	"github.com/labstack/gommon/log"
)

type HomeAssistantConfig struct {
	Address string
	AccessToken string
}

type ConfigValues struct {
	WebUIPort uint16
	HomeAssistant HomeAssistantConfig
	MQTTBroker string
	LogLevel string
}

type Config struct {
	Values ConfigValues
	ConfigPath string
}

func ReadConfig(path string) (*Config, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}

	var conf ConfigValues
	err = json.Unmarshal(data, &conf)
	if err != nil {
		return nil, err
	}
	if err = conf.validate(); err != nil {
		return nil, err
	}

	return &Config{
		ConfigPath: path,
		Values: conf,
	}, nil
}

func (c *ConfigValues) validate() error {
	// HomeAssistantURL
	if len(c.HomeAssistant.Address) < 1 {
		return errors.New("HomeAssistant.Address must be set")
	}
	if !strings.Contains(c.HomeAssistant.Address, ":") {
		return errors.New("HomeAssistant.Address must have a port provided")
	}
	if len(c.HomeAssistant.AccessToken) < 1 {
		return errors.New("HomeAssistant.AccessToken must be set")
	}

	// MQTTBroker
	if len(c.MQTTBroker) < 1 {
		return errors.New("MQTTBroker must be set")
	}
	if !strings.Contains(c.MQTTBroker, ":") {
		return errors.New("MQTTBroker must have a port provided")
	}

	// LogLevel
	if len(c.LogLevel) < 1 {
		c.LogLevel = "info"
	} else {
		c.LogLevel = strings.ToLower(c.LogLevel)
	}

	return nil
}


func (c *Config) SetLogLevel(l echo.Logger) {
	if c.Values.LogLevel == "debug" {
		l.SetLevel(log.DEBUG)
	} else if c.Values.LogLevel == "info" {
		l.SetLevel(log.INFO)
	} else if c.Values.LogLevel == "error" {
		l.SetLevel(log.ERROR)
	} else if c.Values.LogLevel == "warn" {
		l.SetLevel(log.WARN)
	}
}
