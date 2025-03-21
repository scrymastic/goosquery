package aggregation

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/blastrain/vitess-sqlparser/sqlparser"
	"github.com/scrymastic/goosquery/sql/result"
)

// Operations provides methods for aggregation operations
type Operations struct {
	// Utility methods are provided as function types
	// This allows for dependency injection and easier testing
	Compare   func(a, b interface{}) int
	ToFloat64 func(v interface{}) (float64, bool)
}

// ExtractFromSelect extracts aggregation functions from a SELECT statement
func (o *Operations) ExtractFromSelect(selectExprs sqlparser.SelectExprs) []Info {
	var aggregations []Info

	for _, expr := range selectExprs {
		aliasedExpr, ok := expr.(*sqlparser.AliasedExpr)
		if !ok {
			continue
		}

		// Check if this is an aggregation function
		aggType, column, isDistinct := GetType(aliasedExpr.Expr)
		if aggType == None {
			continue
		}

		// Get alias if available, otherwise use function name as alias
		alias := column
		if !aliasedExpr.As.IsEmpty() {
			alias = aliasedExpr.As.String()
		} else {
			alias = GenerateDefaultAlias(aggType, column)
		}

		aggregations = append(aggregations, Info{
			Type:       aggType,
			Column:     column,
			Alias:      alias,
			IsDistinct: isDistinct,
		})
	}

	return aggregations
}

// HasAggregations checks if the select statement contains any aggregation functions
func (o *Operations) HasAggregations(selectExprs sqlparser.SelectExprs) bool {
	for _, expr := range selectExprs {
		aliasedExpr, ok := expr.(*sqlparser.AliasedExpr)
		if !ok {
			continue
		}

		aggType, _, _ := GetType(aliasedExpr.Expr)
		if aggType != None {
			return true
		}
	}
	return false
}

// Apply applies aggregation functions to the result set
func (o *Operations) Apply(results *result.QueryResult, aggregations []Info, groupBy sqlparser.GroupBy) (*result.QueryResult, error) {
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
		return o.aggregateAll(results, aggregations)
	}

	// Otherwise, group the results and apply aggregations to each group
	return o.aggregateByGroups(results, aggregations, groupByColumns)
}

// aggregateAll applies aggregations to the entire result set (no GROUP BY)
func (o *Operations) aggregateAll(results *result.QueryResult, aggregations []Info) (*result.QueryResult, error) {
	// Create a single aggregated row
	aggregatedRow := make(map[string]interface{})

	// Apply each aggregation function
	for _, agg := range aggregations {
		value, err := o.calculateAggregation(results, agg)
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
func (o *Operations) aggregateByGroups(results *result.QueryResult, aggregations []Info, groupByColumns []string) (*result.QueryResult, error) {
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
			value, err := o.calculateAggregation(groupResults, agg)
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
func (o *Operations) calculateAggregation(results *result.QueryResult, agg Info) (interface{}, error) {
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
		return o.Sum(values)

	case Avg:
		sum, err := o.Sum(values)
		if err != nil {
			return nil, err
		}
		sumFloat, ok := sum.(float64)
		if !ok {
			return nil, fmt.Errorf("cannot calculate average: sum is not a number")
		}
		return sumFloat / float64(len(values)), nil

	case Min:
		return o.Min(values)

	case Max:
		return o.Max(values)

	default:
		return nil, fmt.Errorf("unsupported aggregation type: %v", agg.Type)
	}
}

// Sum calculates the sum of a set of values
func (o *Operations) Sum(values []interface{}) (interface{}, error) {
	var sum float64

	for _, val := range values {
		// Convert to float64 if possible
		numVal, isNum := o.ToFloat64(val)
		if !isNum {
			return nil, fmt.Errorf("cannot sum non-numeric value: %v", val)
		}
		sum += numVal
	}

	return sum, nil
}

// Min finds the minimum value in a set
func (o *Operations) Min(values []interface{}) (interface{}, error) {
	if len(values) == 0 {
		return nil, nil
	}

	minVal := values[0]

	for _, val := range values[1:] {
		if o.Compare(val, minVal) < 0 {
			minVal = val
		}
	}

	return minVal, nil
}

// Max finds the maximum value in a set
func (o *Operations) Max(values []interface{}) (interface{}, error) {
	if len(values) == 0 {
		return nil, nil
	}

	maxVal := values[0]

	for _, val := range values[1:] {
		if o.Compare(val, maxVal) > 0 {
			maxVal = val
		}
	}

	return maxVal, nil
}

// ParseInt parses a string to int, with proper error handling
func ParseInt(s string) (int, error) {
	return strconv.Atoi(s)
}
