package time_info

import (
	"github.com/scrymastic/goosquery/sql/result"
)

var TableName = "time"
var Description = "Track current date and time in UTC."
var Schema = result.Schema{
	result.Column{Name: "weekday", Type: "TEXT", Description: "Current weekday in UTC"},
	result.Column{Name: "year", Type: "INTEGER", Description: "Current year in UTC"},
	result.Column{Name: "month", Type: "INTEGER", Description: "Current month in UTC"},
	result.Column{Name: "day", Type: "INTEGER", Description: "Current day in UTC"},
	result.Column{Name: "hour", Type: "INTEGER", Description: "Current hour in UTC"},
	result.Column{Name: "minutes", Type: "INTEGER", Description: "Current minutes in UTC"},
	result.Column{Name: "seconds", Type: "INTEGER", Description: "Current seconds in UTC"},
	result.Column{Name: "timezone", Type: "TEXT", Description: "Timezone for reported time (hardcoded to UTC)"},
	result.Column{Name: "local_timezone", Type: "TEXT", Description: "Current local timezone in of the system"},
	result.Column{Name: "unix_time", Type: "INTEGER", Description: "Current UNIX time in UTC"},
	result.Column{Name: "timestamp", Type: "TEXT", Description: "Current timestamp (log format) in UTC"},
	result.Column{Name: "datetime", Type: "TEXT", Description: "Current date and time (ISO format) in UTC"},
	result.Column{Name: "iso_8601", Type: "TEXT", Description: "Current time (ISO format) in UTC"},

	result.Column{Name: "win_timestamp", Type: "BIGINT", Description: "Timestamp value in 100 nanosecond units"},
}
