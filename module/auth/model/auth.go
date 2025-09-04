package model

type (
	RegisterRequest struct {
		Name     string `json:"name" validate:"required"`
		Username string `json:"username" validate:"required"`
		Email    string `json:"email"`
		Password string `json:"password" validate:"required"`
	}

	LoginRequest struct {
		Username string `json:"username" validate:"required"`
		Password string `json:"password" validate:"required"`
	}
)
