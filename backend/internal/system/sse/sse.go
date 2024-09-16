package sse

import (
	"Magaz/backend/internal/repository"

	"go.uber.org/zap"
	"gorm.io/gorm"
)

type SSEHub struct {
	Clients    map[chan interface{}]bool
	DB         *gorm.DB
	Logger     *zap.Logger
	Register   chan chan interface{}
	Unregister chan chan interface{}
	Broadcast  chan interface{}
}

func NewSSEHub(db *gorm.DB, logger *zap.Logger) *SSEHub {
	return &SSEHub{
		Clients:    make(map[chan interface{}]bool),
		DB:         db,
		Logger:     logger,
		Register:   make(chan chan interface{}),
		Unregister: make(chan chan interface{}),
		Broadcast:  make(chan interface{}),
	}
}

func (hub *SSEHub) Run() {
	for {
		select {
		case client := <-hub.Register:
			hub.Clients[client] = true
		case client := <-hub.Unregister:
			if _, ok := hub.Clients[client]; ok {
				delete(hub.Clients, client)
				close(client)
			}
		case message := <-hub.Broadcast:
			for client := range hub.Clients {
				client <- message
			}
		}
	}
}

func (hub *SSEHub) BroadcastMessage(message interface{}) {
	hub.Logger.Info("Broadcasting message to clients", zap.Any("message", message))
	for client := range hub.Clients {
		client <- message
	}
}

func (hub *SSEHub) BroadcastOrderUpdate(order repository.OrderView) {
	message := map[string]interface{}{
		"ID":          order.ID,
		"CityName":    order.CityName,
		"ProductName": order.ProductName,
		"Quantity":    order.Quantity,
		"Due":         order.Due,
		"Username":    order.Client.Username,
		"CreatedAt":   order.CreatedAt.Format("2006-01-02 15:04:05"),
	}

	for client := range hub.Clients {
		client <- message
	}
}
