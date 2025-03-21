package executor

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/blastrain/vitess-sqlparser/sqlparser"
	"github.com/scrymastic/goosquery/sql/context"
	"github.com/scrymastic/goosquery/sql/result"
	"github.com/scrymastic/goosquery/tables/networking"
	"github.com/scrymastic/goosquery/tables/system"
	"github.com/scrymastic/goosquery/tables/utility"
)

// Executor is the interface for query executors
type Executor interface {
	Execute(stmt *sqlparser.Select) (*result.QueryResult, error)
}

// GetExecutor returns the appropriate executor for a given table
func GetExecutor(tableName string) (Executor, error) {
	switch tableName {
	case "arp_cache":
		return &TableExecutor{
			BaseExecutor: BaseExecutor{},
			tableName:    "arp_cache",
			genarator:    networking.GenARPCache,
		}, nil
	case "connectivity":
		return &TableExecutor{
			BaseExecutor: BaseExecutor{},
			tableName:    "connectivity",
			genarator:    networking.GenConnectivity,
		}, nil
	case "etc_hosts":
		return &TableExecutor{
			BaseExecutor: BaseExecutor{},
			tableName:    "etc_hosts",
			genarator:    networking.GenEtcHosts,
		}, nil
	case "curl":
		return &TableExecutor{
			BaseExecutor: BaseExecutor{},
			tableName:    "curl",
			genarator:    networking.GenCurl,
		}, nil
	case "listening_ports":
		return &TableExecutor{
			BaseExecutor: BaseExecutor{},
			tableName:    "listening_ports",
			genarator:    networking.GenListeningPorts,
		}, nil
	case "process_open_sockets":
		return &TableExecutor{
			BaseExecutor: BaseExecutor{},
			tableName:    "process_open_sockets",
			genarator:    networking.GenProcessOpenSockets,
		}, nil
	case "routes":
		return &TableExecutor{
			BaseExecutor: BaseExecutor{},
			tableName:    "routes",
			genarator:    networking.GenRoutes,
		}, nil
	case "windows_firewall_rules":
		return &TableExecutor{
			BaseExecutor: BaseExecutor{},
			tableName:    "windows_firewall_rules",
			genarator:    networking.GenWindowsFirewallRules,
		}, nil
	case "appcompat_shims":
		return &TableExecutor{
			BaseExecutor: BaseExecutor{},
			tableName:    "appcompat_shims",
			genarator:    system.GenAppCompatShims,
		}, nil
	case "bitlocker_info":
		return &TableExecutor{
			BaseExecutor: BaseExecutor{},
			tableName:    "bitlocker_info",
			genarator:    system.GenBitlockerInfo,
		}, nil
	case "background_activities_moderator":
		return &TableExecutor{
			BaseExecutor: BaseExecutor{},
			tableName:    "background_activities_moderator",
			genarator:    system.GenBackgroundActivitiesModerator,
		}, nil
	case "authenticode":
		return &TableExecutor{
			BaseExecutor: BaseExecutor{},
			tableName:    "authenticode",
			genarator:    system.GenAuthenticode,
		}, nil
	case "processes":
		return &TableExecutor{
			BaseExecutor: BaseExecutor{},
			tableName:    "processes",
			genarator:    system.GenProcesses,
		}, nil
	case "uptime":
		return &TableExecutor{
			BaseExecutor: BaseExecutor{},
			tableName:    "uptime",
			genarator:    system.GenUptime,
		}, nil
	case "default_environment":
		return &TableExecutor{
			BaseExecutor: BaseExecutor{},
			tableName:    "default_environment",
			genarator:    system.GenDefaultEnvironments,
		}, nil
	case "os_version":
		return &TableExecutor{
			BaseExecutor: BaseExecutor{},
			tableName:    "os_version",
			genarator:    system.GenOSVersion,
		}, nil
	case "file":
		return &TableExecutor{
			BaseExecutor: BaseExecutor{},
			tableName:    "file",
			genarator:    utility.GenFile,
		}, nil
	case "time":
		return &TableExecutor{
			BaseExecutor: BaseExecutor{},
			tableName:    "time",
			genarator:    utility.GenTime,
		}, nil
	default:
		return nil, fmt.Errorf("unsupported table: %s", tableName)
	}
}

// TableExecutor is a generic executor for tables that return []map[string]interface{}
type TableExecutor struct {
	BaseExecutor
	tableName string
	genarator func(context.Context) ([]map[string]interface{}, error)
}

// Execute executes a query against the table using the provided data function
func (e *TableExecutor) Execute(stmt *sqlparser.Select) (*result.QueryResult, error) {
	// Check if the query uses aggregation functions
	hasAggregations := e.HasAggregations(stmt.SelectExprs)

	// Get aggregation information if needed
	var aggregations []AggregationInfo
	if hasAggregations {
		aggregations = e.GetAggregations(stmt.SelectExprs)
	}

	// Get all required columns for this query - these are the columns we need to fetch
	requiredColumns := e.GetAllRequiredColumns(stmt)

	// Create context for query execution
	ctx := context.Context{}

	// Get constants from WHERE clause
	if stmt.Where != nil {
		e.GetConstants(stmt.Where.Expr, &ctx)
	}

	// Set the columns in the context to ensure all required data is fetched
	ctx.SetColumns(requiredColumns)

	// Fetch data with all necessary columns
	data, err := e.genarator(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to get %s data: %w", e.tableName, err)
	}

	// Create result
	res := result.NewQueryResult()

	// Apply WHERE clause if present
	for _, itemMap := range data {
		if stmt.Where == nil || e.MatchesWhereClause(itemMap, stmt.Where.Expr) {
			// Add all columns at this stage - we'll project down later
			res.AddRecord(itemMap)
		}
	}

	// Apply aggregations if needed
	if hasAggregations {
		res, err = e.ApplyAggregations(res, aggregations, stmt.GroupBy)
		if err != nil {
			return nil, fmt.Errorf("failed to apply aggregations: %w", err)
		}
	}

	// Apply post-query operations (ORDER BY, LIMIT, etc.)
	res, err = e.ApplyPostQueryOperations(res, stmt)
	if err != nil {
		return nil, err
	}

	// Apply final projection to get only the requested columns with proper aliases
	return e.ProjectFinalResults(res, stmt), nil
}

// ApplyPostQueryOperations applies operations that work on the result set (ORDER BY, LIMIT, etc.)
func (e *TableExecutor) ApplyPostQueryOperations(results *result.QueryResult, stmt *sqlparser.Select) (*result.QueryResult, error) {
	// Apply ORDER BY if present
	if len(stmt.OrderBy) > 0 {
		err := e.SortResults(results, stmt.OrderBy)
		if err != nil {
			return nil, fmt.Errorf("failed to sort results: %w", err)
		}
	}

	// Apply LIMIT and OFFSET if present
	if stmt.Limit != nil {
		// Process limit value
		count, ok := stmt.Limit.Rowcount.(*sqlparser.SQLVal)
		if !ok || count.Type != sqlparser.IntVal {
			return nil, fmt.Errorf("invalid LIMIT value: %v", stmt.Limit.Rowcount)
		}

		limitNum, err := strconv.Atoi(string(count.Val))
		if err != nil {
			return nil, fmt.Errorf("failed to parse LIMIT value: %v", err)
		}

		// Process offset value if present
		offsetNum := 0
		if stmt.Limit.Offset != nil {
			offset, ok := stmt.Limit.Offset.(*sqlparser.SQLVal)
			if !ok || offset.Type != sqlparser.IntVal {
				return nil, fmt.Errorf("invalid OFFSET value: %v", stmt.Limit.Offset)
			}

			offsetNum, err = strconv.Atoi(string(offset.Val))
			if err != nil {
				return nil, fmt.Errorf("failed to parse OFFSET value: %v", err)
			}

			if offsetNum < 0 {
				offsetNum = 0
			}
		}

		// Apply limit and offset
		if offsetNum < len(*results) {
			// Create a new result with the limited rows
			limitedResult := result.NewQueryResult()

			// Calculate end position considering both limit and total results
			endPos := len(*results)
			if limitNum > 0 && offsetNum+limitNum < endPos {
				endPos = offsetNum + limitNum
			}

			// Copy the relevant subset of rows
			for i := offsetNum; i < endPos; i++ {
				limitedResult.AddRecord((*results)[i])
			}

			return limitedResult, nil
		} else if offsetNum > 0 {
			// Offset is beyond available rows, return empty result
			return result.NewQueryResult(), nil
		}
	}

	return results, nil
}

// ApplyAggregations applies aggregation functions to the result set
func (e *TableExecutor) ApplyAggregations(results *result.QueryResult, aggregations []AggregationInfo, groupBy sqlparser.GroupBy) (*result.QueryResult, error) {
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
		return e.aggregateAll(results, aggregations)
	}

	// Otherwise, group the results and apply aggregations to each group
	return e.aggregateByGroups(results, aggregations, groupByColumns)
}

// aggregateAll applies aggregations to the entire result set (no GROUP BY)
func (e *TableExecutor) aggregateAll(results *result.QueryResult, aggregations []AggregationInfo) (*result.QueryResult, error) {
	// Create a single aggregated row
	aggregatedRow := make(map[string]interface{})

	// Apply each aggregation function
	for _, agg := range aggregations {
		value, err := e.calculateAggregation(results, agg)
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
func (e *TableExecutor) aggregateByGroups(results *result.QueryResult, aggregations []AggregationInfo, groupByColumns []string) (*result.QueryResult, error) {
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
			value, err := e.calculateAggregation(groupResults, agg)
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
func (e *TableExecutor) calculateAggregation(results *result.QueryResult, agg AggregationInfo) (interface{}, error) {
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
		return e.Sum(values)

	case Avg:
		sum, err := e.Sum(values)
		if err != nil {
			return nil, err
		}
		sumFloat, ok := sum.(float64)
		if !ok {
			return nil, fmt.Errorf("cannot calculate average: sum is not a number")
		}
		return sumFloat / float64(len(values)), nil

	case Min:
		return e.Min(values)

	case Max:
		return e.Max(values)

	default:
		return nil, fmt.Errorf("unsupported aggregation type: %v", agg.Type)
	}
}

// Sum calculates the sum of a set of values
func (e *TableExecutor) Sum(values []interface{}) (interface{}, error) {
	var sum float64

	for _, val := range values {
		// Convert to float64 if possible
		numVal, isNum := e.BaseExecutor.toFloat64(val)
		if !isNum {
			return nil, fmt.Errorf("cannot sum non-numeric value: %v", val)
		}
		sum += numVal
	}

	return sum, nil
}

// Min finds the minimum value in a set
func (e *TableExecutor) Min(values []interface{}) (interface{}, error) {
	if len(values) == 0 {
		return nil, nil
	}

	minVal := values[0]

	for _, val := range values[1:] {
		if e.Compare(val, minVal) < 0 {
			minVal = val
		}
	}

	return minVal, nil
}

// Max finds the maximum value in a set
func (e *TableExecutor) Max(values []interface{}) (interface{}, error) {
	if len(values) == 0 {
		return nil, nil
	}

	maxVal := values[0]

	for _, val := range values[1:] {
		if e.Compare(val, maxVal) > 0 {
			maxVal = val
		}
	}

	return maxVal, nil
}

// contains checks if a string is in a slice of strings
func contains(slice []string, s string) bool {
	for _, item := range slice {
		if item == s {
			return true
		}
	}
	return false
}

// ProjectFinalResults applies final projection to the result to ensure only the requested columns are returned
// It handles column selection, aliases from SELECT, and aggregation function results
func (e *TableExecutor) ProjectFinalResults(results *result.QueryResult, stmt *sqlparser.Select) *result.QueryResult {
	// If there are no results, return empty result
	if len(*results) == 0 {
		return result.NewQueryResult()
	}

	// If it's a SELECT * with no aliases, return the results as is
	if e.isSelectStar(stmt.SelectExprs) && !e.hasAliases(stmt.SelectExprs) {
		return results
	}

	// Create a new result with only the requested columns, including aliases
	projectedResult := result.NewQueryResult()

	// Build a map of original column names to their aliases
	aliasMap := e.buildAliasMap(stmt.SelectExprs)

	// Get the list of columns to include in the final result
	finalColumns := e.getFinalProjectionColumns(stmt.SelectExprs)

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
func (e *TableExecutor) buildAliasMap(selectExprs sqlparser.SelectExprs) map[string]string {
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
func (e *TableExecutor) getFinalProjectionColumns(selectExprs sqlparser.SelectExprs) []string {
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
func (e *TableExecutor) isSelectStar(selectExprs sqlparser.SelectExprs) bool {
	for _, expr := range selectExprs {
		if _, ok := expr.(*sqlparser.StarExpr); ok {
			return true
		}
	}
	return false
}

// hasAliases checks if any expressions have aliases
func (e *TableExecutor) hasAliases(selectExprs sqlparser.SelectExprs) bool {
	for _, expr := range selectExprs {
		aliasedExpr, ok := expr.(*sqlparser.AliasedExpr)
		if ok && !aliasedExpr.As.IsEmpty() {
			return true
		}
	}
	return false
}
