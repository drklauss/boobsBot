package handlers

import (
	"bytes"
	"context"
	"fmt"

	"github.com/drklauss/boobsBot/bot"
	"github.com/drklauss/boobsBot/model"
	"github.com/drklauss/boobsBot/telegram"
	log "github.com/sirupsen/logrus"
)

func TopViewers(ctx context.Context, u *telegram.Update) {
	db, err := bot.GetDB(ctx)
	if err != nil {
		log.Warnln(err)
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
					count(c.id) AS count
			FROM %s v
			LEFT JOIN %s c on
				c.id = v.chatId
			GROUP BY c.id
			ORDER BY 3 DESC
			LIMIT 10`, (new(model.View)).TableName(), (new(model.Chat)).TableName())
	if err := db.Raw(sql).Scan(&tVs).Error; err != nil {
		log.Errorf("could not get top viewers stat: %s", err)
		return
	}

	buf := new(bytes.Buffer)
	buf.WriteString("\xF0\x9F\x93\x8A Top Viewers Report:\n")
	for k, v := range tVs {
		s := fmt.Sprintf("%d. %s (%s) - %d views\n", k+1, v.Title, v.Type, v.Count)
		buf.WriteString(s)
	}

	mes := telegram.MessageSend{
		ChatID:         u.Message.Chat.ID,
		Text:           buf.String(),
		KeyboardMarkup: telegram.ReplyKeyboardRemove{RemoveKeyboard: true},
	}

	if err := mes.Send(); err != nil {
		log.Warnf("could not send message: %v", err)
	}
}
