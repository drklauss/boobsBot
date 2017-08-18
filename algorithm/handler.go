package algorithm

import (
	"boobsBot/entities"
	"boobsBot/config"
)

// Обрабатывет входящий update
func handleUpdate(update entities.Update) {
	switch update.Message.Text {
	case config.HelloCom :

	}
}

