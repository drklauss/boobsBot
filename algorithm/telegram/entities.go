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
