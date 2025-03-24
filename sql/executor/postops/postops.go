package postops

import (
	"fmt"
	"strconv"

	"github.com/blastrain/vitess-sqlparser/sqlparser"
	"github.com/scrymastic/goosquery/sql/executor/operations"
	"github.com/scrymastic/goosquery/sql/result"
)

// ApplyPostQueryOperations applies operations that work on the result set (ORDER BY, LIMIT, etc.)
func ApplyPostQueryOperations(results *result.QueryResult, stmt *sqlparser.Select) (*result.QueryResult, error) {
	// Apply ORDER BY if present
	if len(stmt.OrderBy) > 0 {
		if err := operations.SortResults(results, stmt.OrderBy); err != nil {
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
		if offsetNum >= len(*results) {
			// If offset is beyond the result set, return empty results
			*results = result.QueryResult{}
			return results, nil
		}

		// Apply the offset and limit
		if offsetNum+limitNum > len(*results) {
			limitNum = len(*results) - offsetNum
		}
		*results = (*results)[offsetNum : offsetNum+limitNum]
	}

	return results, nil
}
