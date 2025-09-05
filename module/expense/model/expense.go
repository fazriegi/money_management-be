package model

type Expense struct {
	ID         uint        `db:"id"`
	CategoryId uint        `db:"category_id"`
	Date       interface{} `db:"date"`
	Value      string      `db:"value"`
	UserId     uint        `db:"user_id"`
	Notes      string      `db:"notes"`
}

type AddRequest struct {
	CategoryId uint        `json:"category_id" validate:"required"`
	Date       interface{} `json:"date" validate:"required"`
	Value      float64     `json:"value"`
	Notes      string      `json:"notes"`
}

type ExpenseCategory struct {
	ID   uint   `db:"id" json:"id"`
	Name string `db:"name" json:"name"`
}
