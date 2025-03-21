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
			case "string":
				result[col.Name] = ""
			case "int32":
				result[col.Name] = int32(-1)
			case "int64":
				result[col.Name] = int64(-1)
			case "float64":
				result[col.Name] = float64(-1)
			default:
				result[col.Name] = nil
			}
		}
	}

	return result
}
