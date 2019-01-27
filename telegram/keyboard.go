package telegram

import (
	"github.com/drklauss/boobsBot/config"
)

var adminKeydoard, defaultKeyboard *ReplyKeyboardMarkup

// ReplyKeyboardMarkup is reply keyboard markup
type ReplyKeyboardMarkup struct {
	Keyboard        [][]KeyboardButton `json:"keyboard"`
	ResizeKeyboard  bool               `json:"resize_keyboard"`
	OneTimeKeyboard bool               `json:"one_time_keyboard"`
	Selective       bool               `json:"selective"`
}

// ReplyKeyboardRemove removes ReplyKeyboard
type ReplyKeyboardRemove struct {
	RemoveKeyboard bool `json:"remove_keyboard"`
	Selective      bool `json:"selective"`
}

// KeyboardButton is a button for ReplyKeyboardMarkup
type KeyboardButton struct {
	Text            string `json:"text"`
	RequestContact  bool   `json:"request_contact"`
	RequestLocation bool   `json:"request_location"`
}

// InlineKeyboardMarkup is a inline keyboard
type InlineKeyboardMarkup struct {
	InlineKeyboard [][]InlineKeyboardButton `json:"inline_keyboard"`
}

// InlineKeyboardButton is a button for inline keyboard
type InlineKeyboardButton struct {
	Text                         string `json:"text"`
	URL                          string `json:"url"`
	CallbackData                 string `json:"callback_data"`
	SwitchInlineQuery            string `json:"switch_inline_query"`
	SwitchInlineQueryCurrentChat string `json:"switch_inline_query_current_chat"`
}

func loadKeyboards() {
	// default
	var btns []KeyboardButton
	for _, c := range config.Get().Categories {
		btn := KeyboardButton{
			Text:            c.Name,
			RequestContact:  false,
			RequestLocation: false,
		}
		btns = append(btns, btn)
	}
	keyboards := [][]KeyboardButton{btns}
	defaultKeyboard = &ReplyKeyboardMarkup{
		Keyboard:        keyboards,
		OneTimeKeyboard: true,
		ResizeKeyboard:  true,
		Selective:       false,
	}
	// admin
	btnDebugStart := KeyboardButton{
		Text:            "/debugStart",
		RequestContact:  false,
		RequestLocation: false,
	}
	btnDebugStop := KeyboardButton{
		Text:            "/debugStop",
		RequestContact:  false,
		RequestLocation: false,
	}
	btnUpdate := KeyboardButton{
		Text:            "/update",
		RequestContact:  false,
		RequestLocation: false,
	}
	adminKeydoard = &ReplyKeyboardMarkup{
		Keyboard:        [][]KeyboardButton{[]KeyboardButton{btnDebugStart, btnDebugStop, btnUpdate}},
		OneTimeKeyboard: true,
		ResizeKeyboard:  true,
		Selective:       false,
	}
}

// GetDefaultKeayboard returns default user keyboard
func GetDefaultKeayboard() *ReplyKeyboardMarkup {
	return defaultKeyboard
}

// GetAdminKeayboard returns admin keyboard
func GetAdminKeayboard() *ReplyKeyboardMarkup {
	return adminKeydoard
}
