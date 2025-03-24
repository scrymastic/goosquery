package engine

import (
	"fmt"

	"github.com/blastrain/vitess-sqlparser/sqlparser"
	execintf "github.com/scrymastic/goosquery/sql/executor/interface"
	"github.com/scrymastic/goosquery/sql/parser"
	"github.com/scrymastic/goosquery/sql/result"
)

// Engine provides SQL query capabilities
type Engine struct{}

// NewEngine creates a new SQL engine
func NewEngine() *Engine {
	return &Engine{}
}

// Execute executes a SQL query and returns the result
func (e *Engine) Execute(query string) (*result.QueryResult, error) {
	// Parse the SQL query
	parsedQuery, err := parser.Parse(query)
	if err != nil {
		return nil, err
	}

	// Get the table name
	tableName, err := parser.GetTableName(parsedQuery.Statement)
	if err != nil {
		return nil, err
	}

	// Get the executor for this table
	exec, err := execintf.GetExecutor(tableName)
	if err != nil {
		return nil, err
	}

	// Execute the query
	selectStmt, ok := parsedQuery.Statement.(*sqlparser.Select)
	if !ok {
		return nil, fmt.Errorf("only SELECT statements are supported")
	}

	return exec.Execute(selectStmt)
}
