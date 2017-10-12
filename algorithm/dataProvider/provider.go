package dataProvider

import (
	"log"

	"fmt"

	"time"

	"github.com/drklauss/boobsBot/algorithm/config"
	"github.com/drklauss/boobsBot/algorithm/dataProvider/dbEntities"
	"github.com/drklauss/boobsBot/algorithm/dataProvider/stat"
	"github.com/drklauss/boobsBot/algorithm/telegram"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/sqlite"
)

// Класс является черным ящиком для получения данных из БД
type Provider struct {
	db       *gorm.DB
	cacheIds []int
}

// Run инициализирует провайдер
func (p *Provider) Init(logSql bool) Provider {

	db, err := gorm.Open("sqlite3", config.DbFileName)
	if err != nil {
		log.Fatal("Cannot Open db connection")
	}
	db.DB().SetMaxOpenConns(10)
	db.DB().SetMaxIdleConns(5)
	db.LogMode(logSql)
	p.db = db

	return *p
}

// CacheChatIds кэширует все ChatIds
func (p *Provider) CacheChatIds() {
	rows, err := p.db.Select("id").
		Table("Chats").
		Rows()
	if err != nil {
		log.Fatal(err)
	}
	for rows.Next() {
		var id int
		rows.Scan(&id)
		p.cacheIds = append(p.cacheIds, id)
	}
}

// Возвращает один URL по категории, пишем в стату просмотр
func (p *Provider) GetUrl(chat telegram.Chat, cat string) dbEntities.Url {
	u, err := p.getOneUrl(chat.Id, cat)
	if err != nil {
		log.Println(err)
	}
	go p.createViewEntry(chat.Id, u.Id)

	return u
}

// Обновляет данные
func (p *Provider) UpdateAll() string {
	updater := new(ItemUpdater)
	nsfwLog := updater.Run(p.db, config.TmNSFWCmd)
	realGLog := updater.Run(p.db, config.TmRealGirlsCmd)
	celebLog := updater.Run(p.db, config.TmCelebCmd)

	return fmt.Sprintf("%s%s%s", nsfwLog, realGLog, celebLog)
}

// GetTopViewers4Tm запрашивает TopViewers отчет отформатированный для Telegram
func (p *Provider) GetTopViewers4Tm() string {
	return new(stat.ReportGenerator).SetDb(p.db).GetTopViewers(stat.TelegramFmt)
}

// GetTopViewers4Cl запрашивает TopViewers отчет отформатированный для Console
func (p *Provider) GetTopViewers4Cl() string {
	return new(stat.ReportGenerator).SetDb(p.db).GetTopViewers(stat.ConsoleFmt)
}

// Создает запись в просмотрах для чата
func (p *Provider) createViewEntry(chatId, urlId int) {
	v := dbEntities.View{
		ChatId:      chatId,
		UrlId:       urlId,
		RequestDate: time.Now().Unix(),
	}
	p.db.Create(v)
}

// Создает запись о чате
func (p *Provider) CreateChatEntry(mes telegram.Message) {
	var contain bool
	for _, v := range p.cacheIds {
		if v == mes.Chat.Id {
			contain = true
			break
		}
	}
	if !contain {
		c := dbEntities.Chat{
			Id:    mes.Chat.Id,
			Type:  mes.Chat.Type,
			Title: mes.Chat.Title,
		}
		if len(c.Title) == 0 {
			c.Title = fmt.Sprintf("%s (%s)", mes.From.FirstName, mes.From.UserName)
		}
		p.db.Create(c)
	}
}

// Очищает просмотры для чата
func (p *Provider) clearChatViews(chatId int, cat string) error {
	sql := fmt.Sprintf(`
	DELETE FROM Views
	WHERE Views.urlId IN (
		SELECT u.id FROM Urls AS u
		INNER JOIN Views AS v ON u.id = v.urlId
		WHERE u.category = "%s" AND v.chatId = %d
	)`, cat, chatId)
	err := p.db.Exec(sql).Error

	return err
}

// Возвращает одну ссылку
func (p *Provider) getOneUrl(chatId int, cat string) (dbEntities.Url, error) {
	var u dbEntities.Url
	sql := fmt.Sprintf(`
	SELECT *
	FROM Urls
	WHERE category="%s" AND id NOT IN
		(SELECT u.id
		FROM Urls as u
		LEFT JOIN Views as v
			ON v.urlId = u.id
		WHERE v.chatId = %d)
	ORDER BY RANDOM()
	LIMIT 1`, cat, chatId)
	err := p.db.Raw(sql).Scan(&u).Error
	if err == gorm.ErrRecordNotFound {
		p.clearChatViews(chatId, cat)
		u, err = p.getOneUrl(chatId, cat)
	}

	return u, err
}
