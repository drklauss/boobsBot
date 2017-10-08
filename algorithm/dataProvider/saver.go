package dataProvider

import (
	"crypto/md5"
	"fmt"
	"strings"

	"github.com/drklauss/boobsBot/algorithm/dataProvider/dbEntities"
	"github.com/jinzhu/gorm"
)

type ItemSaver struct {
	db           *gorm.DB
	insertValues []string
	catType      string
}

// Init инициализирует БД для работы ItemSaver-а
func (s *ItemSaver) Init(db *gorm.DB, catType string) *ItemSaver {
	s.db = db
	s.catType = catType

	return s
}

// Add добавляет новый item для последующего сохранения
func (s *ItemSaver) Add(item dbEntities.UrlGetter) *ItemSaver {
	if item.GetValue() == "" {
		return s
	}
	h := md5.New()
	h.Write([]byte(item.GetValue()))
	md := h.Sum(nil)
	newV := fmt.Sprintf("(\"%s\", \"%s\", \"%x\", \"%s\")", s.catType, item.GetValue(), md, item.GetCaption())
	s.insertValues = append(s.insertValues, newV)

	return s
}

// Сохраняет данные в БД
func (s *ItemSaver) Save() {
	insertStr := strings.Join(s.insertValues, ",")
	sql := fmt.Sprintf("INSERT OR IGNORE INTO \"Urls\" (\"category\",\"value\", \"urlHash\", \"caption\") VALUES %s", insertStr)
	s.db.Exec(sql)
	s.insertValues = []string{}
}
