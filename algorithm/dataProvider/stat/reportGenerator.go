package stat

import "github.com/jinzhu/gorm"

const (
	TelegramFmt = 1
	ConsoleFmt  = 2
)

type ReportGenerator struct {
	db *gorm.DB
}

// SetDb устанавливает соединение для содания отчетов
func (p *ReportGenerator) SetDb(db *gorm.DB) *ReportGenerator {
	p.db = db

	return p
}

