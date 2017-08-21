package reddit

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"strings"
)

var RedditToken string

func GetUrls(uType string) ([]string, error ){

}

func checkToken() {
	client := new(http.Client)
	req, _ := http.NewRequest("GET", GetMeUrl, nil)
	req.Header.Set("User-Agent", "My private BoobsBot")
	resp, err := client.Do(req)
	if err != nil {
		fmt.Printf("%s\n", err)
		log.Println(err)
	}
	defer resp.Body.Close()
	responseBody, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Printf("%s\n", err)
		log.Println(err)
	}
	if resp.
	fmt.Printf("%s\n", responseBody)
	log.Fatal("Stopped Manualy")
}

func getNewToken() (string, error) {
	client := new(http.Client)
	params := url.Values{}
	params.Set("grant_type", "password")
	params.Set("username", UserName)
	params.Set("password", Password)
	req, _ := http.NewRequest("POST", GetTokenUrl, strings.NewReader(params.Encode()))
	req.SetBasicAuth(ClientId, ClientSecret)
	req.Header.Set("User-Agent", "My private BoobsBot")
	resp, err := client.Do(req)
	defer resp.Body.Close()
	responseBody, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}
	fmt.Printf("%s\n", responseBody)
}

// Получает новые URL-ы гифок
// todo: пока написана тестовая опреация получения данных о пользователе. Переделать!
// todo: к тому же еще нужно обрабатывать все gyficat url-ы, для получения роликов в mp4
func getNewUrls() {

	client := new(http.Client)
	req, _ := http.NewRequest("GET", "https://oauth.reddit.com/api/v1/me", nil)
	req.Header.Set("Authorization", "bearer "+RedditToken)
	req.Header.Set("User-Agent", "My private BoobsBot")
	resp, err := client.Do(req)
	defer resp.Body.Close()
	responseBody, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Printf("%s\n", err)
	}
	fmt.Printf("%s\n", responseBody)
	log.Fatal("Stopped Manualy")

}
