package operations

import (
	"fmt"
	"regexp"
	"sort"
	"strconv"
	"strings"

	"github.com/blastrain/vitess-sqlparser/sqlparser"
	"github.com/scrymastic/goosquery/sql/result"
)

// Compare compares two values and returns -1, 0, or 1
func Compare(a, b interface{}) int {
	// Handle special cases
	if a == nil && b == nil {
		return 0
	}
	if a == nil {
		return -1
	}
	if b == nil {
		return 1
	}

	// Try to convert both to float64 for numeric comparison
	aFloat, aIsNum := ToFloat64(a)
	bFloat, bIsNum := ToFloat64(b)

	if aIsNum && bIsNum {
		// Both are numeric, compare as numbers
		if aFloat < bFloat {
			return -1
		}
		if aFloat > bFloat {
			return 1
		}
		return 0
	}

	// Convert both to string and compare lexically
	aStr := fmt.Sprintf("%v", a)
	bStr := fmt.Sprintf("%v", b)

	// Additional attempt for numeric string comparison
	// This helps with sorting columns that store numbers as strings
	aFloatFromStr, aIsNumStr := ToFloat64(aStr)
	bFloatFromStr, bIsNumStr := ToFloat64(bStr)

	if aIsNumStr && bIsNumStr {
		// Both strings can be parsed as numbers, compare numerically
		if aFloatFromStr < bFloatFromStr {
			return -1
		}
		if aFloatFromStr > bFloatFromStr {
			return 1
		}
		return 0
	}

	// Fallback to lexical comparison
	if aStr < bStr {
		return -1
	}
	if aStr > bStr {
		return 1
	}
	return 0
}

// ToFloat64 converts a value to float64 if possible
func ToFloat64(v interface{}) (float64, bool) {
	switch vt := v.(type) {
	case int:
		return float64(vt), true
	case int8:
		return float64(vt), true
	case int16:
		return float64(vt), true
	case int32:
		return float64(vt), true
	case int64:
		return float64(vt), true
	case uint:
		return float64(vt), true
	case uint8:
		return float64(vt), true
	case uint16:
		return float64(vt), true
	case uint32:
		return float64(vt), true
	case uint64:
		return float64(vt), true
	case float32:
		return float64(vt), true
	case float64:
		return vt, true
	case string:
		// Attempt to parse as float
		f, err := strconv.ParseFloat(vt, 64)
		if err == nil {
			return f, true
		}
	}
	return 0, false
}

// SortResults sorts the results based on ORDER BY clause
func SortResults(results *result.QueryResult, orderBy sqlparser.OrderBy) error {
	if len(orderBy) == 0 || len(*results) == 0 {
		return nil
	}

	// Create a copy of the results for sorting
	sortedResults := append(result.QueryResult{}, *results...)

	// Sort the results
	sort.SliceStable(sortedResults, func(i, j int) bool {
		return CompareRows(sortedResults[i], sortedResults[j], orderBy)
	})

	// Update the original results
	*results = sortedResults
	return nil
}

// CompareRows compares two rows based on ORDER BY expressions
func CompareRows(a, b map[string]interface{}, orderBy sqlparser.OrderBy) bool {
	for _, order := range orderBy {
		var aVal, bVal interface{}
		var exists bool

		// Extract column values
		switch expr := order.Expr.(type) {
		case *sqlparser.ColName:
			colName := expr.Name.String()
			aVal, exists = a[colName]
			if !exists {
				aVal = nil
			}
			bVal, exists = b[colName]
			if !exists {
				bVal = nil
			}
		default:
			// Unsupported expression type
			continue
		}

		// Compare the values
		cmp := Compare(aVal, bVal)
		if cmp != 0 {
			// Return result based on sort direction
			// In vitess-sqlparser, order.Direction is a string representation
			direction := string(order.Direction)
			isAsc := direction == "" || strings.ToLower(direction) == "asc"
			return (isAsc && cmp < 0) || (!isAsc && cmp > 0)
		}
		// If equal, continue to next ORDER BY expression
	}
	return false
}

// CalculateMin finds the minimum value in a set
func CalculateMin(values []interface{}) (interface{}, error) {
	if len(values) == 0 {
		return nil, nil
	}

	minVal := values[0]

	for _, val := range values[1:] {
		if Compare(val, minVal) < 0 {
			minVal = val
		}
	}

	return minVal, nil
}

// CalculateMax finds the maximum value in a set
func CalculateMax(values []interface{}) (interface{}, error) {
	if len(values) == 0 {
		return nil, nil
	}

	maxVal := values[0]

	for _, val := range values[1:] {
		if Compare(val, maxVal) > 0 {
			maxVal = val
		}
	}

	return maxVal, nil
}

// CalculateSum calculates the sum of a set of values
func CalculateSum(values []interface{}) (interface{}, error) {
	var sum float64

	for _, val := range values {
		// Convert to float64 if possible
		numVal, isNum := ToFloat64(val)
		if !isNum {
			return nil, fmt.Errorf("cannot sum non-numeric value: %v", val)
		}
		sum += numVal
	}

	return sum, nil
}

// ExtractLiteralValue extracts a literal value from a SQL expression
func ExtractLiteralValue(expr sqlparser.Expr) interface{} {
	if sqlVal, ok := expr.(*sqlparser.SQLVal); ok {
		switch sqlVal.Type {
		case sqlparser.IntVal:
			intVal := string(sqlVal.Val)
			// Convert to integer for proper numeric comparison
			if val, err := strconv.Atoi(intVal); err == nil {
				return val
			}
			return intVal
		case sqlparser.FloatVal:
			floatVal := string(sqlVal.Val)
			// Convert to float for proper numeric comparison
			if val, err := strconv.ParseFloat(floatVal, 64); err == nil {
				return val
			}
			return floatVal
		case sqlparser.StrVal:
			strVal := string(sqlVal.Val)
			return strVal
		}
	}
	return nil
}

// MatchesLike checks if a value matches a LIKE pattern
func MatchesLike(a, b interface{}) bool {
	// Convert both to strings
	aStr := fmt.Sprintf("%v", a)
	bStr := fmt.Sprintf("%v", b)

	// Escape special regex characters in the pattern
	pattern := strings.Replace(strings.Replace(bStr, "%", ".*", -1), "_", ".", -1)
	pattern = "^" + pattern + "$" // Ensure full string match

	// Only use case-insensitive matching if the compiler supports it
	reStr := "(?i)" + pattern

	// Check if the string matches the pattern
	matched, err := regexp.MatchString(reStr, aStr)
	if err != nil {
		return false
	}
	return matched
}
