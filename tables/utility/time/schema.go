package time_info

import (
	"github.com/scrymastic/goosquery/tables/specs"
)

var TableName = "time"
var Description = "Track current date and time in UTC."
var Schema = specs.Schema{
	specs.Column{Name: "weekday", Type: "string", Description: "Current weekday in UTC"},
	specs.Column{Name: "year", Type: "int32", Description: "Current year in UTC"},
	specs.Column{Name: "month", Type: "int32", Description: "Current month in UTC"},
	specs.Column{Name: "day", Type: "int32", Description: "Current day in UTC"},
	specs.Column{Name: "hour", Type: "int32", Description: "Current hour in UTC"},
	specs.Column{Name: "minutes", Type: "int32", Description: "Current minutes in UTC"},
	specs.Column{Name: "seconds", Type: "int32", Description: "Current seconds in UTC"},
	specs.Column{Name: "timezone", Type: "string", Description: "Timezone for reported time (hardcoded to UTC)"},
	specs.Column{Name: "local_timezone", Type: "string", Description: "Current local timezone in of the system"},
	specs.Column{Name: "unix_time", Type: "int32", Description: "Current UNIX time in UTC"},
	specs.Column{Name: "timestamp", Type: "string", Description: "Current timestamp (log format) in UTC"},
	specs.Column{Name: "datetime", Type: "string", Description: "Current date and time (ISO format) in UTC"},
	specs.Column{Name: "iso_8601", Type: "string", Description: "Current time (ISO format) in UTC"},

	specs.Column{Name: "win_timestamp", Type: "int64", Description: "Timestamp value in 100 nanosecond units"},
}
