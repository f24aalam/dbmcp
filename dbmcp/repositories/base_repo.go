package repositories

import (
	"context"
	"database/sql"
	"strings"
)

type BaseRepository struct {}

func (r *BaseRepository) RunSelectQuery(ctx context.Context, db *sql.DB, query string) ([]map[string]interface{}, int, error) {
	rows, err := db.QueryContext(ctx, strings.TrimSpace(query))
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()

	columns, err := rows.Columns()
	if err != nil {
		return nil, 0, err
	}

	var result []map[string]interface{}
	for rows.Next() {
		values := make([]interface{}, len(columns))
		valuesPtr := make([]interface{}, len(columns))

		for i := range columns {
			valuesPtr[i] = &values[i]
		}

		err = rows.Scan(valuesPtr...)
		if err != nil {
			return nil, 0, err
		}

		entry := make(map[string]interface{})
		for i, col := range columns {
			var v interface{}

			val := values[i]
			b, ok := val.([]byte)

			if ok {
				v = string(b)
			} else {
				v = val
			}

			entry[col] = v
		}

		result = append(result, entry)
	}

	return result, len(result), nil
}
