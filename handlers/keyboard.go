package handlers

import (
	"github.com/drklauss/boobsBot/config"
	"github.com/drklauss/boobsBot/telegram"
)

var (
	adminKeyboard            *telegram.ReplyKeyboardMarkup
	categoriesInlineKeyboard *telegram.InlineKeyboardMarkup
	categoriesKeyboard       *telegram.ReplyKeyboardMarkup
)

func InitKeyboards() {
	// categories
	categoriesInlineKeyboard = createCategoriesInlineKeyboard()
	categoriesKeyboard = createCategoriesKeyboard()
	// admin
	adminKeyboard = createAdminKeyboard()
}

// GetAdminKeyboard returns admin keyboard.
func GetAdminKeyboard() *telegram.ReplyKeyboardMarkup {
	return adminKeyboard
}

// GetCategoriesInlineKeyboard returns categories keyboard.
func GetCategoriesInlineKeyboard() *telegram.InlineKeyboardMarkup {
	return categoriesInlineKeyboard
}

// GetCategoriesKeyboard returns categories keyboard.
func GetCategoriesKeyboard() *telegram.ReplyKeyboardMarkup {
	return categoriesKeyboard
}

func createCategoriesInlineKeyboard() *telegram.InlineKeyboardMarkup {
	var keyboard [][]telegram.InlineKeyboardButton
	btnsInRow := 3
	for i := 0; i < len(config.Get().Categories); i += btnsInRow {
		var btnsRow []telegram.InlineKeyboardButton
		for k := 0; k < btnsInRow; k++ {
			c := config.Get().Categories[i+k]
			btn := telegram.InlineKeyboardButton{
				Text:         c.Name,
				CallbackData: "/" + c.Name,
			}
			btnsRow = append(btnsRow, btn)
		}
		keyboard = append(keyboard, btnsRow)
	}

	return &telegram.InlineKeyboardMarkup{
		InlineKeyboard: keyboard,
	}
}

func createCategoriesKeyboard() *telegram.ReplyKeyboardMarkup {
	var keyboard [][]telegram.KeyboardButton
	btnsInRow := 3
	for i := 0; i < len(config.Get().Categories); i += btnsInRow {
		var btnsRow []telegram.KeyboardButton
		for k := 0; k < btnsInRow; k++ {
			c := config.Get().Categories[i+k]
			btn := telegram.KeyboardButton{
				Text: "/" + c.Name,
			}
			btnsRow = append(btnsRow, btn)
		}
		keyboard = append(keyboard, btnsRow)
	}

	return &telegram.ReplyKeyboardMarkup{
		Keyboard:        keyboard,
		OneTimeKeyboard: true,
		ResizeKeyboard:  true,
		Selective:       false,
	}
}

func createAdminKeyboard() *telegram.ReplyKeyboardMarkup {
	btnDebugStart := telegram.KeyboardButton{
		Text:            DebugStart,
		RequestContact:  false,
		RequestLocation: false,
	}
	btnDebugStop := telegram.KeyboardButton{
		Text:            DebugStop,
		RequestContact:  false,
		RequestLocation: false,
	}
	btnUpdate := telegram.KeyboardButton{
		Text:            Update,
		RequestContact:  false,
		RequestLocation: false,
	}
	btnTopViewers := telegram.KeyboardButton{
		Text:            TopViewers,
		RequestContact:  false,
		RequestLocation: false,
	}
	btnCategoriesStat := telegram.KeyboardButton{
		Text:            CategoriesStat,
		RequestContact:  false,
		RequestLocation: false,
	}

	return &telegram.ReplyKeyboardMarkup{
		Keyboard:        [][]telegram.KeyboardButton{{btnDebugStart, btnDebugStop, btnTopViewers, btnCategoriesStat, btnUpdate}},
		OneTimeKeyboard: true,
		ResizeKeyboard:  true,
		Selective:       false,
	}
}
