package model

type Income struct {
	PeriodCode string      `db:"period_code" json:"period_code"`
	Type       string      `db:"type" json:"type"`
	Name       string      `db:"name" json:"name"`
	Value      interface{} `db:"value" json:"value"`
	OrderNo    int         `db:"order_no" json:"order_no"`
}

type IncomeResponse struct {
	PeriodCode string `json:"period_code"`
	Type       string `json:"type"`
	Name       string `json:"name"`
	Value      int    `json:"value"`
	OrderNo    int    `json:"order_no"`
}

type GetIncomeRequest struct {
	PaginationRequest
	Search string `json:"search"`
}

type GetIncomeResponse []IncomeResponse

type UpdateIncomeRequest struct {
	PeriodCode string   `json:"period_code" validate:"required"`
	Data       []Income `json:"data"`
}
