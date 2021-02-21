package model

import (
	"crypto/md5"
	"errors"
	"fmt"
	"strings"

	"github.com/jinzhu/gorm"

	"github.com/drklauss/boobsBot/reddit"
)

// Item contains main info - links to videofiles and its captions.
type Item struct {
	db       *gorm.DB
	ID       int    `gorm:"primary_key;column:id"`
	Category string `gorm:"column:category"`
	URL      string `gorm:"column:url"`
	Hash     string `gorm:"column:hash"`
	Caption  string `gorm:"column:caption"`
}

// TableName returns name of the Item table.
func (Item) TableName() string {
	return "Items"
}

// NewItem returns new Item entity.
func NewItem(db *gorm.DB) *Item {
	return &Item{db: db}
}

// List returns a list of items by categories.
func (i *Item) List(cat string) ([]Item, error) {
	var items []Item
	return items, i.db.Where("category=?", cat).Find(&items).Limit(100).Error
}

// Save into DB
// Returns quantity of items were inserted and an error.
func (i *Item) Save(cat string, els []*reddit.Element) (int, error) {
	insertRows := prepareInsertValues(cat, els)
	if len(insertRows) == 0 {
		return 0, errors.New("empty items")
	}

	beforeCount, err := i.Count(cat)
	if err != nil {
		return 0, err
	}

	insertStr := strings.Join(insertRows, ",")
	err = i.db.Exec(fmt.Sprintf("INSERT OR IGNORE INTO \"Items\" (\"category\",\"url\", \"hash\", \"caption\") VALUES %s", insertStr)).Error
	if err != nil {
		return 0, fmt.Errorf("could not save items: %v", err)
	}

	afterCount, err := i.Count(cat)
	if err != nil {
		return 0, err
	}

	return afterCount - beforeCount, nil
}

// Count counts items by category.
func (i *Item) Count(cat string) (int, error) {
	var c int
	if err := i.db.Table("Items").Where("category=?", cat).Count(&c).Error; err != nil {
		return 0, fmt.Errorf("could not count items by categories: %v", err)
	}

	return c, nil
}

// Fill fills item data.
func (i *Item) Fill(chatID int) error {
	sql := fmt.Sprintf(`
	SELECT *
	FROM Items
	WHERE category="%s" AND id NOT IN
		(SELECT i.id
		FROM Items as i
		LEFT JOIN Views as v
			ON v.itemId = i.id
		WHERE v.chatId = %d)
	ORDER BY RANDOM()
	LIMIT 1`, i.Category, chatID)
	err := i.db.Raw(sql).Find(&i).Error
	if err != nil {
		return fmt.Errorf("could not fill item: %v", err)
	}

	return nil
}

func prepareInsertValues(cat string, els []*reddit.Element) []string {
	var insertRows []string
	for _, el := range els {
		if el.URL == "" {
			continue
		}
		h := md5.New()
		h.Write([]byte(el.URL))
		md := h.Sum(nil)
		cap := strings.Replace(el.Caption, `"`, `'`, -1)
		insertRows = append(insertRows, fmt.Sprintf("(\"%s\", \"%s\", \"%x\", \"%s\")", cat, el.URL, md, cap))
	}

	return insertRows
}
