package algorithm

import (
	"encoding/json"
	"fmt"
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
	updateResp     []telegram.Update
	urlProvider    dataProvider.Provider
	motions        []string
	lastUpdateId   int
	lastUpdateTime int64
}

func (d *Dispatcher) Run() {
	d.urlProvider = new(dataProvider.Provider).Init()
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
	var response telegram.Response
	err = json.Unmarshal(responseBody, &response)
	d.updateResp = response.Result

	return err
}

// Обрабатывает полученные обновления
func (d *Dispatcher) processUpdates() {
	fmt.Printf("%#v\n", len(d.updateResp))
	fmt.Printf("%#v\n", d.lastUpdateId)
	upLen := len(d.updateResp)
	if upLen > 0 {
		// todo: поставить время обновления
		d.lastUpdateId = d.updateResp[upLen-1].UpdateId
		d.lastUpdateTime = d.updateResp[upLen-1].Message.Date
		for i := 0; i < upLen; i++ {
			d.handleUpdate(d.updateResp[i])
		}
	}
}
