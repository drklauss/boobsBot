package telegram

import (
	"errors"
	"strings"
)

// Response is a telegram response.
type Response struct {
	Ok     bool     `json:"ok"`
	Result []Update `json:"result"`
}

// ResponseError is a telegram error response.
//type ResponseError struct {
//	Ok          bool   `json:"ok"`
//	ErrorCode   int    `json:"error_code"`
//	Description string `json:"description"`
//}

// Update object represents an incoming update.
// At most one of the optional parameters can be present in any given update.
type Update struct {
	UpdateID      int           `json:"update_id"`
	Message       Message       `json:"message"`
	CallBackQuery CallbackQuery `json:"callback_query"`
	//InlineQuery        InlineQuery        `json:"inline_query"`
	//ChosenInlineResult ChosenInlineResult `json:"chosen_inline_result"`
}

// InlineQuery This object represents an incoming inline query. When the user sends an empty query, your bot could return some default or trending results.
type InlineQuery struct {
	ID       string   `json:"id"`       // Unique identifier for this query
	From     User     `json:"from"`     // Sender
	Location Location `json:"location"` // Опционально. Sender location, only for bots that request user location
	Query    string   `json:"query"`    // Text of the query
	Offset   string   `json:"offset"`   // Offset of the results to be returned, can be controlled by the bot
}

// ChosenInlineResult Represents a result of an inline query that was chosen by the user and sent to their chat partner.
type ChosenInlineResult struct {
	ResultID        string   `json:"result_id"`         // The unique identifier for the result that was chosen
	From            User     `json:"from"`              // The user that chose the result
	Location        Location `json:"location"`          // Опционально. Sender location, only for bots that require user location
	InlineMessageID string   `json:"inline_message_id"` // Опционально. Identifier of the sent inline message. Available only if there is an inline keyboard attached to the message. Will be also received in callback queries and can be used to edit the message.
	Query           string   `json:"query"`             // The query that was used to obtain the result
}

// Location Этот объект представляет точку на карте.
type Location struct {
	Longitude float64 `json:"longitude"` // Долгота, заданная отправителем
	Latitude  float64 `json:"latitude"`  // Широта, заданная отправителем
}

// CallbackQuery object represents an incoming callback query from a callback button in an inline keyboard.
// If the button that originated the query was attached to a message sent by the bot, the field message will be present.
// If the button was attached to a message sent via the bot (in inline mode), the field inline_message_id will be present.
// Exactly one of the fields data or game_short_name will be present.
type CallbackQuery struct {
	ID      string  `json:"id"`
	From    User    `json:"from"`
	Message Message `json:"message"`
	Data    string  `json:"data"`
}

// AnswerCallbackQuery is method to send answers to callback queries sent from inline keyboards.
// The answer will be displayed to the user as a notification at the top of the chat screen or as an alert.
// On success, True is returned.
type AnswerCallbackQuery struct {
	CallbackQueryID string `json:"callback_query_id"`
	Text            string `json:"text"`
	ShowAlert       bool   `json:"show_alert"`
	URL             string `json:"url"`
}

// Message represents message.
type Message struct {
	ID          int         `json:"message_id"`
	From        User        `json:"from"`
	Chat        Chat        `json:"chat"`
	Date        int64       `json:"date"`
	Text        string      `json:"text"`
	ReplyMarkup interface{} `json:"reply_markup"`
}

// User represents a Telegram user or bot.
type User struct {
	ID        int    `json:"id"`
	FirstName string `json:"first_name"`
	IsBot     bool   `json:"is_bot"`
	UserName  string `json:"username"`
}

// Chat represents chat.
type Chat struct {
	ID                          int    `json:"id"`
	Title                       string `json:"title"`
	Type                        string `json:"type"`
	AllMembersAreAdministrators bool   `json:"all_members_are_administrators"`
}

// MessageSend is message entity.
type MessageSend struct {
	ChatID         int
	Text           string
	KeyboardMarkup interface{}
}

// Send sends simple text message.
func (ms *MessageSend) Send() error {
	return SendMessage(*ms)
}

// MediaSend is message entity with media content.
type MediaSend struct {
	ChatID         int
	URL            string
	Caption        string
	KeyboardMarkup interface{}
}

// Send sends message with media content.
func (c *MediaSend) Send() error {
	uSplit := strings.Split(c.URL, ".")
	if len(uSplit) <= 0 {
		return errors.New("empty URL")
	}
	ext := uSplit[len(uSplit)-1]
	for _, photoExt := range []string{"png", "jpg", "jpeg"} {
		if strings.Contains(ext, photoExt) {
			return SendImage(*c)
		}
	}
	for _, videoExt := range []string{"mp4", "gif"} {
		if strings.Contains(ext, videoExt) {
			return SendDocument(*c)
		}
	}
	return nil
}
