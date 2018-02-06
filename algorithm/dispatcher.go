package algorithm

import (
	"log"
	"time"

	"strings"

	"bytes"

	"github.com/drklauss/boobsBot/algorithm/config"
	"github.com/drklauss/boobsBot/algorithm/dataProvider"
	"github.com/drklauss/boobsBot/algorithm/dataProvider/stat"
	"github.com/drklauss/boobsBot/algorithm/telegram"
	"sync"
)

type Dispatcher struct {
	upResp    []telegram.Update
	dataProv  dataProvider.Provider
	debug     bool // todo в гитхабе лежит многопоточная версия, которая некорректно проставляет этот флаг
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
		}
		d.processUpdates()
		time.Sleep(config.TmUpdateTime * time.Second)
	}
}

// Обрабатывает полученные обновления
func (d *Dispatcher) processUpdates() {
	upLen := len(d.upResp)
	if upLen <= 0 {
		return
	}
	d.lastUpdId = d.upResp[upLen-1].UpdateId
	var wg sync.WaitGroup
	for i := 0; i < upLen; i++ {
		if time.Now().Unix() > d.upResp[i].Message.Date+config.TmSkipMessagesTime {
			continue
		}
		wg.Add(1)
		go func(mes telegram.Message) {
			d.handleUpdate(mes)
			wg.Done()
		}(d.upResp[i].Message)
	}
	wg.Wait()
}

// Обрабатывет команду
func (d *Dispatcher) handleUpdate(mes telegram.Message) {
	log.Printf("New Update: %+v \n", mes)
	comName := strings.Split(mes.Text, config.TmFullBotName)
	command := strings.Replace(comName[0], "/", "", -1)
	d.dataProv.CreateChatEntry(mes)
	if d.isHandledDebugCommand(mes, command) {
		return
	}
	switch command {
	case config.TmHelpCmd:
		mes := telegram.MessageSend{
			ChatId:         mes.Chat.Id,
			Text:           d.getHelpMsg(),
			KeyboardMarkup: telegram.ReplyKeyboardRemove{RemoveKeyboard: true},
		}
		mes.Send()
	case config.TmRateCmd:
		mes := telegram.MessageSend{
			ChatId:         mes.Chat.Id,
			Text:           "Please, give us 5 \xE2\xAD\x90\xE2\xAD\x90\xE2\xAD\x90\xE2\xAD\x90\xE2\xAD\x90 at \n https://telegram.me/storebot?start=DornBot",
			KeyboardMarkup: telegram.ReplyKeyboardRemove{RemoveKeyboard: true},
		}
		mes.Send()
	case config.TmNSFWCmd:
		u := d.dataProv.GetUrl(mes.Chat, config.TmNSFWCmd)
		content := telegram.MediaSend{
			ChatId:         mes.Chat.Id,
			Caption:        u.Caption,
			KeyboardMarkup: telegram.ReplyKeyboardRemove{RemoveKeyboard: true},
			Url:            u.Value,
		}
		log.Printf("Send Content: %+v \n", content)
		content.Send()
	case config.TmRealGirlsCmd:
		u := d.dataProv.GetUrl(mes.Chat, config.TmRealGirlsCmd)
		content := telegram.MediaSend{
			ChatId:         mes.Chat.Id,
			Caption:        u.Caption,
			KeyboardMarkup: telegram.ReplyKeyboardRemove{RemoveKeyboard: true},
			Url:            u.Value,
		}
		log.Printf("Send Content: %+v \n", content)
		content.Send()
	case config.TmCelebCmd:
		u := d.dataProv.GetUrl(mes.Chat, config.TmCelebCmd)
		content := telegram.MediaSend{
			ChatId:         mes.Chat.Id,
			Caption:        u.Caption,
			KeyboardMarkup: telegram.ReplyKeyboardRemove{RemoveKeyboard: true},
			Url:            u.Value,
		}
		log.Printf("Send Content: %+v \n", content)
		content.Send()
		// =-_-= Админские команды =-_-=
	case config.TmUpdateCmd:
		if mes.From.Id != config.TmDevUserId {
			return
		}
		s := d.dataProv.UpdateAll()
		mes := telegram.MessageSend{
			ChatId:         mes.Chat.Id,
			Text:           s,
			KeyboardMarkup: telegram.ReplyKeyboardRemove{RemoveKeyboard: false},
		}
		log.Printf("Update Info: %+v \n", strings.Replace(s, "\n", " ", -1))
		mes.Send()
	case config.TmTopViewersCmd:
		if mes.From.Id != config.TmDevUserId {
			return
		}
		s := d.dataProv.GetTopViewers(stat.TelegramFmt)
		mes := telegram.MessageSend{
			ChatId:         mes.Chat.Id,
			Text:           s,
			KeyboardMarkup: telegram.ReplyKeyboardRemove{RemoveKeyboard: false},
		}
		mes.Send()
	case config.TmTotalLinksCmd:
		if mes.From.Id != config.TmDevUserId {
			return
		}
		s := d.dataProv.GetTotalLinks(stat.TelegramFmt)
		mes := telegram.MessageSend{
			ChatId:         mes.Chat.Id,
			Text:           s,
			KeyboardMarkup: telegram.ReplyKeyboardRemove{RemoveKeyboard: false},
		}

		mes.Send()
	}
}

// Обработана ли команда дебага
func (d *Dispatcher) isHandledDebugCommand(mes telegram.Message, com string) bool {
	if mes.From.Id != config.TmDevUserId {
		return false
	}
	switch com {
	case config.TmDebugStartCmd:
		d.debug = true
	case config.TmDebugStopCmd:
		d.debug = false
	case config.TmAdmin:
		mes := telegram.MessageSend{
			ChatId:         mes.Chat.Id,
			KeyboardMarkup: d.getAdminKeyboard(),
			Text:           "Админская клавиатура",
		}
		mes.Send()
		return true
	}

	return d.debug
}

// Возвращает админскую инлайновую клавиатуру
func (d *Dispatcher) getAdminKeyboard() telegram.ReplyKeyboardMarkup {
	btnTopViewers := telegram.KeyboardButton{
		Text:            config.TmTopViewersCmd,
		RequestContact:  false,
		RequestLocation: false,
	}
	btnTotalLinks := telegram.KeyboardButton{
		Text:            config.TmTotalLinksCmd,
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
	btnUpdate := telegram.KeyboardButton{
		Text:            config.TmUpdateCmd,
		RequestContact:  false,
		RequestLocation: false,
	}
	btns := []telegram.KeyboardButton{btnTopViewers, btnTotalLinks, btnDebugStart, btnDebugStop, btnUpdate}
	keyboards := [][]telegram.KeyboardButton{btns}
	rm := telegram.ReplyKeyboardMarkup{
		Keyboard:        keyboards,
		OneTimeKeyboard: true,
		ResizeKeyboard:  true,
		Selective:       false,
	}
	return rm
}

func (d *Dispatcher) getHelpMsg() string {
	buf := new(bytes.Buffer)
	buf.WriteString("We offer you hot girls from popular categories \xF0\x9F\x98\x9C \n")
	buf.WriteString("Available commands:\n")
	buf.WriteString("/nsfw Gives you not safe for work content \n")
	buf.WriteString("/real_girls Gives real girls \n")
	buf.WriteString("/celeb Gives you naked celebrities \n")
	buf.WriteString("/rate Gives a link to rate for us \n")
	buf.WriteString("/help Bot help \n")
	buf.WriteString("\n\n Do not forget to rate us 5 \xE2\xAD\x90\xE2\xAD\x90\xE2\xAD\x90\xE2\xAD\x90\xE2\xAD\x90 at \n")
	buf.WriteString("https://telegram.me/storebot?start=DornBot")

	return buf.String()
}
