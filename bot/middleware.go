package bot

import (
	"context"
	"time"

	"github.com/drklauss/boobsBot/config"

	"github.com/drklauss/boobsBot/telegram"
	log "github.com/sirupsen/logrus"
)

// Middleware func allows to create some middlewares to use them in handlers.
type Middleware func(ctx context.Context, next HandlerFunc, u *telegram.Update) HandlerFunc

// CheckAdmin interrupts special admin requests and apply special things.
func CheckAdmin(ctx context.Context, next HandlerFunc, u *telegram.Update) HandlerFunc {
	defaultHandler := func(ctx context.Context, u *telegram.Update) {
		next(ctx, u)
	}
	isAdmin := false
	for _, adminID := range config.Get().Telegram.Admin {
		if u.Message.From.ID == adminID {
			isAdmin = true
			break
		}
	}
	if !isAdmin {
		return defaultHandler
	}

	var ms telegram.MessageSend
	switch u.Message.Text {
	case "/debugStart":
		SetDebug(true)
		ms = telegram.MessageSend{
			ChatID: u.Message.Chat.ID,
			Text:   "debug enabled",
		}
	case "/debugStop":
		SetDebug(false)
		ms = telegram.MessageSend{
			ChatID: u.Message.Chat.ID,
			Text:   "debug disabled",
		}
	case "/admin":
		ms = telegram.MessageSend{
			ChatID:         u.Message.Chat.ID,
			Text:           "admin keyboard open",
			KeyboardMarkup: telegram.GetAdminKeyboard(),
		}
	default:
		if isAdmin && Debug() {
			log.Infoln("all non-admin commands will be ignored during debug")
			return func(ctx context.Context, u *telegram.Update) {}
		}
		return defaultHandler
	}
	err := ms.Send()
	if err != nil {
		log.Warnln(err)
	}
	return func(ctx context.Context, u *telegram.Update) {}
}

// LogRequest is a dev middleware, that just logs the request data.
func LogRequest(ctx context.Context, next HandlerFunc, u *telegram.Update) HandlerFunc {
	return func(ctx context.Context, u *telegram.Update) {
		t := time.Now()
		log.Infof("process %d message: %s", u.UpdateID, u.Message.Text)
		next(ctx, u)
		log.Infof("update %d handled in: %v", u.UpdateID, time.Since(t))
	}
}
