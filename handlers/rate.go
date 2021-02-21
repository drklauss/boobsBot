package handlers

import (
	"context"

	"github.com/drklauss/boobsBot/telegram"
	log "github.com/sirupsen/logrus"
)

func Rate(ctx context.Context, u *telegram.Update) {
	mes := telegram.MessageSend{
		ChatID:         u.Message.Chat.ID,
		Text:           "Please, give us 5 \xE2\xAD\x90\xE2\xAD\x90\xE2\xAD\x90\xE2\xAD\x90\xE2\xAD\x90 at \n https://telegram.me/storebot?start=DornBot",
		KeyboardMarkup: telegram.ReplyKeyboardRemove{RemoveKeyboard: true},
	}
	if err := mes.Send(); err != nil {
		log.Warnf("could not send message: %v", err)
	}
}
