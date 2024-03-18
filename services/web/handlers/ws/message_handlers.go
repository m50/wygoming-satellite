package ws

import (
	"github.com/gorilla/websocket"
	"github.com/labstack/echo/v4"
)

func EchoHandler(message WebSocketMessage, ws *websocket.Conn, c echo.Context) error {
	c.Logger().Debugf("attempting to write %s '%s'", message.Type, message.Message)
	disp, err := message.Render(c)
	if err != nil {
		return err
	}
	var messageType int
	if message.Type == "binary" {
		messageType = websocket.BinaryMessage
	} else {
		messageType = websocket.TextMessage
	}

	if err := ws.WriteMessage(messageType, disp); err != nil {
		// Handle page refreshes
		// if strings.Contains(err.Error(), "connection has been hijacked") || strings.Contains(err.Error(), "websocket: close sent") {
		// 	c.Logger().Warnf("failed to write %s '%s' due to closed connection, retrying...", message.Type, message.Message)
		// 	return nil
		// }
		return err
	}
	c.Logger().Infof("wrote %s '%s'", message.Type, message.Message)
	return nil
}
