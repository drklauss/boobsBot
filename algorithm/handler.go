package algorithm

import (
	"boobsBot/entities"
	"boobsBot/config"
	"boobsBot/algorithm/sender"
)

// Обрабатывет входящий update
func handleUpdate(update entities.Update) {
	switch update.Message.Text {
	case config.HelloCom:
		sendHello(update)
	case config.JokeCom:
		sendJoke(update)
	}
}

func sendHello(update entities.Update) {
	sender.SendMessage(update.Message.Chat.Id, "Same to you")
}

func sendJoke(update entities.Update) {
	sender.SendMessage(update.Message.Chat.Id, "Ha-Ha-Ha")
}
