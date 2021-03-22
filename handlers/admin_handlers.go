package handlers

import (
	"bytes"
	"fmt"

	"github.com/drklauss/boobsBot/bot"
	"github.com/drklauss/boobsBot/model"
	"github.com/drklauss/boobsBot/telegram"
	log "github.com/sirupsen/logrus"
)

const (
	Admin          = "/admin"
	DebugStart     = "/debugStart"
	DebugStop      = "/debugStop"
	TopViewers     = "/topViewers"
	MonthlyStat    = "/monthlyStat"
	CategoriesStat = "/categoriesStat"
	Update         = "/update"
)

func checkAdmin(req bot.HandlerRequest) bool {
	isAdmin := false
	for _, adminID := range req.Config.Telegram.Admin {
		if req.Update.Message.From.ID == adminID {
			isAdmin = true
			break
		}
	}

	return isAdmin
}

// DebugStartHandler handles start debug request.
func DebugStartHandler(req bot.HandlerRequest) {
	if !checkAdmin(req) {
		return
	}
	bot.SetDebug(true)
	mes := telegram.MessageConfig{
		ChatID: req.Update.Message.Chat.ID,
		Text:   "debug enabled",
	}
	if err := mes.Send(); err != nil {
		log.Warnf("could not send message: %v", err)
	}
}

// DebugStopHandler handles stop debug request.
func DebugStopHandler(req bot.HandlerRequest) {
	if !checkAdmin(req) {
		return
	}
	bot.SetDebug(false)
	mes := telegram.MessageConfig{
		ChatID: req.Update.Message.Chat.ID,
		Text:   "debug disabled",
	}
	if err := mes.Send(); err != nil {
		log.Warnf("could not send message: %v", err)
	}
}

// DebugStopHandler handles show admin menu.
func AdminHandler(req bot.HandlerRequest) {
	if !checkAdmin(req) {
		return
	}
	bot.SetDebug(false)
	mes := telegram.MessageConfig{
		ChatID:         req.Update.Message.Chat.ID,
		Text:           "admin keyboard open",
		KeyboardMarkup: GetAdminKeyboard(),
	}
	if err := mes.Send(); err != nil {
		log.Warnf("could not send message: %v", err)
	}
}

// DebugStopHandler handles show categories statistics.
func CategoriesStatHandler(req bot.HandlerRequest) {
	if !checkAdmin(req) {
		return
	}
	var categoriesStat []struct {
		Category string
		Count    int
	}

	sql := fmt.Sprintf(`
			SELECT category, COUNT(id) AS count
			FROM %s
			GROUP BY category
			ORDER BY 2 DESC`, (new(model.Item)).TableName())
	if err := req.DB.Raw(sql).Scan(&categoriesStat).Error; err != nil {
		log.Errorf("could not get categories stat: %s", err)
		return
	}

	buf := new(bytes.Buffer)
	buf.WriteString("\xF0\x9F\x93\x8A Categories Items Count:\n")
	for _, v := range categoriesStat {
		s := fmt.Sprintf("%s - %d items\n", v.Category, v.Count)
		buf.WriteString(s)
	}

	mes := telegram.MessageConfig{
		ChatID:         req.Update.Message.Chat.ID,
		Text:           buf.String(),
		KeyboardMarkup: telegram.ReplyKeyboardRemove{RemoveKeyboard: true},
	}
	if err := mes.Send(); err != nil {
		log.Warnf("could not send message: %v", err)
	}
}

// DebugStopHandler handles show top viewers request.
func TopViewersHandler(req bot.HandlerRequest) {
	if !checkAdmin(req) {
		return
	}
	var tVs []struct {
		Title string
		Type  string
		Count int
	}

	sql := fmt.Sprintf(`
			SELECT c.title,
					c.type,
					COUNT(c.id) AS count
			FROM %s v
			LEFT JOIN %s c ON
				c.id = v.chatId
			GROUP BY c.id
			ORDER BY 3 DESC
			LIMIT 10`, (new(model.View)).TableName(), (new(model.Chat)).TableName())
	if err := req.DB.Raw(sql).Scan(&tVs).Error; err != nil {
		log.Errorf("could not get top viewers stat: %s", err)
		return
	}

	buf := new(bytes.Buffer)
	buf.WriteString("\xF0\x9F\x93\x8A Top Viewers Report:\n")
	for k, v := range tVs {
		s := fmt.Sprintf("%d. %d views by %s '%s'\n", k+1, v.Count, v.Type, v.Title)
		buf.WriteString(s)
	}

	mes := telegram.MessageConfig{
		ChatID:         req.Update.Message.Chat.ID,
		Text:           buf.String(),
		KeyboardMarkup: telegram.ReplyKeyboardRemove{RemoveKeyboard: true},
	}

	if err := mes.Send(); err != nil {
		log.Warnf("could not send message: %v", err)
	}
}

// MonthlyStatHandler handles show monthly statistics request.
func MonthlyStatHandler(req bot.HandlerRequest) {
	if !checkAdmin(req) {
		return
	}
	var stat []struct {
		Month int
		Day   int
		Count int
	}

	sql := fmt.Sprintf(`
		SELECT
		
			strftime('%%m', datetime(v.requestDate, 'unixepoch')) AS month,
			strftime('%%d', datetime(v.requestDate, 'unixepoch')) AS day,
			COUNT(v.itemId) AS count
		FROM %s AS v
		GROUP BY strftime('%%Y-%%m-%%d', datetime(v.requestDate, 'unixepoch'))
		ORDER BY 
			strftime('%%Y', datetime(v.requestDate, 'unixepoch')) DESC, 
			strftime('%%m', datetime(v.requestDate, 'unixepoch')) DESC, 
			strftime('%%d', datetime(v.requestDate, 'unixepoch')) DESC
		LIMIT 30`, (new(model.View)).TableName())
	if err := req.DB.Raw(sql).Scan(&stat).Error; err != nil {
		log.Errorf("could not get monthly statistics stat: %s", err)
		return
	}

	buf := new(bytes.Buffer)
	buf.WriteString("\xF0\x9F\x93\x8A Monthly Statistics:\n")
	for _, v := range stat {
		s := fmt.Sprintf("%d.%d - %d views\n", v.Day, v.Month, v.Count)
		buf.WriteString(s)
	}

	mes := telegram.MessageConfig{
		ChatID:         req.Update.Message.Chat.ID,
		Text:           buf.String(),
		KeyboardMarkup: telegram.ReplyKeyboardRemove{RemoveKeyboard: true},
	}

	if err := mes.Send(); err != nil {
		log.Warnf("could not send message: %v", err)
	}
}
