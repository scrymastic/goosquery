package impl

import (
	"fmt"

	"github.com/blastrain/vitess-sqlparser/sqlparser"
	"github.com/scrymastic/goosquery/sql/context"
	"github.com/scrymastic/goosquery/sql/executor/aggregation"
	"github.com/scrymastic/goosquery/sql/executor/postops"
	"github.com/scrymastic/goosquery/sql/executor/projection"
	"github.com/scrymastic/goosquery/sql/result"
)

// DataGenerator is a function that generates data for a table
type DataGenerator func(context.Context) ([]map[string]interface{}, error)

// TableExecutor is a generic executor for tables that return []map[string]interface{}
type TableExecutor struct {
	TableName string
	Generator DataGenerator
	BaseExecutor
}

// Execute executes a query against the table using the provided data function
func (e *TableExecutor) Execute(stmt *sqlparser.Select) (*result.QueryResult, error) {
	// Check if the query uses aggregation functions
	hasAggregations := aggregation.HasAggregations(stmt.SelectExprs)

	// Get aggregation information if needed
	var aggs []aggregation.AggregationInfo
	if hasAggregations {
		aggs = aggregation.ExtractAggregations(stmt.SelectExprs)
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
	data, err := e.Generator(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to get %s data: %w", e.TableName, err)
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
		res, err = aggregation.ApplyAggregations(res, aggs, stmt.GroupBy)
		if err != nil {
			return nil, fmt.Errorf("failed to apply aggregations: %w", err)
		}
	}

	// Apply post-query operations (ORDER BY, LIMIT, etc.)
	res, err = postops.ApplyPostQueryOperations(res, stmt)
	if err != nil {
		return nil, err
	}

	// Apply final projection to get only the requested columns with proper aliases
	return projection.ProjectFinalResults(res, stmt), nil
}
