package model

type Asset struct {
	ID         uint   `db:"id"`
	PeriodCode string `db:"period_code"`
	Name       string `db:"name"`
	Amount     string `db:"amount"`
	Value      string `db:"value"`
	OrderNo    int    `db:"order_no"`
	UserID     uint   `db:"user_id"`
}

type AssetRequest struct {
	ID         uint        `json:"id"`
	PeriodCode string      `json:"period_code"`
	Name       string      `json:"name"`
	Amount     interface{} `json:"amount"`
	Value      interface{} `json:"value"`
	OrderNo    int         `json:"order_no"`
}

type AssetResponse struct {
	ID         uint   `json:"id"`
	PeriodCode string `json:"period_code"`
	Name       string `json:"name"`
	Amount     int    `json:"amount"`
	Value      int    `json:"value"`
	OrderNo    int    `json:"order_no"`
}

type GetAssetRequest struct {
	PaginationRequest
	Search     string `query:"search"`
	PeriodCode string `query:"period_code"`
}

type InsertAssetRequest struct {
	PeriodCode string         `json:"period_code" validate:"required"`
	Data       []AssetRequest `json:"data" validate:"required"`
}
