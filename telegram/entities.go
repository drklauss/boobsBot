package telegram

import (
	"errors"
	"strings"
)

// Response is a telegram response
type Response struct {
	Ok     bool     `json:"ok"`
	Result []Update `json:"result"`
}

// Update contains main update info from bot
type Update struct {
	UpdateID int     `json:"update_id"`
	Message  Message `json:"message"`
}

// Message is a message message from chat
type Message struct {
	ID   int    `json:"message_id"`
	From User   `json:"from"`
	Chat Chat   `json:"chat"`
	Date int64  `json:"date"`
	Text string `json:"text"`
}

// User is info about who sent message
type User struct {
	ID        int    `json:"id"`
	FirstName string `json:"first_name"`
	UserName  string `json:"username"`
}

// Chat is where message was sent to
type Chat struct {
	ID                          int    `json:"id"`
	Title                       string `json:"title"`
	Type                        string `json:"type"`
	AllMembersAreAdministrators bool   `json:"all_members_are_administrators"`
}

// MessageSend is message entity
type MessageSend struct {
	ChatID         int
	Text           string
	KeyboardMarkup interface{}
}

// Send sends simple text message
func (ms *MessageSend) Send() error {
	return SendMessage(*ms)
}

// MediaSend is message entity with media content
type MediaSend struct {
	ChatID         int
	URL            string
	Caption        string
	KeyboardMarkup interface{}
}

// Send sends message with media content
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
