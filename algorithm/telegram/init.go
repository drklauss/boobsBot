package telegram

import (
	"log"
	"net/http"
	"net/url"
	"strconv"

	"github.com/boobsBot/algorithm/config"
)

// Отправляет сообщение в чат
func SendMessage(chatId int, text string) {
	u, _ := url.ParseRequestURI(config.TmApiUrl + config.TmToken)
	u.Path += "/sendMessage"
	params := url.Values{}
	params.Set("chat_id", strconv.Itoa(chatId))
	params.Set("text", text)
	u.RawQuery = params.Encode()
	_, err := http.Get(u.String())
	if err != nil {
		log.Println(err)
	}
}

// https://core.telegram.org/bots/api#senddocument
func SendDocument(chatId int, docUrl string) {
	u, _ := url.ParseRequestURI(config.TmApiUrl + config.TmToken)
	u.Path += "/sendDocument"
	params := url.Values{}
	params.Set("chat_id", strconv.Itoa(chatId))
	params.Set("document", docUrl)
	u.RawQuery = params.Encode()
	_, err := http.Get(u.String())
	if err != nil {
		log.Println(err)
	}
}
