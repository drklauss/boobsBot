package algorithm

import (
	"boobsBot/entities"
	"boobsBot/config"
	"net/url"
	"strconv"
	"fmt"
	"net/http"
)

// Обрабатывет входящий update
func handleUpdate(update entities.Update) {
	switch update.Message.Text {
	case config.HelloCom :
		sendHello(update)
	case config.JokeCom:
		sendJoke(update)
	}
}

func sendHello(update entities.Update){
	u, _ := url.ParseRequestURI(config.ApiUrl + config.Token)
	u.Path += "/sendMessage"
	params := url.Values{}
	params.Set("chat_id", strconv.Itoa(update.Message.Chat.Id))
	params.Set("text", "Hello to you, too")
	u.RawQuery = params.Encode()
	fmt.Println(u.String())
	resp, err := http.Get(u.String())
	fmt.Printf("%v \n %v", resp, err)
}


func sendJoke(update entities.Update){
	u, _ := url.ParseRequestURI(config.ApiUrl + config.Token)
	u.Path += "/sendMessage"
	params := url.Values{}
	params.Set("chat_id", strconv.Itoa(update.Message.Chat.Id))
	params.Set("text", "HAHA")
	u.RawQuery = params.Encode()
	fmt.Println(u.String())
	resp, err := http.Get(u.String())
	fmt.Printf("%v \n %v", resp, err)
}