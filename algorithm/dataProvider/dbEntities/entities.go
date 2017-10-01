package dbEntities

type Url struct {
	Id      int    `gorm:"primary_key;column:id"`
	Value   string `gorm:"column:value"`
	UrlHash string `gorm:"column:urlHash"`
}

type View struct {
	UrlId       int   `gorm:"column:urlId"`
	ChatId      int   `gorm:"column:chatId"`
	RequestDate int64 `gorm:"column:requestDate"`
}

type Chat struct {
	Id    int    `gorm:"primary_key;column:id"`
	Title string `gorm:"column:title"`
	Type  string `gorm:"column:type"`
}

func (Url) TableName() string {
	return "Urls"
}

func (View) TableName() string {
	return "Views"
}

func (Chat) TableName() string {
	return "Chats"
}
