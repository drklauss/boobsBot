package telegram

import (
	"github.com/drklauss/boobsBot/config"
)

var (
	adminKeyboard      *ReplyKeyboardMarkup
	categoriesKeyboard *InlineKeyboardMarkup
)

// ReplyKeyboardMarkup is reply keyboard markup.
type ReplyKeyboardMarkup struct {
	Keyboard        [][]KeyboardButton `json:"keyboard"`
	ResizeKeyboard  bool               `json:"resize_keyboard"`
	OneTimeKeyboard bool               `json:"one_time_keyboard"`
	Selective       bool               `json:"selective"`
}

// ReplyKeyboardRemove removes ReplyKeyboard.
type ReplyKeyboardRemove struct {
	RemoveKeyboard bool `json:"remove_keyboard"`
	Selective      bool `json:"selective"`
}

// KeyboardButton is a button for ReplyKeyboardMarkup.
type KeyboardButton struct {
	Text            string `json:"text"`
	RequestContact  bool   `json:"request_contact"`
	RequestLocation bool   `json:"request_location"`
}

// InlineKeyboardMarkup is a inline keyboard.
type InlineKeyboardMarkup struct {
	InlineKeyboard [][]InlineKeyboardButton `json:"inline_keyboard"`
}

// InlineKeyboardButton is a button for inline keyboard.
type InlineKeyboardButton struct {
	Text                         string `json:"text"`
	URL                          string `json:"url"`
	CallbackData                 string `json:"callback_data"`
	SwitchInlineQuery            string `json:"switch_inline_query"`
	SwitchInlineQueryCurrentChat string `json:"switch_inline_query_current_chat"`
}

func loadKeyboards() {
	// default
	var btns []InlineKeyboardButton
	for _, c := range config.Get().Categories {
		btn := InlineKeyboardButton{
			Text:         c.Name,
			CallbackData: "/" + c.Name,
		}
		btns = append(btns, btn)
	}
	keyboards := [][]InlineKeyboardButton{btns}
	categoriesKeyboard = &InlineKeyboardMarkup{
		InlineKeyboard: keyboards,
	}
	// admin
	adminKeyboard = createAdminKeyboard()
}

// GetAdminKeyboard returns admin keyboard.
func GetAdminKeyboard() *ReplyKeyboardMarkup {
	return adminKeyboard
}

// GetCategoriesInlineKeyboard returns categories keyboard.
func GetCategoriesInlineKeyboard() *InlineKeyboardMarkup {
	return categoriesKeyboard
}

func createAdminKeyboard() *ReplyKeyboardMarkup {
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
	btnTopViewers := KeyboardButton{
		Text:            "/topViewers",
		RequestContact:  false,
		RequestLocation: false,
	}
	btnCategoriesStat := KeyboardButton{
		Text:            "/categoriesStat",
		RequestContact:  false,
		RequestLocation: false,
	}

	return &ReplyKeyboardMarkup{
		Keyboard:        [][]KeyboardButton{{btnDebugStart, btnDebugStop, btnTopViewers, btnCategoriesStat, btnUpdate}},
		OneTimeKeyboard: true,
		ResizeKeyboard:  true,
		Selective:       false,
	}
}
