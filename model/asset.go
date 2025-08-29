package model

type Asset struct {
	PeriodCode string      `db:"period_code" json:"period_code"`
	Name       string      `db:"name" json:"name"`
	Amount     interface{} `db:"amount" json:"amount"`
	Value      interface{} `db:"value" json:"value"`
	OrderNo    int         `db:"order_no" json:"order_no"`
	UserID     uint        `db:"user_id"`
}

type AssetRequest struct {
	PaginationRequest
	Search     string `query:"search"`
	PeriodCode string `query:"period_code"`
}

type InsertAssetRequest struct {
	PeriodCode string  `json:"period_code" validate:"required"`
	Data       []Asset `json:"data" validate:"required"`
}
