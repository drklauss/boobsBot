package dataProvider

import (
	"log"

	"fmt"

	"crypto/md5"

	"strings"

	"github.com/boobsBot/algorithm/config"
	"github.com/boobsBot/algorithm/dataProvider/gfycat"
	gorm2 "github.com/boobsBot/algorithm/dataProvider/gorm"
	"github.com/boobsBot/algorithm/dataProvider/reddit"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/sqlite"
)

// Класс является черным ящиком для получения нужного URL-а по категории
type Provider struct {
	db             *gorm.DB
	totalUrlsCount int
}

// Init инициализирует провайдер
func (p *Provider) Init() Provider {

	db, err := gorm.Open("sqlite3", config.DbFileName)
	if err != nil {
		log.Fatal("Cannot Open db connection")
	}
	db.DB().SetMaxOpenConns(10)
	db.DB().SetMaxIdleConns(5)
	db.LogMode(true)
	p.db = db

	return *p
}

// Возвращает один URL гифки, уменьшаем срез на один урл
// Если в срезе URL последний - обновляем срез
func (p *Provider) GetUrl(chatId int) string {
	url, err := p.getOneUrl(chatId)
	if err != nil {
		log.Println(err)
	}
	p.createViewEntry(chatId, url.Id)

	return url.Value
}

// Создает запись в просмотрах для чата
func (p *Provider) createViewEntry(chatId int, urlId int) error {
	view := gorm2.Views{ChatId: chatId, UrlId: urlId}
	err := p.db.Create(view).Error

	return err
}

// Очищает просмотры для чата
func (p *Provider) clearChatViews(chatId int) error {
	err := p.db.Where("chatId = ?", chatId).
		Delete(gorm2.Views{}).
		Error

	return err
}

// Возвращает одну ссылку
func (p *Provider) getOneUrl(chatId int) (gorm2.Urls, error) {
	var url gorm2.Urls
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

// Обновляет URLs
func (p *Provider) Update(up int) {
	log.Printf("Starting update for %d urls", up)
	if p.totalUrlsCount == 0 {
		p.totalUrlsCount = p.getTotalEntriesCount()
	}
	var (
		totalUp  int // количество обновленных ссылок
		errCount int // количество ошибок при обновлении
	)
	for totalUp < up || errCount == 10 {
		redditUrls, err := reddit.GetNames(config.HotCategory)
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

// Сохраняет картинки в БД
func (p *Provider) saveUrls(urls []string) {
	insertStr := p.prepareInsertString(urls)
	sql := fmt.Sprintf("INSERT OR IGNORE INTO \"Urls\" (\"value\", \"urlHash\") VALUES %s", insertStr)
	p.db.Exec(sql)
}

// Подготавливяет строку для INSERT
func (p *Provider) prepareInsertString(urls []string) string {
	var values []string
	for _, v := range urls {
		if v == "" {
			continue
		}
		hasher := md5.New()
		hasher.Write([]byte(v))
		md := hasher.Sum(nil)
		values = append(values, fmt.Sprintf("(\"%s\", \"%s\")", v, fmt.Sprintf("%x", md)))
	}

	return strings.Join(values, ",")
}

// Возарвщает общее количество записей в таблице Urls
func (p *Provider) getTotalEntriesCount() int {
	var count int
	p.db.Table("Urls").Count(&count)

	return count
}
