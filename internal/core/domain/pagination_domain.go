package domain

type Pagination struct {
	Page      int         `json:"page"`
	Limit     int         `json:"limit"`
	Total     int64       `json:"total"`
	TotalPage int         `json:"total_page"`
	Data      interface{} `json:"data"`
}
