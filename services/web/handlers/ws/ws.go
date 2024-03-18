package ws

import (
	"crypto/sha256"
	"fmt"
	"strings"

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
		if ok := w.readMessage(ws, c); !ok {
			return nil
		}
	}
}

func (w *WS) processMessages(feed chan WebSocketMessage, ws *websocket.Conn, c echo.Context) {
	for {
		select {
		case msg := <-feed:
			if msg.Type == Close || msg.Type == "" || len(strings.TrimSpace(msg.Message)) < 1 {
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

func (w *WS) readMessage(ws *websocket.Conn, c echo.Context) bool {
	var chatMessage WebSocketMessage
	if err := ws.ReadJSON(&chatMessage); err != nil {
		if strings.Contains(err.Error(), "websocket: close 1001") {
			c.Logger().Warn("connection broken, closing reader...")
			return false
		}
		c.Logger().Error(err)
		return true
	}
	c.Logger().Info(chatMessage)
	if chatMessage.Type != Close && len(strings.TrimSpace(chatMessage.Message)) < 1 {
		// Ignore empty messages
		return true
	}
	w.feed <- chatMessage
	if chatMessage.Type == Close {
		return false
	}

	return true
}

func hashStr(in string) string {
	h := sha256.New()
	h.Write([]byte(in))
	return string(h.Sum(nil))
}
