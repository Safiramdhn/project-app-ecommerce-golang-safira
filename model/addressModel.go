package model

import "database/sql"

type Address struct {
	ID         int            `json:"id,omitempty"`
	UserID     string         `json:"user_id,omitempty"`
	Name       string         `json:"name,omitempty"`
	IsDefault  bool           `json:"is_default,omitempty"`
	Street     string         `json:"street,omitempty"`
	District   sql.NullString `json:"district,omitempty"`
	City       sql.NullString `json:"city,omitempty"`
	State      sql.NullString `json:"state,omitempty"`
	Country    string         `json:"country,omitempty"`
	PostalCode string         `json:"zipcode,omitempty"`
	Detail     `json:"-"`
}

type AddressDTO struct {
	ID         int    `json:"id"`                              // No validation for ID
	Name       string `json:"name" validate:"required"`        // Name is required
	IsDefault  bool   `json:"is_default"`                      // No validation for boolean
	Street     string `json:"street" validate:"required"`      // Street is required
	District   string `json:"district"`                        // District is required
	City       string `json:"city"`                            // City is required
	State      string `json:"state"`                           // State is required
	Country    string `json:"country" validate:"required"`     // Country is required
	PostalCode string `json:"postal_code" validate:"required"` // Required, numeric, 5 digits
}
