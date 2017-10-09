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

type SubRedditResponse struct {
	Data CommonData `json:"data"`
}

type CommonData struct {
	Children []Children `json:"children"`
}

type Children struct {
	Data Data `json:"data"`
}

type Data struct {
	Domain string `json:"domain"`
	Url    string `json:"url"`
	Name   string `json:"name"`  // Используется при вытаскивания с gfycat
	Title  string `json:"title"` // Используется для именования картинок
}

type ImageItem struct {
	UrlValue string
	Caption   string
}

func (item *ImageItem) GetValue() string {
	return item.UrlValue
}

func (item *ImageItem) GetCaption() string {
	return item.Caption
}
