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
	itemsTable := (new(Item)).TableName()
	viewsTable := v.TableName()
	sql := fmt.Sprintf(`
	DELETE FROM %s AS v
	WHERE v.itemId IN (
		SELECT u.id 
		FROM %s AS u
			INNER JOIN %s AS v1 ON u.id = v1.itemId
		WHERE u.category = "%s" AND v1.chatId = %d
	)`, viewsTable, itemsTable, viewsTable, category, chatID)

	return v.db.Exec(sql).Error
}
