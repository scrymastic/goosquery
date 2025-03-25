package sqlctx

import "slices"

// Context provides information about the SQL query execution context
// It stores metadata, constants, and other information relevant to query execution
type Context struct {
	// Constants extracted from the WHERE clause (e.g., id = 1, name = 'foo')
	Constants map[string]string

	// Columns requested in the query
	Columns []string

	// // Additional query metadata
	// Metadata map[string]interface{}
}

// NewContext creates a new query execution context
func NewContext() *Context {
	return &Context{
		Constants: make(map[string]string),
		Columns:   []string{},
		// Metadata:  make(map[string]interface{}),
	}
}

// AddConstant adds a constant value extracted from the query
// For example, in WHERE id = 1, this would add "id" -> "1"
func (c *Context) AddConstant(name string, value string) {
	if c.Constants == nil {
		c.Constants = make(map[string]string)
	}
	c.Constants[name] = value
}

// HasConstant checks if a constant exists in the context
func (c *Context) HasConstant(name string) bool {
	if c.Constants == nil {
		return false
	}
	_, exists := c.Constants[name]
	return exists
}

// GetAllConstantNames returns the names of all constants
func (c *Context) GetAllConstantNames() []string {
	if c.Constants == nil {
		return []string{}
	}

	names := make([]string, 0, len(c.Constants))
	for name := range c.Constants {
		names = append(names, name)
	}
	return names
}

// AddColumn adds a column to the list of requested columns
func (c *Context) AddColumn(column string) {
	if c.Columns == nil {
		c.Columns = []string{}
	}
	c.Columns = append(c.Columns, column)
}

// SetColumns sets the list of requested columns
func (c *Context) SetColumns(columns []string) {
	c.Columns = columns
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
	return []string{c.Constants[key]}
}
