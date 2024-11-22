package model

import "database/sql"

type User struct {
	ID             string         `json:"id,omitempty"`
	Name           string         `json:"name,omitempty"`
	Email          sql.NullString `json:"email"`
	PasswordHashed string         `json:"-"`
	PhoneNumber    sql.NullString `json:"phone_number"`
	// Wishlist       Wishlist       `json:"wishlist,omitempty"`
	Detail `json:"-"`
}

type UserDTO struct {
	Name               string `json:"name"`
	Password           string `json:"password" validate:"min=8"`
	EmailOrPhoneNumber string `json:"email_or_phone_number"`
}
