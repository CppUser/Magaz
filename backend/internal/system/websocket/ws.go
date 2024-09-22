package ws

import (
	"Magaz/backend/internal/repository"
	"errors"
	"go.uber.org/zap"
	"gorm.io/gorm"
	"sync"
)

type Manager struct {
	logger       *zap.Logger
	clients      ClientList
	sync.RWMutex //many people concurently connected to api , protecting with ReadWriteMutex
	DB           *gorm.DB
	handlers     map[string]EventHandler
}

func NewManager(logger *zap.Logger, db *gorm.DB) *Manager {
	m := &Manager{
		logger:   logger,
		clients:  make(ClientList),
		DB:       db,
		handlers: make(map[string]EventHandler),
	}

	m.setupEventHandlers()
	return m
}

func (m *Manager) setupEventHandlers() {
	m.handlers[EventSendMessage] = SendMessageHandler
	m.handlers[EventOrderRelease] = OrderReleaseHandler
	m.handlers[EventOutAssignAddress] = AssignAddressHandler
	m.handlers[EventInUpdateAddress] = UpdateAddressHandler
}

func (m *Manager) routeEvents(event Event, c *Client) error {
	if handler, ok := m.handlers[event.Type]; ok {
		if err := handler(event, c); err != nil {
			return err
		}
		return nil
	} else {
		return errors.New("no handler for " + event.Type)
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

// TODO: Refactor so that any caller can use providing type and payload
func (m *Manager) BroadcastOrder(ordView repository.OrderView) {
	event := Event{
		Type:    "new_order",
		Payload: ordView,
	}

	// Iterate over all connected clients
	for client := range m.clients {
		select {
		case client.egress <- event: // Send event to client's egress channel
			// Successfully sent to egress
		default:
			// Handle case where sending fails (e.g., client connection issue)
			close(client.egress)      // Close the connection
			delete(m.clients, client) // Remove client from the manager
		}
	}
}
