package algorithm

import (
	"boobsBot/algorithm/gfycat"
	"boobsBot/config"
	"boobsBot/entities"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"strconv"
	"time"
)

type Dispatcher struct {
	updateResp     []entities.Update
	urlProvider    gfycat.Provider
	motions        []string
	lastUpdateId   int
	lastUpdateTime int64
}

func (d *Dispatcher) Run() {
	d.urlProvider = new(gfycat.Provider).Init()
	for {
		time.Sleep(config.TmUpdateTime * time.Second)
		err := d.initUpdateEntities()
		if err != nil {
			log.Println(err)
			continue
		}
		d.processUpdates()
	}
}

// Получение обновлений
func (d *Dispatcher) initUpdateEntities() error {
	d.updateResp = []entities.Update{}
	u, _ := url.ParseRequestURI(config.TmApiUrl + config.TmToken)
	u.Path += "/getUpdates"
	params := url.Values{}
	params.Set("offset", strconv.Itoa(d.lastUpdateId+1))
	u.RawQuery = params.Encode()
	resp, err := http.Get(u.String())
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	responseBody, err := ioutil.ReadAll(resp.Body)
	var response entities.Response
	err = json.Unmarshal(responseBody, &response)
	d.updateResp = response.Result

	return err
}

// Обрабатывает обновления
// todo: сделать проверку на время, т.к. если приложение выключено, и написать много сообщений боту - жестко спаммит
// todo: здесь же сделать проверку на
func (d *Dispatcher) processUpdates() {
	fmt.Printf("%#v\n", len(d.updateResp))
	fmt.Printf("%#v\n", d.lastUpdateId)
	upLen := len(d.updateResp)
	if upLen > 0 {
		d.lastUpdateId = d.updateResp[upLen-1].UpdateId
		d.lastUpdateTime = d.updateResp[upLen-1].Message.Date
		for i := 0; i < upLen; i++ {
			d.handleUpdate(d.updateResp[i])
		}
	}
}
