package dataProvider

import (
	"log"

	"crypto/md5"
	"fmt"
	"strings"

	"bytes"

	"time"

	"github.com/drklauss/boobsBot/algorithm/config"
	"github.com/drklauss/boobsBot/algorithm/dataProvider/dbEntities"
	"github.com/drklauss/boobsBot/algorithm/dataProvider/reddit"
	"github.com/jinzhu/gorm"
)

type ItemUpdater struct {
	db           *gorm.DB
	totalUp      int      // количество обновленных ссылок
	errCount     int      // количество ошибок при обновлении
	insertValues []string // вносимые значения
	catType      string   // категория
	log          *bytes.Buffer
}

// Run инициализирует БД для работы ItemUpdater-а
func (upd *ItemUpdater) Run(db *gorm.DB, catType string) []byte {
	upd.db = db
	upd.catType = catType
	upd.log = bytes.NewBufferString("")
	upd.updateItems()

	return upd.log.Bytes()
}

// Обновляет Items
func (upd *ItemUpdater) updateItems() {
	log.Println(fmt.Sprintf("Starting update %s \n", upd.catType))
	upStart := time.Now()
	totalEntriesCount := upd.getTotalEntriesCount()
	var lastNameId string
	for upd.totalUp < config.GetMaxUpdates && upd.errCount < 5 {
		execTime := time.Since(upStart).Seconds()
		if execTime > config.UpdateTimeOut {
			upd.log.WriteString(fmt.Sprintf("Update \"%s\" timeout\n", upd.catType))
			break
		}

		switch upd.catType {
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
		upd.insertValues = []string{}
		afterInsertCount := upd.getTotalEntriesCount()
		upd.totalUp += afterInsertCount - totalEntriesCount
		totalEntriesCount = afterInsertCount
	}
	endUpd := fmt.Sprintf("Updated %d %s entries\n", upd.totalUp, upd.catType)
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
			continue
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
