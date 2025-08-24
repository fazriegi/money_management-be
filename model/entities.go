package model

type Response struct {
	Status
	Data any `json:"data"`
}

type Status struct {
	Code      int    `json:"code"`
	Status    string `json:"status"`
	Message   string `json:"message"`
	IsSuccess bool   `json:"is_success"`
}

type PaginationRequest struct {
	Page  *uint   `query:"page"`
	Limit *uint   `query:"limit"`
	Sort  *string `query:"sort"`
}

type PaginationResponse struct {
	Page             int   `json:"page"`
	TotalPages       int   `json:"total_pages"`
	TotalRows        int64 `json:"total_rows"`
	CurrentRowsCount int   `json:"current_rows_count"`
}
