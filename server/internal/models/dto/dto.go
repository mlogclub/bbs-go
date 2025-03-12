package dto

// 站点导航
type ActionLink struct {
	Title string `json:"title"`
	Url   string `json:"url"`
}

type ApiRoute struct {
	Method string `json:"method"`
	Path   string `json:"path"`
	Name   string `json:"name"`
}
