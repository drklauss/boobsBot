package handler

import (
	"context"

	"github.com/drklauss/boobsBot/telegram"
	"github.com/leesper/holmes"
)

// Categories returns all cateories available for request
func Categories(ctx context.Context, u *telegram.Update) {
	ms := telegram.MessageSend{
		ChatID:         u.Message.Chat.ID,
		Text:           "Choose category from list",
		KeyboardMarkup: telegram.GetDefaultKeayboard(),
	}
	err := ms.Send()
	if err != nil {
		holmes.Errorln(err)
		return
	}
}
