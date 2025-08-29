package model

type Expense struct {
	PeriodCode  string      `db:"period_code" json:"period_code"`
	Name        string      `db:"name" json:"name"`
	Value       interface{} `db:"value" json:"value"`
	OrderNo     int         `db:"order_no" json:"order_no"`
	LiabilityID *int        `db:"liability_id" json:"liability_id"`
	UserID      uint        `db:"user_id"`
}

type ExpenseResponse struct {
	PeriodCode  string `json:"period_code"`
	Name        string `json:"name"`
	Value       int    `json:"value"`
	OrderNo     int    `json:"order_no"`
	LiabilityID *int   `json:"liability_id"`
}

type GetExpenseRequest struct {
	PaginationRequest
	Search     string `query:"search"`
	PeriodCode string `query:"period_code"`
}

type GetExpenseResponse []ExpenseResponse

type UpdateExpenseRequest struct {
	PeriodCode string    `json:"period_code" validate:"required"`
	Data       []Expense `json:"data"`
}
