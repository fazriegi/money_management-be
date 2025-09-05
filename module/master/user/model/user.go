package model

type (
	User struct {
		ID       uint   `json:"id" db:"id"`
		Name     string `json:"name" db:"name"`
		Username string `json:"username" db:"username"`
		Email    string `json:"email" db:"email"`
		Password string `json:"password" db:"password"`
	}

	UserResponse struct {
		Name     string `json:"name"`
		Username string `json:"username"`
		Email    string `json:"email"`
	}
)
