package sender

import (
	"boobsBot/config"
	"log"
	"net/http"
	"net/url"
	"strconv"
)

// https://core.telegram.org/bots/api#senddocument
func SendDocument(chatId int, docUrl string) {
	u, _ := url.ParseRequestURI(config.ApiUrl + config.Token)
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
