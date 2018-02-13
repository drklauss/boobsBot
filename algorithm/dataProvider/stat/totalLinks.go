package stat

import (
	"log"

	"strconv"

	"bytes"

	"fmt"

	"github.com/olekukonko/tablewriter"
)

// GetTopViewers возвращает отформатированнный topViewers отчет
func (p *ReportGenerator) GetCategories(format int) string {
	var categoriesStat []struct {
		Category string
		Count    int
	}
	sql := `SELECT category, 
					count(id) as count
			FROM Urls
			GROUP BY category`
	err := p.db.Raw(sql).Scan(&categoriesStat).Error
	if err != nil {
		log.Println(err)
		return ""
	}

	buf := new(bytes.Buffer)
	switch format {
	case ConsoleFmt:
		var data [][]string
		for _, v := range categoriesStat {
			data = append(data, []string{v.Category, strconv.Itoa(v.Count)})
		}
		table := tablewriter.NewWriter(buf)
		table.SetHeader([]string{"Title", "Count"})
		table.SetBorder(false)
		table.AppendBulk(data)
		table.Render()
	case TelegramFmt:
		buf.WriteString("\xF0\x9F\x93\x8A Categories Items Count:\n")
		for _, v := range categoriesStat {
			s := fmt.Sprintf("%s -- %d \n", v.Category, v.Count)
			buf.WriteString(s)
		}
	}

	return buf.String()
}
