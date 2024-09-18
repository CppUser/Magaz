package ws

import (
	"github.com/gorilla/websocket"
	"go.uber.org/zap"
)

type ClientList map[*Client]bool

type Client struct {
	conn *websocket.Conn
	mng  *Manager
	//egress is used to avoid concurrent writes on ws connection
	egress chan []byte
}

func NewClient(conn *websocket.Conn, mng *Manager) *Client {
	return &Client{
		conn:   conn,
		mng:    mng,
		egress: make(chan []byte),
	}
}

func (c *Client) ReadMessage() {
	defer func() {
		// cleanup connection
		c.mng.RemoveClient(c)
	}()
	for {
		messageType, p, err := c.conn.ReadMessage()
		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				c.mng.logger.Error("websocket connection closed unexpectedly", zap.Error(err))
			}
			break
		}

		//TODO. hack for testing
		for wsclient := range c.mng.clients {
			wsclient.egress <- p
		}

		c.mng.logger.Debug("websocket client read msg", zap.String("type", string(int32(messageType))))
		c.mng.logger.Debug("websocket client read msg", zap.String("message", string(p)))

	}

}

func (c *Client) WriteMessage() {
	defer func() {
		// cleanup connection
		c.mng.RemoveClient(c)
	}()

	for {
		select {
		case message, ok := <-c.egress:
			if !ok {
				if err := c.conn.WriteMessage(websocket.CloseMessage, []byte{}); err != nil {
					c.mng.logger.Error("websocket connection closed unexpectedly", zap.Error(err))

				}
				return
			}

			if err := c.conn.WriteMessage(websocket.TextMessage, message); err != nil {
				c.mng.logger.Error("failed to send message", zap.Error(err))
			}
		}
	}

}
