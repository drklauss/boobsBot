package handlers

import (
	"github.com/drklauss/boobsBot/bot"
	"github.com/drklauss/boobsBot/telegram"
	log "github.com/sirupsen/logrus"
)

func Rate(req bot.HandlerRequest) {
	mes := telegram.MessageConfig{
		ChatID:         req.Update.Message.Chat.ID,
		Text:           "Please, give us 5 \xE2\xAD\x90\xE2\xAD\x90\xE2\xAD\x90\xE2\xAD\x90\xE2\xAD\x90 at \n https://telegram.me/storebot?start=" + req.Config.Telegram.BotName,
		KeyboardMarkup: telegram.ReplyKeyboardRemove{RemoveKeyboard: true},
	}
	if err := mes.Send(); err != nil {
		log.Warnf("could not send message: %v", err)
	}
}
