package gfycat

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"

	"sync"

	"github.com/drklauss/boobsBot/algorithm/config"
)

// ConvertNamesToUrls преобразовывает названия в ссылки
func ConvertNamesToUrls(names []string) ([]string, error) {
	var validUrls []string
	mutex := sync.Mutex{}
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
		mutex.Lock()
		validUrls = append(validUrls, gfy.MobileUrl)
		mutex.Unlock()
	}
	close(gfyCh)

	return validUrls, nil
}

func gfyWorker(namesCh <-chan string, gfyCh chan<- GfyItem) {
	for name := range namesCh {
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
		gfyCh <- response.GfyItem
	}

}
