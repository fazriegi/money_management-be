package model

import "github.com/fazriegi/money_management-be/module/common"

type Asset struct {
	ID         uint   `db:"id"`
	CategoryId uint   `db:"category_id"`
	Value      string `db:"value"`
	Amount     string `db:"amount"`
	UserId     uint   `db:"user_id"`
	Notes      string `db:"notes"`
}

type AddRequest struct {
	CategoryId uint    `json:"category_id" validate:"required"`
	Value      float64 `json:"value" validate:"required"`
	Amount     float64 `json:"amount" validate:"required"`
	Notes      string  `json:"notes" validate:"required"`
}

type AssetCategory struct {
	ID   uint   `db:"id" json:"id"`
	Name string `db:"name" json:"name"`
}

type ListRequest struct {
	common.PaginationRequest
	Keyword string `query:"keyword"`
	UserId  uint
}

type GetAsset struct {
	ID         uint   `db:"id"`
	CategoryId uint   `db:"category_id"`
	Category   string `db:"category"`
	Value      string `db:"value"`
	Amount     string `db:"amount"`
	UserId     uint   `db:"user_id"`
	Notes      string `db:"notes"`
}

type ListResponse struct {
	ID         uint    `json:"id"`
	CategoryId uint    `json:"category_id"`
	Category   string  `json:"category"`
	Value      float64 `json:"value"`
	Amount     float64 `json:"amount"`
	Notes      string  `json:"notes"`
}

type UpdateRequest struct {
	ID         uint
	CategoryId uint        `json:"category_id" validate:"required"`
	Amount     interface{} `json:"amount" validate:"required"`
	Value      float64     `json:"value" validate:"required"`
	Notes      string      `json:"notes" validate:"required"`
}
