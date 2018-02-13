package gfycat

import (
	"encoding/json"
	"io/ioutil"
	"net/http"

	"sync"

	"strings"

	"github.com/drklauss/boobsBot/algorithm/config"
	"github.com/fatih/camelcase"
)

// ConvertNamesToUrls преобразовывает названия в ссылки
func ConvertNamesToUrls(names []string) ([]*GfyItem, error) {
	var gfyItems []*GfyItem
	gfyCh := make(chan *GfyItem, config.ConvertNamesThreads)
	namesCh := make(chan string, len(names))
	for i := 1; i <= config.ConvertNamesThreads; i++ {
		go gfyWorker(namesCh, gfyCh)
	}

	for _, name := range names {
		namesCh <- name
	}
	close(namesCh)

	mutex := sync.Mutex{}
	for a := 0; a < len(names); a++ {
		mutex.Lock()
		gfyItems = append(gfyItems, <-gfyCh)
		mutex.Unlock()
	}
	close(gfyCh)
	convertCamelCaseNames(gfyItems)

	return gfyItems, nil
}

// Отдельный процесс для обработки данных
func gfyWorker(namesCh <-chan string, gfyCh chan<- *GfyItem) {
	for name := range namesCh {
		client := new(http.Client)
		req, _ := http.NewRequest("GET", config.GfycatUrl+name, nil)
		resp, err := client.Do(req)
		if err != nil {
			gfyCh <- &GfyItem{}
			return
		}
		respBody, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			gfyCh <- &GfyItem{}
			return
		}
		var response Response
		err = json.Unmarshal(respBody, &response)
		if err != nil {
			gfyCh <- &GfyItem{}
			return
		}
		resp.Body.Close()
		gfyCh <- &response.GfyItem
	}

}

// Преоборазовывает camelCase названия в обычные
func convertCamelCaseNames(gfyItems []*GfyItem) {
	for _, gfyItem := range gfyItems {
		splittedWords := camelcase.Split(gfyItem.GfyName)
		gfyItem.GfyName = strings.Join(splittedWords, " ")
	}
}
