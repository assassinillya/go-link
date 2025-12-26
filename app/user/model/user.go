package model

import "time"

type User struct {
	ID        uint      `gorm:"primaryKey" json:"id"`
	Username  string    `gorm:"unique;not null" json:"username"`
	Password  string    `gorm:"not null" json:"-"` // json:"-" 表示返回 JSON 时隐藏密码字段
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

// TableName 指定表名（GORM 默认是用复数，为了保险显式指定）
func (User) TableName() string {
	// 返回数据库中的表名
	return "users"
}
