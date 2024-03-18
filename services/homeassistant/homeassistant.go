package homeassistant

import (
	"encoding/json"
	"fmt"
	"net/url"
	"strings"

	"github.com/gorilla/websocket"
	"github.com/labstack/gommon/log"
	"github.com/m50/wygoming-satellite/services/config"
)

type HomeAssistant struct {
	nextRequestId int
	conf          *config.Config
	logger        *log.Logger
	conn          *websocket.Conn
	requestEvents map[int]chan []byte
}

var singleton *HomeAssistant

func GetHomeAssistantConnect() *HomeAssistant {
	return singleton
}

func NewHomeAssistantConnection(conf *config.Config, logger *log.Logger) *HomeAssistant {
	if singleton == nil {
		singleton = &HomeAssistant{
			nextRequestId: 0,
			conf:   conf,
			logger: logger,
			requestEvents: make(map[int]chan []byte),
		}
	}

	return singleton
}

func tryAndUpgradeConnect(u url.URL) (*websocket.Conn, error) {
	c, _, err := websocket.DefaultDialer.Dial(u.String(), nil)
	if err != nil && err.Error() == "websocket: bad handshake" {
		u.Scheme = "wss"
		return tryAndUpgradeConnect(u)
	}
	return c, err
}

func auth(accessToken string) map[string]string {
	return map[string]string{"type": "auth", "access_token": accessToken}
}

func (ha *HomeAssistant) Run() {
	if ha.conn != nil {
		ha.logger.Error("Attempted to connected to HomeAssisstant with an already open connection")
		return
	}
	wsPath := "/api/websocket"
	u := url.URL{Scheme: "ws", Host: ha.conf.Values.HomeAssistant.Address, Path: wsPath}
	c, err := tryAndUpgradeConnect(u)
	if err != nil {
		ha.logger.Error(err)
		return
	}
	ha.conn = c
	defer c.Close()
	defer func() { ha.conn = nil }()

	ha.readMessages()
}

func (ha *HomeAssistant) readMessages() {
	for {
		_, msg, err := ha.conn.ReadMessage()
		if err != nil {
			if strings.Contains(err.Error(), "unexpected EOF") {
				return
			}
			ha.logger.Error(err)
			continue
		}
		ha.logger.Info("<- ", string(msg))
		if strings.Contains(string(msg), "auth_required") {
			ha.Request(0, auth(ha.conf.Values.HomeAssistant.AccessToken))
			ha.Done(0)
			continue
		}

		var m struct{
			ID int `json:"id"`
		}
		if err := json.Unmarshal(msg, &m); err != nil {
			ha.logger.Error(err)
			continue
		}
		if _, ok := ha.requestEvents[m.ID]; ok {
			ha.requestEvents[m.ID] <- msg
		}
	}
}

func (ha *HomeAssistant) Done(reqID int) error {
	if _, ok := ha.requestEvents[reqID]; !ok {
		return fmt.Errorf("Channel for request %d events not open", reqID)
	}
	close(ha.requestEvents[reqID])
	delete(ha.requestEvents, reqID)

	return nil
}

func (ha *HomeAssistant) Request(reqID int, evt interface{}) (chan []byte, error) {
	if ha.conn == nil {
		return nil, fmt.Errorf("Unable to write message, HomeAssistant connection not open")
	}

	ha.logger.Info("-> ", evt)

	if err := ha.conn.WriteJSON(evt); err != nil {
		return nil, err
	}
	respChan := make(chan []byte, 2)
	ha.requestEvents[reqID] = respChan

	return respChan, nil
}

func (ha *HomeAssistant) AwaitResponse(reqID int) ([]byte, error) {
	if _, ok := ha.requestEvents[reqID]; !ok {
		return []byte{}, fmt.Errorf("Request %d is already closed, unable to await responses", reqID)
	}
	defer ha.Done(reqID)
	resp := <-ha.requestEvents[reqID]
	return resp, nil
}

func (ha *HomeAssistant) NextRequestId() int {
	ha.nextRequestId++
	return ha.nextRequestId
}
