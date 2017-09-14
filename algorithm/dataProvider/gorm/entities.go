package gorm

type Urls struct {
	Id      int    `gorm:"primary_key;column:id"`
	Value   string `gorm:"column:value"`
	UrlHash string `gorm:"column:urlHash"`
}

type Views struct {
	UrlId  int `gorm:"column:urlId"`
	ChatId int `gorm:"column:chatId"`
}

func (Urls) TableName() string {
	return "Urls"
}

func (Views) TableName() string {
	return "Views"
}
