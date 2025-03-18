package executor

import (
	"fmt"

	"github.com/blastrain/vitess-sqlparser/sqlparser"
	"github.com/scrymastic/goosquery/sql/result"
	"github.com/scrymastic/goosquery/tables/system"
)

// ProcessesExecutor handles queries for the processes table
type ProcessesExecutor struct {
	BaseExecutor
}

// NewProcessesExecutor creates a new processes executor
func NewProcessesExecutor() *ProcessesExecutor {
	return &ProcessesExecutor{}
}

// Execute executes a query against the processes table
func (e *ProcessesExecutor) Execute(stmt *sqlparser.Select) (*result.QueryResult, error) {
	// Fetch all processes
	processes, err := system.GenProcesses()
	if err != nil {
		return nil, fmt.Errorf("failed to get processes: %w", err)
	}

	// Create result
	res := result.NewQueryResult()

	// Get selected columns
	selectedColumns := e.GetSelectedColumns(stmt.SelectExprs)

	// Apply WHERE clause if present and project columns
	for _, proc := range processes {
		// Convert process to map for filtering
		procMap := e.processToMap(proc)

		if stmt.Where == nil || e.MatchesWhereClause(procMap, stmt.Where.Expr) {
			// Project only the selected columns
			projectedRow := e.ProjectRow(procMap, selectedColumns)
			res.AddRow(projectedRow)
		}
	}

	// Set the result columns explicitly if specific columns were selected
	if selectedColumns != nil && len(selectedColumns) > 0 && len(res.Rows) > 0 {
		res.SetColumns(selectedColumns)
	}

	return res, nil
}

// processToMap converts a Process struct to a map for output
func (e *ProcessesExecutor) processToMap(proc system.Process) map[string]interface{} {
	return map[string]interface{}{
		"pid":                proc.PID,
		"name":               proc.Name,
		"path":               proc.Path,
		"cmdline":            proc.CMDLine,
		"state":              proc.State,
		"cwd":                proc.CWD,
		"root":               proc.Root,
		"on_disk":            proc.OnDisk,
		"wired_size":         proc.WiredSize,
		"resident_size":      proc.ResidentSize,
		"total_size":         proc.TotalSize,
		"disk_bytes_read":    proc.DiskBytesRead,
		"disk_bytes_written": proc.DiskBytesWritten,
		"parent":             proc.Parent,
		"threads":            proc.Threads,
		"nice":               proc.Nice,
		"elevated_token":     proc.ElevatedToken,
		"elapsed_time":       proc.ElapsedTime,
		"handle_count":       proc.HandleCount,
	}
}
