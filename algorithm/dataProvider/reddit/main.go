package reddit

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"strings"
	"encoding/json"
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
	req, _ := http.NewRequest("GET", GetMeUrl, nil)
	req.Header.Set("User-Agent", UserAgent)
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
	params.Set("username", UserName)
	params.Set("password", Password)
	req, _ := http.NewRequest("POST", GetTokenUrl, strings.NewReader(params.Encode()))
	req.SetBasicAuth(ClientId, ClientSecret)
	req.Header.Set("User-Agent", UserAgent)
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
	var urls[]string
	var err error
	client := new(http.Client)
	fmt.Println(ApiUrl+NSFW+uType)
	fmt.Printf("%+v\n",TokenSample)
	req, _ := http.NewRequest("GET", ApiUrl+NSFW+uType, nil)
	req.Header.Set("Authorization", "bearer "+TokenSample.Token)
	req.Header.Set("User-Agent", UserAgent)
	resp, err := client.Do(req)
	defer resp.Body.Close()
	respBody, err := ioutil.ReadAll(resp.Body) // todo: разобрать тело в структурку и вернуть срез URL-ов
	if err != nil {
		fmt.Printf("%s\n", err)
		return urls, err
	}

	return urls, nil
}
