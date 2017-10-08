package gfycat

type Response struct {
	GfyItem GfyItem `json:"gfyItem"`
}
type GfyItem struct {
	MobileUrl string `json:"mobileUrl"`
	GfyName   string `json:"gfyName"`
}

func (item GfyItem) GetValue() string {
	return item.MobileUrl
}

func (item GfyItem) GetCaption() string {
	return item.GfyName
}
func (item GfyItem) GetCategoryId() string {
	return item.GfyName
}