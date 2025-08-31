package model

type Liability struct {
	ID          uint   `db:"id"`
	PeriodCode  string `db:"period_code"`
	Name        string `db:"name"`
	Value       string `db:"value"`
	Installment string `db:"installment"`
	OrderNo     int    `db:"order_no"`
	UserID      uint   `db:"user_id"`
}

type LiabilityRequest struct {
	ID          uint        `json:"id"`
	PeriodCode  string      `json:"period_code"`
	Name        string      `json:"name"`
	Value       interface{} `json:"value"`
	Installment interface{} `json:"installment"`
	OrderNo     int         `json:"order_no"`
}

type LiabilityResponse struct {
	ID          uint   `json:"id"`
	PeriodCode  string `json:"period_code"`
	Name        string `json:"name"`
	Value       int    `json:"value"`
	Installment int    `json:"installment"`
	OrderNo     int    `json:"order_no"`
}

type GetLiabilityRequest struct {
	PaginationRequest
	Search     string `query:"search"`
	PeriodCode string `query:"period_code"`
}

type GetLiabilityResponse []LiabilityResponse

type UpdateLiabilityRequest struct {
	PeriodCode string             `json:"period_code" validate:"required"`
	Data       []LiabilityRequest `json:"data"`
}

type ValidateDeleteRequest struct {
	ID uint `query:"id"`
}
