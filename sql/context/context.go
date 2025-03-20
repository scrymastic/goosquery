package context

import "slices"

// Context holds information about query execution.
// It contains both column information and constants/parameters
// extracted from the query (like WHERE clause conditions).
type Context struct {
	// Columns contains the list of columns used in the query
	Columns []string

	// Constants is a map of parameter name to values.
	// For SQL queries, this typically contains constants from
	// the WHERE clause, such as key=value conditions.
	// For example, in "SELECT * FROM curl WHERE url='https://example.com'",
	// Constants would contain: {"url": ["https://example.com"]}
	Constants map[string][]string
}

// IsColumnUsed checks if a column is used in the query.
// If no columns are specified (empty Columns slice),
// it returns true for all columns.
func (c *Context) IsColumnUsed(column string) bool {
	if len(c.Columns) == 0 {
		return true
	}
	return slices.Contains(c.Columns, column)
}

// IsAnyOfColumnsUsed checks if any of the specified columns are used.
// Returns true if at least one of the columns in the input list
// is used in the query.
func (c *Context) IsAnyOfColumnsUsed(columns []string) bool {
	for _, column := range columns {
		if c.IsColumnUsed(column) {
			return true
		}
	}
	return false
}

// GetConstants returns all values for a specific key.
// For example, if multiple URLs are specified in WHERE clauses
// like "url='https://example1.com' OR url='https://example2.com'",
// GetConstants("url") would return ["https://example1.com", "https://example2.com"]
func (c *Context) GetConstants(key string) []string {
	if c.Constants == nil {
		return nil
	}
	return c.Constants[key]
}

// AddConstant adds a constant value for a key.
// This is typically used by executors that extract parameters
// from WHERE clauses and other parts of the query.
func (c *Context) AddConstant(key, value string) {
	if c.Constants == nil {
		c.Constants = make(map[string][]string)
	}
	c.Constants[key] = append(c.Constants[key], value)
}

// HasConstants checks if there are any constants with the given key.
// Returns true if the key exists and has at least one value.
func (c *Context) HasConstants(key string) bool {
	if c.Constants == nil {
		return false
	}
	values, exists := c.Constants[key]
	return exists && len(values) > 0
}

// GetAllConstantNames returns all constant names in the context.
func (c *Context) GetAllConstantNames() []string {
	names := make([]string, 0, len(c.Constants))
	for name := range c.Constants {
		names = append(names, name)
	}
	return names
}

// GetAllConstants returns all constants in the context.
func (c *Context) GetAllConstants() map[string][]string {
	return c.Constants
}
