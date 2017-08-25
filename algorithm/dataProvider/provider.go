package dataProvider

import (
	"log"

	"boobsBot/algorithm/config"
	"boobsBot/algorithm/dataProvider/reddit"
)

// Класс является черным ящиком для получения нужного URL-а по категории
type Provider struct {
	Urls        map[string][]string
	RedditToken string
}

// Init инициализирует провайдер
// Парсит урлы и складывает их в срез
func (p *Provider) Init() Provider {
	var err error
	p.Urls[config.New], err = reddit.GetUrls(config.New) // todo: обновить надо все, сделать горутину
	if err != nil {
		log.Println(err)
	}
	log.Fatal("Stopped Manualy")

	return *p
}

// Возвращает один URL гифки, уменьшаем срез на один урл
func (p *Provider) GetUrl(uType string) string {
	u := p.Urls[uType][len(p.Urls)-1]
	p.Urls[uType] = p.Urls[uType][:len(p.Urls[uType])-1]
	p.updateUrls(uType)

	return u
}

// Обновляет URL-ы
// Если срез пустой, то запрашиваем новые урлы
func (p *Provider) updateUrls(uType string) {
	if len(p.Urls[uType]) != 0 {
		return
	}
	urls, err := reddit.GetUrls(uType)
	if err != nil {
		log.Println(err)
	}
	p.Urls[uType] = urls
}
