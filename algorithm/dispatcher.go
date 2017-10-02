package algorithm

import (
	"log"
	"time"

	"fmt"
	"strings"

	"github.com/drklauss/boobsBot/algorithm/config"
	"github.com/drklauss/boobsBot/algorithm/dataProvider"
	"github.com/drklauss/boobsBot/algorithm/telegram"
)

type Dispatcher struct {
	upResp    []telegram.Update
	dataProv  dataProvider.Provider
	debug     bool
	lastUpdId int
}

// Run запускает бота
func (d *Dispatcher) Run() {
	d.dataProv = new(dataProvider.Provider).Init(true)
	d.dataProv.CacheChatIds()

	for {
		var err error
		d.upResp, err = telegram.GetUpdateEntities(d.lastUpdId)
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
	upLen := len(d.upResp)
	if upLen > 0 {
		d.lastUpdId = d.upResp[upLen-1].UpdateId
		for i := 0; i < upLen; i++ {
			if time.Now().Unix() > d.upResp[i].Message.Date+config.TmSkipMessagesTime {
				continue
			}
			d.handleUpdate(d.upResp[i].Message)
		}
	}
}

// Обрабатывет команду
func (d *Dispatcher) handleUpdate(mes telegram.Message) {
	comName := strings.Split(mes.Text, config.TmFullBotName)
	command := strings.Replace(comName[0], "/", "", -1)
	d.dataProv.CreateChatEntry(mes)
	if d.isDebug(mes.From.Id, command) {
		return
	}
	switch command {
	case config.TmHelloCmd:
		telegram.SendMessage(mes.Chat.Id, fmt.Sprintf("Hello, %s!", mes.From.FirstName))
	case config.TmHotCmd:
		u := d.dataProv.GetUrl(mes.Chat)
		telegram.SendDocument(mes.Chat.Id, u)
	case config.TmTopViewersCmd:
		if mes.From.Id != config.TmDevUserId {
			return
		}
		b := d.dataProv.GetTopViewers4Tm()
		telegram.SendMessage(mes.Chat.Id, fmt.Sprintf("%s", b))

	}
}

// Проверяет включен ли дебаг для разработчика
func (d *Dispatcher) isDebug(userId int, com string) bool {
	if userId != config.TmDevUserId {
		return false
	}
	switch com {
	case config.TmDebugStartCmd:
		d.debug = true
	case config.TmDebugEndCmd:
		d.debug = false
	}

	return d.debug
}
