package filtering

import (
	"fmt"
	"strings"

	"github.com/blastrain/vitess-sqlparser/sqlparser"
	"github.com/scrymastic/goosquery/sql/result"
)

// Filtering contains methods for filtering query results
type Filtering struct {
	// Function types for utility methods (dependency injection)
	ToFloat64          func(val interface{}) (float64, error)
	Compare            func(a, b interface{}) int
	MatchesLike        func(val interface{}, pattern string) bool
	CaseInsensitiveGet func(m map[string]interface{}, key string) (interface{}, bool)
}

// ApplyWhere filters query results based on the WHERE clause
func (f *Filtering) ApplyWhere(results *result.QueryResult, where *sqlparser.Where) (*result.QueryResult, error) {
	if where == nil {
		return results, nil
	}

	// Create a new result to hold filtered records
	filtered := result.NewQueryResult()

	// Filter each record
	for _, record := range *results {
		matches, err := f.evaluateExpr(record, where.Expr)
		if err != nil {
			return nil, fmt.Errorf("error evaluating WHERE: %w", err)
		}

		if matches {
			filtered.AddRecord(record)
		}
	}

	return filtered, nil
}

// evaluateExpr evaluates a SQL expression against a record
func (f *Filtering) evaluateExpr(record map[string]interface{}, expr sqlparser.Expr) (bool, error) {
	switch expr := expr.(type) {
	case *sqlparser.ComparisonExpr:
		return f.evaluateComparison(record, expr)
	case *sqlparser.AndExpr:
		return f.evaluateAnd(record, expr)
	case *sqlparser.OrExpr:
		return f.evaluateOr(record, expr)
	case *sqlparser.NotExpr:
		return f.evaluateNot(record, expr)
	case *sqlparser.ParenExpr:
		return f.evaluateExpr(record, expr.Expr)
	case *sqlparser.IsExpr:
		return f.evaluateIsExpr(record, expr)
	case *sqlparser.RangeCond:
		return f.evaluateRange(record, expr)
	case *sqlparser.ExistsExpr:
		// Not supporting subqueries in this implementation
		return false, fmt.Errorf("EXISTS expressions not supported")
	case *sqlparser.SQLVal:
		// A direct value is typically a boolean TRUE or FALSE
		if expr.Type == sqlparser.IntVal {
			val := string(expr.Val)
			return val != "0", nil
		}
		return false, fmt.Errorf("unsupported direct value in WHERE: %v", expr)
	}

	return false, fmt.Errorf("unsupported expression in WHERE: %T", expr)
}

// evaluateComparison evaluates a comparison expression
func (f *Filtering) evaluateComparison(record map[string]interface{}, expr *sqlparser.ComparisonExpr) (bool, error) {
	leftVal, err := f.getExprValue(record, expr.Left)
	if err != nil {
		return false, err
	}

	// Handle NULL values in comparisons
	if leftVal == nil && expr.Operator != "is" && expr.Operator != "is not" {
		return false, nil
	}

	// Special case for LIKE operator
	if expr.Operator == sqlparser.LikeStr {
		rightVal, err := f.getExprValue(record, expr.Right)
		if err != nil {
			return false, err
		}

		if rightStr, ok := rightVal.(string); ok {
			return f.MatchesLike(leftVal, rightStr), nil
		}
		return false, fmt.Errorf("right side of LIKE must be a string")
	}

	// Handle IN operator
	if expr.Operator == sqlparser.InStr {
		return f.evaluateIn(leftVal, expr.Right)
	}

	// Handle NOT IN operator
	if expr.Operator == sqlparser.NotInStr {
		result, err := f.evaluateIn(leftVal, expr.Right)
		return !result, err
	}

	// For regular comparisons
	rightVal, err := f.getExprValue(record, expr.Right)
	if err != nil {
		return false, err
	}

	// Handle NULL values in right side
	if rightVal == nil && expr.Operator != "is" && expr.Operator != "is not" {
		return false, nil
	}

	return f.compareValues(leftVal, rightVal, expr.Operator)
}

// evaluateAnd evaluates an AND expression
func (f *Filtering) evaluateAnd(record map[string]interface{}, expr *sqlparser.AndExpr) (bool, error) {
	leftResult, err := f.evaluateExpr(record, expr.Left)
	if err != nil {
		return false, err
	}

	// Short-circuit evaluation
	if !leftResult {
		return false, nil
	}

	return f.evaluateExpr(record, expr.Right)
}

// evaluateOr evaluates an OR expression
func (f *Filtering) evaluateOr(record map[string]interface{}, expr *sqlparser.OrExpr) (bool, error) {
	leftResult, err := f.evaluateExpr(record, expr.Left)
	if err != nil {
		return false, err
	}

	// Short-circuit evaluation
	if leftResult {
		return true, nil
	}

	return f.evaluateExpr(record, expr.Right)
}

// evaluateNot evaluates a NOT expression
func (f *Filtering) evaluateNot(record map[string]interface{}, expr *sqlparser.NotExpr) (bool, error) {
	result, err := f.evaluateExpr(record, expr.Expr)
	if err != nil {
		return false, err
	}

	return !result, nil
}

// evaluateNullCheck evaluates a NULL check (IS NULL, IS NOT NULL)
func (f *Filtering) evaluateNullCheck(record map[string]interface{}, expr *sqlparser.ComparisonExpr) (bool, error) {
	val, err := f.getExprValue(record, expr.Left)
	if err != nil {
		return false, err
	}

	isNull := val == nil
	return (expr.Operator == "is") == isNull, nil
}

// evaluateIsExpr evaluates an IS expression (IS TRUE, IS FALSE)
func (f *Filtering) evaluateIsExpr(record map[string]interface{}, expr *sqlparser.IsExpr) (bool, error) {
	leftVal, err := f.getExprValue(record, expr.Expr)
	if err != nil {
		return false, err
	}

	// The IsExpr just checks if expr is null, true, or false
	var result bool

	// Check if the value is null, true, or false
	if leftVal == nil {
		result = expr.Operator == "is null"
	} else if boolVal, ok := leftVal.(bool); ok {
		if boolVal {
			result = expr.Operator == "is true"
		} else {
			result = expr.Operator == "is false"
		}
	} else {
		// Try converting to a boolean
		result = false // default if can't convert
	}

	// Apply NOT if needed
	return result, nil
}

// evaluateRange evaluates a BETWEEN or NOT BETWEEN expression
func (f *Filtering) evaluateRange(record map[string]interface{}, expr *sqlparser.RangeCond) (bool, error) {
	val, err := f.getExprValue(record, expr.Left)
	if err != nil {
		return false, err
	}

	from, err := f.getExprValue(record, expr.From)
	if err != nil {
		return false, err
	}

	to, err := f.getExprValue(record, expr.To)
	if err != nil {
		return false, err
	}

	// Check if value is between from and to
	inRange := f.Compare(val, from) >= 0 && f.Compare(val, to) <= 0

	if expr.Operator == sqlparser.BetweenStr {
		return inRange, nil
	}
	return !inRange, nil
}

// evaluateIn evaluates an IN expression
func (f *Filtering) evaluateIn(leftVal interface{}, right sqlparser.Expr) (bool, error) {
	// Handle IN clause with a tuple
	if list, ok := right.(sqlparser.ValTuple); ok {
		for _, item := range list {
			// Convert each value in the list
			val, err := f.getConstValue(item)
			if err != nil {
				return false, err
			}

			// Check if leftVal matches any value in the list
			if f.Compare(leftVal, val) == 0 {
				return true, nil
			}
		}
		return false, nil
	}

	return false, fmt.Errorf("unsupported IN expression: %T", right)
}

// getExprValue gets the value of an expression in the context of a record
func (f *Filtering) getExprValue(record map[string]interface{}, expr sqlparser.Expr) (interface{}, error) {
	switch expr := expr.(type) {
	case *sqlparser.ColName:
		columnName := strings.ToLower(expr.Name.String())

		// Check case insensitive for the column
		val, exists := f.CaseInsensitiveGet(record, columnName)
		if !exists {
			return nil, nil // Missing column treated as NULL
		}
		return val, nil

	case sqlparser.ValTuple:
		// For tuple expressions like (1, 2, 3)
		values := make([]interface{}, len(expr))
		for i, val := range expr {
			v, err := f.getExprValue(record, val)
			if err != nil {
				return nil, err
			}
			values[i] = v
		}
		return values, nil

	case *sqlparser.SQLVal:
		return f.getConstValue(expr)

	case *sqlparser.NullVal:
		return nil, nil

	case *sqlparser.FuncExpr:
		// Not supporting functions in WHERE clause in this implementation
		return nil, fmt.Errorf("functions in WHERE clause not supported")

	case *sqlparser.Subquery:
		// Not supporting subqueries in this implementation
		return nil, fmt.Errorf("subqueries not supported")

	case *sqlparser.BinaryExpr:
		return f.evaluateBinaryExpr(record, expr)
	}

	return nil, fmt.Errorf("unsupported expression type in WHERE: %T", expr)
}

// getConstValue converts a SQL value to its Go equivalent
func (f *Filtering) getConstValue(expr sqlparser.Expr) (interface{}, error) {
	sqlVal, ok := expr.(*sqlparser.SQLVal)
	if !ok {
		return nil, fmt.Errorf("expected SQL value, got %T", expr)
	}

	switch sqlVal.Type {
	case sqlparser.StrVal:
		return string(sqlVal.Val), nil
	case sqlparser.IntVal:
		// First try parsing as an integer
		return string(sqlVal.Val), nil
	case sqlparser.FloatVal:
		return string(sqlVal.Val), nil
	case sqlparser.HexNum, sqlparser.HexVal:
		return nil, fmt.Errorf("hex values not supported")
	}

	return nil, fmt.Errorf("unsupported SQL value type: %v", sqlVal.Type)
}

// evaluateBinaryExpr evaluates a binary expression (e.g., a + b, a - b)
func (f *Filtering) evaluateBinaryExpr(record map[string]interface{}, expr *sqlparser.BinaryExpr) (interface{}, error) {
	leftVal, err := f.getExprValue(record, expr.Left)
	if err != nil {
		return nil, err
	}

	rightVal, err := f.getExprValue(record, expr.Right)
	if err != nil {
		return nil, err
	}

	// If either value is nil, the result is nil
	if leftVal == nil || rightVal == nil {
		return nil, nil
	}

	// Convert to float64 for arithmetic operations
	leftFloat, err := f.ToFloat64(leftVal)
	if err != nil {
		return nil, err
	}

	rightFloat, err := f.ToFloat64(rightVal)
	if err != nil {
		return nil, err
	}

	switch expr.Operator {
	case "+":
		return leftFloat + rightFloat, nil
	case "-":
		return leftFloat - rightFloat, nil
	case "*":
		return leftFloat * rightFloat, nil
	case "/":
		if rightFloat == 0 {
			return nil, fmt.Errorf("division by zero")
		}
		return leftFloat / rightFloat, nil
	}

	return nil, fmt.Errorf("unsupported binary operator: %s", expr.Operator)
}

// compareValues compares two values based on the comparison operator
func (f *Filtering) compareValues(left, right interface{}, operator string) (bool, error) {
	// Special handling for NULL values
	if left == nil || right == nil {
		switch operator {
		case "is":
			return left == right, nil
		case "is not":
			return left != right, nil
		default:
			return false, nil
		}
	}

	// Get comparison result
	compResult := f.Compare(left, right)

	// Apply comparison based on operator
	switch operator {
	case "=":
		return compResult == 0, nil
	case "!=", "<>":
		return compResult != 0, nil
	case "<":
		return compResult < 0, nil
	case "<=":
		return compResult <= 0, nil
	case ">":
		return compResult > 0, nil
	case ">=":
		return compResult >= 0, nil
	}

	return false, fmt.Errorf("unsupported comparison operator: %s", operator)
}
