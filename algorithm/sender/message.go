package sender

import (
	"net/url"
	"strconv"
	"fmt"
	"net/http"
	"boobsBot/config"
	"log"
)

// Отправляет сообщение в чат
func SendMessage(chatId int, text string) {
	u, _ := url.ParseRequestURI(config.ApiUrl + config.Token)
	u.Path += "/sendMessage"
	params := url.Values{}
	params.Set("chat_id", strconv.Itoa(chatId))
	params.Set("text", text)
	u.RawQuery = params.Encode()
	fmt.Println(u.String())
	_, err := http.Get(u.String())
	if err != nil {
		log.Println(err)
	}
}
