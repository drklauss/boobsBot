package dataProvider

import (
	"log"

	"fmt"
	"sync"

	"github.com/boobsBot/algorithm/config"
	"github.com/boobsBot/algorithm/dataProvider/gfycat"
	"github.com/boobsBot/algorithm/dataProvider/reddit"
)

// Класс является черным ящиком для получения нужного URL-а по категории
// https://gfycat.com/cajax/get/gfycatUrl - не посеять этот кусок кода
type Provider struct {
	Urls map[string][]string
}

// Init инициализирует провайдер
// Парсит урлы и складывает их в срез
func (p *Provider) Init() Provider {
	var wg sync.WaitGroup
	p.Urls = make(map[string][]string)
	wg.Add(3)
	go func() {
		p.updateUrls(config.New)
		wg.Done()
	}()
	go func() {
		p.updateUrls(config.Hot)
		wg.Done()
	}()
	go func() {
		p.updateUrls(config.Top)
		wg.Done()
	}()
	wg.Wait()

	return *p
}

// Возвращает один URL гифки, уменьшаем срез на один урл
// Если в срезе URL последний - обновляем срез
func (p *Provider) GetUrl(uType string) string {
	var u string
	ln := len(p.Urls[uType])
	if ln == 1 {
		fmt.Printf("equals to 1")
		go p.updateUrls(uType)
	}
	if ln > 0 {
		u = p.Urls[uType][len(p.Urls[uType])-1]
		p.Urls[uType] = p.Urls[uType][:len(p.Urls[uType])-1]
		fmt.Printf("+%v\n", p.Urls[uType])
	}
	log.Printf("Remain %d %s url(s)", len(p.Urls[uType]), uType)

	return u
}

// Обновляет URLs
func (p *Provider) updateUrls(uType string) {
	log.Printf("Starting update for %s urls", uType)
	redditUrls, err := reddit.GetNames(uType)
	fmt.Printf("Updated %d %v url(s)", len(redditUrls), uType)
	if err != nil {
		log.Println(err.Error())
	}
	convertedUrls, err := gfycat.ConvertNamesToUrls(redditUrls, uType)
	if err != nil {
		log.Println(err.Error())
	}
	log.Printf("Updated %d %s url(s)", len(convertedUrls), uType)
	p.Urls[uType] = convertedUrls
}
