package reddit

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"strings"

	log "github.com/sirupsen/logrus"
)

const gfycatApiURL = "https://api.gfycat.com/v1/gfycats/"

type gfyItem struct {
	GfyItem struct {
		MobileURL string `json:"mobileUrl"`
		GfyName   string `json:"gfyName"`
	} `json:"gfyItem"`
}

func (c *Converter) processingGfycat(d Data) (*Element, error) {
	client := new(http.Client)
	paths := strings.Split(d.URL, "/")
	hash := paths[len(paths)-1]
	req, _ := http.NewRequest("GET", gfycatApiURL+hash, nil)
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer func() {
		if err := resp.Body.Close(); err != nil {
			log.Errorf("error close body: %v", err)
		}
	}()
	respBody, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var item gfyItem
	err = json.Unmarshal(respBody, &item)
	if err != nil {
		return nil, err
	}

	return &Element{
		URL:     item.GfyItem.MobileURL,
		Caption: d.Title,
	}, nil
}
