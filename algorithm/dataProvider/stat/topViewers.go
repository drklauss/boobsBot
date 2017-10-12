package stat

import (
	"log"

	"strconv"

	"bytes"

	"fmt"

	"github.com/olekukonko/tablewriter"
)

type TopViewer struct {
	Title string
	Type  string
	Count int
}

// GetTopViewers возвращает отформатированнный topViewers отчет
func (p *ReportGenerator) GetTopViewers(format int) string {
	var tVs []TopViewer
	sql := `SELECT c.title,
					c.type,
					count(c.id) AS count
			FROM Views v
			LEFT JOIN Chats c on
				c.id = v.chatId
			GROUP BY c.id
			ORDER BY count(c.id) DESC
			LIMIT 10`
	err := p.db.Raw(sql).Scan(&tVs).Error
	if err != nil {
		log.Println(err)
		return ""
	}

	buf := new(bytes.Buffer)
	switch format {
	case ConsoleFmt:
		var data [][]string
		for k, v := range tVs {
			data = append(data, []string{strconv.Itoa(k + 1), v.Title, v.Type, strconv.Itoa(v.Count)})
		}
		table := tablewriter.NewWriter(buf)
		table.SetHeader([]string{"№", "Title", "Type", "Count"})
		table.SetBorder(false)
		table.AppendBulk(data)
		table.Render()
	case TelegramFmt:
		buf.WriteString("\xF0\x9F\x93\x8A Top Viewers Report:\n")
		for k, v := range tVs {
			s := fmt.Sprintf("%d. %s (%s) -- %d \n", k+1, v.Title, v.Type, v.Count)
			buf.WriteString(s)
		}
	}

	return buf.String()
}
