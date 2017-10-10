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
		mes := telegram.MessageSend{
			ChatId:         mes.Chat.Id,
			Text:           "Hello",
			KeyboardMarkup: telegram.ReplyKeyboardRemove{RemoveKeyboard: true},
		}
		telegram.SendMessage(mes)
	case config.TmNSFWCmd:
		u := d.dataProv.GetUrl(mes.Chat, config.TmNSFWCmd)
		doc := telegram.DocumentSend{
			ChatId:         mes.Chat.Id,
			Caption:        u.Caption,
			KeyboardMarkup: telegram.ReplyKeyboardRemove{RemoveKeyboard: true},
			Url:            u.Value,
		}
		fmt.Printf("%+v\n", doc)
		telegram.SendDocument(doc)
	case config.TmRealGirlsCmd:
		u := d.dataProv.GetUrl(mes.Chat, config.TmRealGirlsCmd)
		doc := telegram.DocumentSend{
			ChatId:         mes.Chat.Id,
			Caption:        u.Caption,
			KeyboardMarkup: telegram.ReplyKeyboardRemove{RemoveKeyboard: true},
			Url:            u.Value,
		}
		telegram.SendDocument(doc)
	case config.TmCelebCmd:
		u := d.dataProv.GetUrl(mes.Chat, config.TmCelebCmd)
		doc := telegram.DocumentSend{
			ChatId:         mes.Chat.Id,
			Caption:        u.Caption,
			KeyboardMarkup: telegram.ReplyKeyboardRemove{RemoveKeyboard: true},
			Url:            u.Value,
		}
		telegram.SendDocument(doc)
		// =-_-= Админские команды =-_-=
	case config.TmAdmin:
		if mes.From.Id != config.TmDevUserId {
			return
		}
		mes := telegram.MessageSend{
			ChatId:         mes.Chat.Id,
			KeyboardMarkup: d.getAdminKeyboard(),
			Text:           fmt.Sprintf("Hello, %s!", mes.From.FirstName),
		}
		telegram.SendMessage(mes)
	case config.TmTopViewersCmd:
		if mes.From.Id != config.TmDevUserId {
			return
		}
		b := d.dataProv.GetTopViewers4Tm()
		mes := telegram.MessageSend{
			ChatId:         mes.Chat.Id,
			Text:           fmt.Sprintf("%s", b),
			KeyboardMarkup: telegram.ReplyKeyboardRemove{RemoveKeyboard: true},
		}
		telegram.SendMessage(mes)
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
	case config.TmDebugStopCmd:
		d.debug = false
	}

	return d.debug
}

func (d *Dispatcher) getAdminKeyboard() telegram.ReplyKeyboardMarkup {
	btnTopViewers := telegram.KeyboardButton{
		Text:            config.TmTopViewersCmd,
		RequestContact:  false,
		RequestLocation: false,
	}
	btnDebugStart := telegram.KeyboardButton{
		Text:            config.TmDebugStartCmd,
		RequestContact:  false,
		RequestLocation: false,
	}
	btnDebugStop := telegram.KeyboardButton{
		Text:            config.TmDebugStopCmd,
		RequestContact:  false,
		RequestLocation: false,
	}
	btns := []telegram.KeyboardButton{btnTopViewers, btnDebugStart, btnDebugStop}
	keyboards := [][]telegram.KeyboardButton{btns}
	rm := telegram.ReplyKeyboardMarkup{
		Keyboard:        keyboards,
		OneTimeKeyboard: true,
		ResizeKeyboard:  true,
		Selective:       false,
	}
	return rm
}
