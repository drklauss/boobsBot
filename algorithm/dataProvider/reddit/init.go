package reddit

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"strings"

	"strconv"

	"github.com/drklauss/boobsBot/algorithm/config"
	"github.com/drklauss/boobsBot/algorithm/dataProvider/dbEntities"
	"github.com/drklauss/boobsBot/algorithm/dataProvider/gfycat"
)

var tokenSample TokenResponse

// GetItems возвращает срез ImageItems
func GetItems(uType string, lastNameId string) ([]dbEntities.Url, string, error) {
	if !isGoodToken() {
		err := refreshToken()
		if err != nil {
			return []dbEntities.Url{}, "", err
		}
	}
	subResp, err := fetchRdtResp(uType, lastNameId)
	if err != nil {
		return []dbEntities.Url{}, "", err
	}
	gfyUrls, imgurItems := sortUrls(&subResp)
	gfyItems, err := gfycat.ConvertNamesToUrls(gfyUrls)
	if err != nil {
		return []dbEntities.Url{}, "", err
	}
	for _, gfyItem := range gfyItems {
		item := dbEntities.Url{
			Value:   gfyItem.MobileUrl,
			Caption: gfyItem.GfyName,
		}
		imgurItems = append(imgurItems, item)
	}

	return imgurItems, getLastNameId(&subResp), nil
}

// Возвращает lastNameId для последующей пагинации
func getLastNameId(subResp *SubRedditResponse) string {
	if len(subResp.Data.Children) > 0 {
		return subResp.Data.Children[len(subResp.Data.Children)-1].Data.Name
	}
	return ""
}

// Проверяет действителен ли токен
func isGoodToken() bool {
	isGoodToken := true
	client := new(http.Client)
	req, _ := http.NewRequest("GET", config.RdtGetMeUrl, nil)
	req.Header.Set("User-Agent", config.RdtUserAgent)
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
	params.Set("username", config.RdtUserName)
	params.Set("password", config.RdtPassword)
	req, _ := http.NewRequest("POST", config.RdtGetTokenUrl, strings.NewReader(params.Encode()))
	req.SetBasicAuth(config.RdtClientId, config.RdtClientSecret)
	req.Header.Set("User-Agent", config.RdtUserAgent)
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

// Возвращает SubRedditResponse
func fetchRdtResp(uType string, lastNameId string) (SubRedditResponse, error) {
	var err error
	client := new(http.Client)
	req, _ := http.NewRequest("GET", config.RdtApiUrl+uType, nil)
	data := req.URL.Query()
	data.Set("limit", strconv.Itoa(config.RdtUrlsLimit))
	data.Set("after", lastNameId)
	req.URL.RawQuery = data.Encode()
	req.Header.Set("Authorization", "bearer "+tokenSample.Token)
	req.Header.Set("User-Agent", config.RdtUserAgent)
	resp, err := client.Do(req)
	if err != nil {
		return SubRedditResponse{}, err
	}
	defer resp.Body.Close()
	respBody, _ := ioutil.ReadAll(resp.Body)
	var subResp SubRedditResponse
	err = json.Unmarshal(respBody, &subResp)
	if err != nil {
		return SubRedditResponse{}, err
	}

	return subResp, nil
}

// Сортирует по картинкам и видосикам
func sortUrls(subResp *SubRedditResponse) ([]string, []dbEntities.Url) {
	var (
		gfyUrls []string
		items   []dbEntities.Url
	)
	for _, child := range subResp.Data.Children {
		if child.Data.Domain == config.GfycatDomain {
			namesSlice := strings.Split(child.Data.Url, "/")
			gfyUrls = append(gfyUrls, namesSlice[len(namesSlice)-1])
		}
		if child.Data.Domain == config.ImgurDomain {
			item := dbEntities.Url{
				Value:   child.Data.Url,
				Caption: child.Data.Title,
			}
			items = append(items, item)
		}
	}

	return gfyUrls, items
}
