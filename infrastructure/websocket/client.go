package websocket

import (
	"log"
	"time"

	"github.com/gorilla/websocket"
)

// Client 代表一個 WebSocket 連線的客戶端
type Client struct {
	Conn          *websocket.Conn
	UserID        string
	Send          chan []byte
	Hub           *Hub
	LastHeartbeat time.Time
	IsClosed      bool
}

// ReadPump 處理從客戶端讀取訊息
func (c *Client) ReadPump() {
	defer func() {
		c.Hub.Unregister <- c
		c.Conn.Close()
	}()

	c.Conn.SetReadLimit(512)
	c.Conn.SetReadDeadline(time.Now().Add(60 * time.Second))
	c.Conn.SetPongHandler(func(string) error {
		c.Conn.SetReadDeadline(time.Now().Add(60 * time.Second))
		return nil
	})

	for {
		_, message, err := c.Conn.ReadMessage()
		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				log.Printf("error: %v", err)
			}
			break
		}
		c.Hub.Broadcast <- message
	}
}

// WritePump 處理向客戶端寫入訊息
func (c *Client) WritePump() {
	ticker := time.NewTicker(54 * time.Second)
	defer func() {
		ticker.Stop()
		c.Conn.Close()
	}()

	for {
		select {
		case message, ok := <-c.Send:
			c.Conn.SetWriteDeadline(time.Now().Add(10 * time.Second))
			if !ok {
				c.Conn.WriteMessage(websocket.CloseMessage, []byte{})
				return
			}

			w, err := c.Conn.NextWriter(websocket.TextMessage)
			if err != nil {
				return
			}
			w.Write(message)

			if err := w.Close(); err != nil {
				return
			}
		case <-ticker.C:
			c.Conn.SetWriteDeadline(time.Now().Add(10 * time.Second))
			if err := c.Conn.WriteMessage(websocket.PingMessage, nil); err != nil {
				return
			}
		}
	}
}

// 處理 WebSocket 訊息
func (c *Client) HandleMessage(msg map[string]interface{}) {
	// 更新心跳時間
	if msg["type"] == "heartbeat" {
		c.LastHeartbeat = time.Now()
		return
	}

	// 處理其他類型的訊息
	// TODO: 根據實際需求處理不同類型的訊息
}

// 發送訊息給特定客戶端
func (c *Client) SendMessage(msg []byte) {
	if !c.IsClosed {
		c.Send <- msg
	}
}
