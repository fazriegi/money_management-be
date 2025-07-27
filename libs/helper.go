package libs

import "github.com/fazriegi/money_management-be/model"

func CustomResponse(code int, message string) model.Status {
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

	status := model.Status{
		Code:      code,
		Message:   message,
		Status:    statuses[code],
		IsSuccess: code >= 200 && code <= 299,
	}

	return status
}
