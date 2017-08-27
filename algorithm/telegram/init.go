package telegram

import (
	"log"
	"net/http"
	"net/url"
	"strconv"

	"encoding/json"
	"io/ioutil"

	"github.com/boobsBot/algorithm/config"
)

// SendMessage отправляет сообщение в чат
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

// SendDocument отправляет документ в чат
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

// GetUpdateEntities возвращает обновления
func GetUpdateEntities(lastUpdateId int) ([]Update, error) {
	var response Response
	u, _ := url.ParseRequestURI(config.TmApiUrl + config.TmToken)
	u.Path += "/getUpdates"
	params := url.Values{}
	params.Set("offset", strconv.Itoa(lastUpdateId+1))
	u.RawQuery = params.Encode()
	resp, err := http.Get(u.String())
	if err != nil {
		return response.Result, err
	}
	defer resp.Body.Close()
	responseBody, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return response.Result, err
	}
	if err = json.Unmarshal(responseBody, &response); err != nil {
		return response.Result, err
	}

	return response.Result, nil
}
