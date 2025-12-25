package model

import "time"

// Link 对应数据库中的 links 表
type Link struct {
	ID          uint      `gorm:"primaryKey" json:"id"`
	UserID      uint      `gorm:"index" json:"user_id"` // 哪个用户创建的
	OriginalURL string    `gorm:"not null" json:"original_url"`
	ShortCode   string    `gorm:"unique;index" json:"short_code"` // 短码，如 aB3d
	VisitCount  int       `gorm:"default:0" json:"visit_count"`
	CreatedAt   time.Time `json:"created_at"`
}

func (Link) TableName() string {
	return "links"
}
