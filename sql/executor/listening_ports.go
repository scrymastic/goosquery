package executor

import (
	"fmt"

	"github.com/blastrain/vitess-sqlparser/sqlparser"
	"github.com/scrymastic/goosquery/sql/result"
	"github.com/scrymastic/goosquery/tables/networking"
)

// ListeningPortsExecutor handles queries for the listening_ports table
type ListeningPortsExecutor struct {
	BaseExecutor
}

// NewListeningPortsExecutor creates a new executor for listening ports
func NewListeningPortsExecutor() *ListeningPortsExecutor {
	return &ListeningPortsExecutor{}
}

// Execute executes a query against the listening_ports table
func (e *ListeningPortsExecutor) Execute(stmt *sqlparser.Select) (*result.QueryResult, error) {
	// Fetch all listening ports
	ports, err := networking.GenListeningPorts()
	if err != nil {
		return nil, fmt.Errorf("failed to get listening ports: %w", err)
	}

	// Create result
	res := result.NewQueryResult()

	// Get selected columns
	selectedColumns := e.GetSelectedColumns(stmt.SelectExprs)

	// Apply WHERE clause if present and project columns
	for _, port := range ports {
		// Convert port to map for filtering
		portMap := e.portToMap(port)

		if stmt.Where == nil || e.MatchesWhereClause(portMap, stmt.Where.Expr) {
			// Project only the selected columns
			projectedRow := e.ProjectRow(portMap, selectedColumns)
			res.AddRow(projectedRow)
		}
	}

	// Set the result columns explicitly if specific columns were selected
	if selectedColumns != nil && len(selectedColumns) > 0 && len(res.Rows) > 0 {
		res.SetColumns(selectedColumns)
	}

	return res, nil
}

// portToMap converts a ListeningPort struct to a map for output
func (e *ListeningPortsExecutor) portToMap(port networking.ListeningPort) map[string]interface{} {
	return map[string]interface{}{
		"pid":      port.PID,
		"port":     port.Port,
		"protocol": port.Protocol,
		"family":   port.Family,
		"address":  port.Address,
		"fd":       port.FD,
		"socket":   port.Socket,
		"path":     port.Path,
	}
}
