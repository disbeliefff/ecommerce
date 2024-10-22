package model

import (
	"github.com/google/uuid"
)

type Image struct {
	ID    uuid.UUID `json:"id" db:"id"`
	Name  string    `json:"name" db:"name"`
	Size  uint      `json:"size" db:"size"`
	Bytes []byte    `json:"bytes" db:"bytes"`
}
