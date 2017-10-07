package telegram

type Response struct {
	Ok     bool     `json:"ok"`
	Result []Update `json:"result"`
}

type Update struct {
	UpdateId int     `json:"update_id"`
	Message  Message `json:"message"`
}

type Message struct {
	Id   int    `json:"message_id"`
	From User   `json:"from"`
	Chat Chat   `json:"chat"`
	Date int64  `json:"date"`
	Text string `json:"text"`
}

type User struct {
	Id        int    `json:"id"`
	FirstName string `json:"first_name"`
	UserName  string `json:"username"`
}

type Chat struct {
	Id                          int    `json:"id"`
	Title                       string `json:"title"`
	Type                        string `json:"type"`
	AllMembersAreAdministrators bool   `json:"all_members_are_administrators"`
}

// Разметка клавиатуры
type ReplyKeyboardMarkup struct {
	Keyboard        [][]KeyboardButton `json:"keyboard"`
	ResizeKeyboard  bool               `json:"resize_keyboard"`
	OneTimeKeyboard bool               `json:"one_time_keyboard"`
	Selective       bool               `json:"selective"`
}

// Удаляет кастомуную ReplyKeyboard
type ReplyKeyboardRemove struct {
	RemoveKeyboard bool `json:"remove_keyboard"`
	Selective      bool `json:"selective"`
}

type KeyboardButton struct {
	Text            string `json:"text"`
	RequestContact  bool   `json:"request_contact"`
	RequestLocation bool   `json:"request_location"`
}

// Размека инлайновой клавиатуры
type InlineKeyboardMarkup struct {
	InlineKeyboard [][]InlineKeyboardButton `json:"inline_keyboard"`
}

type InlineKeyboardButton struct {
	Text                         string `json:"text"`
	Url                          string `json:"url"`
	CallbackData                 string `json:"callback_data"`
	SwitchInlineQuery            string `json:"switch_inline_query"`
	SwitchInlineQueryCurrentChat string `json:"switch_inline_query_current_chat"`
}

type MessageSend struct {
	ChatId         int
	Text           string
	KeyboardMarkup interface{}
}

type DocumentSend struct {
	ChatId         int
	Url            string
	Caption        string
	KeyboardMarkup interface{}
}
