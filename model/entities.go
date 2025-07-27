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
	Page  *uint   `json:"page"`
	Limit *uint   `json:"limit"`
	Sort  *string `json:"sort"`
}

type PaginationResponse struct {
	Page             int   `json:"page"`
	TotalPages       int   `json:"total_pages"`
	TotalRows        int64 `json:"total_rows"`
	CurrentRowsCount int   `json:"current_rows_count"`
}
