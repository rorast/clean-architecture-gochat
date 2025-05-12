package websocket

import (
	"sync"
)

// Hub 負責管理所有 WebSocket 連線
type Hub struct {
	Clients    map[string]*Client
	Register   chan *Client
	Unregister chan *Client
	Broadcast  chan []byte
	lock       sync.RWMutex
}

// NewHub 創建一個新的 Hub 實例
func NewHub() *Hub {
	return &Hub{
		Clients:    make(map[string]*Client),
		Register:   make(chan *Client),
		Unregister: make(chan *Client),
		Broadcast:  make(chan []byte),
	}
}

// Run 啟動 WebSocket Hub
func (h *Hub) Run() {
	for {
		select {
		case client := <-h.Register:
			h.lock.Lock()
			h.Clients[client.UserID] = client
			h.lock.Unlock()
		case client := <-h.Unregister:
			h.lock.Lock()
			if _, ok := h.Clients[client.UserID]; ok {
				delete(h.Clients, client.UserID)
				close(client.Send)
			}
			h.lock.Unlock()
		case message := <-h.Broadcast:
			h.lock.RLock()
			for _, client := range h.Clients {
				select {
				case client.Send <- message:
				default:
					close(client.Send)
					delete(h.Clients, client.UserID)
				}
			}
			h.lock.RUnlock()
		}
	}
}
