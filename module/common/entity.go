package common

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

func (s Response) CustomResponse(code int, message string, data any) Response {
	statuses := map[int]string{
		500: "internal server error",
		422: "unprocessable content",
		415: "unsupported media type",
		413: "request entity too large",
		404: "not found",
		401: "unauthorized",
		400: "bad request",
		303: "redirect",
		204: "no content",
		201: "created",
		200: "success",
	}

	return Response{
		Status: Status{
			Code:      code,
			Message:   message,
			Status:    statuses[code],
			IsSuccess: code >= 200 && code <= 299,
		},
		Data: data,
	}
}
