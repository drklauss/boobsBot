package model

import (
	"fmt"

	"github.com/jinzhu/gorm"
)

// View holds item views from chats.
type View struct {
	db          *gorm.DB
	ItemID      int   `gorm:"column:itemId"`
	ChatID      int   `gorm:"column:chatId"`
	RequestDate int64 `gorm:"column:requestDate"`
}

// TableName returns name of the View table.
func (View) TableName() string {
	return "Views"
}

// NewView return new View entity.
func NewView(db *gorm.DB) *View {
	return &View{db: db}
}

// Clear erases all views for user by category.
func (v *View) Clear(chatID int, category string) error {
	sql := fmt.Sprintf(`
	DELETE FROM Views
	WHERE Views.urlId IN (
		SELECT u.id FROM Urls AS u
		INNER JOIN Views AS v ON u.id = v.urlId
		WHERE u.category = "%s" AND v.chatId = %d
	)`, category, chatID)
	return v.db.Exec(sql).Error
}
