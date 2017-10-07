package dataProvider

import (
	"log"

	"fmt"

	"crypto/md5"

	"strings"

	"time"

	"github.com/drklauss/boobsBot/algorithm/config"
	"github.com/drklauss/boobsBot/algorithm/dataProvider/dbEntities"
	"github.com/drklauss/boobsBot/algorithm/dataProvider/gfycat"
	"github.com/drklauss/boobsBot/algorithm/dataProvider/reddit"
	"github.com/drklauss/boobsBot/algorithm/dataProvider/stat"
	"github.com/drklauss/boobsBot/algorithm/telegram"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/sqlite"
)

const getUpdates = 100

// Класс является черным ящиком для получения данных из БД
type Provider struct {
	db             *gorm.DB
	totalUrlsCount int
	cacheIds       []int
}

// Init инициализирует провайдер
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

// Возвращает один URL гифки, пишем в стату просмотр
func (p *Provider) GetUrl(chat telegram.Chat) dbEntities.Url {
	u, err := p.getOneUrl(chat.Id)
	if err != nil {
		log.Println(err)
	}
	go p.createViewEntry(chat.Id, u.Id)

	return u
}

// Обновляет URLs
func (p *Provider) UpdateVideoUrls() {
	log.Printf("Starting update for %d urls", 100)
	if p.totalUrlsCount == 0 {
		p.totalUrlsCount = p.getTotalEntriesCount()
	}
	var (
		totalUp  int // количество обновленных ссылок
		errCount int // количество ошибок при обновлении
	)
	for totalUp < getUpdates || errCount >= 10 {
		redditUrls, err := reddit.GetNames(config.RdtHotCategory)
		if err != nil {
			log.Println(err.Error())
			errCount++
		}
		if len(redditUrls) == 0 {
			log.Println("There is no urls")
			break
		}
		convertedUrls, err := gfycat.ConvertNamesToUrls(redditUrls)
		if err != nil {
			log.Println(err.Error())
			errCount++
		}
		p.saveUrls(convertedUrls)
		afterInsertCount := p.getTotalEntriesCount()
		totalUp += afterInsertCount - p.totalUrlsCount
		p.totalUrlsCount = afterInsertCount
	}
	log.Printf("Updated %d urls", totalUp)
}

// GetTopViewers4Tm запрашивает TopViewers отчет отформатированный для Telegram
func (p *Provider) GetTopViewers4Tm() []byte {
	return new(stat.ReportGenerator).SetDb(p.db).GetTopViewers(stat.TelegramFmt)
}

// GetTopViewers4Cl запрашивает TopViewers отчет отформатированный для Console
func (p *Provider) GetTopViewers4Cl() []byte {
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
func (p *Provider) clearChatViews(chatId int) error {
	err := p.db.Where("chatId = ?", chatId).
		Delete(dbEntities.View{}).
		Error

	return err
}

// Возвращает одну ссылку
func (p *Provider) getOneUrl(chatId int) (dbEntities.Url, error) {
	var url dbEntities.Url
	sql := fmt.Sprintf(`
	SELECT *
	FROM Urls
	WHERE id NOT IN
		(SELECT u.id
		FROM Urls as u
		LEFT JOIN Views as v
			ON v.urlId = u.id
		WHERE v.chatId = %d)
	ORDER BY RANDOM()
	LIMIT 1`, chatId)
	err := p.db.Raw(sql).Scan(&url).Error
	if err == gorm.ErrRecordNotFound {
		p.clearChatViews(chatId)
		url, err = p.getOneUrl(chatId)
	}

	return url, err
}

// Сохраняет данные в БД
func (p *Provider) saveUrls(urls []gfycat.GfyItem) {
	insertStr := p.prepareInsertString(urls)
	sql := fmt.Sprintf("INSERT OR IGNORE INTO \"Urls\" (\"value\", \"urlHash\", \"caption\") VALUES %s", insertStr)
	p.db.Exec(sql)
}

// Подготавливяет строку для INSERT
func (p *Provider) prepareInsertString(gfyItems []gfycat.GfyItem) string {
	var values []string
	for _, item := range gfyItems {
		if item.MobileUrl == "" {
			continue
		}
		hasher := md5.New()
		hasher.Write([]byte(item.MobileUrl))
		md := hasher.Sum(nil)
		values = append(values, fmt.Sprintf("(\"%s\", \"%x\", \"%s\")", item.MobileUrl, md, item.GfyName))
	}

	return strings.Join(values, ",")
}

// Возарвщает общее количество записей в таблице Urls
func (p *Provider) getTotalEntriesCount() int {
	var count int
	p.db.Table("Urls").Count(&count)

	return count
}
