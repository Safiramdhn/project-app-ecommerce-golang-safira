package model

type Recommendation struct {
	ID            int     `json:"id,omitempty"`
	Product       Product `json:"product"`
	IsRecommended bool    `json:"is_recommended,omitempty"`
	SetInBanner   bool    `json:"set_in_banner,omitempty"`
	Title         string  `json:"title"`
	Subtitle      string  `json:"subtitle"`
	PhotoUrl      string  `json:"photo_url"`
	PathUrl       string  `json:"path_url"`
	Detail        `json:"-"`
}

type RecommendationDTO struct {
	IsRecommended bool `json:"is_recommended,omitempty"`
	SetInBanner   bool `json:"set_in_banner,omitempty"`
}
