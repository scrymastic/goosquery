package aggregation

// AggregationType represents the type of aggregation function
type AggregationType int

const (
	// Count represents COUNT aggregation function
	Count AggregationType = iota
	// Sum represents SUM aggregation function
	Sum
	// Avg represents AVG aggregation function
	Avg
	// Min represents MIN aggregation function
	Min
	// Max represents MAX aggregation function
	Max
)

// AggregationInfo represents an aggregation operation
type AggregationInfo struct {
	Type       AggregationType
	Column     string
	Alias      string
	IsDistinct bool
}
