package model

import (
	"github.com/google/uuid"
)

type Currency struct {
	ID     uuid.UUID `json:"id" db:"id"`
	Name   string    `json:"name" db:"name"`
	Symbol string    `json:"symbol" db:"symbol"`
}
