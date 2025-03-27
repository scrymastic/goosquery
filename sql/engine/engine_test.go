package engine

import (
	"testing"
)

// Test query execution
func TestExecute(t *testing.T) {
	engine := NewEngine()
	query := "select count(name) from processes where name = 'Cursor.exe';"
	result, err := engine.Execute(query)

	if err != nil {
		t.Fatalf("Failed to execute query: %v", err)
	}

	if result == nil {
		t.Fatalf("Expected non-nil result, got nil")
	}
}
