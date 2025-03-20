package executor

import (
	"fmt"

	"github.com/blastrain/vitess-sqlparser/sqlparser"
	"github.com/scrymastic/goosquery/sql/context"
	"github.com/scrymastic/goosquery/sql/result"
	"github.com/scrymastic/goosquery/tables/networking"
	"github.com/scrymastic/goosquery/tables/system/default_environment"
	"github.com/scrymastic/goosquery/tables/system/os_version"
	"github.com/scrymastic/goosquery/tables/system/processes"
	"github.com/scrymastic/goosquery/tables/system/uptime"
)

// Executor is the interface for query executors
type Executor interface {
	Execute(stmt *sqlparser.Select) (*result.QueryResult, error)
}

// GetExecutor returns the appropriate executor for a given table
func GetExecutor(tableName string) (Executor, error) {
	switch tableName {
	case "processes":
		return &TableExecutor{
			BaseExecutor: BaseExecutor{},
			tableName:    "processes",
			dataFunc:     processes.GenProcesses,
		}, nil
	case "arp_cache":
		return &TableExecutor{
			BaseExecutor: BaseExecutor{},
			tableName:    "arp_cache",
			dataFunc:     networking.GenARPCache,
		}, nil
	case "connectivity":
		return &TableExecutor{
			BaseExecutor: BaseExecutor{},
			tableName:    "connectivity",
			dataFunc:     networking.GenConnectivity,
		}, nil
	case "etc_hosts":
		return &TableExecutor{
			BaseExecutor: BaseExecutor{},
			tableName:    "etc_hosts",
			dataFunc:     networking.GenEtcHosts,
		}, nil
	case "uptime":
		return &TableExecutor{
			BaseExecutor: BaseExecutor{},
			tableName:    "uptime",
			dataFunc:     uptime.GenUptime,
		}, nil
	case "default_environment":
		return &TableExecutor{
			BaseExecutor: BaseExecutor{},
			tableName:    "default_environment",
			dataFunc:     default_environment.GenDefaultEnvironments,
		}, nil
	case "os_version":
		return &TableExecutor{
			BaseExecutor: BaseExecutor{},
			tableName:    "os_version",
			dataFunc:     os_version.GenOSVersion,
		}, nil
	case "curl":
		return &TableExecutor{
			BaseExecutor: BaseExecutor{},
			tableName:    "curl",
			dataFunc:     networking.GenCurl,
		}, nil
	case "listening_ports":
		return &TableExecutor{
			BaseExecutor: BaseExecutor{},
			tableName:    "listening_ports",
			dataFunc:     networking.GenListeningPorts,
		}, nil
	case "process_open_sockets":
		return &TableExecutor{
			BaseExecutor: BaseExecutor{},
			tableName:    "process_open_sockets",
			dataFunc:     networking.GenProcessOpenSockets,
		}, nil
	case "routes":
		return &TableExecutor{
			BaseExecutor: BaseExecutor{},
			tableName:    "routes",
			dataFunc:     networking.GenRoutes,
		}, nil
	case "windows_firewall_rules":
		return &TableExecutor{
			BaseExecutor: BaseExecutor{},
			tableName:    "windows_firewall_rules",
			dataFunc:     networking.GenWindowsFirewallRules,
		}, nil
	default:
		return nil, fmt.Errorf("unsupported table: %s", tableName)
	}
}

// TableExecutor is a generic executor for tables that return []map[string]interface{}
type TableExecutor struct {
	BaseExecutor
	tableName string
	dataFunc  func(context.Context) ([]map[string]interface{}, error)
}

// Execute executes a query against the table using the provided data function
func (e *TableExecutor) Execute(stmt *sqlparser.Select) (*result.QueryResult, error) {
	// Get selected columns for display (empty slice means all columns)
	selectedColumns := e.GetSelectedColumns(stmt.SelectExprs)

	// Create context for query execution
	ctx := context.Context{}

	// Get constants from WHERE clause
	if stmt.Where != nil {
		e.GetConstants(stmt.Where.Expr, &ctx)
	}

	// Merge constants with selected columns to get all necessary columns if selectedColumns is empty
	allColumns := []string{}
	if len(selectedColumns) != 0 {
		allColumns = append(selectedColumns, ctx.GetAllConstantNames()...)
	}

	// Add constants to context
	ctx.Columns = allColumns

	// Fetch data with all necessary columns
	data, err := e.dataFunc(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to get %s data: %w", e.tableName, err)
	}

	// Create result
	res := result.NewQueryResult()

	// Apply WHERE clause if present and project columns
	for _, itemMap := range data {
		if stmt.Where == nil || e.MatchesWhereClause(itemMap, stmt.Where.Expr) {
			// Project only the selected columns
			projectedRow := e.ProjectRow(itemMap, selectedColumns)
			res.AddRecord(projectedRow)
		}
	}

	return res, nil
}

func (e *TableExecutor) GetConstants(expr sqlparser.Expr, ctx *context.Context) {
	if expr == nil {
		return
	}

	switch expr := expr.(type) {
	case *sqlparser.ComparisonExpr:
		if colName, ok := expr.Left.(*sqlparser.ColName); ok {
			fieldName := colName.Name.String()

			if sqlVal, ok := expr.Right.(*sqlparser.SQLVal); ok && expr.Operator == "=" {
				value := string(sqlVal.Val)

				ctx.AddConstant(fieldName, value)
			}
		}
	case *sqlparser.AndExpr:
		e.GetConstants(expr.Left, ctx)
		e.GetConstants(expr.Right, ctx)
	case *sqlparser.OrExpr:
		e.GetConstants(expr.Left, ctx)
		e.GetConstants(expr.Right, ctx)
	case *sqlparser.ParenExpr:
		e.GetConstants(expr.Expr, ctx)
	}
}
