package reddit

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"strconv"
	"strings"

	"github.com/pkg/errors"

	"github.com/leesper/holmes"

	"github.com/drklauss/boobsBot/config"
)

var rClient *client

const (
	getTokenURL = "https://www.reddit.com/api/v1/access_token"
	getMeURL    = "https://www.reddit.com/api/v1/me"
	apiURL      = "https://oauth.reddit.com"
	agent       = "My Boobs Bot v2.0 (by dr.klauss)"
)

// client is a reddit client
type client struct {
	config *config.Reddit
	token  *TokenResponse
	sender *http.Client
	last   map[string]string
}

// Init reddit client
func Init(c *config.Reddit) error {
	if rClient != nil {
		return nil
	}
	rClient = &client{
		config: c,
		sender: &http.Client{},
		last:   make(map[string]string),
	}
	if !isGoodToken() {
		err := refreshToken()
		if err != nil {
			return err
		}
	}
	return nil
}

// GetItems gets the SubRedditResponse
func GetItems(catPath string) (*SubRedditResponse, error) {
	if !isGoodToken() {
		err := refreshToken()
		if err != nil {
			return nil, err
		}
	}
	var err error

	req, _ := http.NewRequest("GET", apiURL+catPath, nil)
	data := req.URL.Query()
	data.Set("limit", strconv.Itoa(rClient.config.Limit))
	data.Set("after", "")
	if last, ok := rClient.last[catPath]; ok {
		data.Set("after", last)
	}
	req.URL.RawQuery = data.Encode()
	req.Header.Set("Authorization", "bearer "+rClient.token.Token)
	req.Header.Set("User-Agent", agent)
	resp, err := rClient.sender.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	respBody, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	subResp := SubRedditResponse{Category: catPath}
	err = json.Unmarshal(respBody, &subResp)
	if err != nil {
		return nil, err
	}
	if len(subResp.Data.Children) == 0 {
		return nil, fmt.Errorf("subreddit children is empty")
	}
	rClient.last[catPath] = subResp.Data.Children[len(subResp.Data.Children)-1].Data.Name

	return &subResp, nil
}

func refreshToken() error {
	params := url.Values{}
	params.Set("grant_type", "password")
	params.Set("username", rClient.config.Username)
	params.Set("password", rClient.config.Password)
	req, _ := http.NewRequest("POST", getTokenURL, strings.NewReader(params.Encode()))
	req.SetBasicAuth(rClient.config.ClientID, rClient.config.Secret)
	req.Header.Set("User-Agent", agent)
	resp, err := rClient.sender.Do(req)
	if err != nil {
		return errors.Wrap(err, "could not send token request")
	}
	defer resp.Body.Close()
	respBody, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return errors.Wrap(err, "could not read body with token")
	}
	err = json.Unmarshal(respBody, &rClient.token)
	if err != nil {
		return errors.Wrap(err, "could not unmarshall token")
	}

	return nil
}

// Проверяет действителен ли токен
func isGoodToken() bool {
	if rClient.token == nil {
		holmes.Warnln("token is nil")
		return false
	}
	req, _ := http.NewRequest("GET", getMeURL, nil)
	req.Header.Set("User-Agent", agent)
	resp, err := rClient.sender.Do(req)
	if err != nil {
		holmes.Warnf("could not check token: %v", err)
		return false
	}
	if resp.StatusCode != http.StatusOK {
		holmes.Warnf("it seems token is bad, cause status code is %d", resp.StatusCode)
		return false
	}

	return true
}
