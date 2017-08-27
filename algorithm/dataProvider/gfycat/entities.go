package gfycat

type Response struct {
	GfyItem GfyItem `json:"gfyItem"`
}
type GfyItem struct {
	MobileUrl string `json:"mobileUrl"`
}
