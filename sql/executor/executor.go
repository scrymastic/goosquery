package executor

import (
	"fmt"

	"github.com/blastrain/vitess-sqlparser/sqlparser"
	"github.com/scrymastic/goosquery/sql/result"
)

// Executor is the interface for query executors
type Executor interface {
	Execute(stmt *sqlparser.Select) (*result.QueryResult, error)
}

// GetExecutor returns the appropriate executor for a given table
func GetExecutor(tableName string) (Executor, error) {
	switch tableName {
	case "processes":
		return NewProcessesExecutor(), nil
	case "listening_ports":
		return NewListeningPortsExecutor(), nil
	// Add other table executors here
	default:
		return nil, fmt.Errorf("unsupported table: %s", tableName)
	}
}
