package bot

import (
	"context"
	"fmt"
	"io/ioutil"
	"os"
	"os/signal"
	"time"

	"github.com/drklauss/boobsBot/config"
	"github.com/drklauss/boobsBot/model"
	"github.com/drklauss/boobsBot/telegram"
	"github.com/jinzhu/gorm"
	log "github.com/sirupsen/logrus"
)

const (
	dbPath = "./db.sqlite3"
	dbSQL  = "./deployments/db.sql"
)

// Bot is a BOT =).
type Bot struct {
	middlewares []Middleware
	handlers    map[string]HandlerFunc
	config      *config.Config
	db          *gorm.DB
}

type HandlerRequest struct {
	Ctx    context.Context
	Update *telegram.Update
	DB     *gorm.DB
	Config *config.Config
}

// HandlerFunc is a command handler.
type HandlerFunc func(req HandlerRequest)

// New returns a bot Bot
func New(c *config.Config, debug bool) (*Bot, error) {
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
	if debug {
		db = db.Debug()
	}

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
	var updates = make(chan telegram.Update)
	defer close(updates)
	go b.getUpdates(ctx, updates)
	b.runWorkers(ctx, updates)

	var signalCh = make(chan os.Signal, 1)
	signal.Notify(signalCh, os.Interrupt)
	<-signalCh
	cancel()
	log.Warnln("received sigkill, waiting for 3 seconds and ending...")
	time.Sleep(3 * time.Second)
}

func (b *Bot) getUpdates(ctx context.Context, updates chan<- telegram.Update) {
	var attempts = 0
	var restPeriod = 20 * b.config.Telegram.Time.Update
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
			attempts++
			if attempts >= 10 {
				log.Warnf("gone for sleep for %d seconds", restPeriod)
				time.Sleep(time.Duration(restPeriod) * time.Second)
			}
			continue
		}
		for _, oneUp := range u {
			updates <- oneUp
		}
		attempts = 0
	}
}

func (b *Bot) runWorkers(ctx context.Context, updates <-chan telegram.Update) {
	for i := 1; i <= b.config.Telegram.Workers; i++ {
		go b.worker(i, updates, ctx)
	}
}

func (b *Bot) worker(worker int, updates <-chan telegram.Update, ctx context.Context) {
	var text string
	for upd := range updates { // iterate over incoming message
		if isTooOldUpdate(&upd, b.config) {
			log.Debugf("sorry, it is a very old update")
			continue
		}
		if upd.Message.ID > 0 {
			text = upd.Message.Text
			log.Debugf("worker %d processing %d message from %s with text '%s'", worker, upd.UpdateID, upd.Message.From.UserName, text)
		} else if upd.CallBackQuery.ID != "" {
			text = upd.CallBackQuery.Data
			log.Debugf("worker %d processing %d callback from %s with text '%s'", worker, upd.UpdateID, upd.Message.From.UserName, text)
		} else {
			text = ""
		}
		if hCommand, okCommand := b.handlers[text]; okCommand {
			req := HandlerRequest{
				Update: &upd,
				Config: b.config,
				Ctx:    ctx,
				DB:     b.db,
			}
			for _, middleware := range b.middlewares {
				hCommand = middleware(req, hCommand)
			}
			hCommand(req)
		} else {
			log.Errorf("unknown command: %s", upd.Message.Text)
		}
	}
}

func tryCreateDb() (*gorm.DB, error) {
	log.Debug("create db...")
	if _, err := os.Create(dbPath); err != nil {
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

	if err = db.Exec(string(b)).Error; err != nil {
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
