package aggregation

import (
	"fmt"
	"strings"

	"github.com/blastrain/vitess-sqlparser/sqlparser"
)

// Type represents the type of aggregation function
type Type int

const (
	None Type = iota
	Count
	Sum
	Avg
	Min
	Max
)

// Info stores information about an aggregation function
type Info struct {
	Type       Type
	Column     string
	Alias      string
	IsDistinct bool
}

// GetType determines the type of aggregation from an expression
func GetType(expr sqlparser.Expr) (Type, string, bool) {
	funcExpr, ok := expr.(*sqlparser.FuncExpr)
	if !ok {
		return None, "", false
	}

	funcName := strings.ToUpper(funcExpr.Name.String())
	isDistinct := funcExpr.Distinct

	// Check if any argument exists
	if len(funcExpr.Exprs) == 0 && funcName != "COUNT" {
		return None, "", false
	}

	// Get column name for aggregation
	var columnName string
	if funcName == "COUNT" && len(funcExpr.Exprs) == 0 {
		// COUNT(*) case
		columnName = "*"
	} else if len(funcExpr.Exprs) > 0 {
		// Normal column aggregation case
		aliasedExpr, ok := funcExpr.Exprs[0].(*sqlparser.AliasedExpr)
		if !ok {
			return None, "", false
		}

		colName, ok := aliasedExpr.Expr.(*sqlparser.ColName)
		if !ok {
			return None, "", false
		}
		columnName = colName.Name.String()
	}

	switch funcName {
	case "COUNT":
		return Count, columnName, isDistinct
	case "SUM":
		return Sum, columnName, isDistinct
	case "AVG":
		return Avg, columnName, isDistinct
	case "MIN":
		return Min, columnName, isDistinct
	case "MAX":
		return Max, columnName, isDistinct
	default:
		return None, "", false
	}
}

// GenerateDefaultAlias creates a default alias for an aggregation function
func GenerateDefaultAlias(aggType Type, column string) string {
	switch aggType {
	case Count:
		return fmt.Sprintf("COUNT(%s)", column)
	case Sum:
		return fmt.Sprintf("SUM(%s)", column)
	case Avg:
		return fmt.Sprintf("AVG(%s)", column)
	case Min:
		return fmt.Sprintf("MIN(%s)", column)
	case Max:
		return fmt.Sprintf("MAX(%s)", column)
	default:
		return column
	}
}
