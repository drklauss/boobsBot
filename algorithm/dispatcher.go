package algorithm

import (
	"log"
	"time"

	"fmt"
	"strings"

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

// Run запускает бота
func (d *Dispatcher) Run() {
	d.urlProvider = new(dataProvider.Provider).Init()
	for {
		var err error
		d.updateResp, err = telegram.GetUpdateEntities(d.lastUpdateId)
		if err != nil {
			log.Println(err)
			continue
		}
		d.processUpdates()
		time.Sleep(config.TmUpdateTime * time.Second)
	}
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

// Обрабатывет команду
func (d *Dispatcher) handleUpdate(update telegram.Update) {
	comName := strings.Split(update.Message.Text, config.TmFullBotName)
	command := strings.Replace(comName[0], "/", "", -1)
	switch command {
	case config.Hello:
		telegram.SendMessage(update.Message.Chat.Id, fmt.Sprintf("Well, Hello, %s!", update.Message.From.FirstName))
	case config.New:
		telegram.SendDocument(update.Message.Chat.Id, d.urlProvider.GetUrl(config.New))
	case config.Hot:
		telegram.SendDocument(update.Message.Chat.Id, d.urlProvider.GetUrl(config.Hot))
	case config.Top:
		telegram.SendDocument(update.Message.Chat.Id, d.urlProvider.GetUrl(config.Top))
	}
}
