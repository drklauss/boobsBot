package algorithm

import (
	"strings"

	"boobsBot/algorithm/sender"
	"boobsBot/config"
	"boobsBot/entities"
)

// Обрабатывет входящий update
func (d *Dispatcher) handleUpdate(update entities.Update) {
	comName := strings.Split(update.Message.Text, config.FullBotName)
	switch comName[0] {
	case config.HelloCom:
		sendHello(update)
	case config.JokeCom:
		sendJoke(update)
	case config.DornCom:
		sendDorn(update)
	}
}

func sendHello(update entities.Update) {
	sender.SendMessage(update.Message.Chat.Id, "Same to you")
}

func sendJoke(update entities.Update) {
	sender.SendMessage(update.Message.Chat.Id, "Ha-Ha-Ha")
}

func sendDorn(update entities.Update) {

	sender.SendDocument(update.Message.Chat.Id, "https://thumbs.gfycat.com/ZigzagAbleDuckling-mobile.mp4")
	// Удаляем отправленный элемент из массива
}
