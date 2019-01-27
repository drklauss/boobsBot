package handler

import (
	"context"
	"fmt"

	"github.com/drklauss/boobsBot/bot"
	"github.com/drklauss/boobsBot/model"

	"github.com/drklauss/boobsBot/config"
	"github.com/drklauss/boobsBot/reddit"
	"github.com/drklauss/boobsBot/telegram"
	"github.com/leesper/holmes"
)

const (
	updateCount = 100
)

// Update gets links from reddit by categories and save them
func Update(ctx context.Context, u *telegram.Update) {
	conn, err := bot.GetDB(ctx)

	if err != nil {
		holmes.Warnln(err)
		return
	}
	if err = reddit.Init(config.Get().Reddit); err != nil {
		holmes.Warnf("could not initialize reddit client: %v", err)
		return
	}
	cc := config.Get().Reddit.Categories
	item := model.NewItem(conn)
	for _, c := range cc {
		go func(c config.Category) { // goroutine for each category
			pathSaved := make(map[string]int)
			for _, urlPath := range c.Source {
				pathSaved[urlPath] = 0 // count saved items for each path for category
				countErrs := 0
				for pathSaved[urlPath] < updateCount && countErrs < 3 {
					resp, err := reddit.GetItems(urlPath)
					if err != nil {
						countErrs++
						holmes.Warnf("could not fetch items for %s with path %s: %v", c.Name, urlPath, err)
						continue
					}
					els := resp.Convert()
					if len(els) == 0 {
						continue
					}
					count, err := item.Save(c.Name, els)
					if err != nil {
						countErrs++
						holmes.Warnf("could not save items: %v", err)
					}
					pathSaved[urlPath] += count
					holmes.Infof("fetched %d items for \"%s\" via \"%s\"", count, c.Name, urlPath)
				}
			}
			var t string
			for path, count := range pathSaved {
				t = fmt.Sprintf("total fetched %d items for \"%s\" by %s", count, c, path)
				holmes.Infoln(t)
				ms := telegram.MessageSend{
					ChatID: config.Get().Telegram.Admin,
					Text:   t,
				}
				ms.Send()
			}
		}(c)
	}
}
