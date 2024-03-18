package ws

import (
	"fmt"

	"github.com/labstack/echo/v4"
	"github.com/m50/wygoming-satellite/views"
	"github.com/m50/wygoming-satellite/views/components"
)

// WebSocketMessageType can be "echo", "message", "binary"
type WebSocketMessageType string

// WebSocketMessage is the shape of messages coming in from the Websocket
type WebSocketMessage struct {
	Message    string `json:"message"`
	Close      bool `json:"close"`
	BinaryData []byte
	Type       WebSocketMessageType
}

// Render will return either a []byte if WebSocketMessage Type is binary, otherwise it will return a string
func (m *WebSocketMessage) Render(c echo.Context) ([]byte, error) {
	if m.Type == "echo" {
		disp, err := views.ToStr(c, components.WSEcho(m.Message))
		return []byte(disp), err
	} else if m.Type == "message" {
		disp, err := views.ToStr(c, components.WSMessage(m.Message))
		return []byte(disp), err
	} else if m.Type == "binary" {
		return m.BinaryData, nil
	}

	return []byte{}, fmt.Errorf("Not a valid WebSocketMessage Type")
}
