package ws

import (
	"strings"

	"github.com/gorilla/websocket"
	"github.com/labstack/echo/v4"

	"github.com/m50/wygoming-satellite/services/homeassistant"
)

func EchoHandler(msg WebSocketMessage, ws *websocket.Conn, c echo.Context) error {
	if msg.Type == "" || len(strings.TrimSpace(msg.Message)) < 1 {
		return nil
		// return fmt.Errorf("Empty message, not echoing...")
	}
	c.Logger().Debugf("attempting to write %s '%s'", msg.Type, msg.Message)
	disp, err := msg.Render(c)
	if err != nil {
		return err
	}
	var messageType int
	if msg.Type == Binary {
		messageType = websocket.BinaryMessage
	} else {
		messageType = websocket.TextMessage
	}

	if err := ws.WriteMessage(messageType, disp); err != nil {
		return err
	}
	c.Logger().Infof("wrote %s '%s'", msg.Type, msg.Message)
	return nil
}

func PipelineHandler(msg WebSocketMessage, c *websocket.Conn, ctx echo.Context) error {
		if msg.Type != Echo {
			return nil
		}

		ha := homeassistant.GetHomeAssistantConnection()
		pm := ha.GetPipelineManager()
		pm.ListPipelines()
		r, err := pm.RunTextPipeline(pm.PreferredPipeline, msg.Message)
		if err != nil {
			return err
		}

		rMsg := WebSocketMessage{
			Message: r,
			Type: Message,
		}

		return EchoHandler(rMsg, c, ctx)
	}
