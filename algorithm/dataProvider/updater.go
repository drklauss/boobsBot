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
	"bytes"
)

const getUpdates = 100

type ItemUpdater struct {
	db                *gorm.DB
	totalEntriesCount int      // общее количество записей
	totalUp           int      // количество обновленных ссылок
	errCount          int      // количество ошибок при обновлении
	insertValues      []string // вносимые значения
	catType           string   // категория
	log               *bytes.Buffer
}

// Run инициализирует БД для работы ItemUpdater-а
func (upd *ItemUpdater) Run(db *gorm.DB, catType string) []byte {
	upd.db = db
	upd.catType = catType
	upd.totalUp = 0
	upd.insertValues = []string{}
	upd.log = bytes.NewBufferString("")
	upd.updateItems(catType)

	return upd.log.Bytes()
}

// Обновляет Items
func (upd *ItemUpdater) updateItems(catType string) {
	log.Println(fmt.Sprintf("Starting update %s \n", catType))
	upd.totalEntriesCount = upd.getTotalEntriesCount()
	var lastNameId string

	for upd.totalUp < getUpdates || upd.errCount > 5 {
		switch catType {
		case config.TmNSFWCmd:
			items, last, err := reddit.GetItems(config.RdtNSFWHot, lastNameId)
			if err != nil {
				log.Println(err.Error())
				upd.errCount++
				continue
			}
			lastNameId = last
			upd.add(items...)
		case config.TmCelebCmd:
			items, last, err := reddit.GetItems(config.RdtCelebHot, lastNameId)
			if err != nil {
				log.Println(err.Error())
				upd.errCount++
				continue
			}
			lastNameId = last
			upd.add(items...)
		case config.TmRealGirlsCmd:
			items, last, err := reddit.GetItems(config.RdtRealGirlsHot, lastNameId)
			if err != nil {
				log.Println(err.Error())
				upd.errCount++
				continue
			}
			lastNameId = last
			upd.add(items...)
		}
		upd.save()
		afterInsertCount := upd.getTotalEntriesCount()
		upd.totalUp += afterInsertCount - upd.totalEntriesCount
		upd.totalEntriesCount = afterInsertCount
	}
	endUpd := fmt.Sprintf("Updated %d %s entries\n", upd.totalUp, catType)
	upd.log.WriteString(endUpd)
	log.Println(endUpd)
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
	err := upd.db.Exec(fmt.Sprintf("INSERT OR IGNORE INTO \"Urls\" (\"category\",\"value\", \"urlHash\", \"caption\") VALUES %s", insertStr)).Error
	if err != nil {
		log.Println(err.Error())
	}
}
