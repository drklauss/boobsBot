package telegram

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net"
	"net/http"
	"net/url"
	"strconv"

	"github.com/drklauss/boobsBot/config"
	log "github.com/sirupsen/logrus"
	"golang.org/x/net/proxy"
)

var tClient *client

// client is a telegram client.
type client struct {
	config       *config.Telegram
	sender       *http.Client
	lastUpdateID int
}

// Init initialize client.
func Init(c *config.Telegram) error {
	httpClient := &http.Client{}
	loadKeyboards()
	if c.Proxy != nil {
		var a *proxy.Auth
		if c.Proxy.User != "" {
			a = &proxy.Auth{
				User:     c.Proxy.User,
				Password: c.Proxy.Password,
			}
		}
		address := fmt.Sprintf("%s:%d", c.Proxy.Server, c.Proxy.Port)
		dialer, err := proxy.SOCKS5("tcp", address, a, proxy.Direct)
		dialContext := func(ctx context.Context, network, address string) (net.Conn, error) {
			return dialer.Dial(network, address)
		}
		if err != nil {
			return fmt.Errorf("could not dial to proxy %v: %v", c.Proxy, err)
		}
		httpClient.Transport = &http.Transport{
			DialContext: dialContext,
			//Dial: dialer.Dial,
		}
		log.Debugf("proxy is used %s:%d", c.Proxy.Server, c.Proxy.Port)
	}
	tClient = &client{
		config: c,
		sender: httpClient,
	}

	return nil
}

// GetUpdateEntities returns telegram updates.
func GetUpdateEntities() ([]Update, error) {
	var response Response
	u, _ := url.ParseRequestURI(tClient.config.API + tClient.config.Token)
	u.Path += "/getUpdates"
	params := url.Values{}
	params.Set("offset", strconv.Itoa(tClient.lastUpdateID+1))
	u.RawQuery = params.Encode()
	req, err := http.NewRequest("GET", u.String(), nil)
	if err != nil {
		return nil, fmt.Errorf("could not create request: %w", err)
	}
	resp, err := tClient.sender.Do(req)
	if err != nil {
		return response.Result, fmt.Errorf("could not make request: %w", err)
	}
	defer resp.Body.Close()
	responseBody, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return response.Result, fmt.Errorf("could not read body: %w", err)
	}
	if err = json.Unmarshal(responseBody, &response); err != nil {
		return response.Result, fmt.Errorf("could not unmarshall body: %w", err)
	}
	updLen := len(response.Result)
	if updLen > 0 {
		tClient.lastUpdateID = response.Result[updLen-1].UpdateID
	}
	return response.Result, nil
}

// SendMessage sends text message into chat.
func SendMessage(mes MessageSend) error {
	u, _ := url.ParseRequestURI(tClient.config.API + tClient.config.Token)
	u.Path += "/sendMessage"
	params := url.Values{}
	params.Set("chat_id", strconv.Itoa(mes.ChatID))
	params.Set("text", mes.Text)
	if mes.KeyboardMarkup != nil {
		b, err := json.Marshal(mes.KeyboardMarkup)
		if err != nil {
			return fmt.Errorf("could not marshall keyboard: %w", err)
		}
		params.Set("reply_markup", fmt.Sprintf("%s", b))
	}
	u.RawQuery = params.Encode()
	req, err := http.NewRequest("POST", u.String(), nil)
	if err != nil {
		return fmt.Errorf("could not create request: %w", err)
	}
	resp, err := tClient.sender.Do(req)
	if err != nil {
		return fmt.Errorf("could not make request: %w", err)
	}
	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("could not make request: %w", err)
	}
	return nil
}

// SendImage sends image into chat.
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
			return fmt.Errorf("could not marshall keyboard: %w", err)
		}
		params.Set("reply_markup", fmt.Sprintf("%s", b))
	}
	u.RawQuery = params.Encode()
	req, err := http.NewRequest("POST", u.String(), nil)
	if err != nil {
		return fmt.Errorf("could not create request: %w", err)
	}
	resp, err := tClient.sender.Do(req)
	if err != nil {
		return fmt.Errorf("could not make request: %w", err)
	}

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("could not make request: %w", err)
	}

	return nil
}

// SendDocument sends document into chat.
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
			return fmt.Errorf("could not marshall keyboard: %w", err)
		}
		params.Set("reply_markup", fmt.Sprintf("%s", b))
	}
	u.RawQuery = params.Encode()
	req, err := http.NewRequest("POST", u.String(), nil)
	if err != nil {
		return fmt.Errorf("could not create request: %w", err)
	}
	resp, err := tClient.sender.Do(req)
	if err != nil {
		return fmt.Errorf("could not make request: %w", err)
	}

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("could not make request: %w", err)
	}
	return nil
}

// SendAnswerCallbackQuery sends query answer.
func SendAnswerCallbackQuery(acq AnswerCallbackQuery) error {
	u, _ := url.ParseRequestURI(tClient.config.API + tClient.config.Token)
	u.Path += "/answerCallbackQuery"
	params := url.Values{}
	params.Set("callback_query_id", acq.CallbackQueryID)
	params.Set("text", acq.Text)
	if acq.ShowAlert {
		params.Set("show_alert", "true")
	}
	// params.Set("url", acq.URL)
	u.RawQuery = params.Encode()
	req, err := http.NewRequest("POST", u.String(), nil)
	if err != nil {
		return fmt.Errorf("could not create request: %w", err)
	}
	resp, err := tClient.sender.Do(req)
	if err != nil {
		return fmt.Errorf("could not make request: %w", err)
	}

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("could not make request: %w", err)
	}
	return nil
}
