package model

type Product struct {
	ID           int       `json:"id,omitempty"`
	Name         string    `json:"name,omitempty"`
	Description  string    `json:"description,omitempty"`
	CategoryID   int       `json:"category_id,omitempty"`
	Category     Category  `json:"category,omitempty"`
	Price        float64   `json:"price,omitempty"`
	Discount     float64   `json:"discount,omitempty"`
	PhotoURL     string    `json:"photo_url,omitempty"`
	IsNewProduct bool      `json:"is_new_product,omitempty"`
	HasVariant   bool      `json:"has_variant,omitempty"`
	Rating       float64   `json:"rating,omitempty"`
	Variant      []Variant `json:"variants,omitempty"`
	Detail       `json:"-"`
}

type ProductDTO struct {
	Name       string  `json:"name"`
	CategoryID int     `json:"category_id"`
	Price      float64 `json:"price"`
	Discount   int     `json:"discount"`
	PhotoUrl   string  `json:"photo_url"`
	HasVariant bool    `json:"has_variant"`
	VariantDTO `json:"variant"`
}
