package captcha

type CaptchaData struct {
	Id          string `json:"id"`
	ImageBase64 string `json:"imageBase64"`
	ThumbBase64 string `json:"thumbBase64"`
	ThumbSize   int    `json:"thumbSize"`
}
