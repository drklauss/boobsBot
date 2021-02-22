package handlers

import (
	"bytes"

	"github.com/drklauss/boobsBot/bot"
	"github.com/drklauss/boobsBot/config"
	"github.com/drklauss/boobsBot/telegram"
	log "github.com/sirupsen/logrus"
)

// Help sends a helping command.
func Help(req bot.HandlerRequest) {
	buf := new(bytes.Buffer)
	buf.WriteString("We offer you hot girls from popular categories \xF0\x9F\x98\x9C \n")
	buf.WriteString("Available commands:\n")
	buf.WriteString("/start shows menu\n")
	buf.WriteString("/cats sends all available categories\n")
	buf.WriteString("/rate gives a link to rate for us\n")
	buf.WriteString("/help Bot help \n")
	for _, c := range config.Get().Categories {
		buf.WriteString("/" + c.Name + "\n")
	}
	buf.WriteString("\n\n Do not forget to rate us 5 \xE2\xAD\x90\xE2\xAD\x90\xE2\xAD\x90\xE2\xAD\x90\xE2\xAD\x90 at \n")
	buf.WriteString("https://telegram.me/storebot?start=" + req.Config.Telegram.BotName)

	mes := telegram.MessageConfig{
		ChatID: req.Update.Message.Chat.ID,
		Text:   buf.String(),
	}
	if err := mes.Send(); err != nil {
		log.Warnf("could not send message: %v", err)
	}
}

func Start(req bot.HandlerRequest) {
	mes := telegram.MessageConfig{
		ChatID:         req.Update.Message.Chat.ID,
		Text:           "Hi!",
		KeyboardMarkup: GetCategoriesKeyboard(),
	}
	if err := mes.Send(); err != nil {
		log.Warnf("could not send message: %v", err)
	}
}
