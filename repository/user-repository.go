package repository

import (
	"fmt"

	"github.com/doug-martin/goqu/v9"
	"github.com/fazriegi/money_management-be/libs"
	"github.com/fazriegi/money_management-be/model"
	"github.com/jmoiron/sqlx"
)

type IUserRepository interface {
	GetUserByUsername(username string, db *sqlx.DB) (model.User, error)
	InsertUser(data *model.User, db *sqlx.DB) error
}

type UserRepository struct {
}

func NewUserRepository() IUserRepository {
	return &UserRepository{}
}

func (r *UserRepository) GetUserByUsername(username string, db *sqlx.DB) (result model.User, err error) {
	dialect := libs.GetDialect()

	dataset := dialect.From("users").Select(goqu.I("username"), goqu.I("password"), goqu.I("email"), goqu.I("id"), goqu.I("name")).Where(goqu.I("username").Eq(username))

	query, val, err := dataset.ToSQL()
	if err != nil {
		return result, fmt.Errorf("failed to build SQL query: %w", err)
	}

	err = db.Get(&result, query, val...)
	if err != nil {
		return result, err
	}

	return
}

func (r *UserRepository) InsertUser(data *model.User, db *sqlx.DB) error {
	dialect := libs.GetDialect()

	dataset := dialect.Insert("users").Rows(*data)
	sql, val, err := dataset.ToSQL()
	if err != nil {
		return fmt.Errorf("failed to build SQL query: %w", err)
	}

	_, err = db.Exec(sql, val...)
	if err != nil {
		return fmt.Errorf("failed to execute query: %w", err)
	}

	return nil
}
