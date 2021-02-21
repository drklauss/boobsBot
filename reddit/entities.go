package reddit

import "fmt"

// ErrorResponse is a reddit error response.
type ErrorResponse struct {
	Message   string `json:"message"`
	ErrorCode int    `json:"error"`
}

func (er *ErrorResponse) Error() string {
	return fmt.Sprintf("error request to reddit API: code=%d, messages=%s", er.ErrorCode, er.Message)
}

// TokenResponse needs for token update.
type TokenResponse struct {
	Token     string `json:"access_token"`
	Type      string `json:"token_type"`
	ExpiresIn int    `json:"expires_in"`
	Scope     string `json:"scope"`
}

// SubRedditResponse is a subreddit response, that contains unnecessary info.
type SubRedditResponse struct {
	Data struct {
		Children []struct {
			Data `json:"data"`
		} `json:"children"`
	} `json:"data"`
	Category string
}

// Data is internal subreddit field, that contains necessary information for url convert.
type Data struct {
	Domain string `json:"domain"`
	URL    string `json:"url"`
	Name   string `json:"name"`
	Title  string `json:"title"`
}

// Convert converts part of the reddit request into entities ready for write into db.
func (sr *SubRedditResponse) Convert() []*Element {
	c, _ := NewConverter()
	return c.Run(sr)
}
