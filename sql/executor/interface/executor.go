package execintf

import (
	"github.com/blastrain/vitess-sqlparser/sqlparser"
	"github.com/scrymastic/goosquery/sql/result"
)

// Executor is the interface for query executors
type Executor interface {
	Execute(stmt *sqlparser.Select) (*result.Results, error)
}
