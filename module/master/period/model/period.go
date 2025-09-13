package model

type MonthlyPeriod struct {
	DayOfMonth uint8 `db:"day_of_month" json:"day_of_month"`
}
