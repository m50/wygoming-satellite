package ws

import (
	"strings"

	"clardy.eu/wygoming-satellite/views"
	"clardy.eu/wygoming-satellite/views/components"
	"github.com/labstack/echo/v4"
	"golang.org/x/net/websocket"
)

// WebSocketMessage is the shape of messages coming in from the Websocket
type WebSocketMessage struct {
	Message string `json:"message"`
	Headers map[string]string `json:"HEADERS"`
}

// WS is the representation of the WebSocket handler
type WS struct {
	// Feed is the feed of websocket messages
	Feed chan WebSocketMessage
}

var singleton *WS;

// GetWs retrieves the WebSocket singleton, making it if it hasn't been created yet.
func GetWS() WS {
	if singleton != nil {
		return *singleton
	}
	feed := make(chan WebSocketMessage, 10)

	singleton = &WS{
		Feed: feed,
	}

	return *singleton
}

// Handle is the websocket handler for the echo server
func (w *WS) Handle(c echo.Context) error {
	websocket.Handler(func(ws *websocket.Conn) {
		defer ws.Close()
		w.Feed <- WebSocketMessage{Message: "Hello client!"}
		go w.writeMessage(ws, c)
		for {
			if w.readMessage(ws, c) {
				break
			}
		}
	}).ServeHTTP(c.Response(), c.Request())
	return nil
}

func (w *WS) writeMessage(ws *websocket.Conn, c echo.Context) {
	for {
		select {
		case message := <-w.Feed:
			c.Logger().Infof("writing '%v'", message.Message)
			disp, err := views.ToStr(c, components.Message(message.Message))
			if err != nil {
				c.Logger().Error(err)
				w.Feed <-message
				break;
			}
			if err := websocket.Message.Send(ws, disp); err != nil {
				if strings.Contains(err.Error(), "use of closed network connection") {
					w.Feed <-message
					return
				}
				c.Logger().Error(err)
			}
			break
		}
	}
}


func (w *WS) readMessage(ws *websocket.Conn, c echo.Context) bool {
	var chatMessage WebSocketMessage
	if err := websocket.JSON.Receive(ws, &chatMessage); err != nil {
		if strings.Contains(err.Error(), "use of closed network connection") {
			return true
		}
		if err.Error() == "EOF" {
			ws.Close()
			return true
		}
		c.Logger().Error(err)
		return false
	}
	w.Feed <- chatMessage

	return false
}
