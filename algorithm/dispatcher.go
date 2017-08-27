package algorithm

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"strconv"
	"time"

	"github.com/boobsBot/algorithm/config"
	"github.com/boobsBot/algorithm/dataProvider"
	"github.com/boobsBot/algorithm/telegram"
)

type Dispatcher struct {
	updateResp   []telegram.Update
	urlProvider  dataProvider.Provider
	motions      []string
	lastUpdateId int
}

func (d *Dispatcher) Run() {
	d.urlProvider = new(dataProvider.Provider).Init()
	for {
		err := d.initUpdateEntities()
		if err != nil {
			log.Println(err)
			continue
		}
		d.processUpdates()
		time.Sleep(config.TmUpdateTime * time.Second)
	}
}

// Получение обновлений
func (d *Dispatcher) initUpdateEntities() error {
	d.updateResp = []telegram.Update{}
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
	if err != nil {
		return err
	}
	var response telegram.Response
	if err = json.Unmarshal(responseBody, &response); err != nil {
		return err
	}
	d.updateResp = response.Result

	return nil
}

// Обрабатывает полученные обновления
func (d *Dispatcher) processUpdates() {
	upLen := len(d.updateResp)
	if upLen > 0 {
		d.lastUpdateId = d.updateResp[upLen-1].UpdateId
		for i := 0; i < upLen; i++ {
			if time.Now().Unix() > d.updateResp[i].Message.Date+config.TmSkipMessagesTime {
				continue
			}
			d.handleUpdate(d.updateResp[i])
		}
	}
}
