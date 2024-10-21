package model

import (
	"github.com/google/uuid"
	"time"
)

type Product struct {
	ID            uuid.UUID `json:"id" db:"id"`
	Name          string    `json:"name" db:"name"`
	Description   string    `json:"description" db:"description"`
	ImageID       uuid.UUID `json:"image_id,omitempty" db:"image_id"`
	Price         int64     `json:"price" db:"price"`
	CurrencyID    int       `json:"currency_id" db:"currency_id"`
	Rating        int       `json:"rating" db:"rating"`
	CategoryID    int       `json:"category_id" db:"category_id"`
	Specification any       `json:"specification,omitempty" db:"specification"`
	CreatedAt     time.Time `json:"created_at" db:"created_at"`
	UpdatedAt     time.Time `json:"updated_at" db:"updated_at"`
}
