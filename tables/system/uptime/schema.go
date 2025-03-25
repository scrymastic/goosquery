package uptime

import (
	"github.com/scrymastic/goosquery/sql/result"
)

var TableName = "uptime"
var Description = "Track time passed since last boot. Some systems track this as calendar time, some as runtime."
var Schema = result.Schema{
	result.Column{Name: "days", Type: "INTEGER", Description: "Days of uptime"},
	result.Column{Name: "hours", Type: "INTEGER", Description: "Hours of uptime"},
	result.Column{Name: "minutes", Type: "INTEGER", Description: "Minutes of uptime"},
	result.Column{Name: "seconds", Type: "INTEGER", Description: "Seconds of uptime"},
	result.Column{Name: "total_seconds", Type: "BIGINT", Description: "Total uptime seconds"},
}
