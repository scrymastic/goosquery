package result

// QueryResult represents the result of a SQL query
// It's implemented as a slice of records (maps) for flexibility
type QueryResult []map[string]interface{}

// NewQueryResult creates a new empty query result
func NewQueryResult() *QueryResult {
	res := QueryResult{}
	return &res
}

// AddRecord adds a record (row) to the query result
func (r *QueryResult) AddRecord(record map[string]interface{}) {
	*r = append(*r, record)
}

// GetColumns returns all column names in the result
// This inspects all records and combines their keys to get all possible columns
func (r *QueryResult) GetColumns() []string {
	if len(*r) == 0 {
		return []string{}
	}

	// Use a map to eliminate duplicates
	columnMap := make(map[string]bool)

	// Collect all column names from all records
	for _, record := range *r {
		for column := range record {
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

// IsEmpty returns true if the result contains no records
func (r *QueryResult) IsEmpty() bool {
	return len(*r) == 0
}

// Size returns the number of records in the result
func (r *QueryResult) Size() int {
	return len(*r)
}

// Clone creates a deep copy of the query result
func (r *QueryResult) Clone() *QueryResult {
	if r.IsEmpty() {
		return NewQueryResult()
	}

	clone := NewQueryResult()
	for _, record := range *r {
		// Create a new map for each record
		recordCopy := make(map[string]interface{})
		for k, v := range record {
			recordCopy[k] = v
		}
		clone.AddRecord(recordCopy)
	}

	return clone
}

// GetValue retrieves a value from a specific record and column
// Returns the value and a boolean indicating if the value was found
func (r *QueryResult) GetValue(rowIndex int, columnName string) (interface{}, bool) {
	if rowIndex < 0 || rowIndex >= r.Size() {
		return nil, false
	}

	record := (*r)[rowIndex]
	value, exists := record[columnName]
	return value, exists
}

// SetValue sets a value in a specific record and column
// Returns true if successful, false if the record doesn't exist
func (r *QueryResult) SetValue(rowIndex int, columnName string, value interface{}) bool {
	if rowIndex < 0 || rowIndex >= r.Size() {
		return false
	}

	record := (*r)[rowIndex]
	record[columnName] = value
	return true
}

// GetRow returns a specific record by index
// Returns the record and a boolean indicating if the record was found
func (r *QueryResult) GetRow(rowIndex int) (map[string]interface{}, bool) {
	if rowIndex < 0 || rowIndex >= r.Size() {
		return nil, false
	}

	return (*r)[rowIndex], true
}

// GetColumnValues returns all values for a specific column
func (r *QueryResult) GetColumnValues(columnName string) []interface{} {
	values := make([]interface{}, 0, r.Size())

	for _, record := range *r {
		if value, exists := record[columnName]; exists {
			values = append(values, value)
		} else {
			values = append(values, nil)
		}
	}

	return values
}

// ForEach executes a function for each record in the result
func (r *QueryResult) ForEach(fn func(map[string]interface{}) error) error {
	for _, record := range *r {
		if err := fn(record); err != nil {
			return err
		}
	}
	return nil
}

// Filter returns a new QueryResult containing only records that pass the filter function
func (r *QueryResult) Filter(fn func(map[string]interface{}) bool) *QueryResult {
	filtered := NewQueryResult()

	for _, record := range *r {
		if fn(record) {
			filtered.AddRecord(record)
		}
	}

	return filtered
}

// Map applies a transformation function to each record and returns a new QueryResult
func (r *QueryResult) Map(fn func(map[string]interface{}) map[string]interface{}) *QueryResult {
	mapped := NewQueryResult()

	for _, record := range *r {
		transformedRecord := fn(record)
		mapped.AddRecord(transformedRecord)
	}

	return mapped
}
