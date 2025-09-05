package libs

import (
	"errors"
	"reflect"
	"strings"

	"github.com/doug-martin/goqu/v9"
	_ "github.com/doug-martin/goqu/v9/dialect/mysql"
	"github.com/fazriegi/money_management-be/module/common"
	"github.com/jmoiron/sqlx"
)

func PaginationRequest(dataset *goqu.SelectDataset, req common.PaginationRequest) *goqu.SelectDataset {
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

	sliceType := destVal.Elem().Type().Elem()
	if sliceType.Kind() != reflect.Struct {
		return errors.New("slice elements must be structs")
	}

	columns, err := rows.Columns()
	if err != nil {
		return err
	}

	columnTypes, err := rows.ColumnTypes()
	if err != nil {
		return err
	}

	for rows.Next() {
		elem := reflect.New(sliceType).Elem()
		// Create a DB column → struct field mapping based on the db tag
		fieldMap := make(map[string]reflect.Value)
		for i := 0; i < elem.NumField(); i++ {
			field := sliceType.Field(i)
			tag := field.Tag.Get("db")
			if tag == "" {
				continue
			}
			fieldMap[tag] = elem.Field(i)
		}

		// Prepare a place to scan query results
		scanArgs := make([]interface{}, len(columns))
		for i, col := range columns {
			if field, exists := fieldMap[col]; exists {
				// Handle interface{} fields by scanning into appropriate types
				if field.Kind() == reflect.Interface {
					switch columnTypes[i].DatabaseTypeName() {
					case "INT", "INTEGER", "BIGINT":
						var v int64
						scanArgs[i] = &v
					case "FLOAT", "REAL", "DOUBLE":
						var v float64
						scanArgs[i] = &v
					case "DECIMAL", "NUMERIC":
						var v float64
						scanArgs[i] = &v
					case "VARCHAR", "TEXT", "STRING":
						var v string
						scanArgs[i] = &v
					default:
						// Fallback to interface{} for unknown types
						scanArgs[i] = field.Addr().Interface()
					}
				} else {
					scanArgs[i] = field.Addr().Interface()
				}
			} else {
				// Ignore unused columns
				var dummy interface{}
				scanArgs[i] = &dummy
			}
		}

		if err := rows.Scan(scanArgs...); err != nil {
			return err
		}

		// Convert scanned values to interface{} for target fields
		for i, col := range columns {
			if field, exists := fieldMap[col]; exists && field.Kind() == reflect.Interface {
				val := reflect.ValueOf(scanArgs[i]).Elem().Interface()
				field.Set(reflect.ValueOf(val))
			}
		}

		destVal.Elem().Set(reflect.Append(destVal.Elem(), elem))
	}

	return rows.Err()
}

func GetDialect() goqu.DialectWrapper {
	return goqu.Dialect("mysql")
}
