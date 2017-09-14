package reddit

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"strings"

	"strconv"

	"github.com/boobsBot/algorithm/config"
)

var tokenSample TokenResponse
var lastNameId string

// GetNames возвращает названия новых видео
func GetNames(uType string) ([]string, error) {
	if !isGoodToken() {
		refreshToken()
	}
	return fetchNames(uType)
}

// Проверяет действителен ли токен
func isGoodToken() bool {
	isGoodToken := true
	client := new(http.Client)
	req, _ := http.NewRequest("GET", config.GetMeUrl, nil)
	req.Header.Set("User-Agent", config.UserAgent)
	resp, err := client.Do(req)
	if err != nil {
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
	err = json.Unmarshal(respBody, &tokenSample)
	if err != nil {
		return err
	}

	return nil
}

// Возвращает названия новых видео по типу
func fetchNames(uType string) ([]string, error) {
	var urls []string
	var err error
	client := new(http.Client)
	req, _ := http.NewRequest("GET", config.ApiUrl+config.NSFW+uType, nil)
	data := req.URL.Query()
	data.Set("limit", strconv.Itoa(config.UrlsLimit))
	data.Set("after", lastNameId)
	req.URL.RawQuery = data.Encode()
	req.Header.Set("Authorization", "bearer "+tokenSample.Token)
	req.Header.Set("User-Agent", config.UserAgent)
	resp, err := client.Do(req)
	if err != nil {
		return urls, err
	}
	defer resp.Body.Close()
	respBody, _ := ioutil.ReadAll(resp.Body)
	var subResp SubRedditResponse
	err = json.Unmarshal(respBody, &subResp)
	if err != nil {
		return urls, err
	}
	urls = getOnlyUsefulNames(subResp)

	return urls, nil
}

// Возвращает только названия видео с gfycat
func getOnlyUsefulNames(subResp SubRedditResponse) []string {
	var names []string
	for _, child := range subResp.Data.Children {
		if child.Data.Domain != config.GfycatDomain {
			continue
		}
		namesSlice := strings.Split(child.Data.Url, "/")
		names = append(names, namesSlice[len(namesSlice)-1])
		lastNameId = child.Data.Name
	}

	return names
}
