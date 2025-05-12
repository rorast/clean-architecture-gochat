package entities

import "time"

type User struct {
	ID            uint      `json:"id" gorm:"primaryKey"`
	Name          string    `json:"name"`
	Password      string    `json:"-" gorm:"size:128;not null"`
	Phone         string    `json:"phone"`
	Email         string    `json:"email" gorm:"size:128;unique"`
	Avatar        string    `json:"avatar" gorm:"size:255"`
	Identity      string    `json:"identity"`
	ClientIp      string    `json:"client_ip"`
	ClientPort    string    `json:"client_port"`
	Salt          string    `json:"salt"`
	LoginTime     time.Time `json:"login_time"`
	HeartbeatTime time.Time `json:"heartbeat_time"`
	LogoutTime    time.Time `json:"logout_time"`
	IsLogout      bool      `json:"is_logout"`
	DeviceInfo    string    `json:"device_info"`
	Username      string    `json:"username" gorm:"-"`
	CreatedAt     time.Time `json:"created_at" gorm:"default:CURRENT_TIMESTAMP"`
	UpdatedAt     time.Time `json:"updated_at" gorm:"default:CURRENT_TIMESTAMP"`
}
