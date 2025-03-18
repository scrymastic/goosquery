package result

import (
	"encoding/json"
	"fmt"
	"sort"
)

// QueryResult represents the result of a SQL query
type QueryResult struct {
	Rows    []map[string]interface{} `json:"rows"`
	Columns []string                 `json:"columns"`
}

// NewQueryResult creates a new empty query result
func NewQueryResult() *QueryResult {
	return &QueryResult{
		Rows:    []map[string]interface{}{},
		Columns: []string{},
	}
}

// AddRow adds a row to the result and updates columns if needed
func (r *QueryResult) AddRow(row map[string]interface{}) {
	// If this is the first row and columns are not set, initialize columns from the row
	if len(r.Columns) == 0 && len(r.Rows) == 0 {
		for k := range row {
			r.Columns = append(r.Columns, k)
		}
		// Sort columns for consistent output
		sort.Strings(r.Columns)
	}
	r.Rows = append(r.Rows, row)
}

// SetColumns explicitly sets the result columns, overriding any columns set from rows
func (r *QueryResult) SetColumns(columns []string) {
	r.Columns = columns
}

// ToJSON converts the result to a formatted JSON string
func (r *QueryResult) ToJSON() (string, error) {
	jsonData, err := json.MarshalIndent(r, "", "  ")
	if err != nil {
		return "", fmt.Errorf("failed to convert result to JSON: %w", err)
	}
	return string(jsonData), nil
}
