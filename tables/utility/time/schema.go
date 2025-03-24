package time_info

import (
	"github.com/scrymastic/goosquery/tables/specs"
)

var TableName = "time"
var Description = "Track current date and time in UTC."
var Schema = specs.Schema{
	specs.Column{Name: "weekday", Type: "TEXT", Description: "Current weekday in UTC"},
	specs.Column{Name: "year", Type: "INTEGER", Description: "Current year in UTC"},
	specs.Column{Name: "month", Type: "INTEGER", Description: "Current month in UTC"},
	specs.Column{Name: "day", Type: "INTEGER", Description: "Current day in UTC"},
	specs.Column{Name: "hour", Type: "INTEGER", Description: "Current hour in UTC"},
	specs.Column{Name: "minutes", Type: "INTEGER", Description: "Current minutes in UTC"},
	specs.Column{Name: "seconds", Type: "INTEGER", Description: "Current seconds in UTC"},
	specs.Column{Name: "timezone", Type: "TEXT", Description: "Timezone for reported time (hardcoded to UTC)"},
	specs.Column{Name: "local_timezone", Type: "TEXT", Description: "Current local timezone in of the system"},
	specs.Column{Name: "unix_time", Type: "INTEGER", Description: "Current UNIX time in UTC"},
	specs.Column{Name: "timestamp", Type: "TEXT", Description: "Current timestamp (log format) in UTC"},
	specs.Column{Name: "datetime", Type: "TEXT", Description: "Current date and time (ISO format) in UTC"},
	specs.Column{Name: "iso_8601", Type: "TEXT", Description: "Current time (ISO format) in UTC"},

	specs.Column{Name: "win_timestamp", Type: "BIGINT", Description: "Timestamp value in 100 nanosecond units"},
}
