package uptime

import (
	"github.com/scrymastic/goosquery/tables/specs"
)

var TableName = "uptime"
var Description = "Track time passed since last boot. Some systems track this as calendar time, some as runtime."
var Schema = specs.Schema{
	specs.Column{Name: "days", Type: "int32", Description: "Days of uptime"},
	specs.Column{Name: "hours", Type: "int32", Description: "Hours of uptime"},
	specs.Column{Name: "minutes", Type: "int32", Description: "Minutes of uptime"},
	specs.Column{Name: "seconds", Type: "int32", Description: "Seconds of uptime"},
	specs.Column{Name: "total_seconds", Type: "int64", Description: "Total uptime seconds"},
}
