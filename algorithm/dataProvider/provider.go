package dataProvider

import (
	"log"

	"time"

	"github.com/boobsBot/algorithm/config"
	"github.com/boobsBot/algorithm/dataProvider/gfycat"
	"github.com/boobsBot/algorithm/dataProvider/reddit"
	"fmt"
)

// Класс является черным ящиком для получения нужного URL-а по категории
// https://gfycat.com/cajax/get/gfycatUrl - не посеять этот кусок кода
type Provider struct {
	Urls map[string][]string
}

// Init инициализирует провайдер
// Парсит урлы и складывает их в срез
func (p *Provider) Init() Provider {
	var err error
	p.Urls = make(map[string][]string)
	go p.updateUrls(config.New)
	go p.updateUrls(config.Hot)
	go p.updateUrls(config.Top)
	if err != nil {
		log.Println(err)
	}
	time.Sleep(10 * time.Second)
	fmt.Printf("%v\n", p.Urls)
	log.Fatal("Stopped Manualy")

	return *p
}

// Возвращает один URL гифки, уменьшаем срез на один урл
// Если в срезе URL последний - обновляем срез
func (p *Provider) GetUrl(uType string) string {
	var u string
	ln := len(p.Urls[uType])
	if ln == 1 {
		go p.updateUrls(uType)
	}
	if ln > 0 {
		u = p.Urls[uType][len(p.Urls)-1]
		p.Urls[uType] = p.Urls[uType][:len(p.Urls[uType])-1]
	}
	log.Printf("Remain %d, %s url(s)", p.Urls[uType], uType)

	return u
}

// Обновляет URLs
func (p *Provider) updateUrls(uType string) {
	log.Printf("Starting update for %s urls", uType)
	redditUrls, err := reddit.GetNames(uType)
	if err != nil {
		log.Println(err.Error())
	}
	convertedUrls, err := gfycat.ConvertNamesToUrls(redditUrls)
	if err != nil {
		log.Println(err.Error())
	}
	log.Printf("Updated %d %s url(s)", len(convertedUrls), uType)
	p.Urls[uType] = convertedUrls
}
