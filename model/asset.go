package model

type Asset struct {
	PeriodCode string      `db:"period_code" json:"period_code"`
	Name       string      `db:"name" json:"name"`
	Amount     interface{} `db:"amount" json:"amount"`
	Value      interface{} `db:"value" json:"value"`
	OrderNo    int         `db:"order_no" json:"order_no"`
}

type AssetRequest struct {
	PaginationRequest
	Search string `json:"search"`
}

type InsertAssetRequest struct {
	PeriodCode string  `json:"period_code" validate:"required"`
	Data       []Asset `json:"data"`
}
