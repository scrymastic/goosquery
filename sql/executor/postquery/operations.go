package postquery

import (
	"fmt"
	"sort"
	"strconv"

	"github.com/blastrain/vitess-sqlparser/sqlparser"
	"github.com/scrymastic/goosquery/sql/result"
)

// Operations contains methods for post-query operations
type Operations struct {
	// Utility methods are provided as function types
	// This allows for dependency injection and easier testing
	Compare func(a, b interface{}) int
}

// Apply applies post-query operations like ORDER BY, LIMIT, etc.
func (p *Operations) Apply(results *result.QueryResult, stmt *sqlparser.Select) (*result.QueryResult, error) {
	// Apply ORDER BY if present
	if len(stmt.OrderBy) > 0 {
		err := p.Sort(results, stmt.OrderBy)
		if err != nil {
			return nil, fmt.Errorf("failed to sort results: %w", err)
		}
	}

	// Apply LIMIT and OFFSET if present
	if stmt.Limit != nil {
		limitedResults, err := p.ApplyLimit(results, stmt.Limit)
		if err != nil {
			return nil, err
		}
		return limitedResults, nil
	}

	return results, nil
}

// Sort sorts the query results based on the ORDER BY clause
func (p *Operations) Sort(results *result.QueryResult, orderBy sqlparser.OrderBy) error {
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
			compResult := p.Compare(valI, valJ)

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

// ApplyLimit applies the LIMIT and OFFSET clauses to the results
func (p *Operations) ApplyLimit(results *result.QueryResult, limit *sqlparser.Limit) (*result.QueryResult, error) {
	// Process limit value
	count, ok := limit.Rowcount.(*sqlparser.SQLVal)
	if !ok || count.Type != sqlparser.IntVal {
		return nil, fmt.Errorf("invalid LIMIT value: %v", limit.Rowcount)
	}

	limitNum, err := strconv.Atoi(string(count.Val))
	if err != nil {
		return nil, fmt.Errorf("failed to parse LIMIT value: %v", err)
	}

	// Process offset value if present
	offsetNum := 0
	if limit.Offset != nil {
		offset, ok := limit.Offset.(*sqlparser.SQLVal)
		if !ok || offset.Type != sqlparser.IntVal {
			return nil, fmt.Errorf("invalid OFFSET value: %v", limit.Offset)
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

	return results, nil
}
