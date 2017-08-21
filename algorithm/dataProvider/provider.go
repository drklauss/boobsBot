package dataProvider

import (
	"log"

	"boobsBot/algorithm/dataProvider/reddit"
)

// Класс обеспечивает получение url для gif
// curl -X POST -d "grant_type=password&username=dr_klauss&password=aL4514209"
// 		--user "hmOdEs1gOvXN4w:VwJV78wGCMGD2pvNQGyeaDHlzlk"
// 		--user-agent "My BoobsBot" https://www.reddit.com/api/v1/access_token
type Provider struct {
	urls        map[string][]string
	RedditToken string
}

// Init инициализирует провайдер
// Парсит урлы и складывает их в срез
func (p *Provider) Init() Provider {
	reddit.GetUrls()
	log.Fatal("Stopped Manualy")

	return *p
}

// Возвращает один URL гифки, уменьшаем срез на один урл
func (p *Provider) GetUrl(uType string) string {
	u := p.urls[uType][len(p.urls)-1]
	p.urls[uType] = p.urls[uType][:len(p.urls[uType])-1]
	p.updateUrls(uType)

	return u
}

// Обновляет URL-ы
// Если срез пустой, то запрашиваем новые урлы
func (p *Provider) updateUrls(uType string) {
	if len(p.urls[uType]) != 0 {
		return
	}
	urls, err := reddit.GetUrls(uType)
	if err != nil {
		log.Println(err)
	}
	p.urls[uType] = urls
}
