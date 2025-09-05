package model

import "github.com/fazriegi/money_management-be/module/common"

type GetCashflow struct {
	ID       uint        `db:"id"`
	Category string      `db:"category"`
	Date     interface{} `db:"date"`
	Value    string      `db:"value"`
	UserId   uint        `db:"user_id"`
}

type ListFilter struct {
	UserId    uint
	StartDate string `query:"start_date"`
	EndDate   string `query:"end_date"`
}

type ListRequest struct {
	common.PaginationRequest
	ListFilter
	Category string `query:"category"`
}

type CashflowData struct {
	ID       uint        `json:"id"`
	Category string      `json:"category"`
	Date     interface{} `json:"date"`
	Value    float64     `json:"value"`
}
