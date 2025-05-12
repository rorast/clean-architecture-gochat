package websocket

import (
	"net/http"
	"sync"
	"time"

	"github.com/gorilla/websocket"
)

var (
	hub  *Hub
	once sync.Once
)

var Upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

// GetHub 獲取 Hub 實例
func GetHub() *Hub {
	once.Do(func() {
		hub = NewHub()
		go hub.Run()
	})
	return hub
}

// 升級 HTTP 連線至 WebSocket 連線
// 參數: conn - WebSocket 連線，userID - 用戶ID
// 返回值: *Client - 代表已建立的 WebSocket 連線客戶端
func UpgradeConnection(conn *websocket.Conn, userID string) *Client {
	client := &Client{
		Conn:          conn,
		UserID:        userID,
		Send:          make(chan []byte, 50),
		Hub:           GetHub(),
		LastHeartbeat: time.Now(),
		IsClosed:      false,
	}

	GetHub().Register <- client
	go client.ReadPump()
	go client.WritePump()
	return client
}

// 透過用戶ID取得 WebSocket 連線的客戶端
func GetClient(userID string) (*Client, bool) {
	hub := GetHub()
	hub.lock.RLock()
	defer hub.lock.RUnlock()
	client, exists := hub.Clients[userID]
	return client, exists
}

// 廣播訊息給所有連線的客戶端
func Broadcast(msg []byte) {
	GetHub().Broadcast <- msg
}
