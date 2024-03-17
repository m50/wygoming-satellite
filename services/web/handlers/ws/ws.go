package ws

import (
	"errors"
	"strings"

	"github.com/gorilla/websocket"
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

// WS is the representation of the WebSocket handler
type WS struct {
	// Feed is the feed of websocket messages
	Feed      chan WebSocketMessage
	CloseFeed chan bool
	upgrader  websocket.Upgrader
}

var singleton *WS

// GetWs retrieves the WebSocket singleton, making it if it hasn't been created yet.
func GetWS() *WS {
	if singleton == nil {
		singleton = &WS{
			Feed:      make(chan WebSocketMessage, 10),
			CloseFeed: make(chan bool),
			upgrader:  websocket.Upgrader{},
		}
	}

	return singleton
}

// Handle is the websocket handler for the echo server
func (w *WS) Handle(c echo.Context) error {
	ws, err := w.upgrader.Upgrade(c.Response(), c.Request(), nil)
	if err != nil {
		return err
	}
	defer ws.Close()

	w.Feed <- WebSocketMessage{Message: "Hello! How can I help you today?", Type: "message"}

	go w.writeMessage(ws, c)
	for {
		if ok := w.readMessage(ws, c); !ok {
			return nil
		}
	}
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

	return []byte{}, errors.New("Not a valid WebSocketMessage Type")
}

func (w *WS) writeMessage(ws *websocket.Conn, c echo.Context) {
	for {
		select {
		case <-w.CloseFeed:
			return
		case message := <-w.Feed:
			c.Logger().Debugf("attempting to write %s '%s'", message.Type, message.Message)
			disp, err := message.Render(c)
			if err != nil {
				c.Logger().Error(err)
				w.Feed <- message
				break
			}
			var messageType int
			if message.Type == "binary" {
				messageType = websocket.BinaryMessage
			} else {
				messageType = websocket.TextMessage
			}

			if err := ws.WriteMessage(messageType, disp); err != nil {
				// Handle page refreshes
				if strings.Contains(err.Error(), "connection has been hijacked") || strings.Contains(err.Error(), "websocket: close sent") {
					c.Logger().Warnf("failed to write %s '%s' due to closed connection, retrying...", message.Type, message.Message)
					w.Feed <-message
					return
				}
				c.Logger().Error(err)
			}
			c.Logger().Infof("wrote %s '%s'", message.Type, message.Message)
			break
		}
	}
}

func (w *WS) readMessage(ws *websocket.Conn, c echo.Context) bool {
	var chatMessage WebSocketMessage
	if err := ws.ReadJSON(&chatMessage); err != nil {
		// Handle page refreshes
		if strings.Contains(err.Error(), "(going away)") || strings.Contains(err.Error(), "close 1005") {
			return false
		}
		c.Logger().Error(err)
		return true
	}
	if chatMessage.Close {
		w.CloseFeed <-true
		return false
	}
	if len(strings.TrimSpace(chatMessage.Message)) < 1 {
		// Ignore empty messages
		return true
	}
	chatMessage.Type = "echo"
	w.Feed <- chatMessage

	return true
}
