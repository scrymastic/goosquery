package result

// Results represents the result of a SQL query
// It's implemented as a slice of results (maps) for flexibility
type Results []Result

// NewQueryResult creates a new empty query result
func NewQueryResult() *Results {
	res := Results{}
	return &res
}

// AppendResult adds a result (row) to the query result
func (r *Results) AppendResult(result Result) {
	*r = append(*r, result)
}

func (r *Results) AppendResults(results Results) {
	*r = append(*r, results...)
}

// GetColumns returns all column names in the result
// This inspects all results and combines their keys to get all possible columns
func (r *Results) GetColumns() []string {
	if len(*r) == 0 {
		return []string{}
	}

	// Use a map to eliminate duplicates
	columnMap := make(map[string]bool)

	// Collect all column names from all results
	for _, result := range *r {
		for column := range result {
			columnMap[column] = true
		}
	}

	// Convert map keys to slice
	columns := make([]string, 0, len(columnMap))
	for column := range columnMap {
		columns = append(columns, column)
	}

	return columns
}

// IsEmpty returns true if the result contains no results
func (r *Results) IsEmpty() bool {
	return len(*r) == 0
}

// Size returns the number of results in the result
func (r *Results) Size() int {
	return len(*r)
}

// Clone creates a deep copy of the query result
func (r *Results) Clone() *Results {
	if r.IsEmpty() {
		return NewQueryResult()
	}

	clone := NewQueryResult()
	for _, result := range *r {
		// Create a new map for each result
		resultCopy := make(map[string]interface{})
		for k, v := range result {
			resultCopy[k] = v
		}
		clone.AppendResult(resultCopy)
	}

	return clone
}

// GetValue retrieves a value from a specific result and column
// Returns the value and a boolean indicating if the value was found
func (r *Results) GetValue(rowIndex int, columnName string) (interface{}, bool) {
	if rowIndex < 0 || rowIndex >= r.Size() {
		return nil, false
	}

	result := (*r)[rowIndex]
	value, exists := result[columnName]
	return value, exists
}

// SetValue sets a value in a specific result and column
// Returns true if successful, false if the result doesn't exist
func (r *Results) SetValue(rowIndex int, columnName string, value interface{}) bool {
	if rowIndex < 0 || rowIndex >= r.Size() {
		return false
	}

	result := (*r)[rowIndex]
	result[columnName] = value
	return true
}

// GetRow returns a specific result by index
func (r *Results) GetRow(rowIndex int) Result {
	if rowIndex < 0 || rowIndex >= r.Size() {
		return nil
	}

	return (*r)[rowIndex]
}

// GetColumnValues returns all values for a specific column
func (r *Results) GetColumnValues(columnName string) []interface{} {
	values := make([]interface{}, 0, r.Size())

	for _, result := range *r {
		if value, exists := result[columnName]; exists {
			values = append(values, value)
		} else {
			values = append(values, nil)
		}
	}

	return values
}

// ForEach executes a function for each result in the result
func (r *Results) ForEach(fn func(Result) error) error {
	for _, result := range *r {
		if err := fn(result); err != nil {
			return err
		}
	}
	return nil
}

// Filter returns a new QueryResult containing only results that pass the filter function
func (r *Results) Filter(fn func(Result) bool) *Results {
	filtered := NewQueryResult()

	for _, result := range *r {
		if fn(result) {
			filtered.AppendResult(result)
		}
	}

	return filtered
}

// Map applies a transformation function to each result and returns a new QueryResult
func (r *Results) Map(fn func(Result) Result) *Results {
	mapped := NewQueryResult()

	for _, result := range *r {
		transformedResult := fn(result)
		mapped.AppendResult(transformedResult)
	}

	return mapped
}
