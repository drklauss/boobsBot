package algorithm

import (
	"strings"

	"fmt"

	"github.com/boobsBot/algorithm/config"
	"github.com/boobsBot/algorithm/telegram"
)

// Обрабатывет команду
func (d *Dispatcher) handleUpdate(update telegram.Update) {
	comName := strings.Split(update.Message.Text, config.TmFullBotName)
	command := strings.Replace(comName[0], "/", "", -1)
	switch command {
	case config.Hello:
		telegram.SendMessage(update.Message.Chat.Id, "Well, Hello")
	case config.New:
		u := d.urlProvider.GetUrl(config.New)
		fmt.Printf("Real got %v\n", u)

		telegram.SendDocument(update.Message.Chat.Id, u)
	case config.Hot:
		telegram.SendDocument(update.Message.Chat.Id, d.urlProvider.GetUrl(config.Hot))
	case config.Top:
		telegram.SendDocument(update.Message.Chat.Id, d.urlProvider.GetUrl(config.Top))
	}
}
