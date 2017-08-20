package sender

import (
	"boobsBot/config"
	"log"
	"net/http"
	"net/url"
	"strconv"
)

// Отправляет сообщение в чат
func SendMessage(chatId int, text string) {
	u, _ := url.ParseRequestURI(config.ApiUrl + config.Token)
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
