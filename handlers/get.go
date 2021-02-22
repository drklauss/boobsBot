package handlers

import (
	"context"
	"fmt"
	"time"

	"github.com/drklauss/boobsBot/bot"
	"github.com/drklauss/boobsBot/model"
	"github.com/drklauss/boobsBot/telegram"
	"github.com/jinzhu/gorm"
	log "github.com/sirupsen/logrus"
)

// Get main handler that handles requests from chats.
func Get(ctx context.Context, u *telegram.Update) {
	var (
		db     *gorm.DB
		err    error
		cat    *string
		item   *model.Item
		chatId int
	)

	if db, err = bot.GetDB(ctx); err != nil {
		log.Warnln(err)
		return
	}

	if cat, err = bot.GetCategory(ctx); err != nil {
		log.Warnln(err)
		return
	}

	if u.Message.Chat.ID != 0 {
		chatId = u.Message.Chat.ID
	} else {
		chatId = u.CallBackQuery.Message.Chat.ID
	}

	if item, err = getItem(db, chatId, *cat); err != nil {
		log.Errorln(err)
		return
	}

	if item == nil {
		log.Warnf("could not get item for cat %s", *cat)
		return
	}

	if u.CallBackQuery.ID != "" {
		acq := telegram.AnswerCallbackQuery{
			CallbackQueryID: u.CallBackQuery.ID,
			Text:            "Nice choice, maaaan!",
			URL:             item.URL,
		}
		if err = telegram.SendAnswerCallbackQuery(acq); err != nil {
			log.Errorln(err)
			return
		}
	}

	ms := telegram.MediaSend{
		URL:     item.URL,
		Caption: item.Caption,
		ChatID:  chatId,
	}
	if err = ms.Send(); err != nil {
		log.Errorln(err)
		return
	}

	go writeStat(db, &u.Message.Chat, item)
}

func writeStat(conn *gorm.DB, chat *telegram.Chat, item *model.Item) {
	c := model.Chat{
		ID:    chat.ID,
		Title: chat.Title,
		Type:  chat.Type,
	}
	if err := conn.FirstOrCreate(&c, c).Error; err != nil {
		log.Warnf("could not create chat entry: %v", err)
		return
	}

	v := model.View{
		ChatID:      chat.ID,
		ItemID:      item.ID,
		RequestDate: time.Now().Unix(),
	}
	if err := conn.Create(v).Error; err != nil {
		log.Warnf("could not create view entry: %v", err)
	}

}

func getItem(db *gorm.DB, chatID int, cat string) (*model.Item, error) {
	item := model.NewItem(db)
	item.Category = cat
	err := item.Fill(chatID)
	if err != nil {
		if err == gorm.ErrRecordNotFound { // there is not items to show
			// probably user has watched all items in category
			// clear user's watch history for category
			if err = model.NewView(db).Clear(chatID, cat); err != nil {
				return nil, fmt.Errorf("could not clear views: %v", err)
			}
			// try get item
			if err = item.Fill(chatID); err != nil {
				return nil, err
			}
			return item, nil
		}
		return nil, err
	}
	return item, nil
}
