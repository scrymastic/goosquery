package util

import (
	"github.com/scrymastic/goosquery/sql/context"
)

// InitColumns sets default values for all requested columns in a map
// This ensures all requested columns are present in the result, even if
// we cannot populate them with actual values
func InitColumns(ctx context.Context, columnDefs map[string]string) map[string]interface{} {
	result := make(map[string]interface{})

	// Initialize all requested columns with appropriate default values
	for col, colType := range columnDefs {
		if ctx.IsColumnUsed(col) {
			switch colType {
			case "string":
				result[col] = ""
			case "int32":
				result[col] = int32(-1)
			case "int64":
				result[col] = int64(-1)
			case "float64":
				result[col] = float64(-1)
			default:
				result[col] = nil // Default for unknown types
			}
		}
	}

	return result
}
