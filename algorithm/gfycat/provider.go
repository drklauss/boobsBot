package gfycat

import (
	"fmt"
	"log"

	"boobsBot/config"

	"io/ioutil"
	"net/http"

	"github.com/PuerkitoBio/goquery"
)

type Provider struct {
	urls []string
}

// Init инициализирует провайдер
// Парсит урлы и складывает их в срез
func (p *Provider) Init() Provider {
	parseMainPage()

	return *p
}

func (p *Provider) GetUrl() {
	parseMainPage()
}

func parseMainPage() {
	//u, _ := url.ParseRequestURI(config.RedditUrl)
	//var w http.ResponseWriter
	client := new(http.Client)
	req, _ := http.NewRequest("GET", config.RedditUrl, nil)
	//cookie := http.Cookie{Name: "over18", Value: "1"}
	//req.AddCookie(&cookie)
	resp, _ := client.Do(req)
	body, _ := ioutil.ReadAll(resp.Body)
	defer resp.Body.Close()
	//resp, err := http.Get(u.String())
	//fmt.Fprintf("%v\n", w.Header().Set())
	fmt.Printf("%s\n", body)

	doc, err := goquery.NewDocument(config.RedditUrl)
	if err != nil {
		log.Println(err)
	}
	println("here")

	// Find the review items
	doc.Find(".interstitial").Each(func(i int, s *goquery.Selection) {
		url, exist := s.Attr("data-url")
		fmt.Printf("%v\n", exist)
		if exist {

			fmt.Printf("%v\n", url)
		}

	})
}
