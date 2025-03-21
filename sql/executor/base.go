package executor

import (
	"fmt"
	"regexp"
	"sort"
	"strconv"
	"strings"

	"github.com/blastrain/vitess-sqlparser/sqlparser"
	"github.com/scrymastic/goosquery/sql/context"
	"github.com/scrymastic/goosquery/sql/result"
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
		return e.Compare(fieldValue, targetValue) == 0
	case "!=", "<>":
		return e.Compare(fieldValue, targetValue) != 0
	case ">":
		return e.Compare(fieldValue, targetValue) > 0
	case "<":
		return e.Compare(fieldValue, targetValue) < 0
	case ">=":
		return e.Compare(fieldValue, targetValue) >= 0
	case "<=":
		return e.Compare(fieldValue, targetValue) <= 0
	case "like":
		return e.MatchesLike(fieldValue, targetValue)
	default:
		return false
	}
}

// Compare compares two values, returns:
//
//	1 if a > b
//	0 if a == b
//
// -1 if a < b
func (e *BaseExecutor) Compare(a, b interface{}) int {
	// Handle nil values
	if a == nil && b == nil {
		return 0
	}
	if a == nil {
		return -1
	}
	if b == nil {
		return 1
	}

	// Handle numeric comparison
	aFloat, aIsNum := e.toFloat64(a)
	bFloat, bIsNum := e.toFloat64(b)
	if aIsNum && bIsNum {
		if aFloat < bFloat {
			return -1
		}
		if aFloat > bFloat {
			return 1
		}
		return 0
	}

	// String comparison for non-numeric values
	aStr := fmt.Sprintf("%v", a)
	bStr := fmt.Sprintf("%v", b)
	return strings.Compare(aStr, bStr)
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

// SortResults sorts the query results based on the ORDER BY clause
func (e *BaseExecutor) SortResults(results *result.QueryResult, orderBy sqlparser.OrderBy) error {
	if len(orderBy) == 0 {
		return nil
	}

	// Sort the results based on the ORDER BY clauses
	sort.Slice(*results, func(i, j int) bool {
		// For each ORDER BY clause
		for _, order := range orderBy {
			// Get column name
			colName, ok := order.Expr.(*sqlparser.ColName)
			if !ok {
				continue
			}
			columnName := colName.Name.String()

			// Get values to compare
			rowI := (*results)[i]
			rowJ := (*results)[j]
			valI, existsI := rowI[columnName]
			valJ, existsJ := rowJ[columnName]

			// Handle missing values
			if !existsI && !existsJ {
				continue
			}
			if !existsI {
				return order.Direction == sqlparser.AscScr
			}
			if !existsJ {
				return order.Direction == sqlparser.DescScr
			}

			// Compare values
			compResult := e.Compare(valI, valJ)

			// Apply sort direction
			if compResult != 0 {
				if order.Direction == sqlparser.AscScr {
					return compResult < 0
				}
				return compResult > 0
			}
		}
		return false
	})

	return nil
}

// GetSelectedColumns extracts the column names directly mentioned in the SELECT statement
// Returns an empty slice for SELECT * queries
func (e *BaseExecutor) GetSelectedColumns(selectExprs sqlparser.SelectExprs) []string {
	var columns []string

	for _, expr := range selectExprs {
		switch expr := expr.(type) {
		case *sqlparser.AliasedExpr:
			// Handle simple column reference
			if colName, ok := expr.Expr.(*sqlparser.ColName); ok {
				// If there's an alias, use it for the result but add the column name to the list
				columns = append(columns, colName.Name.String())
			}
		case *sqlparser.StarExpr:
			// Handle * expression (all columns)
			return []string{} // Return empty slice to represent all columns
		}
	}

	return columns
}

// GetWhereColumns extracts column names used in the WHERE clause
func (e *BaseExecutor) GetWhereColumns(whereExpr sqlparser.Expr) []string {
	if whereExpr == nil {
		return []string{}
	}

	columnsMap := make(map[string]bool)
	e.extractWhereColumns(whereExpr, columnsMap)

	// Convert map keys to slice
	columns := make([]string, 0, len(columnsMap))
	for col := range columnsMap {
		columns = append(columns, col)
	}

	return columns
}

// extractWhereColumns recursively extracts column names from WHERE expressions
func (e *BaseExecutor) extractWhereColumns(expr sqlparser.Expr, columnsMap map[string]bool) {
	if expr == nil {
		return
	}

	switch expr := expr.(type) {
	case *sqlparser.ComparisonExpr:
		// Extract column from left side of comparison
		if colName, ok := expr.Left.(*sqlparser.ColName); ok {
			columnsMap[colName.Name.String()] = true
		} else {
			// Left side could be a complex expression, recursively extract
			e.extractWhereColumns(expr.Left, columnsMap)
		}

		// Extract column from right side of comparison (if it's a column)
		if colName, ok := expr.Right.(*sqlparser.ColName); ok {
			columnsMap[colName.Name.String()] = true
		} else {
			// Right side could be a complex expression, recursively extract
			e.extractWhereColumns(expr.Right, columnsMap)
		}

	case *sqlparser.AndExpr:
		e.extractWhereColumns(expr.Left, columnsMap)
		e.extractWhereColumns(expr.Right, columnsMap)

	case *sqlparser.OrExpr:
		e.extractWhereColumns(expr.Left, columnsMap)
		e.extractWhereColumns(expr.Right, columnsMap)

	case *sqlparser.ParenExpr:
		e.extractWhereColumns(expr.Expr, columnsMap)

	case *sqlparser.NotExpr:
		e.extractWhereColumns(expr.Expr, columnsMap)

	case *sqlparser.IsExpr:
		e.extractWhereColumns(expr.Expr, columnsMap)

	case *sqlparser.RangeCond:
		e.extractWhereColumns(expr.Left, columnsMap)
		e.extractWhereColumns(expr.From, columnsMap)
		e.extractWhereColumns(expr.To, columnsMap)

	case *sqlparser.BinaryExpr:
		e.extractWhereColumns(expr.Left, columnsMap)
		e.extractWhereColumns(expr.Right, columnsMap)

	case *sqlparser.FuncExpr:
		// Extract columns from function arguments
		for _, arg := range expr.Exprs {
			if aliasedExpr, ok := arg.(*sqlparser.AliasedExpr); ok {
				e.extractWhereColumns(aliasedExpr.Expr, columnsMap)
			}
		}
	}
}

// GetAggregationColumns extracts column names used in aggregation functions
func (e *BaseExecutor) GetAggregationColumns(selectExprs sqlparser.SelectExprs) []string {
	columnsMap := make(map[string]bool)

	for _, expr := range selectExprs {
		aliasedExpr, ok := expr.(*sqlparser.AliasedExpr)
		if !ok {
			continue
		}

		// Check if this is an aggregation function
		funcExpr, isFuncExpr := aliasedExpr.Expr.(*sqlparser.FuncExpr)
		if !isFuncExpr {
			continue
		}

		funcName := strings.ToUpper(funcExpr.Name.String())
		if funcName != "COUNT" && funcName != "SUM" && funcName != "AVG" &&
			funcName != "MIN" && funcName != "MAX" {
			continue
		}

		// Skip COUNT(*) as it doesn't reference any specific column
		if funcName == "COUNT" && len(funcExpr.Exprs) == 0 {
			continue
		}

		// Extract column names from function arguments
		for _, arg := range funcExpr.Exprs {
			if argAliasedExpr, ok := arg.(*sqlparser.AliasedExpr); ok {
				if colName, ok := argAliasedExpr.Expr.(*sqlparser.ColName); ok {
					columnsMap[colName.Name.String()] = true
				}
			}
		}
	}

	// Convert map keys to slice
	columns := make([]string, 0, len(columnsMap))
	for col := range columnsMap {
		columns = append(columns, col)
	}

	return columns
}

// GetAllRequiredColumns returns a combined list of all columns required for the query:
// - Columns in the SELECT clause
// - Columns used in aggregation functions
// - Columns in the WHERE clause
// - Columns in GROUP BY
// - Columns in ORDER BY
func (e *BaseExecutor) GetAllRequiredColumns(stmt *sqlparser.Select) []string {
	columnsMap := make(map[string]bool)

	// Add columns from SELECT clause
	for _, col := range e.GetSelectedColumns(stmt.SelectExprs) {
		columnsMap[col] = true
	}
	// If e.GetSelectedColumns returns an empty slice, it means we're selecting all columns (*)
	if len(columnsMap) == 0 {
		return []string{}
	}

	// Add columns from aggregation functions
	for _, col := range e.GetAggregationColumns(stmt.SelectExprs) {
		columnsMap[col] = true
	}

	// Add columns from WHERE clause
	if stmt.Where != nil {
		for _, col := range e.GetWhereColumns(stmt.Where.Expr) {
			columnsMap[col] = true
		}
	}

	// Add columns from GROUP BY
	for _, expr := range stmt.GroupBy {
		if colName, ok := expr.(*sqlparser.ColName); ok {
			columnsMap[colName.Name.String()] = true
		}
	}

	// Add columns from ORDER BY
	for _, order := range stmt.OrderBy {
		if colName, ok := order.Expr.(*sqlparser.ColName); ok {
			columnsMap[colName.Name.String()] = true
		}
	}

	// Convert map keys to slice
	columns := make([]string, 0, len(columnsMap))
	for col := range columnsMap {
		columns = append(columns, col)
	}

	// If the result is empty, it likely means we're selecting all columns (*)
	if len(columns) == 0 {
		return []string{}
	}

	return columns
}

// ProjectRow filters a row to only include the specified columns
// If columns is nil or empty, all columns are included (SELECT *)
func (e *BaseExecutor) ProjectRow(row map[string]interface{}, columns []string) map[string]interface{} {
	// If no specific columns are requested or empty list, return the whole row
	if len(columns) == 0 {
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

// GetConstants extracts constants from a WHERE clause expression
func (e *BaseExecutor) GetConstants(expr sqlparser.Expr, ctx *context.Context) {
	if expr == nil {
		return
	}

	switch expr := expr.(type) {
	case *sqlparser.ComparisonExpr:
		if colName, ok := expr.Left.(*sqlparser.ColName); ok {
			fieldName := colName.Name.String()

			if sqlVal, ok := expr.Right.(*sqlparser.SQLVal); ok && expr.Operator == "=" {
				value := string(sqlVal.Val)

				ctx.AddConstant(fieldName, value)
			}
		}
	case *sqlparser.AndExpr:
		e.GetConstants(expr.Left, ctx)
		e.GetConstants(expr.Right, ctx)
	case *sqlparser.OrExpr:
		e.GetConstants(expr.Left, ctx)
		e.GetConstants(expr.Right, ctx)
	case *sqlparser.ParenExpr:
		e.GetConstants(expr.Expr, ctx)
	}
}

// AggregationType represents the type of aggregation function
type AggregationType int

const (
	None AggregationType = iota
	Count
	Sum
	Avg
	Min
	Max
)

// AggregationInfo stores information about an aggregation function
type AggregationInfo struct {
	Type       AggregationType
	Column     string
	Alias      string
	IsDistinct bool
}

// GetAggregationType determines the type of aggregation from an expression
func GetAggregationType(expr sqlparser.Expr) (AggregationType, string, bool) {
	funcExpr, ok := expr.(*sqlparser.FuncExpr)
	if !ok {
		return None, "", false
	}

	funcName := strings.ToUpper(funcExpr.Name.String())
	isDistinct := funcExpr.Distinct

	// Check if any argument exists
	if len(funcExpr.Exprs) == 0 && funcName != "COUNT" {
		return None, "", false
	}

	// Get column name for aggregation
	var columnName string
	if funcName == "COUNT" && len(funcExpr.Exprs) == 0 {
		// COUNT(*) case
		columnName = "*"
	} else if len(funcExpr.Exprs) > 0 {
		// Normal column aggregation case
		aliasedExpr, ok := funcExpr.Exprs[0].(*sqlparser.AliasedExpr)
		if !ok {
			return None, "", false
		}

		colName, ok := aliasedExpr.Expr.(*sqlparser.ColName)
		if !ok {
			return None, "", false
		}
		columnName = colName.Name.String()
	}

	switch funcName {
	case "COUNT":
		return Count, columnName, isDistinct
	case "SUM":
		return Sum, columnName, isDistinct
	case "AVG":
		return Avg, columnName, isDistinct
	case "MIN":
		return Min, columnName, isDistinct
	case "MAX":
		return Max, columnName, isDistinct
	default:
		return None, "", false
	}
}

// GetAggregations extracts aggregation functions from a SELECT statement
func (e *BaseExecutor) GetAggregations(selectExprs sqlparser.SelectExprs) []AggregationInfo {
	var aggregations []AggregationInfo

	for _, expr := range selectExprs {
		aliasedExpr, ok := expr.(*sqlparser.AliasedExpr)
		if !ok {
			continue
		}

		// Check if this is an aggregation function
		aggType, column, isDistinct := GetAggregationType(aliasedExpr.Expr)
		if aggType == None {
			continue
		}

		// Get alias if available, otherwise use function name as alias
		alias := column
		if !aliasedExpr.As.IsEmpty() {
			alias = aliasedExpr.As.String()
		} else {
			// Generate a default alias like COUNT(column)
			switch aggType {
			case Count:
				alias = fmt.Sprintf("COUNT(%s)", column)
			case Sum:
				alias = fmt.Sprintf("SUM(%s)", column)
			case Avg:
				alias = fmt.Sprintf("AVG(%s)", column)
			case Min:
				alias = fmt.Sprintf("MIN(%s)", column)
			case Max:
				alias = fmt.Sprintf("MAX(%s)", column)
			}
		}

		aggregations = append(aggregations, AggregationInfo{
			Type:       aggType,
			Column:     column,
			Alias:      alias,
			IsDistinct: isDistinct,
		})
	}

	return aggregations
}

// HasAggregations checks if the select statement contains any aggregation functions
func (e *BaseExecutor) HasAggregations(selectExprs sqlparser.SelectExprs) bool {
	for _, expr := range selectExprs {
		aliasedExpr, ok := expr.(*sqlparser.AliasedExpr)
		if !ok {
			continue
		}

		aggType, _, _ := GetAggregationType(aliasedExpr.Expr)
		if aggType != None {
			return true
		}
	}
	return false
}
