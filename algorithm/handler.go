package algorithm

import (
	"strings"

	"boobsBot/algorithm/config"
	"boobsBot/algorithm/telegram"
)

// Обрабатывет команду
func (d *Dispatcher) handleUpdate(update telegram.Update) {
	comName := strings.Split(update.Message.Text, config.TmFullBotName)
	command := strings.Replace(comName[0], "/", "", -1)
	switch command {
	case config.Hello:
		sendHello(update)
	case config.New:
		sendNewCorn(update)
	case config.Hot:
		sendHotCorn(update)
	case config.Random:
		sendRandomCorn(update)
	}
}

func sendHello(update telegram.Update) {
	telegram.SendMessage(update.Message.Chat.Id, "Same to you")
}

func sendNewCorn(update telegram.Update) {
	telegram.SendMessage(update.Message.Chat.Id, "Ha-Ha-Ha")
}

func sendHotCorn(update telegram.Update) {

	telegram.SendDocument(update.Message.Chat.Id, "https://thumbs.gfycat.com/ZigzagAbleDuckling-mobile.mp4")
	// Удаляем отправленный элемент из массива
}

func sendRandomCorn(update telegram.Update) {

	telegram.SendDocument(update.Message.Chat.Id, "https://thumbs.gfycat.com/ZigzagAbleDuckling-mobile.mp4")
	// Удаляем отправленный элемент из массива
}
