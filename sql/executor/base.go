package executor

import (
	"fmt"
	"regexp"
	"strconv"
	"strings"

	"github.com/blastrain/vitess-sqlparser/sqlparser"
)

// BaseExecutor provides common functionality for all executors
type BaseExecutor struct{}

// MatchesWhereClause checks if a row matches the WHERE clause
func (e *BaseExecutor) MatchesWhereClause(row map[string]interface{}, expr sqlparser.Expr) bool {
	switch expr := expr.(type) {
	case *sqlparser.ComparisonExpr:
		return e.EvaluateComparison(row, expr)
	case *sqlparser.AndExpr:
		return e.MatchesWhereClause(row, expr.Left) && e.MatchesWhereClause(row, expr.Right)
	case *sqlparser.OrExpr:
		return e.MatchesWhereClause(row, expr.Left) || e.MatchesWhereClause(row, expr.Right)
	case *sqlparser.ParenExpr:
		return e.MatchesWhereClause(row, expr.Expr)
	}
	return false
}

// EvaluateComparison evaluates a comparison expression against a row
func (e *BaseExecutor) EvaluateComparison(row map[string]interface{}, expr *sqlparser.ComparisonExpr) bool {
	// Get left side (field name)
	colName, ok := expr.Left.(*sqlparser.ColName)
	if !ok {
		return false
	}
	fieldName := colName.Name.String()

	// Get field value from row
	fieldValue, exists := row[fieldName]
	if !exists {
		// Try case-insensitive matching
		for k, v := range row {
			if strings.EqualFold(k, fieldName) {
				fieldValue = v
				exists = true
				break
			}
		}
		if !exists {
			return false
		}
	}

	// Get right side (value)
	var targetValue interface{}
	switch val := expr.Right.(type) {
	case *sqlparser.SQLVal:
		// Handle different value types
		switch val.Type {
		case sqlparser.StrVal:
			targetValue = string(val.Val)
		case sqlparser.IntVal:
			intVal, err := strconv.ParseInt(string(val.Val), 10, 64)
			if err != nil {
				return false
			}
			targetValue = intVal
		case sqlparser.FloatVal:
			floatVal, err := strconv.ParseFloat(string(val.Val), 64)
			if err != nil {
				return false
			}
			targetValue = floatVal
		}
	case *sqlparser.NullVal:
		targetValue = nil
	}

	// Perform the comparison based on operator
	switch expr.Operator {
	case "=":
		return e.Equals(fieldValue, targetValue)
	case "!=", "<>":
		return !e.Equals(fieldValue, targetValue)
	case ">":
		return e.GreaterThan(fieldValue, targetValue)
	case "<":
		return e.LessThan(fieldValue, targetValue)
	case ">=":
		return e.GreaterThan(fieldValue, targetValue) || e.Equals(fieldValue, targetValue)
	case "<=":
		return e.LessThan(fieldValue, targetValue) || e.Equals(fieldValue, targetValue)
	case "LIKE":
		return e.MatchesLike(fieldValue, targetValue)
	default:
		return false
	}
}

// Equals checks if two values are equal
func (e *BaseExecutor) Equals(a, b interface{}) bool {
	// Handle nil values
	if a == nil && b == nil {
		return true
	}
	if a == nil || b == nil {
		return false
	}

	// Handle string comparison
	aStr, aIsStr := a.(string)
	bStr, bIsStr := b.(string)
	if aIsStr && bIsStr {
		return strings.EqualFold(aStr, bStr)
	}

	// Handle numeric comparison
	return a == b
}

// GreaterThan checks if a > b
func (e *BaseExecutor) GreaterThan(a, b interface{}) bool {
	// Handle nil values
	if a == nil || b == nil {
		return false
	}

	// Convert to float64 if possible
	aFloat, aIsNum := e.toFloat64(a)
	bFloat, bIsNum := e.toFloat64(b)

	if aIsNum && bIsNum {
		return aFloat > bFloat
	}

	// String comparison
	aStr := fmt.Sprintf("%v", a)
	bStr := fmt.Sprintf("%v", b)
	return aStr > bStr
}

// LessThan checks if a < b
func (e *BaseExecutor) LessThan(a, b interface{}) bool {
	// Handle nil values
	if a == nil || b == nil {
		return false
	}

	// Convert to float64 if possible
	aFloat, aIsNum := e.toFloat64(a)
	bFloat, bIsNum := e.toFloat64(b)

	if aIsNum && bIsNum {
		return aFloat < bFloat
	}

	// String comparison
	aStr := fmt.Sprintf("%v", a)
	bStr := fmt.Sprintf("%v", b)
	return aStr < bStr
}

// MatchesLike implements SQL LIKE operator
func (e *BaseExecutor) MatchesLike(a, b interface{}) bool {
	if a == nil || b == nil {
		return false
	}

	aStr := fmt.Sprintf("%v", a)
	pattern := fmt.Sprintf("%v", b)

	// Convert SQL LIKE pattern to Go regex
	// % becomes .*, _ becomes .
	pattern = strings.ReplaceAll(pattern, "%", ".*")
	pattern = strings.ReplaceAll(pattern, "_", ".")

	// Escape regex metacharacters in the non-wildcard parts
	pattern = strings.ReplaceAll(pattern, "[", "\\[")
	pattern = strings.ReplaceAll(pattern, "]", "\\]")
	pattern = strings.ReplaceAll(pattern, "(", "\\(")
	pattern = strings.ReplaceAll(pattern, ")", "\\)")

	matched, err := regexp.MatchString("^"+pattern+"$", aStr)
	if err != nil {
		return false
	}
	return matched
}

// toFloat64 attempts to convert a value to float64
func (e *BaseExecutor) toFloat64(v interface{}) (float64, bool) {
	switch val := v.(type) {
	case int:
		return float64(val), true
	case int8:
		return float64(val), true
	case int16:
		return float64(val), true
	case int32:
		return float64(val), true
	case int64:
		return float64(val), true
	case uint:
		return float64(val), true
	case uint8:
		return float64(val), true
	case uint16:
		return float64(val), true
	case uint32:
		return float64(val), true
	case uint64:
		return float64(val), true
	case float32:
		return float64(val), true
	case float64:
		return val, true
	case string:
		// Try to convert string to number
		if f, err := strconv.ParseFloat(val, 64); err == nil {
			return f, true
		}
	}
	return 0, false
}

// GetSelectedColumns extracts the column names from a SELECT statement
func (e *BaseExecutor) GetSelectedColumns(selectExprs sqlparser.SelectExprs) []string {
	var columns []string

	for _, expr := range selectExprs {
		switch expr := expr.(type) {
		case *sqlparser.AliasedExpr:
			// Handle simple column reference
			if colName, ok := expr.Expr.(*sqlparser.ColName); ok {
				// If there's an alias, use it
				if !expr.As.IsEmpty() {
					columns = append(columns, expr.As.String())
				} else {
					columns = append(columns, colName.Name.String())
				}
			}
		case *sqlparser.StarExpr:
			// Handle * expression (all columns)
			return nil // nil indicates all columns
		}
	}

	return columns
}

// ProjectRow filters a row to only include the specified columns
// If columns is nil, all columns are included (SELECT *)
func (e *BaseExecutor) ProjectRow(row map[string]interface{}, columns []string) map[string]interface{} {
	// If no specific columns are requested or empty list, return the whole row
	if columns == nil || len(columns) == 0 {
		return row
	}

	// Create a new row with only the selected columns
	projectedRow := make(map[string]interface{})

	// Add only the selected columns to the projected row
	for _, col := range columns {
		if value, exists := row[col]; exists {
			projectedRow[col] = value
		} else {
			// Try case-insensitive matching if exact match not found
			for k, v := range row {
				if strings.EqualFold(k, col) {
					projectedRow[col] = v
					break
				}
			}
		}
	}

	return projectedRow
}
