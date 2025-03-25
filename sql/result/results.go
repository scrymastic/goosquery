package result

import (
	"github.com/scrymastic/goosquery/sql/sqlctx"
)

type Column struct {
	Name        string
	Type        string
	Description string
}

type Schema []Column

type Result map[string]interface{}

func (r *Result) Add(key string, value interface{}) {
	(*r)[key] = value
}

func (r *Result) Get(key string) interface{} {
	return (*r)[key]
}

func (r *Result) Set(key string, value interface{}) {
	// Only set the value if the key already exists
	if _, ok := (*r)[key]; ok {
		(*r)[key] = value
	}
}

func NewResult(ctx *sqlctx.Context, schema Schema) *Result {
	result := Result{}

	for _, col := range schema {
		if ctx.IsColumnUsed(col.Name) {
			switch col.Type {
			case "TEXT":
				result.Add(col.Name, "")
			case "INTEGER":
				result.Add(col.Name, int32(-1))
			case "BIGINT":
				result.Add(col.Name, int64(-1))
			case "FLOAT":
				result.Add(col.Name, float64(-1))
			default:
				result.Add(col.Name, nil)
			}
		}
	}

	return &result
}
