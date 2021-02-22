package handlers

import (
	"fmt"

	"github.com/drklauss/boobsBot/bot"
	"github.com/drklauss/boobsBot/config"
	"github.com/drklauss/boobsBot/model"
	"github.com/drklauss/boobsBot/reddit"
	"github.com/drklauss/boobsBot/telegram"
	"github.com/jinzhu/gorm"
	log "github.com/sirupsen/logrus"
)

const (
	updateCount = 10
)

// UpdateHandler gets links from reddit by categories and save them.
func UpdateHandler(req bot.HandlerRequest) {
	if !checkAdmin(req) {
		return
	}
	if err := reddit.Init(config.Get().Reddit); err != nil {
		log.Errorf("could not initialize reddit client: %v", err)
		return
	}
	cc := config.Get().Reddit.Categories
	for _, c := range cc {
		// goroutine for each category
		// go
		// remove goroutine - dont know reddit api limits for parallel requests
		updateCategory(c, req.DB, req.Update.Message.Chat.ID)
	}
}

type updateSubRedditResult struct {
	fetched   int
	converted int
	inserted  int
	errors    int
}

// updateCategory gets new posts from categories' subreddits.
func updateCategory(c config.Category, db *gorm.DB, chatID int) {
	results := make(map[string]updateSubRedditResult)
	for _, subredditURI := range c.Source {
		results[subredditURI] = updateSubReddit(c.Name, subredditURI, db)
	}
	var message string
	message = fmt.Sprintf("Update has been completed for category '%s'\n", c.Name)
	for subredditURI, result := range results {
		message += fmt.Sprintf("- fetched %d converted %d inserted %d items for '%s' from subreddit '%s'", result.fetched, result.converted, result.inserted, c.Name, subredditURI)
	}
	log.Infoln(message)
	ms := telegram.MessageConfig{
		ChatID: chatID,
		Text:   message,
	}
	if err := ms.Send(); err != nil {
		log.Errorf("error send message with update result to telegram: %v", err)
	}
}

// updateSubReddit gets new posts from subreddit and saves them to db.
func updateSubReddit(cat string, subredditURI string, db *gorm.DB) (result updateSubRedditResult) {
	last := ""
	countErrs := 0
	for result.inserted < updateCount && countErrs < 3 {
		resp, err := reddit.GetItems(subredditURI, last)
		if err != nil {
			countErrs++
			log.Errorf("could not fetch items for '%s' from subreddit '%s': %v", cat, subredditURI, err)
			continue
		}
		last = resp.Data.Children[len(resp.Data.Children)-1].Data.Name
		result.fetched += len(resp.Data.Children)
		els := resp.Convert()
		if len(els) == 0 {
			continue
		}
		result.converted += len(els)
		item := model.NewItem(db)
		insertedCount, err := item.Save(cat, els)
		if err != nil {
			countErrs++
			log.Errorf("could not save items: %v", err)
		}
		result.inserted += insertedCount
	}
	log.Infof("fetched %d converted %d inserted %d items for '%s' from subreddit '%s'", result.fetched, result.converted, result.inserted, cat, subredditURI)

	return result
}
