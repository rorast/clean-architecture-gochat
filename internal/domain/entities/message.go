package entities

import (
	"encoding/json"
	"time"
)

// MessageType 定義不同的訊息類型
type MessageType int

const (
	MessageTypePrivate   MessageType = 1
	MessageTypeGroup     MessageType = 2
	MessageTypeSystem    MessageType = 3
	MessageTypeHeartbeat MessageType = 4
)

// MediaType 定義不同的媒體類型
type MediaType int

const (
	MediaTypeText  MediaType = 1
	MediaTypeImage MediaType = 2
	MediaTypeVoice MediaType = 3
	MediaTypeVideo MediaType = 4
	MediaTypeFile  MediaType = 5
)

// Message 實體
type Message struct {
	ID        uint        `json:"id" gorm:"primaryKey"`
	UserId    uint        `json:"user_id" gorm:"not null"`
	TargetId  uint        `json:"target_id" gorm:"not null"`
	RoomID    uint        `json:"room_id" gorm:"index"` // 聊天室ID
	Type      MessageType `json:"type" gorm:"not null"`
	Media     MediaType   `json:"media" gorm:"not null"`
	Content   string      `json:"content" gorm:"type:text"`
	Metadata  JSON        `json:"metadata" gorm:"type:json"`
	CreatedAt time.Time   `json:"created_at" gorm:"default:CURRENT_TIMESTAMP"`
	UpdatedAt time.Time   `json:"updated_at" gorm:"default:CURRENT_TIMESTAMP"`
}

// JSON 類型用於存儲 JSON 數據
type JSON json.RawMessage

// 實現 GORM 的 Scanner 和 Valuer 接口
func (j *JSON) Scan(value interface{}) error {
	bytes, ok := value.([]byte)
	if !ok {
		return nil
	}
	result := json.RawMessage{}
	err := json.Unmarshal(bytes, &result)
	*j = JSON(result)
	return err
}

func (j JSON) Value() (interface{}, error) {
	if len(j) == 0 {
		return nil, nil
	}
	return json.RawMessage(j).MarshalJSON()
}

// RoomType 定義聊天室類型
type RoomType int

const (
	RoomTypePrivate RoomType = 1
	RoomTypeGroup   RoomType = 2
)

// ChatRoom 代表一個聊天室
type ChatRoom struct {
	ID           uint      `json:"id" gorm:"primaryKey"`
	Name         string    `json:"name" gorm:"size:100;not null"`
	Type         RoomType  `json:"type" gorm:"not null"`
	CreatorID    uint      `json:"creator_id" gorm:"not null"`
	Participants []uint    `json:"participants" gorm:"type:json"`
	CreatedAt    time.Time `json:"created_at" gorm:"default:CURRENT_TIMESTAMP"`
	UpdatedAt    time.Time `json:"updated_at" gorm:"default:CURRENT_TIMESTAMP"`
}
