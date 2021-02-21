package model

import "github.com/jinzhu/gorm"

// Chat holds requester chat info.
type Chat struct {
	db    *gorm.DB
	ID    int    `gorm:"primary_key;column:id"`
	Title string `gorm:"column:title"`
	Type  string `gorm:"column:type"`
}

// TableName returns name of the Chat table.
func (Chat) TableName() string {
	return "Chats"
}

// NewChat returns new Chat entity.
func NewChat(db *gorm.DB) *Chat {
	return &Chat{db: db}
}

// Save creates or updates chat entity.
func (c *Chat) Save() error {
	return c.db.FirstOrCreate(c, c).Error
}
