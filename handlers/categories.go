package handlers

import (
	"github.com/drklauss/boobsBot/bot"
	"github.com/drklauss/boobsBot/telegram"
	log "github.com/sirupsen/logrus"
)

// Categories returns all categories available for request.
func Categories(req bot.HandlerRequest) {
	ms := telegram.MessageConfig{
		ChatID:         req.Update.Message.Chat.ID,
		Text:           "Choose category from list",
		KeyboardMarkup: GetCategoriesInlineKeyboard(),
	}
	err := ms.Send()
	if err != nil {
		log.Errorln(err)
		return
	}
}
