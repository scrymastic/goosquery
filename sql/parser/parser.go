package parser

import (
	"fmt"

	"github.com/blastrain/vitess-sqlparser/sqlparser"
)

// ParsedQuery represents a parsed SQL query
type ParsedQuery struct {
	Statement sqlparser.Statement
	Original  string
}

// Parse parses a SQL query string into a structured form
func Parse(query string) (*ParsedQuery, error) {
	stmt, err := sqlparser.Parse(query)
	if err != nil {
		return nil, fmt.Errorf("SQL parse error: %w", err)
	}

	return &ParsedQuery{
		Statement: stmt,
		Original:  query,
	}, nil
}

// GetTableName extracts the table name from a query
func GetTableName(stmt sqlparser.Statement) (string, error) {
	selectStmt, ok := stmt.(*sqlparser.Select)
	if !ok {
		return "", fmt.Errorf("only SELECT statements are supported")
	}

	if len(selectStmt.From) == 0 {
		return "", fmt.Errorf("no FROM clause found")
	}

	fromExpr, ok := selectStmt.From[0].(*sqlparser.AliasedTableExpr)
	if !ok {
		return "", fmt.Errorf("unsupported FROM expression")
	}

	tableExpr, ok := fromExpr.Expr.(sqlparser.TableName)
	if !ok {
		return "", fmt.Errorf("invalid table name")
	}

	return tableExpr.Name.String(), nil
}
