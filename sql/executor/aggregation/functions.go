package aggregation

import (
	"fmt"
	"strings"

	"github.com/blastrain/vitess-sqlparser/sqlparser"
	"github.com/scrymastic/goosquery/sql/executor/operations"
	"github.com/scrymastic/goosquery/sql/result"
)

// ApplyAggregations applies aggregation functions to the result set
func ApplyAggregations(results *result.QueryResult, aggregations []AggregationInfo, groupBy sqlparser.GroupBy) (*result.QueryResult, error) {
	// If there are no results, return empty result with aggregation columns
	if len(*results) == 0 {
		emptyResult := result.NewQueryResult()
		emptyRow := make(map[string]interface{})

		// Add zero values for all aggregations
		for _, agg := range aggregations {
			switch agg.Type {
			case Count:
				emptyRow[agg.Alias] = 0
			case Sum, Avg:
				emptyRow[agg.Alias] = 0.0
			case Min, Max:
				emptyRow[agg.Alias] = nil
			}
		}

		emptyResult.AddRecord(emptyRow)
		return emptyResult, nil
	}

	// Check if GROUP BY is used
	hasGroupBy := len(groupBy) > 0

	// Get the group by columns if provided
	groupByColumns := []string{}
	for _, expr := range groupBy {
		colName, ok := expr.(*sqlparser.ColName)
		if !ok {
			continue
		}
		groupByColumns = append(groupByColumns, colName.Name.String())
	}

	// If no GROUP BY, apply aggregations to the entire result set
	if !hasGroupBy {
		return aggregateAll(results, aggregations)
	}

	// Otherwise, group the results and apply aggregations to each group
	return aggregateByGroups(results, aggregations, groupByColumns)
}

// aggregateAll applies aggregations to the entire result set (no GROUP BY)
func aggregateAll(results *result.QueryResult, aggregations []AggregationInfo) (*result.QueryResult, error) {
	// Create a single aggregated row
	aggregatedRow := make(map[string]interface{})

	// Apply each aggregation function
	for _, agg := range aggregations {
		value, err := calculateAggregation(results, agg)
		if err != nil {
			return nil, err
		}
		aggregatedRow[agg.Alias] = value
	}

	// Create a new result with just the aggregated row
	aggregatedResult := result.NewQueryResult()
	aggregatedResult.AddRecord(aggregatedRow)

	return aggregatedResult, nil
}

// aggregateByGroups applies aggregations to each group of results
func aggregateByGroups(results *result.QueryResult, aggregations []AggregationInfo, groupByColumns []string) (*result.QueryResult, error) {
	// Check if group columns exist in the results
	if len(*results) > 0 {
		firstRow := (*results)[0]
		for _, col := range groupByColumns {
			if _, exists := firstRow[col]; !exists {
				return nil, fmt.Errorf("group by column %s not found in result set", col)
			}
		}
	}

	// Map to store groups: groupKey -> []rowIndex
	groups := make(map[string][]int)

	// Group the rows
	for i, row := range *results {
		// Create a key for this group (combination of values of group by columns)
		var keyParts []string
		for _, col := range groupByColumns {
			keyParts = append(keyParts, fmt.Sprintf("%v", row[col]))
		}
		// Join the key parts with a colon, if there are multiple group by columns
		groupKey := strings.Join(keyParts, ":")

		// Add this row's index to the appropriate group
		groups[groupKey] = append(groups[groupKey], i)
	}

	// Create a new result with one row per group
	aggregatedResult := result.NewQueryResult()

	// Process each group
	for _, rowIndices := range groups {
		// Create a subset of results containing only rows in this group
		groupResults := result.NewQueryResult()
		for _, idx := range rowIndices {
			groupResults.AddRecord((*results)[idx])
		}

		// Create a row with group by columns and aggregated values
		aggregatedRow := make(map[string]interface{})

		// First, add the group by columns
		for _, col := range groupByColumns {
			aggregatedRow[col] = (*results)[rowIndices[0]][col]
		}

		// Then add the aggregated values
		for _, agg := range aggregations {
			value, err := calculateAggregation(groupResults, agg)
			if err != nil {
				return nil, err
			}
			aggregatedRow[agg.Alias] = value
		}

		aggregatedResult.AddRecord(aggregatedRow)
	}

	return aggregatedResult, nil
}

// calculateAggregation applies a single aggregation function to a set of results
func calculateAggregation(results *result.QueryResult, agg AggregationInfo) (interface{}, error) {
	if len(*results) == 0 {
		return nil, nil
	}

	// Handle COUNT(*) as a special case
	if agg.Type == Count && agg.Column == "*" {
		return len(*results), nil
	}

	// For other aggregations, collect the values to aggregate
	var values []interface{}

	// If distinct, collect unique values
	if agg.IsDistinct {
		uniqueValues := make(map[string]interface{})
		for _, row := range *results {
			if val, exists := row[agg.Column]; exists && val != nil {
				// Use string representation as map key
				uniqueValues[fmt.Sprintf("%v", val)] = val
			}
		}

		for _, val := range uniqueValues {
			values = append(values, val)
		}
	} else {
		// Otherwise collect all values
		for _, row := range *results {
			if val, exists := row[agg.Column]; exists && val != nil {
				values = append(values, val)
			}
		}
	}

	// If no values to aggregate, return appropriate default
	if len(values) == 0 {
		switch agg.Type {
		case Count:
			return 0, nil
		case Sum, Avg:
			return 0.0, nil
		default:
			return nil, nil
		}
	}

	// Apply the aggregation function
	switch agg.Type {
	case Count:
		return len(values), nil

	case Sum:
		return operations.CalculateSum(values)

	case Avg:
		sum, err := operations.CalculateSum(values)
		if err != nil {
			return nil, err
		}
		sumFloat, ok := sum.(float64)
		if !ok {
			return nil, fmt.Errorf("cannot calculate average: sum is not a number")
		}
		return sumFloat / float64(len(values)), nil

	case Min:
		return operations.CalculateMin(values)

	case Max:
		return operations.CalculateMax(values)

	default:
		return nil, fmt.Errorf("unsupported aggregation type: %v", agg.Type)
	}
}

// ExtractAggregations extracts aggregation functions from SELECT expressions
func ExtractAggregations(selectExprs sqlparser.SelectExprs) []AggregationInfo {
	var aggregations []AggregationInfo

	for _, expr := range selectExprs {
		aliasedExpr, ok := expr.(*sqlparser.AliasedExpr)
		if !ok {
			continue
		}

		funcExpr, ok := aliasedExpr.Expr.(*sqlparser.FuncExpr)
		if !ok {
			continue
		}

		funcName := strings.ToUpper(funcExpr.Name.String())
		var aggType AggregationType
		switch funcName {
		case "COUNT":
			aggType = Count
		case "SUM":
			aggType = Sum
		case "AVG":
			aggType = Avg
		case "MIN":
			aggType = Min
		case "MAX":
			aggType = Max
		default:
			continue // Not an aggregation function
		}

		// Get the column and check for DISTINCT
		columnName := "*"
		isDistinct := funcExpr.Distinct

		if funcName != "COUNT" || len(funcExpr.Exprs) > 0 {
			// For non-COUNT(*) functions, get the column name
			if len(funcExpr.Exprs) > 0 {
				argExpr, ok := funcExpr.Exprs[0].(*sqlparser.AliasedExpr)
				if !ok {
					continue
				}

				colName, ok := argExpr.Expr.(*sqlparser.ColName)
				if ok {
					columnName = colName.Name.String()
				}
			}
		}

		// Get the alias if provided, otherwise use function name
		alias := funcName
		if len(columnName) > 0 && columnName != "*" {
			alias = fmt.Sprintf("%s(%s)", funcName, columnName)
		} else if columnName == "*" {
			alias = fmt.Sprintf("%s(*)", funcName)
		}

		if !aliasedExpr.As.IsEmpty() {
			alias = aliasedExpr.As.String()
		}

		aggregations = append(aggregations, AggregationInfo{
			Type:       aggType,
			Column:     columnName,
			Alias:      alias,
			IsDistinct: isDistinct,
		})
	}

	return aggregations
}

// HasAggregations checks if the query contains aggregation functions
func HasAggregations(selectExprs sqlparser.SelectExprs) bool {
	for _, expr := range selectExprs {
		aliasedExpr, ok := expr.(*sqlparser.AliasedExpr)
		if !ok {
			continue
		}

		funcExpr, ok := aliasedExpr.Expr.(*sqlparser.FuncExpr)
		if !ok {
			continue
		}

		funcName := strings.ToUpper(funcExpr.Name.String())
		switch funcName {
		case "COUNT", "SUM", "AVG", "MIN", "MAX":
			return true
		}
	}
	return false
}
