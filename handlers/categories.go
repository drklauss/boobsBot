package handlers

import (
	"context"

	"github.com/drklauss/boobsBot/telegram"
	log "github.com/sirupsen/logrus"
)

// Categories returns all categories available for request.
func Categories(ctx context.Context, u *telegram.Update) {
	ms := telegram.MessageSend{
		ChatID:         u.Message.Chat.ID,
		Text:           "Choose category from list",
		KeyboardMarkup: telegram.GetCategoriesInlineKeyboard(),
	}
	err := ms.Send()
	if err != nil {
		log.Errorln(err)
		return
	}
}
