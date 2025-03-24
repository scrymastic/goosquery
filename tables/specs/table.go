package specs

import (
	"github.com/scrymastic/goosquery/sql/context"
)

type Column struct {
	Name        string
	Type        string
	Description string
}

type Schema []Column

func Init(ctx context.Context, schema Schema) map[string]interface{} {
	result := make(map[string]interface{})

	for _, col := range schema {
		if ctx.IsColumnUsed(col.Name) {
			switch col.Type {
			case "TEXT":
				result[col.Name] = ""
			case "INTEGER":
				result[col.Name] = int32(-1)
			case "BIGINT":
				result[col.Name] = int64(-1)
			case "FLOAT":
				result[col.Name] = float64(-1)
			default:
				result[col.Name] = nil
			}
		}
	}

	return result
}
