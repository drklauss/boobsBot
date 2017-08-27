package gfycat

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"

	"github.com/boobsBot/algorithm/config"
)

var validUrls []string

// ConvertNamesToUrls преобразовывает названия в ссылки
func ConvertNamesToUrls(names []string) ([]string, error) {
	gfyCh := make(chan GfyItem, config.Threads)
	namesCh := make(chan string, len(names))
	for i := 1; i <= config.Threads; i++ {
		go gfyWorker(namesCh, gfyCh)
	}

	for _, name := range names {
		namesCh <- name
	}
	close(namesCh)

	for a := 0; a < len(names); a++ {
		gfy := <-gfyCh
		validUrls = append(validUrls, gfy.MobileUrl)
	}
	close(gfyCh)

	return validUrls, nil
}

func gfyWorker(namesCh <-chan string, gfyCh chan<- GfyItem) {
	for name := range namesCh {
		fmt.Printf("%v\n", name)
		client := new(http.Client)
		req, _ := http.NewRequest("GET", config.GfycatUrl+name, nil)
		resp, err := client.Do(req)
		if err != nil {
			log.Println(err)
			return
		}
		respBody, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			log.Println(err)
			return
		}
		var response Response
		err = json.Unmarshal(respBody, &response)
		if err != nil {
			log.Println(err)
			return
		}
		resp.Body.Close()

		fmt.Printf("%v\n", response.GfyItem)
		gfyCh <- response.GfyItem
	}

}
