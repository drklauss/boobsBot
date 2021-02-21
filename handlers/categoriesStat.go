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

func CategoriesStat(ctx context.Context, u *telegram.Update) {
	db, err := bot.GetDB(ctx)
	if err != nil {
		log.Warnln(err)
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
	if err := db.Raw(sql).Scan(&categoriesStat).Error; err != nil {
		log.Errorf("could not get categories stat: %s", err)
		return
	}

	buf := new(bytes.Buffer)
	buf.WriteString("\xF0\x9F\x93\x8A Categories Items Count:\n")
	for _, v := range categoriesStat {
		s := fmt.Sprintf("%s - %d items\n", v.Category, v.Count)
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
