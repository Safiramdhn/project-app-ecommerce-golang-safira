package model

import (
	"database/sql"
	"time"
)

type Detail struct {
	Status    string       `json:"status,omitempty"`
	CreatedAt time.Time    `json:"created_at,omitempty"`
	UpdatedAt sql.NullTime `json:"updated_at,omitempty"`
	DeletedAt sql.NullTime `json:"deleted_at,omitempty"`
}
