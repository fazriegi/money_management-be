package libs

import (
	"errors"
	"reflect"
	"strings"

	"github.com/doug-martin/goqu/v9"
	"github.com/fazriegi/money_management-be/model"
	"github.com/jmoiron/sqlx"
)

func PaginationRequest(dataset *goqu.SelectDataset, req model.PaginationRequest) *goqu.SelectDataset {
	if req.Sort != nil && *req.Sort != "" {
		parts := strings.Fields(*req.Sort)
		if len(parts) > 0 {
			field := parts[0]
			direction := "ASC"
			if len(parts) >= 2 && strings.ToUpper(parts[1]) == "DESC" {
				direction = "DESC"
			}

			col := goqu.I(field)
			if direction == "DESC" {
				dataset = dataset.Order(col.Desc())
			} else {
				dataset = dataset.Order(col.Asc())
			}
		}
	}

	if req.Page != nil && req.Limit != nil && *req.Page > 0 && *req.Limit > 0 {
		offset := (*req.Page - 1) * *req.Limit
		dataset = dataset.Limit(*req.Limit).Offset(offset)
	}

	return dataset
}

func ScanRowsIntoStructs(rows *sqlx.Rows, destSlice interface{}) error {
	destVal := reflect.ValueOf(destSlice)
	if destVal.Kind() != reflect.Ptr || destVal.Elem().Kind() != reflect.Slice {
		return errors.New("destSlice must be a pointer to a slice")
	}

	columns, err := rows.Columns()
	if err != nil {
		return err
	}

	for rows.Next() {
		// Create a new struct of the slice element type
		elemType := destVal.Elem().Type().Elem()
		elemPtr := reflect.New(elemType)
		elem := elemPtr.Elem()

		fieldMap := make(map[string]reflect.Value)
		for i := 0; i < elem.NumField(); i++ {
			field := elem.Type().Field(i)
			tag := field.Tag.Get("db")
			if tag != "" {
				fieldMap[tag] = elem.Field(i)
			}
		}

		scanArgs := make([]interface{}, len(columns))
		for i, col := range columns {
			if field, ok := fieldMap[col]; ok {
				scanArgs[i] = field.Addr().Interface()
			} else {
				var dummy interface{}
				scanArgs[i] = &dummy
			}
		}

		if err := rows.Scan(scanArgs...); err != nil {
			return err
		}

		destVal.Elem().Set(reflect.Append(destVal.Elem(), elem))
	}

	return rows.Err()
}
