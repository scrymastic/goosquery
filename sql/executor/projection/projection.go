package projection

import (
	"fmt"
	"strings"

	"github.com/blastrain/vitess-sqlparser/sqlparser"
	"github.com/scrymastic/goosquery/sql/result"
)

// ProjectFinalResults applies final projection to the result to ensure only the requested columns are returned
// It handles column selection, aliases from SELECT, and aggregation function results
func ProjectFinalResults(results *result.QueryResult, stmt *sqlparser.Select) *result.QueryResult {
	// If there are no results, return empty result
	if len(*results) == 0 {
		return result.NewQueryResult()
	}

	// If it's a SELECT * with no aliases, return the results as is
	if isSelectStar(stmt.SelectExprs) && !hasAliases(stmt.SelectExprs) {
		return results
	}

	// Create a new result with only the requested columns, including aliases
	projectedResult := result.NewQueryResult()

	// Build a map of original column names to their aliases
	aliasMap := buildAliasMap(stmt.SelectExprs)

	// Get the list of columns to include in the final result
	finalColumns := getFinalProjectionColumns(stmt.SelectExprs)

	// Apply projection for each row
	for _, row := range *results {
		projectedRow := make(map[string]interface{})

		// Handle columns with aliases
		for originalCol, alias := range aliasMap {
			if value, exists := row[originalCol]; exists {
				projectedRow[alias] = value
			}
		}

		// Handle direct columns (no aliases)
		for _, col := range finalColumns {
			if _, isAliased := aliasMap[col]; isAliased {
				// Already handled above
				continue
			}

			// Not an alias, use the column directly
			if value, exists := row[col]; exists {
				projectedRow[col] = value
			}
		}

		projectedResult.AddRecord(projectedRow)
	}

	return projectedResult
}

// buildAliasMap creates a map of original column names to their aliases
func buildAliasMap(selectExprs sqlparser.SelectExprs) map[string]string {
	aliasMap := make(map[string]string)

	for _, expr := range selectExprs {
		aliasedExpr, ok := expr.(*sqlparser.AliasedExpr)
		if !ok || aliasedExpr.As.IsEmpty() {
			continue
		}

		alias := aliasedExpr.As.String()

		// Handle direct column reference with alias
		if colName, ok := aliasedExpr.Expr.(*sqlparser.ColName); ok {
			aliasMap[colName.Name.String()] = alias
			continue
		}

		// Handle aggregation function with alias
		if funcExpr, ok := aliasedExpr.Expr.(*sqlparser.FuncExpr); ok {
			funcName := strings.ToUpper(funcExpr.Name.String())

			// Special case for COUNT(*)
			if funcName == "COUNT" && len(funcExpr.Exprs) == 0 {
				// COUNT(*) is stored using the function name as the original column
				aliasMap["COUNT(*)"] = alias
				continue
			}

			// For other aggregation functions with a single column argument
			if len(funcExpr.Exprs) > 0 {
				if argAliasedExpr, ok := funcExpr.Exprs[0].(*sqlparser.AliasedExpr); ok {
					if colName, ok := argAliasedExpr.Expr.(*sqlparser.ColName); ok {
						// Create aggregate function name as key
						originalCol := fmt.Sprintf("%s(%s)", funcName, colName.Name.String())
						aliasMap[originalCol] = alias
					}
				}
			}
		}
	}

	return aliasMap
}

// getFinalProjectionColumns gets the list of columns to include in the final result
func getFinalProjectionColumns(selectExprs sqlparser.SelectExprs) []string {
	var columns []string

	for _, expr := range selectExprs {
		switch expr := expr.(type) {
		case *sqlparser.AliasedExpr:
			// For direct column reference
			if colName, ok := expr.Expr.(*sqlparser.ColName); ok {
				// If there's no alias, use the column name
				if expr.As.IsEmpty() {
					columns = append(columns, colName.Name.String())
				}
				// If there is an alias, it will be handled by the alias map
			} else if funcExpr, ok := expr.Expr.(*sqlparser.FuncExpr); ok {
				// For aggregation function
				funcName := strings.ToUpper(funcExpr.Name.String())

				// Special case for COUNT(*)
				if funcName == "COUNT" && len(funcExpr.Exprs) == 0 {
					if expr.As.IsEmpty() {
						columns = append(columns, "COUNT(*)")
					}
					continue
				}

				// For other aggregation functions
				if len(funcExpr.Exprs) > 0 && expr.As.IsEmpty() {
					if argAliasedExpr, ok := funcExpr.Exprs[0].(*sqlparser.AliasedExpr); ok {
						if colName, ok := argAliasedExpr.Expr.(*sqlparser.ColName); ok {
							// If no alias, use function name as column name
							originalCol := fmt.Sprintf("%s(%s)", funcName, colName.Name.String())
							columns = append(columns, originalCol)
						}
					}
				}
			}
		case *sqlparser.StarExpr:
			// For SELECT *, return empty to signal all columns
			return []string{}
		}
	}

	return columns
}

// isSelectStar checks if the query is a SELECT * query
func isSelectStar(selectExprs sqlparser.SelectExprs) bool {
	for _, expr := range selectExprs {
		if _, ok := expr.(*sqlparser.StarExpr); ok {
			return true
		}
	}
	return false
}

// hasAliases checks if any expressions have aliases
func hasAliases(selectExprs sqlparser.SelectExprs) bool {
	for _, expr := range selectExprs {
		aliasedExpr, ok := expr.(*sqlparser.AliasedExpr)
		if ok && !aliasedExpr.As.IsEmpty() {
			return true
		}
	}
	return false
}
