package dataProvider

import (
	"log"

	"crypto/md5"
	"fmt"
	"strings"

	"github.com/drklauss/boobsBot/algorithm/config"
	"github.com/drklauss/boobsBot/algorithm/dataProvider/dbEntities"
	"github.com/drklauss/boobsBot/algorithm/dataProvider/reddit"
	"github.com/jinzhu/gorm"
)

const getUpdates = 100

type ItemUpdater struct {
	db                *gorm.DB
	totalEntriesCount int      // общее количество записей
	totalUp           int      // количество обновленных ссылок
	errCount          int      // количество ошибок при обновлении
	insertValues      []string // вносимые значения
	catType           string   // категория
}

// Run инициализирует БД для работы ItemUpdater-а
func (upd *ItemUpdater) Run(db *gorm.DB, catType string) *ItemUpdater {

	upd.db = db
	upd.catType = catType
	upd.totalUp = 0
	upd.insertValues = []string{}
	upd.updateItems(catType)

	return upd
}

// Обновляет Items
func (upd *ItemUpdater) updateItems(catType string) {
	log.Printf("Starting update %s for %d entries", catType, getUpdates)
	upd.totalEntriesCount = upd.getTotalEntriesCount()
	for upd.totalUp < getUpdates || upd.errCount > 10 {
		fmt.Printf("%d %d %d %s\n", upd.totalUp, upd.totalEntriesCount, upd.errCount, catType)
		switch catType {
		case config.TmNSFWCmd:
			items, err := reddit.GetItems(config.RdtNSFWHot)
			if err != nil {
				log.Println(err.Error())
				upd.errCount++
				continue
			}
			upd.add(items...)
		case config.TmCeleb:
			items, err := reddit.GetItems(config.RdtCelebHot)
			if err != nil {
				log.Println(err.Error())
				upd.errCount++
				continue
			}
			upd.add(items...)
		case config.TmRealGirlsCmd:
			fmt.Println("HERERERERERERERERERE")
			items, err := reddit.GetItems(config.RdtRealGirlsHot)
			fmt.Printf("%+v", items)
			if err != nil {
				log.Println(err.Error())
				upd.errCount++
				continue
			}
			fmt.Printf("%+v", items)
			upd.add(items...)
		}
		upd.save()
		afterInsertCount := upd.getTotalEntriesCount()
		upd.totalUp += afterInsertCount - upd.totalEntriesCount
		upd.totalEntriesCount = afterInsertCount
	}
	log.Printf("Updated %d %s entries", upd.totalUp, catType)
}

// Возвращает общее количество записей в таблице Urls
func (upd *ItemUpdater) getTotalEntriesCount() int {
	var count int
	upd.db.Table("Urls").Count(&count)

	return count
}

// Добавляет новые items для последующего сохранения
func (upd *ItemUpdater) add(items ...dbEntities.Url) *ItemUpdater {
	for _, item := range items {
		if item.Value == "" {
			return upd
		}
		h := md5.New()
		h.Write([]byte(item.Value))
		md := h.Sum(nil)
		newV := fmt.Sprintf("(\"%s\", \"%s\", \"%x\", \"%s\")", upd.catType, item.Value, md, strings.Replace(item.Caption, `"`, `'`, -1))
		upd.insertValues = append(upd.insertValues, newV)
	}

	return upd
}

// Сохраняет данные в БД
func (upd *ItemUpdater) save() {
	insertStr := strings.Join(upd.insertValues, ",")
	upd.db.Exec(fmt.Sprintf("INSERT OR IGNORE INTO \"Urls\" (\"category\",\"value\", \"urlHash\", \"caption\") VALUES %s", insertStr))
}
