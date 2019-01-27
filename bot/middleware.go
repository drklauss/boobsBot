package bot

import (
	"context"
	"time"

	"github.com/drklauss/boobsBot/config"

	"github.com/drklauss/boobsBot/telegram"
	"github.com/leesper/holmes"
)

// Middleware func allows to create some middlewares to use them in handlers
type Middleware func(ctx context.Context, next HandlFunc, u *telegram.Update) HandlFunc

// CheckAdmin interrupts special admin requests and apply special things
func CheckAdmin(ctx context.Context, next HandlFunc, u *telegram.Update) HandlFunc {
	defaultHanlder := func(ctx context.Context, u *telegram.Update) {
		next(ctx, u)
	}
	isAdmin := false
	for _, adminID := range config.Get().Telegram.Admin {
		if u.Message.From.ID == adminID {
			isAdmin = true
		}
	}
	if !isAdmin {
		return defaultHanlder
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
			Text:           "admin kb",
			KeyboardMarkup: telegram.GetAdminKeayboard(),
		}
	default:
		return defaultHanlder
	}
	err := ms.Send()
	if err != nil {
		holmes.Warnln(err)
	}
	return func(ctx context.Context, u *telegram.Update) {}
}

// LogRequest is a dev middleware, that just logs the request data
func LogRequest(ctx context.Context, next HandlFunc, u *telegram.Update) HandlFunc {
	return func(ctx context.Context, u *telegram.Update) {
		t := time.Now()
		next(ctx, u)
		holmes.Infof("update %d handled in: %v", u.UpdateID, time.Since(t))
	}
}
