package model

import "database/sql"

type Cart struct {
	ID          int        `json:"id"`
	UserID      string     `json:"user_id"`
	TotalAmount int        `json:"total_amount"`
	TotalPrice  float64    `json:"total_price"`
	Items       []CartItem `json:"cart_items"`
	CartStatus  string     `json:"-"`
}

type CartItem struct {
	ID          int              `json:"id"`
	CartID      int              `json:"cart_id,omitempty"`
	ProductID   int              `json:"product_id,omitempty"`
	Product     Product          `json:"product"`
	Amount      int              `json:"amount"`
	ItemVariant []CarttemVariant `json:"cart_item_variants,omitempty"`
	SubTotal    float64          `json:"subtotal"`
}

type CarttemVariant struct {
	ID              int           `json:"id"`
	CartItemID      int           `json:"-"`
	VariantID       sql.NullInt64 `json:"-"`
	Variant         Variant       `json:"variant"`
	OptionID        sql.NullInt64 `json:"-"`
	Option          VariantOption `json:"option"`
	AdditionalPrice float64       `json:"-"`
}

type CartItemDTO struct {
	ProductID int                  `json:"product_id"`
	Variant   []CartItemVariantDTO `json:"variant,omitempty"`
	Amount    int                  `json:"amount,omitempty"`
}

type CartItemVariantDTO struct {
	VariantID       int     `json:"variant_id,omitempty"`
	VariantOptionID int     `json:"variant_option_id,omitempty"`
	AdditionalPrice float64 `json:"-"`
}
