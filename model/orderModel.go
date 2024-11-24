package model

import "database/sql"

type Order struct {
	ID            int         `json:"id"`
	UserID        string      `json:"-"`
	AddressID     int         `json:"-"`
	Address       Address     `json:"address"`
	ShippingType  string      `json:"shipping_type"`
	ShippingCost  float64     `json:"shipping_cost"`
	PaymentMethod string      `json:"payment_method"`
	TotalAmount   int         `json:"total_amount"`
	TotalPrice    float64     `json:"total_price"`
	OrderItems    []OrderItem `json:"order_items"`
	OrderStatus   string      `json:"order_status"`
	CartID        int         `json:"cart_id"`
}

type OrderDTO struct {
	CartID        int     `json:"cart_id"`
	AddressID     int     `json:"address_id"`
	ShippingType  string  `json:"shipping_type"`
	ShippingCost  float64 `json:"shipping_cost"`
	PaymentMethod string  `json:"payment_method"`
	TotalPrice    float64 `json:"total_price"`
	TotalAmount   int     `json:"total_amount"`
}

type OrderItem struct {
	ID         int                `json:"-"`
	OrderID    int                `json:"-"`
	ProductID  int                `json:"-"`
	Product    Product            `json:"product"`
	Variants   []OrderItemVariant `json:"item_variants"`
	Amount     int                `json:"amount"`
	SubTotal   float64            `json:"subtotal"`
	Review     string             `json:"review"`
	Rating     float64            `json:"rating"`
	Photos     []string           `json:"photos"`
	CartItemID int                `json:"-"`
}

type OrderItemVariant struct {
	ID          int           `json:"id"`
	OrderItemID int           `json:"-"`
	VariantID   sql.NullInt64 `json:"-"`
	Variant     Variant       `json:"variant"`
	OptionID    sql.NullInt64 `json:"-"`
	Option      VariantOption
}
