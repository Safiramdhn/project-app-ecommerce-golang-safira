package model

type Pagination struct {
	Page      int `json:"page,omitempty"`
	PerPage   int `json:"per_page,omitempty"`
	CountData int `json:"countData,omitempty"`
}
