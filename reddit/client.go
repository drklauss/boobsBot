package reddit

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"strconv"
	"strings"

	"github.com/drklauss/boobsBot/config"
	log "github.com/sirupsen/logrus"
)

var rClient *client

const (
	getTokenURL = "https://www.reddit.com/api/v1/access_token"
	getMeURL    = "https://www.reddit.com/api/v1/me"
	apiURL      = "https://oauth.reddit.com"
	agent       = "My Boobs Bot v2.0 (by dr.klauss)"
)

// client is a reddit client.
type client struct {
	config *config.Reddit
	token  *TokenResponse
	sender *http.Client
}

// Init reddit client.
func Init(c *config.Reddit) error {
	if rClient != nil {
		return nil
	}
	rClient = &client{
		config: c,
		sender: &http.Client{},
	}

	return checkToken()
}

// checkToken checks and refreshes token.
func checkToken() error {
	if err := isGoodToken(); err != nil {
		log.Warnf("error check token: %v", err)
		err := refreshToken()
		if err != nil {
			return err
		}
	}
	return nil
}

// GetItems gets the SubRedditResponse.
func GetItems(catPath string, last string) (*SubRedditResponse, error) {
	var err error
	if err = checkToken(); err != nil {
		return nil, err
	}

	req, _ := http.NewRequest("GET", apiURL+catPath, nil)
	data := req.URL.Query()
	data.Set("limit", strconv.Itoa(rClient.config.Limit))
	data.Set("after", last)
	req.URL.RawQuery = data.Encode()
	req.Header.Set("Authorization", "bearer "+rClient.token.Token)
	req.Header.Set("User-Agent", agent)
	resp, err := rClient.sender.Do(req)
	if err != nil {
		return nil, err
	}
	defer func() {
		if err := resp.Body.Close(); err != nil {
			log.Errorf("error close body: %v", err)
		}
	}()
	respBody, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	if resp.StatusCode != http.StatusOK {
		var errResp = ErrorResponse{}
		if err = json.Unmarshal(respBody, &errResp); err != nil {
			return nil, fmt.Errorf("could not unmarshall ErrorResponse: %w", err)
		}
		return nil, &errResp
	}
	var subResp = SubRedditResponse{Category: catPath}
	if err = json.Unmarshal(respBody, &subResp); err != nil {
		return nil, err
	}
	if len(subResp.Data.Children) == 0 {
		return nil, fmt.Errorf("subreddit children is empty")
	}

	return &subResp, nil
}

// refreshToken gets new token.
func refreshToken() error {
	log.Debug("refresh reddit token...")
	params := url.Values{}
	params.Set("grant_type", "password")
	params.Set("username", rClient.config.Username)
	params.Set("password", rClient.config.Password)
	req, _ := http.NewRequest("POST", getTokenURL, strings.NewReader(params.Encode()))
	req.SetBasicAuth(rClient.config.ClientID, rClient.config.Secret)
	req.Header.Set("User-Agent", agent)
	resp, err := rClient.sender.Do(req)
	if err != nil {
		return fmt.Errorf("could not send token request: %w", err)
	}
	defer func() {
		if err := resp.Body.Close(); err != nil {
			log.Errorf("error close body: %v", err)
		}
	}()
	respBody, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("could not read body with token: %w", err)
	}
	if resp.StatusCode != http.StatusOK {
		var errResp = ErrorResponse{}
		if err = json.Unmarshal(respBody, &errResp); err != nil {
			return fmt.Errorf("could not unmarshall ErrorResponse: %w", err)
		}
		return &errResp
	}

	if err = json.Unmarshal(respBody, &rClient.token); err != nil {
		return fmt.Errorf("could not unmarshall token: %w", err)
	}
	log.Debugf("new toket is: %v", rClient.token)

	return nil
}

// isGoodToken checks if token is alive.
func isGoodToken() error {
	if rClient.token == nil {
		return errors.New("token is nil")
	}
	req, _ := http.NewRequest("GET", getMeURL, nil)
	req.Header.Set("User-Agent", agent)
	resp, err := rClient.sender.Do(req)
	if err != nil {
		return fmt.Errorf("could not check token: %w", err)
	}
	respBody, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("could not read body: %w", err)
	}

	if resp.StatusCode != http.StatusOK {
		var errResp = ErrorResponse{}
		if err = json.Unmarshal(respBody, &errResp); err != nil {
			return fmt.Errorf("could not unmarshall ErrorResponse: %w", err)
		}
		return &errResp
	}

	return nil
}
