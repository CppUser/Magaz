package ws

import (
	"go.uber.org/zap"
	"sync"
)

type Manager struct {
	logger       *zap.Logger
	clients      ClientList
	sync.RWMutex //many people concurently connected to api , protecting with ReadWriteMutex
}

func NewManager(logger *zap.Logger) *Manager {
	return &Manager{
		logger:  logger,
		clients: make(ClientList),
	}
}

func (m *Manager) AddClient(client *Client) {
	m.Lock() // locking manager to avoid modifying client list map at same time when multiple entries
	defer m.Unlock()

	m.clients[client] = true // client connected

}

func (m *Manager) RemoveClient(client *Client) {
	m.Lock()
	defer m.Unlock()

	if _, ok := m.clients[client]; ok {
		err := client.conn.Close()
		if err != nil {
			m.logger.Warn("Failed to close websocket connection", zap.Error(err))
			return
		}
		delete(m.clients, client)
	}
}
