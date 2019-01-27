package handler

import (
	"context"
	"fmt"
	"time"

	"github.com/jinzhu/gorm"

	"github.com/drklauss/boobsBot/bot"
	"github.com/drklauss/boobsBot/model"
	"github.com/leesper/holmes"

	"github.com/drklauss/boobsBot/telegram"
)

// Get main handler that handles requests from chats
func Get(ctx context.Context, u *telegram.Update) {
	conn, err := bot.GetDB(ctx)
	if err != nil {
		holmes.Warnln(err)
	}
	cat, err := bot.GetCategory(ctx)
	if err != nil {
		holmes.Warnln(err)
	}
	item, err := getItem(conn, u.Message.Chat.ID, *cat)
	if err != nil {
		holmes.Errorln(err)
		return
	}
	ms := telegram.MediaSend{
		ChatID:         u.Message.Chat.ID,
		URL:            item.URL,
		Caption:        item.Caption,
		KeyboardMarkup: telegram.GetDefaultKeayboard(),
	}
	holmes.Debugf("ready for send for %d \"%v\"", ms.ChatID, ms.URL)
	err = ms.Send()
	if err != nil {
		holmes.Errorln(err)
		return
	}
	go writeStat(conn, &u.Message.Chat, item)
}

func writeStat(conn *gorm.DB, chat *telegram.Chat, item *model.Item) {
	c := model.Chat{
		ID:    chat.ID,
		Title: chat.Title,
		Type:  chat.Type,
	}
	if err := conn.FirstOrCreate(&c, c).Error; err != nil {
		holmes.Warnf("could not create chat entry: %v", err)
		return
	}

	v := model.View{
		ChatID:      chat.ID,
		ItemID:      item.ID,
		RequestDate: time.Now().Unix(),
	}
	if err := conn.Create(v).Error; err != nil {
		holmes.Warnf("could not create view entry: %v", err)
	}

}

func getItem(db *gorm.DB, chatID int, cat string) (*model.Item, error) {
	item := model.NewItem(db)
	item.Category = cat
	err := item.Fill(chatID)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			if err = model.NewView(db).Clear(chatID, cat); err != nil {
				return nil, fmt.Errorf("could not clear views: %v", err)
			}
			item, err = getItem(db, chatID, cat)
			if err != nil {
				return nil, err
			}
			return nil, nil
		}
		return nil, err
	}
	return item, nil
}
