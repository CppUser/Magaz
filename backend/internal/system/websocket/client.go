package ws

import (
	"encoding/json"
	"github.com/gorilla/websocket"
	"go.uber.org/zap"
	"time"
)

var (
	pongWait   = 60 * time.Second
	pingPeriod = (pongWait * 9) / 10
)

type ClientList map[*Client]bool

type Client struct {
	conn *websocket.Conn
	mng  *Manager
	//egress is used to avoid concurrent writes on ws connection
	egress chan Event
}

func NewClient(conn *websocket.Conn, mng *Manager) *Client {
	return &Client{
		conn:   conn,
		mng:    mng,
		egress: make(chan Event),
	}
}

func (c *Client) ReadMessage() {
	defer func() {
		// cleanup connection
		c.mng.RemoveClient(c)
	}()

	c.conn.SetReadLimit(1000) //TODO:Adjust limit on messages size to be send

	if err := c.conn.SetReadDeadline(time.Now().Add(pongWait)); err != nil {
		c.mng.logger.Error("failed to set read deadline", zap.Error(err))
		return
	}

	c.conn.SetPongHandler(c.pongHandler)

	for {
		_, p, err := c.conn.ReadMessage()
		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				c.mng.logger.Error("ReadMessage websocket connection closed unexpectedly", zap.Error(err))
			}
			break
		}

		var request Event
		if err := json.Unmarshal(p, &request); err != nil {
			c.mng.logger.Error("error marshaling event", zap.Error(err))
			//break
		}

		if err := c.mng.routeEvents(request, c); err != nil {
			c.mng.logger.Error("error handling event", zap.Error(err))
		}
	}

}

func (c *Client) WriteMessage() {
	defer func() {
		// cleanup connection
		c.mng.RemoveClient(c)
	}()

	ticker := time.NewTicker(pingPeriod)

	for {
		select {
		case message, ok := <-c.egress:
			if !ok {
				if err := c.conn.WriteMessage(websocket.CloseMessage, []byte{}); err != nil {
					c.mng.logger.Error("WriteMessage websocket connection closed unexpectedly", zap.Error(err))

				}
				return
			}

			data, err := json.Marshal(message)
			if err != nil {
				c.mng.logger.Error("error marshaling event", zap.Error(err))
				return
			}

			if err := c.conn.WriteMessage(websocket.TextMessage, data); err != nil {
				c.mng.logger.Error("failed to send message", zap.Error(err))
			}
		case <-ticker.C:
			if err := c.conn.WriteMessage(websocket.PingMessage, []byte{}); err != nil {
				c.mng.logger.Error("websocket connection closed unexpectedly", zap.Error(err))
				return
			}
		}

	}

}

func (c *Client) pongHandler(pongMsg string) error {
	return c.conn.SetReadDeadline(time.Now().Add(pongWait))
}
