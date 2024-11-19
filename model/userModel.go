package model

import "database/sql"

type User struct {
	ID             string         `json:"id,omitempty"`
	Name           string         `json:"name,omitempty"`
	Email          sql.NullString `json:"email," validate:"email"`
	PasswordHashed string         `json:"-"`
	PhoneNumber    sql.NullString `json:"phone_number"`
	Status         string         `json:"status,omitempty"`
	CreatedAt      string         `json:"created_at,omitempty"`
	UpdatedAt      sql.NullString `json:"updated_at,omitempty"`
	DeletedAt      sql.NullString `json:"deleted_at,omitempty"`
}

type UserDTO struct {
	Name               string `json:"name"`
	Password           string `json:"password" validate:"min=8"`
	PhoneNumber        string `json:"phone_number"`
	EmailOrPhoneNumber string `json:"email_or_phone_number"`
}
