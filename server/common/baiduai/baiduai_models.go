package baiduai

// AiData 要获取描述的数据
type AiData struct {
	Title string
	Desc  string
}

type AiTags struct {
	LogID     int64   `json:"log_id"`
	Items     []AiTag `json:"items"`
	ErrorMSG  string  `json:"error_msg"`
	ErrorCode int     `json:"error_code"`
}

type AiTag struct {
	Score float64 `json:"score"`
	Tag   string  `json:"tag"`
}

type AiTagParam struct {
	Title   string `json:"title"`
	Content string `json:"content"`
}

type AiCategories struct {
	LogID     int64          `json:"log_id"`
	Item      AiCategoryItem `json:"item"`
	ErrorMSG  string         `json:"error_msg"`
	ErrorCode int            `json:"error_code"`
}

type AiCategoryItem struct {
	TopCategory    []AiTag `json:"lv1_tag_list"`
	SecondCatrgory []AiTag `json:"lv2_tag_list"`
}

type AiAnalyzeRet struct {
	Tags    []string
	Summary string
}
