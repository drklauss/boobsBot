package telegram

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"strconv"

	"github.com/drklauss/boobsBot/config"
	"github.com/pkg/errors"
	"golang.org/x/net/proxy"
)

var tClient *client

// client is a telegram client
type client struct {
	config       *config.Telegram
	sender       *http.Client
	lastUpdateID int
}

// Init initialize client
func Init(c *config.Telegram) error {
	httpClient := &http.Client{}
	loadKeyboards()
	if c.Proxy != nil {
		a := &proxy.Auth{
			User:     c.Proxy.User,
			Password: c.Proxy.Password,
		}
		address := fmt.Sprintf("%s:%d", c.Proxy.Server, c.Proxy.Port)
		dialer, err := proxy.SOCKS5("tcp", address, a, proxy.Direct)
		if err != nil {
			return fmt.Errorf("could not dial to proxy %v: %v", c.Proxy, err)
		}
		httpClient.Transport = &http.Transport{
			Dial: dialer.Dial,
		}
	}
	tClient = &client{
		config: c,
		sender: httpClient,
	}

	return nil
}

// GetUpdateEntities returns telegram updates
func GetUpdateEntities() ([]Update, error) {
	var response Response
	u, _ := url.ParseRequestURI(tClient.config.API + tClient.config.Token)
	u.Path += "/getUpdates"
	params := url.Values{}
	params.Set("offset", strconv.Itoa(tClient.lastUpdateID+1))
	u.RawQuery = params.Encode()
	req, err := http.NewRequest("GET", u.String(), nil)
	if err != nil {
		return nil, errors.Wrap(err, "could not create request")
	}
	resp, err := tClient.sender.Do(req)
	if err != nil {
		return response.Result, errors.Wrap(err, "could not make request")
	}
	defer resp.Body.Close()
	responseBody, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return response.Result, errors.Wrap(err, "could not read body")
	}
	if err = json.Unmarshal(responseBody, &response); err != nil {
		return response.Result, errors.Wrap(err, "could not unmarshall body")
	}
	updLen := len(response.Result)
	if updLen > 0 {
		tClient.lastUpdateID = response.Result[updLen-1].UpdateID
	}
	return response.Result, nil
}

// SendMessage sends text message into chat
func SendMessage(mes MessageSend) error {
	u, _ := url.ParseRequestURI(tClient.config.API + tClient.config.Token)
	u.Path += "/sendMessage"
	params := url.Values{}
	params.Set("chat_id", strconv.Itoa(mes.ChatID))
	params.Set("text", mes.Text)
	if mes.KeyboardMarkup != nil {
		b, err := json.Marshal(mes.KeyboardMarkup)
		if err != nil {
			return errors.Wrap(err, "could not marshall keyboard")
		}
		params.Set("reply_markup", fmt.Sprintf("%s", b))
	}
	u.RawQuery = params.Encode()
	req, err := http.NewRequest("POST", u.String(), nil)
	if err != nil {
		return errors.Wrap(err, "could not create request")
	}
	resp, err := tClient.sender.Do(req)
	if err != nil {
		return errors.Wrap(err, "could not make request")
	}
	if resp.StatusCode != http.StatusOK {
		return errors.Wrap(err, "could not make request")
	}
	return nil
}

// SendImage sends image into chat
func SendImage(photo MediaSend) error {
	u, _ := url.ParseRequestURI(tClient.config.API + tClient.config.Token)
	u.Path += "/sendPhoto"
	params := url.Values{}
	params.Set("chat_id", strconv.Itoa(photo.ChatID))
	params.Set("photo", photo.URL)
	params.Set("caption", photo.Caption)
	if photo.KeyboardMarkup != nil {
		b, err := json.Marshal(photo.KeyboardMarkup)
		if err != nil {
			return errors.Wrap(err, "could not marshall keyboard")
		}
		params.Set("reply_markup", fmt.Sprintf("%s", b))
	}
	u.RawQuery = params.Encode()
	req, err := http.NewRequest("POST", u.String(), nil)
	if err != nil {
		return errors.Wrap(err, "could not create request")
	}
	resp, err := tClient.sender.Do(req)
	if err != nil {
		return errors.Wrap(err, "could not make request")
	}

	if resp.StatusCode != http.StatusOK {
		return errors.Wrap(err, "could not make request")
	}

	return nil
}

// SendDocument sends document into chat
func SendDocument(doc MediaSend) error {
	u, _ := url.ParseRequestURI(tClient.config.API + tClient.config.Token)
	u.Path += "/sendDocument"
	params := url.Values{}
	params.Set("chat_id", strconv.Itoa(doc.ChatID))
	params.Set("document", doc.URL)
	params.Set("caption", doc.Caption)
	if doc.KeyboardMarkup != nil {
		b, err := json.Marshal(doc.KeyboardMarkup)
		if err != nil {
			return errors.Wrap(err, "could not marshall keyboard")
		}
		params.Set("reply_markup", fmt.Sprintf("%s", b))
	}
	u.RawQuery = params.Encode()
	req, err := http.NewRequest("POST", u.String(), nil)
	if err != nil {
		return errors.Wrap(err, "could not create request")
	}
	resp, err := tClient.sender.Do(req)
	if err != nil {
		return errors.Wrap(err, "could not make request")
	}

	if resp.StatusCode != http.StatusOK {
		return errors.Wrap(err, "could not make request")
	}
	return nil
}
