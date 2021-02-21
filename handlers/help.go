package handlers

import (
	"bytes"
	"context"

	"github.com/drklauss/boobsBot/telegram"
	log "github.com/sirupsen/logrus"
)

// Help sends a helping command.
func Help(ctx context.Context, u *telegram.Update) {
	buf := new(bytes.Buffer)
	buf.WriteString("We offer you hot girls from popular categories \xF0\x9F\x98\x9C \n")
	buf.WriteString("Available commands:\n")
	buf.WriteString("/categories contains all avilable categories \n")
	buf.WriteString("/rate Gives a link to rate for us \n")
	buf.WriteString("/help Bot help \n")
	buf.WriteString("\n\n Do not forget to rate us 5 \xE2\xAD\x90\xE2\xAD\x90\xE2\xAD\x90\xE2\xAD\x90\xE2\xAD\x90 at \n")
	buf.WriteString("https://telegram.me/storebot?start=DornBot")

	mes := telegram.MessageSend{
		ChatID: u.Message.Chat.ID,
		Text:   buf.String(),
	}
	if err := mes.Send(); err != nil {
		log.Warnf("could not send message: %v", err)
	}
}
