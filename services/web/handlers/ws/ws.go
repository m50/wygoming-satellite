package ws

import (
	"crypto/sha256"
	"fmt"

	"github.com/gorilla/websocket"
	"github.com/labstack/echo/v4"
)
type MessageHandler func(WebSocketMessage, *websocket.Conn, echo.Context) error

// WS is the representation of the WebSocket handler
type WS struct {
	// feed is the feed of websocket messages
	feed        chan WebSocketMessage
	connections map[string]chan WebSocketMessage
	handlers    []MessageHandler
	upgrader    websocket.Upgrader
}

var singleton *WS

// GetWs retrieves the WebSocket singleton, making it if it hasn't been created yet.
func GetWS() *WS {
	if singleton == nil {
		singleton = &WS{
			feed:        make(chan WebSocketMessage, 5),
			connections: map[string]chan WebSocketMessage{},
			handlers:    []MessageHandler{},
			upgrader:    websocket.Upgrader{},
		}

		go singleton.fanOut()
	}

	return singleton
}

func (w *WS) fanOut() {
	for {
		select {
		case msg, ok := <-w.feed:
			if !ok {
				return
			}
			for _, conn := range(w.connections) {
				conn <-msg
			}
		}
	}
}

func (w *WS) RegisterHandle(handler MessageHandler) {
	w.handlers = append(w.handlers, handler)
}

// Handle is the websocket handler for the echo server
func (w *WS) Handle(c echo.Context) error {
	key := hashStr(c.Request().RemoteAddr)
	if _, exists := w.connections[key]; exists {
		return fmt.Errorf("A connection from this address already exists, only one connection per use allowed")
	}
	ws, err := w.upgrader.Upgrade(c.Response(), c.Request(), nil)
	if err != nil {
		return err
	}
	defer ws.Close()

	feed := make(chan WebSocketMessage, 5)
	feed <- WebSocketMessage{Message: "Hello! How can I help you today?", Type: "message"}
	defer close(feed)

	w.connections[key] = feed
	defer delete(w.connections, key)

	go w.processMessages(feed, ws, c)
	for {
		// keep alive
	}
	// go w.processMessages(ws, c)
	// for {
	// 	if ok := w.readMessage(ws, c); !ok {
	// 		return nil
	// 	}
	// }
}

func (w *WS) processMessages(feed chan WebSocketMessage, ws *websocket.Conn, c echo.Context) {
	for {
		select {
		case msg := <-feed:
			if msg.Close {
				return
			}
			for _, h := range(w.handlers) {
				if err := h(msg, ws, c); err != nil {
					c.Logger().Error(err)
				}
			}
		}
	}
}

// func (w *WS) processMessages(ws *websocket.Conn, c echo.Context) {
// 	for {
// 		select {
// 		case message := <-w.feed:
// 			c.Logger().Debugf("attempting to write %s '%s'", message.Type, message.Message)
// 			disp, err := message.Render(c)
// 			if err != nil {
// 				c.Logger().Error(err)
// 				break
// 			}
// 			var messageType int
// 			if message.Type == "binary" {
// 				messageType = websocket.BinaryMessage
// 			} else {
// 				messageType = websocket.TextMessage
// 			}

// 			if err := ws.WriteMessage(messageType, disp); err != nil {
// 				// Handle page refreshes
// 				if strings.Contains(err.Error(), "connection has been hijacked") || strings.Contains(err.Error(), "websocket: close sent") {
// 					c.Logger().Warnf("failed to write %s '%s' due to closed connection, retrying...", message.Type, message.Message)
// 					w.feed <-message
// 					return
// 				}
// 				c.Logger().Error(err)
// 			}
// 			c.Logger().Infof("wrote %s '%s'", message.Type, message.Message)
// 			break
// 		}
// 	}
// }

// func (w *WS) readMessage(ws *websocket.Conn, c echo.Context) bool {
// 	var chatMessage WebSocketMessage
// 	if err := ws.ReadJSON(&chatMessage); err != nil {
// 		// Handle page refreshes
// 		if strings.Contains(err.Error(), "(going away)") || strings.Contains(err.Error(), "close 1005") {
// 			return false
// 		}
// 		c.Logger().Error(err)
// 		return true
// 	}
// 	if !chatMessage.Close && len(strings.TrimSpace(chatMessage.Message)) < 1 {
// 		// Ignore empty messages
// 		return true
// 	}
// 	chatMessage.Type = "echo"
// 	w.feed <- chatMessage

// 	return true
// }

func hashStr(in string) string {
	h := sha256.New()
	h.Write([]byte(in))
	return string(h.Sum(nil))
}
