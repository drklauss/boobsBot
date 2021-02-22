package handlers

import (
	"fmt"
	"strings"
	"time"

	"github.com/drklauss/boobsBot/bot"
	"github.com/drklauss/boobsBot/model"
	"github.com/drklauss/boobsBot/telegram"
	"github.com/jinzhu/gorm"
	log "github.com/sirupsen/logrus"
)

// Get main handler that handles requests from chats.
func Get(req bot.HandlerRequest) {
	var (
		err    error
		item   *model.Item
		chatId int
	)

	var category string
	if req.Update.Message.ID > 0 {
		category = req.Update.Message.Text
	} else if req.Update.CallBackQuery.ID != "" {
		category = req.Update.CallBackQuery.Data
	}
	if len(category) > 0 && category[:1] == "/" { // cats saved in db without "/"
		category = category[1:]
	}
	if strings.Index(category, "@"+req.Config.Telegram.BotName) > 0 {
		category = category[:strings.Index(category, "@"+req.Config.Telegram.BotName)]
	}

	if req.Update.Message.Chat.ID != 0 {
		chatId = req.Update.Message.Chat.ID
	} else {
		chatId = req.Update.CallBackQuery.Message.Chat.ID
	}

	if item, err = getItem(req.DB, chatId, category); err != nil {
		log.Errorln(err)
		return
	}
	if item == nil {
		log.Warnf("could not get item for cat %s", category)
		return
	}

	if req.Update.CallBackQuery.ID != "" {
		//acq := telegram.AnswerCallbackQuery{
		//	CallbackQueryID: req.Update.CallBackQuery.ID,
		//	Text:            "Nice choice, maaaan!",
		//	URL:             item.URL,
		//}
		//if err = telegram.SendAnswerCallbackQuery(acq); err == nil {
		//	return
		//}
	}

	ms := telegram.MediaSend{
		URL:     item.URL,
		Caption: item.Caption,
		ChatID:  chatId,
	}
	if req.Update.Message.Chat.Type == "private" {
		ms.KeyboardMarkup = GetCategoriesKeyboard()
	}
	if err = ms.Send(); err != nil {
		log.Errorln(err)
		return
	}

	go writeStat(req.DB, &req.Update.Message.Chat, item)
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
