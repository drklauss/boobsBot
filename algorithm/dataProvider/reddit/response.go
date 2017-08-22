package reddit

type ErrorResponse struct {
	Message string `json:"message"`
	Error   int    `json:"error"`
}

type TokenResponse struct {
	Token     string `json:"access_token"`
	Type      string `json:"token_type"`
	ExpiresIn int    `json:"expires_in"`
	Scope     string `json:"scope"`
}

type NSFWResponse struct {

}