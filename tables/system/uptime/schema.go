package uptime

import (
	"github.com/scrymastic/goosquery/tables/specs"
)

var TableName = "uptime"
var Description = "Track time passed since last boot. Some systems track this as calendar time, some as runtime."
var Schema = specs.Schema{
	specs.Column{Name: "days", Type: "INTEGER", Description: "Days of uptime"},
	specs.Column{Name: "hours", Type: "INTEGER", Description: "Hours of uptime"},
	specs.Column{Name: "minutes", Type: "INTEGER", Description: "Minutes of uptime"},
	specs.Column{Name: "seconds", Type: "INTEGER", Description: "Seconds of uptime"},
	specs.Column{Name: "total_seconds", Type: "BIGINT", Description: "Total uptime seconds"},
}
