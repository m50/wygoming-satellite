package ws

import (
	"fmt"

	"github.com/labstack/echo/v4"
	"github.com/m50/wygoming-satellite/views"
	"github.com/m50/wygoming-satellite/views/components"
)

// WebSocketMessageType can be "echo", "message", "binary"
type WebSocketMessageType string

var Echo WebSocketMessageType = "echo"
var Message WebSocketMessageType = "message"
var Binary WebSocketMessageType = "binary"
var Close WebSocketMessageType = "close"

// WebSocketMessage is the shape of messages coming in from the Websocket
type WebSocketMessage struct {
	Message    string               `json:"message"`
	BinaryData []byte               `json:"binary"`
	Type       WebSocketMessageType `json:"type"`
}

// Render will return either a []byte if WebSocketMessage Type is binary, otherwise it will return a string
func (m *WebSocketMessage) Render(c echo.Context) ([]byte, error) {
	if m.Type == Echo {
		disp, err := views.ToStr(c, components.WSEcho(m.Message))
		return []byte(disp), err
	} else if m.Type == Message {
		disp, err := views.ToStr(c, components.WSMessage(m.Message))
		return []byte(disp), err
	} else if m.Type == Binary {
		return m.BinaryData, nil
	} else if m.Type == Close {
		return []byte{}, nil
	}

	return []byte{}, fmt.Errorf("Not a valid WebSocketMessage Type %v", m.Type)
}
