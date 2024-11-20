package model

type PaginatedResponse struct {
	StandardResponse
	Page       int `json:"page,omitempty"`
	Limit      int `json:"limit,omitempty"`
	TotalItems int `json:"total_items,omitempty"`
	TotalPages int `json:"total_pages,omitempty"`
}
