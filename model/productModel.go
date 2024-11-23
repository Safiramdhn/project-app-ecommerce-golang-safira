package model

type Product struct {
	ID                 int            `json:"id,omitempty"`
	Name               string         `json:"name,omitempty"`
	Description        string         `json:"description,omitempty"`
	CategoryID         int            `json:"category_id,omitempty"`
	Category           Category       `json:"category,omitempty"`
	Price              float64        `json:"price,omitempty"`
	Discount           float64        `json:"discount,omitempty"`
	PriceAfterDiscount float64        `json:"price_after_discount"`
	PhotoURL           string         `json:"photo_url,omitempty"`
	IsNewProduct       bool           `json:"is_new_product,omitempty"`
	HasVariant         bool           `json:"has_variant,omitempty"`
	Rating             float64        `json:"rating,omitempty"`
	TotalStock         int            `json:"total_stock,omitempty"`
	Variant            []Variant      `json:"variants,omitempty"`
	SpecialProduct     SpecialProduct `json:"special_products,omitempty"`
	Detail             `json:"-"`
}

type ProductDTO struct {
	Name       string     `json:"name"`
	CategoryID int        `json:"category_id" validate:"required,gt 0"`
	Price      float64    `json:"price" validate:"required,gt 0.0"`
	Discount   float64    `json:"discount"`
	PhotoUrl   string     `json:"photo_url"`
	HasVariant bool       `json:"has_variant"`
	VariantDTO VariantDTO `json:"variant"`
}

type SpecialProduct struct {
	IsBestSelling bool `json:"is_best_selling,omitempty"`
	IsNewProduct  bool `json:"is_new_product,omitempty"`
}

type WeeklyPromo struct {
	ID            int     `json:"id,omitempty"`
	ProductID     int     `json:"-"`
	Product       Product `json:"product,omitempty"`
	PromoDiscount float64 `json:"promo_discount,omitempty"`
	PromoPrice    float64 `json:"promo_price,omitempty"`
	StartDate     string  `json:"start_date"`
	EndDate       string  `json:"end_date"`
}
