package result

// QueryResult represents the result of a SQL query
type QueryResult []map[string]interface{}

// NewQueryResult creates a new empty query result
func NewQueryResult() *QueryResult {
	return &QueryResult{}
}

func (r *QueryResult) AddRecord(record map[string]interface{}) {
	*r = append(*r, record)
}
