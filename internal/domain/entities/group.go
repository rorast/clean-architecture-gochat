package entities

import "time"

type Group struct {
	ID        uint      `json:"id" gorm:"primaryKey"`
	Name      string    `json:"name" gorm:"size:100;not null"`
	OwnerId   uint      `json:"owner_id" gorm:"not null"`
	Icon      string    `json:"icon" gorm:"size:255"`
	Type      int       `json:"type" gorm:"default:0"`      // 群組類型：0-默認，1-興趣愛好，2-行業交流，3-生活休閒，4-學習考試
	Desc      string    `json:"desc" gorm:"size:200"`       // 群組描述
	Size      int       `json:"size" gorm:"default:50"`     // 群組規模：0-小群(50人)，1-中群(200人)，2-大群(500人)
	JoinType  int       `json:"join_type" gorm:"default:0"` // 入群方式：0-自由加入，1-需要驗證，2-不允許加入
	CreatedAt time.Time `json:"created_at" gorm:"default:CURRENT_TIMESTAMP"`
	UpdatedAt time.Time `json:"updated_at" gorm:"default:CURRENT_TIMESTAMP"`
	DeletedAt time.Time `json:"deleted_at" gorm:"index"`
}

// TableName 指定表名
func (Group) TableName() string {
	return "groups"
}
