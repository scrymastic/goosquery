package context

import "slices"

// Context provides information about the SQL query execution context
// It stores metadata, constants, and other information relevant to query execution
type Context struct {
	// Constants extracted from the WHERE clause (e.g., id = 1, name = 'foo')
	Constants map[string]string

	// Columns requested in the query
	Columns []string

	// Additional query metadata
	Metadata map[string]interface{}
}

// NewContext creates a new query execution context
func NewContext() *Context {
	return &Context{
		Constants: make(map[string]string),
		Columns:   []string{},
		Metadata:  make(map[string]interface{}),
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

// GetConstant retrieves a constant value by name
// Returns the value and a boolean indicating if the constant exists
func (c *Context) GetConstant(name string) (string, bool) {
	if c.Constants == nil {
		return "", false
	}
	value, exists := c.Constants[name]
	return value, exists
}

// HasConstant checks if a constant exists in the context
func (c *Context) HasConstant(name string) bool {
	if c.Constants == nil {
		return false
	}
	_, exists := c.Constants[name]
	return exists
}

// GetAllConstants returns all constants as a map
func (c *Context) GetAllConstants() map[string]string {
	if c.Constants == nil {
		return make(map[string]string)
	}
	return c.Constants
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

// SetMetadata sets a metadata value by key
func (c *Context) SetMetadata(key string, value interface{}) {
	if c.Metadata == nil {
		c.Metadata = make(map[string]interface{})
	}
	c.Metadata[key] = value
}

// GetMetadata retrieves a metadata value by key
// Returns the value and a boolean indicating if the key exists
func (c *Context) GetMetadata(key string) (interface{}, bool) {
	if c.Metadata == nil {
		return nil, false
	}
	value, exists := c.Metadata[key]
	return value, exists
}

// GetAllMetadata returns all metadata as a map
func (c *Context) GetAllMetadata() map[string]interface{} {
	if c.Metadata == nil {
		return make(map[string]interface{})
	}
	return c.Metadata
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

// GetColumns returns the list of requested columns
func (c *Context) GetColumns() []string {
	if c.Columns == nil {
		return []string{}
	}
	return c.Columns
}

// Clone creates a deep copy of the context
func (c *Context) Clone() *Context {
	clone := NewContext()

	// Copy constants
	for name, value := range c.Constants {
		clone.Constants[name] = value
	}

	// Copy columns
	clone.Columns = make([]string, len(c.Columns))
	copy(clone.Columns, c.Columns)

	// Copy metadata
	for key, value := range c.Metadata {
		clone.Metadata[key] = value
	}

	return clone
}

// Merge combines this context with another context
// Values from the other context will override values in this context if they conflict
func (c *Context) Merge(other *Context) {
	// Merge constants
	for name, value := range other.Constants {
		c.Constants[name] = value
	}

	// Merge columns (deduplicate)
	existingColumns := make(map[string]bool)
	for _, column := range c.Columns {
		existingColumns[column] = true
	}

	for _, column := range other.Columns {
		if !existingColumns[column] {
			c.Columns = append(c.Columns, column)
			existingColumns[column] = true
		}
	}

	// Merge metadata
	for key, value := range other.Metadata {
		c.Metadata[key] = value
	}
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

// HasConstants checks if there are any constants with the given key.
// Returns true if the key exists and has at least one value.
func (c *Context) HasConstants(key string) bool {
	if c.Constants == nil {
		return false
	}
	_, exists := c.Constants[key]
	return exists
}
