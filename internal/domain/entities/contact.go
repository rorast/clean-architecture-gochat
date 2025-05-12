package entities

import "time"

// 人員關係
type Contact struct {
	ID        uint      `json:"id"`
	OwnerId   uint      `json:"owner_id"`  // 誰的關係信息
	TargetId  uint      `json:"target_id"` // 對應的誰 /群 ID
	Type      int       `json:"type"`      // 對應的類型  1好友  2群  3xx
	Desc      string    `json:"desc"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	DeletedAt time.Time `json:"deleted_at"`
}
