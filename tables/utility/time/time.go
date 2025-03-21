package time_info

import (
	"fmt"
	"runtime"
	"time"

	"github.com/scrymastic/goosquery/sql/context"
	"github.com/scrymastic/goosquery/tables/specs"
)

// Time represents the system time information matching osquery's time table schema
type Time struct {
	Weekday       string `json:"weekday"`
	Year          int32  `json:"year"`
	Month         int32  `json:"month"`
	Day           int32  `json:"day"`
	Hour          int32  `json:"hour"`
	Minutes       int32  `json:"minutes"`
	Seconds       int32  `json:"seconds"`
	Timezone      string `json:"timezone"`
	LocalTimezone string `json:"local_timezone"`
	UnixTime      int64  `json:"unix_time"`
	Timestamp     string `json:"timestamp"`
	Datetime      string `json:"datetime"`
	ISO8601       string `json:"iso_8601"`
	WinTimestamp  *int64 `json:"win_timestamp,omitempty"`
}

const (
	WINDOWS_TICK      = 100         // nanoseconds
	SEC_TO_UNIX_EPOCH = 11644473600 // seconds between 1601 and 1970
)

// GenTime returns the current system time information
func GenTime(ctx context.Context) ([]map[string]interface{}, error) {
	utcNow := time.Now().UTC()

	entry := specs.Init(ctx, Schema)

	if ctx.IsColumnUsed("weekday") {
		entry["weekday"] = utcNow.Weekday().String()
	}

	if ctx.IsColumnUsed("year") {
		entry["year"] = int32(utcNow.Year())
	}

	if ctx.IsColumnUsed("month") {
		entry["month"] = int32(utcNow.Month())
	}

	if ctx.IsColumnUsed("day") {
		entry["day"] = int32(utcNow.Day())
	}

	if ctx.IsColumnUsed("hour") {
		entry["hour"] = int32(utcNow.Hour())
	}

	if ctx.IsColumnUsed("minutes") {
		entry["minutes"] = int32(utcNow.Minute())
	}

	if ctx.IsColumnUsed("seconds") {
		entry["seconds"] = int32(utcNow.Second())
	}

	if ctx.IsColumnUsed("timezone") {
		entry["timezone"] = "UTC"
	}

	if ctx.IsColumnUsed("local_timezone") {
		_, offset := time.Now().Zone()
		entry["local_timezone"] = fmt.Sprintf("UTC%+d", offset/3600)
	}

	if ctx.IsColumnUsed("unix_time") {
		entry["unix_time"] = int32(utcNow.Unix())
	}

	if ctx.IsColumnUsed("timestamp") {
		entry["timestamp"] = utcNow.Format("Mon Jan _2 15:04:05 2006 UTC")
	}

	if ctx.IsColumnUsed("datetime") {
		entry["datetime"] = utcNow.Format("2006-01-02 15:04:05")
	}

	if ctx.IsColumnUsed("iso_8601") {
		entry["iso_8601"] = utcNow.Format("2006-01-02T15:04:05Z")
	}

	// Windows-specific timestamp
	if runtime.GOOS == "windows" && ctx.IsColumnUsed("win_timestamp") {
		// Convert Unix timestamp to Windows timestamp (100-nanosecond intervals since Jan 1, 1601)
		winTime := (utcNow.Unix() + SEC_TO_UNIX_EPOCH) * 1e7
		entry["win_timestamp"] = int64(winTime)
	}

	return []map[string]interface{}{entry}, nil
}
