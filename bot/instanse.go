package bot

import (
	"context"
	"fmt"
	"io/ioutil"
	"os"
	"os/signal"
	"strings"
	"time"

	"github.com/drklauss/boobsBot/config"
	"github.com/drklauss/boobsBot/model"
	"github.com/drklauss/boobsBot/telegram"
	"github.com/jinzhu/gorm"
	log "github.com/sirupsen/logrus"
)

const (
	dbPath = "./db.sqlite3"
	dbSQL  = "./db.sql"
)

// Bot is a BOT =).
type Bot struct {
	middlewares []Middleware
	handlers    map[string]HandlerFunc
	config      *config.Config
	db          *gorm.DB
}

// HandlerFunc is a command handler.
type HandlerFunc func(ctx context.Context, u *telegram.Update)

// New returns a bot Bot
func New(c *config.Config) (*Bot, error) {
	log.Debug("initialize bot...")
	if err := telegram.Init(c.Telegram); err != nil {
		return nil, err
	}

	log.Debug("connect to db...")
	db, err := gorm.Open("sqlite3", dbPath)
	if err != nil || !db.HasTable(&model.Item{}) {
		db, err = tryCreateDb()
		if err != nil {
			return nil, err
		}
	}
	log.Debug("db is connected")
	db = db.Debug()

	//if err != nil {
	//	fmt.Println(err)
	//	os.Exit(1)
	//}

	//if err != nil {
	//	return nil, err
	//}

	b := &Bot{
		config:   c,
		handlers: make(map[string]HandlerFunc),
		db:       db,
	}
	log.Debug("bot is created")
	return b, nil
}

// UseMiddlewares allows to use middlewares inside of the request
func (b *Bot) UseMiddlewares(ms ...Middleware) {
	b.middlewares = ms
}

// Handle associate path, i.e. incoming command with handler func
func (b *Bot) Handle(path string, h HandlerFunc) {
	b.handlers[path] = h
	// for presaved command
	b.handlers[path+"@"+b.config.Telegram.BotName] = h
}

// Run starts the bot
func (b *Bot) Run() {
	log.Debug("run bot...")
	ctx := context.Background()
	ctx, cancel := context.WithCancel(ctx)
	ctx = SetDB(ctx, b.db)
	updates := make(chan telegram.Update)
	defer close(updates)
	upTries := 0
	go func(ctx context.Context) {
		for {
			time.Sleep(time.Duration(b.config.Telegram.Time.Update) * time.Second)
			select {
			case <-ctx.Done():
				return
			default:
			}
			u, err := telegram.GetUpdateEntities()

			if err != nil {
				log.Warnf("could not get updates: %v", err)
				upTries++
				if upTries >= 10 {
					restPeriod := 20 * b.config.Telegram.Time.Update
					log.Warnf("gone for sleep for %d seconds", restPeriod)
					time.Sleep(time.Duration(restPeriod) * time.Second)
				}
				continue
			}
			for _, oneUp := range u {
				updates <- oneUp
			}
			upTries = 0
		}
	}(ctx)
	b.workerPool(ctx, &updates)
	signalCh := make(chan os.Signal, 1)
	signal.Notify(signalCh, os.Interrupt)
	<-signalCh
	cancel()
	log.Warnln("received sigkill, waiting for 3 seconds and ending...")
	time.Sleep(3 * time.Second)
}

func (b *Bot) workerPool(ctx context.Context, updates *chan telegram.Update) {
	for i := 1; i <= b.config.Telegram.Workers; i++ {
		go func(i int, ctx context.Context) {
			for upd := range *updates { // iterate over income message
				if isTooOldUpdate(&upd, b.config) {
					log.Debugf("sorry, it is a very old update")
					continue
				}

				log.Debugf("worker %d processing %d from %+v with text \"%s\"", i, upd.UpdateID, upd.Message.From, upd.Message.Text)
				hCommand, okCommand := b.handlers[upd.Message.Text]
				hCallback, okCallback := b.handlers[upd.CallBackQuery.Data]
				// simple text command handler
				if okCommand {
					cat := upd.Message.Text
					if len(cat) > 0 && cat[:1] == "/" { // cats saved in db without "/"
						cat = cat[1:]
					}
					if strings.Index(cat, "@"+b.config.Telegram.BotName) > 0 {
						cat = cat[:strings.Index(cat, "@"+b.config.Telegram.BotName)]
					}
					ctx = SetCategory(ctx, &cat)
					for _, middleware := range b.middlewares {
						hCommand = middleware(ctx, hCommand, &upd)
					}
					hCommand(ctx, &upd)
					continue
				}
				// inline callbacks handler
				if okCallback {
					cat := upd.CallBackQuery.Data
					if len(cat) > 0 && cat[:1] == "/" { // cats saved in db without "/"
						cat = cat[1:]
					}
					ctx = SetCategory(ctx, &cat)
					for _, m := range b.middlewares {
						hCallback = m(ctx, hCallback, &upd)
					}
					hCallback(ctx, &upd)
					continue
				}
				log.Infof("incorrect command processed: %v", upd.Message.Text)
			}
		}(i, ctx)
	}
}

func tryCreateDb() (*gorm.DB, error) {
	log.Debug("create db...")
	_, err := os.Create(dbPath)
	if err != nil {
		return nil, err
	}
	db, err := gorm.Open("sqlite3", dbPath)
	if err != nil {
		return nil, fmt.Errorf("could not open sqlite3: %v", err)
	}
	b, err := ioutil.ReadFile(dbSQL)
	if err != nil {
		return nil, fmt.Errorf("could not read sql file: %v", err)
	}
	err = db.Exec(string(b)).Error
	if err != nil {
		return nil, fmt.Errorf("could not execute sql: %v", err)
	}

	return db, nil
}

func isTooOldUpdate(upd *telegram.Update, config *config.Config) bool {
	if upd.Message.Date != 0 {
		log.Debugf("current time %d, message time %d, skip %d", time.Now().Unix(), upd.Message.Date, config.Telegram.Time.SkipMessages)
		return time.Now().Unix() > upd.Message.Date+config.Telegram.Time.SkipMessages
	}

	if upd.CallBackQuery.Message.Date != 0 {
		log.Debugf("current time %d, CallBackQuery message time %d, skip %d", time.Now().Unix(), upd.CallBackQuery.Message.Date, config.Telegram.Time.QuerySkipMessages)
		return time.Now().Unix() > upd.CallBackQuery.Message.Date+config.Telegram.Time.QuerySkipMessages
	}
	return false
}
