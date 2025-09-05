package model

import "github.com/fazriegi/money_management-be/module/common"

type Income struct {
	ID         uint        `db:"id"`
	CategoryId uint        `db:"category_id"`
	Date       interface{} `db:"date"`
	Value      string      `db:"value"`
	UserId     uint        `db:"user_id"`
	Notes      string      `db:"notes"`
}

type GetIncome struct {
	ID         uint        `db:"id"`
	CategoryId uint        `db:"category_id"`
	Category   string      `db:"category"`
	Date       interface{} `db:"date"`
	Value      string      `db:"value"`
	UserId     uint        `db:"user_id"`
	Notes      string      `db:"notes"`
}

type IncomeData struct {
	ID         uint        `json:"id"`
	CategoryId uint        `json:"category_id"`
	Category   string      `json:"category"`
	Date       interface{} `json:"date"`
	Value      float64     `json:"value"`
	Notes      string      `json:"notes"`
}

type AddRequest struct {
	CategoryId uint        `json:"category_id" validate:"required"`
	Date       interface{} `json:"date" validate:"required"`
	Value      float64     `json:"value" validate:"required"`
	Notes      string      `json:"notes"`
}

type IncomeCategory struct {
	ID   uint   `db:"id" json:"id"`
	Name string `db:"name" json:"name"`
}

type ListRequest struct {
	common.PaginationRequest
	Keyword string `json:"keyword"`
	UserId  uint
}

type UpdateRequest struct {
	ID         uint
	CategoryId uint        `json:"category_id" validate:"required"`
	Date       interface{} `json:"date" validate:"required"`
	Value      float64     `json:"value"`
	Notes      string      `json:"notes"`
}
