package model

type Wishlist struct {
	ID        int     `json:"id,omitempty"`
	UserID    string  `json:"user_id,omitempty"`
	ProductID int     `json:"-"`
	Product   Product `json:"product"`
}

type WishlistDTO struct {
	UserID    string `json:"user_id"`
	ProductID int    `json:"product_id"`
}
