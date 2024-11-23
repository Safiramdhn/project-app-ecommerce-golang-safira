package model

import "database/sql"

type Cart struct {
	ID               int             `json:"id"`
	UserID           string          `json:"user_id"`
	ProductID        int             `json:"-"`
	VariantIDs       []sql.NullInt64 `json:"-"`
	VariantOptionIDs []sql.NullInt64 `json:"-"`
	Product          Product         `json:"product"`
	Variants         []Variant       `json:"variants"`
	VariantOptions   []VariantOption `json:"variant_options"`
	Amount           int             `json:"amount"`
	TotalPrice       float64         `json:"total_price"`
	Detail           `json:"-"`      //
}

type CartDTO struct {
	ProductID        int   `json:"product_id"`
	VariantIDs       []int `json:"variant_id"`
	VariantOptionIDs []int `json:"variant_option_id"`
	Amount           int   `json:"amount,omitempty"`
}
