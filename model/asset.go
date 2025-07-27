package model

import (
	"github.com/google/uuid"
)

type Asset struct {
	PeriodCode string  `db:"period_code"`
	Name       string  `db:"name"`
	Amount     int     `db:"amount"`
	Value      float64 `db:"value"`
}

type AssetRequest struct {
	PaginationRequest
	Search string `json:"search"`
}

type AssetResponse struct {
	ID    uuid.UUID `json:"id"`
	Email string    `json:"email"`
	Name  string    `json:"name"`
}
