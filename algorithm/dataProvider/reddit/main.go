package reddit

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"strings"

	"boobsBot/algorithm/config"
)

var TokenSample TokenResponse

func GetUrls(uType string) ([]string, error) {
	if !isGoodToken() {
		refreshToken()
	}
	return getNewUrls(uType)
}

// Проверяет действителен ли токен
func isGoodToken() bool {
	isGoodToken := true
	client := new(http.Client)
	req, _ := http.NewRequest("GET", config.GetMeUrl, nil)
	req.Header.Set("User-Agent", config.UserAgent)
	resp, err := client.Do(req)
	if err != nil {
		fmt.Printf("%s\n", err)
		log.Println(err)
	}
	if resp.StatusCode != http.StatusOK {
		isGoodToken = false
	}

	return isGoodToken
}

// Получает новый токен и записывает его
func refreshToken() error {
	client := new(http.Client)
	params := url.Values{}
	params.Set("grant_type", "password")
	params.Set("username", config.UserName)
	params.Set("password", config.Password)
	req, _ := http.NewRequest("POST", config.GetTokenUrl, strings.NewReader(params.Encode()))
	req.SetBasicAuth(config.ClientId, config.ClientSecret)
	req.Header.Set("User-Agent", config.UserAgent)
	resp, err := client.Do(req)
	defer resp.Body.Close()
	respBody, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return err
	}
	err = json.Unmarshal(respBody, &TokenSample)
	if err != nil {
		return err
	}

	return nil
}

// Получает новые URL-ы по типу
func getNewUrls(uType string) ([]string, error) {
	var urls []string
	var err error
	client := new(http.Client)
	fmt.Println(config.ApiUrl + config.NSFW + uType)
	fmt.Printf("%+v\n", TokenSample)

	req, _ := http.NewRequest("GET", config.ApiUrl+config.NSFW+uType, nil)
	data := req.URL.Query()
	data.Set("limit", "50")
	req.URL.RawQuery = data.Encode()
	req.Header.Set("Authorization", "bearer "+TokenSample.Token)
	req.Header.Set("User-Agent", config.UserAgent)
	resp, err := client.Do(req)
	if err != nil {
		return urls, err
	}
	fmt.Printf("%v\n", req.RequestURI)
	defer resp.Body.Close()
	respBody, _ := ioutil.ReadAll(resp.Body) // todo: разобрать тело в структурку и вернуть срез URL-ов
	var subResp SubRedditResponse
	err = json.Unmarshal(respBody, &subResp)
	fmt.Printf("%v\n", subResp)
	if err != nil {
		fmt.Printf("%s\n", err)
		return urls, err
	}
	log.Fatal()

	return urls, nil
}
