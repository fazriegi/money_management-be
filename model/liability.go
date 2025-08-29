package model

type Liability struct {
	PeriodCode string      `db:"period_code" json:"period_code"`
	Name       string      `db:"name" json:"name"`
	Value      interface{} `db:"value" json:"value"`
	OrderNo    int         `db:"order_no" json:"order_no"`
	UserID     uint        `db:"user_id"`
}

type LiabilityResponse struct {
	PeriodCode string `json:"period_code"`
	Name       string `json:"name"`
	Value      int    `json:"value"`
	OrderNo    int    `json:"order_no"`
}

type GetLiabilityRequest struct {
	PaginationRequest
	Search     string `query:"search"`
	PeriodCode string `query:"period_code"`
}

type GetLiabilityResponse []LiabilityResponse

type UpdateLiabilityRequest struct {
	PeriodCode string      `json:"period_code" validate:"required"`
	Data       []Liability `json:"data"`
}
