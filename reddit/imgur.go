package reddit

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"

	"github.com/drklauss/boobsBot/config"
)

const imgurApiURL = "https://api.imgur.com/3/"

type imgurItem struct {
	Data struct {
		Mp4  string `json:"mp4"`
		Link string `json:"link"`
	} `json:"data"`
	Success bool `json:"success"`
}

func (c *Converter) processingImgur(d Data) (*Element, error) {
	if reMedia.MatchString(d.URL) {
		el := &Element{
			URL:     d.URL,
			Caption: d.Title,
		}
		return el, nil
	}
	apiURL := imgurApiURL + "image/"
	if strings.Contains(d.URL, "/gallery/") {
		apiURL = imgurApiURL + "gallery/"
	}
	var item imgurItem
	client := new(http.Client)
	paths := strings.Split(d.URL, "/")
	hash := paths[len(paths)-1]
	// first try to get link by image hash
	req, _ := http.NewRequest("GET", apiURL+hash, nil)
	req.Header.Add("Authorization", "Client-ID "+config.Get().Imgur.ClientID)
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	respBody, err := ioutil.ReadAll(resp.Body)
	defer resp.Body.Close()
	if err != nil {
		return nil, err
	}
	err = json.Unmarshal(respBody, &item)
	if err != nil {
		return nil, err
	}
	if item.Data.Link == "" {
		return nil, fmt.Errorf("could not fill URL for element: empty link")
	}
	el := &Element{
		URL:     item.Data.Link,
		Caption: d.Title,
	}
	if item.Data.Mp4 != "" {
		el.URL = item.Data.Mp4
	}

	return el, nil
}
