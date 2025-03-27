package impl

import (
	"github.com/blastrain/vitess-sqlparser/sqlparser"
	"github.com/scrymastic/goosquery/sql/executor/operations"
	"github.com/scrymastic/goosquery/sql/result"
	"github.com/scrymastic/goosquery/sql/sqlctx"
)

// BaseExecutor provides common functionality for all executors
type BaseExecutor struct{}

// MatchesWhereClause checks if a row matches the WHERE clause
func (e *BaseExecutor) MatchesWhereClause(row result.Result, expr sqlparser.Expr) bool {
	switch expr := expr.(type) {
	case *sqlparser.ComparisonExpr:
		return e.EvaluateComparison(row, expr)
	case *sqlparser.AndExpr:
		return e.MatchesWhereClause(row, expr.Left) && e.MatchesWhereClause(row, expr.Right)
	case *sqlparser.OrExpr:
		return e.MatchesWhereClause(row, expr.Left) || e.MatchesWhereClause(row, expr.Right)
	case *sqlparser.NotExpr:
		return !e.MatchesWhereClause(row, expr.Expr)
	case *sqlparser.ParenExpr:
		return e.MatchesWhereClause(row, expr.Expr)
	}
	return false
}

// EvaluateComparison evaluates a comparison expression
func (e *BaseExecutor) EvaluateComparison(row result.Result, expr *sqlparser.ComparisonExpr) bool {
	// Get left operand
	var leftValue interface{}
	if colName, ok := expr.Left.(*sqlparser.ColName); ok {
		columnName := colName.Name.String()
		if val, exists := row[columnName]; exists {
			leftValue = val
		} else {
			return false // Column doesn't exist in this row
		}
	} else {
		// Try to extract a literal value
		leftValue = operations.ExtractLiteralValue(expr.Left)
	}

	// Get right operand
	var rightValue interface{}
	if colName, ok := expr.Right.(*sqlparser.ColName); ok {
		columnName := colName.Name.String()
		if val, exists := row[columnName]; exists {
			rightValue = val
		} else {
			return false // Column doesn't exist in this row
		}
	} else {
		// Try to extract a literal value
		rightValue = operations.ExtractLiteralValue(expr.Right)
	}

	// Handle NULL values
	if leftValue == nil || rightValue == nil {
		// NULL comparison semantics
		if expr.Operator == "is" {
			return rightValue == nil
		}
		if expr.Operator == "is not" {
			return rightValue != nil
		}
		return false // Any comparison with NULL yields false
	}

	// Perform the comparison
	switch expr.Operator {
	case "=":
		return operations.Compare(leftValue, rightValue) == 0
	case "!=", "<>":
		return operations.Compare(leftValue, rightValue) != 0
	case "<":
		return operations.Compare(leftValue, rightValue) < 0
	case "<=":
		return operations.Compare(leftValue, rightValue) <= 0
	case ">":
		return operations.Compare(leftValue, rightValue) > 0
	case ">=":
		return operations.Compare(leftValue, rightValue) >= 0
	case "like":
		return operations.MatchesLike(leftValue, rightValue)
	case "not like":
		return !operations.MatchesLike(leftValue, rightValue)
	case "in":
		// TODO: Implement IN comparison
		return false
	case "not in":
		// TODO: Implement NOT IN comparison
		return false
	}

	return false
}

// GetAllRequiredColumns returns all columns required for the query
func (e *BaseExecutor) GetAllRequiredColumns(stmt *sqlparser.Select) []string {
	columnsMap := make(map[string]bool)

	// Add columns from SELECT clause
	for _, col := range e.GetSelectedColumns(stmt.SelectExprs) {
		if col != "*" {
			columnsMap[col] = true
		} else {
			return []string{"*"}
		}
	}

	// Add columns from WHERE clause
	if stmt.Where != nil {
		for _, col := range e.GetWhereColumns(stmt.Where.Expr) {
			columnsMap[col] = true
		}
	}

	// Add columns from ORDER BY clause
	for _, order := range stmt.OrderBy {
		switch expr := order.Expr.(type) {
		case *sqlparser.ColName:
			columnsMap[expr.Name.String()] = true
		}
	}

	// Add columns from GROUP BY clause
	for _, expr := range stmt.GroupBy {
		switch expr := expr.(type) {
		case *sqlparser.ColName:
			columnsMap[expr.Name.String()] = true
		}
	}

	// Convert map to slice
	result := make([]string, 0, len(columnsMap))
	for col := range columnsMap {
		result = append(result, col)
	}

	return result
}

// GetSelectedColumns extracts column names from SELECT expressions
func (e *BaseExecutor) GetSelectedColumns(selectExprs sqlparser.SelectExprs) []string {
	columns := []string{}
	for _, expr := range selectExprs {
		switch expr := expr.(type) {
		case *sqlparser.AliasedExpr:
			switch exprType := expr.Expr.(type) {
			case *sqlparser.ColName:
				columns = append(columns, exprType.Name.String())
			case *sqlparser.FuncExpr:
				// Add columns used in function arguments
				columns = append(columns, e.GetAggregationColumns(selectExprs)...)
			}
		case *sqlparser.StarExpr:
			// SELECT * means all columns
			return []string{"*"}
		}
	}
	return columns
}

// GetWhereColumns extracts column names from WHERE expression
func (e *BaseExecutor) GetWhereColumns(whereExpr sqlparser.Expr) []string {
	columnsMap := make(map[string]bool)
	e.extractWhereColumns(whereExpr, columnsMap)

	columns := make([]string, 0, len(columnsMap))
	for col := range columnsMap {
		columns = append(columns, col)
	}
	return columns
}

// extractWhereColumns recursively extracts column names from an expression
func (e *BaseExecutor) extractWhereColumns(expr sqlparser.Expr, columnsMap map[string]bool) {
	switch expr := expr.(type) {
	case *sqlparser.ComparisonExpr:
		// Handle left side
		if colName, ok := expr.Left.(*sqlparser.ColName); ok {
			columnsMap[colName.Name.String()] = true
		} else {
			e.extractWhereColumns(expr.Left, columnsMap)
		}

		// Handle right side
		if colName, ok := expr.Right.(*sqlparser.ColName); ok {
			columnsMap[colName.Name.String()] = true
		} else {
			e.extractWhereColumns(expr.Right, columnsMap)
		}

	case *sqlparser.AndExpr:
		e.extractWhereColumns(expr.Left, columnsMap)
		e.extractWhereColumns(expr.Right, columnsMap)

	case *sqlparser.OrExpr:
		e.extractWhereColumns(expr.Left, columnsMap)
		e.extractWhereColumns(expr.Right, columnsMap)

	case *sqlparser.NotExpr:
		e.extractWhereColumns(expr.Expr, columnsMap)

	case *sqlparser.ParenExpr:
		e.extractWhereColumns(expr.Expr, columnsMap)

	case *sqlparser.FuncExpr:
		for _, arg := range expr.Exprs {
			if aliasedExpr, ok := arg.(*sqlparser.AliasedExpr); ok {
				e.extractWhereColumns(aliasedExpr.Expr, columnsMap)
			}
		}
	}
}

// GetAggregationColumns extracts column names used in aggregation functions
func (e *BaseExecutor) GetAggregationColumns(selectExprs sqlparser.SelectExprs) []string {
	columns := []string{}
	for _, expr := range selectExprs {
		aliasedExpr, ok := expr.(*sqlparser.AliasedExpr)
		if !ok {
			continue
		}

		funcExpr, ok := aliasedExpr.Expr.(*sqlparser.FuncExpr)
		if !ok {
			continue
		}

		// Extract column from first arg of aggregation function
		if len(funcExpr.Exprs) > 0 {
			argExpr, ok := funcExpr.Exprs[0].(*sqlparser.AliasedExpr)
			if !ok {
				continue
			}

			colName, ok := argExpr.Expr.(*sqlparser.ColName)
			if ok {
				columns = append(columns, colName.Name.String())
			}
		}
	}
	return columns
}

// GetConstants extracts constants from WHERE expressions and adds them to the context
func (e *BaseExecutor) GetConstants(expr sqlparser.Expr, ctx *sqlctx.Context) {
	switch expr := expr.(type) {
	case *sqlparser.ComparisonExpr:
		if colName, ok := expr.Left.(*sqlparser.ColName); ok {
			if sqlVal, ok := expr.Right.(*sqlparser.SQLVal); ok {
				strVal := string(sqlVal.Val)
				// Store column name and value in context (without operator)
				ctx.AddConstant(colName.Name.String(), strVal)
			}
		} else if colName, ok := expr.Right.(*sqlparser.ColName); ok {
			if sqlVal, ok := expr.Left.(*sqlparser.SQLVal); ok {
				strVal := string(sqlVal.Val)
				// Store column name and value in context (without operator)
				// Note: original code was inverting operator but AddConstant doesn't take an operator
				ctx.AddConstant(colName.Name.String(), strVal)
			}
		}
	case *sqlparser.AndExpr:
		e.GetConstants(expr.Left, ctx)
		e.GetConstants(expr.Right, ctx)
	case *sqlparser.OrExpr:
		// OR expressions can't be used for optimization in this simple implementation
		// but you could handle them more sophisticatedly
	case *sqlparser.ParenExpr:
		e.GetConstants(expr.Expr, ctx)
	}
}
