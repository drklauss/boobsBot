package dataProvider

import (
	"log"

	"github.com/drklauss/boobsBot/algorithm/config"
	"github.com/drklauss/boobsBot/algorithm/dataProvider/gfycat"
	"github.com/drklauss/boobsBot/algorithm/dataProvider/reddit"
)

// Обновляет Video Items
func (p *Provider) updateVideoItems() {
	log.Printf("Starting update for %d urls", 100)
	if p.totalUrlsCount == 0 {
		p.totalUrlsCount = p.getTotalEntriesCount()
	}
	var (
		totalUp  int // количество обновленных ссылок
		errCount int // количество ошибок при обновлении
	)
	for totalUp < getUpdates || errCount >= 10 {
		redditUrls, err := reddit.GetVideoTitles(config.RdtNSFWHot)
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

		uSaver := new(ItemSaver).Init(p.db, config.TmNSFWVideo)
		for _, v := range convertedUrls {
			uSaver.Add(v)
		}
		uSaver.Save()
		afterInsertCount := p.getTotalEntriesCount()
		totalUp += afterInsertCount - p.totalUrlsCount
		p.totalUrlsCount = afterInsertCount
	}
	log.Printf("Updated %d urls", totalUp)
}

// Обновляет Image Items по названию категории
func (p *Provider) updateImageItems(catType string) {

}