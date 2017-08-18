package algorithm

import (
	"boobsBot/entities"
	"net/http"
	"log"
	"io/ioutil"
	"encoding/json"
	"net/url"
	"strconv"
	"boobsBot/config"
	"time"
	"fmt"
)

type Dispatcher struct {
	updateResp []entities.Update
	lastUpdate int
}

func (d *Dispatcher) Run() {
	for {
		time.Sleep(config.RequestTime * time.Second)
		log.Println("New request arrived")
		err := d.initUpdateEntities()
		if err != nil {
			log.Println(err)
			continue
		}
		d.processUpdates()
	}
}

// Получение обнвлений
func (d *Dispatcher) initUpdateEntities() error {
	d.updateResp = []entities.Update{}
	u, _ := url.ParseRequestURI(config.ApiUrl + config.Token)
	u.Path += "/getUpdates"
	params := url.Values{}
	params.Set("offset", strconv.Itoa(d.lastUpdate+1))
	u.RawQuery = params.Encode()
	fmt.Println(u.String())
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

func (d *Dispatcher) processUpdates() {
	fmt.Printf("%#v\n", len(d.updateResp))
	fmt.Printf("%#v\n", d.lastUpdate)
	upLen := len(d.updateResp)
	if upLen > 0 {
		d.lastUpdate = d.updateResp[upLen-1].UpdateId
		for i := 0; i < upLen; i++ {
			handleUpdate(d.updateResp[i])
		}
	}
}
