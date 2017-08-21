package gfycat

import (
	"fmt"
	"boobsBot/config"

	"io/ioutil"
	"net/http"
)

// Класс обеспечивает получение url для gif
// curl -X POST -d "grant_type=password&username=dr_klauss&password=aL4514209"
// 		--user "hmOdEs1gOvXN4w:VwJV78wGCMGD2pvNQGyeaDHlzlk"
// 		--user-agent "My BoobsBot" https://www.reddit.com/api/v1/access_token
type Provider struct {
	urls []string
}

// Init инициализирует провайдер
// Парсит урлы и складывает их в срез
func (p *Provider) Init() Provider {
	getNewUrls()

	return *p
}

// Возвращает один URL гифки
// Берем один url и возвращаем его. Уменьшаем срез на один урл и проверяем длину среза.
// Если срез пустой, то снова запрашиваем новые урлы
func (p *Provider) GetUrl() string {
	url := p.urls[len(p.urls)-1]
	p.urls = p.urls[:len(p.urls)-1]
	p.updateUrls()

	return url
}

// Получает новые URL-ы гифок
// todo: пока написана тестовая опреация получения данных о пользователе. Переделать!
// todo: к тому же еще нужно обрабатывать все gyficat url-ы, для получения роликов в mp4
func getNewUrls() {

	client := new(http.Client)
	req, _ := http.NewRequest("GET", "https://oauth.reddit.com/api/v1/me", nil)
	req.Header.Set("Authorization", "bearer "+config.RedditToken)
	req.Header.Set("User-Agent", "My private BoobsBot")
	resp, err := client.Do(req)
	defer resp.Body.Close()
	responseBody, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Printf("%s\n", err)
	}
	fmt.Printf("%s\n", responseBody)

}

// Проверяем урлы, если иъ уже нет, то запрашиваем в срез новые
func (p *Provider) updateUrls() {

}
