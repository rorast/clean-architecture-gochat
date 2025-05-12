package entities

import "time"

type GroupMember struct {
	GroupID   uint      `json:"group_id" gorm:"primaryKey"`
	UserID    uint      `json:"user_id" gorm:"primaryKey"`
	CreatedAt time.Time `json:"created_at" gorm:"default:CURRENT_TIMESTAMP"`
}

// TableName 指定表名
func (GroupMember) TableName() string {
	return "group_members"
}
