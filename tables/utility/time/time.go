package time_info

import (
	"fmt"
	"time"

	"github.com/scrymastic/goosquery/sql/result"
	"github.com/scrymastic/goosquery/sql/sqlctx"
)

const (
	WINDOWS_TICK      = 100         // nanoseconds
	SEC_TO_UNIX_EPOCH = 11644473600 // seconds between 1601 and 1970
)

// GenTime returns the current system time information
func GenTime(ctx *sqlctx.Context) (*result.Results, error) {
	utcNow := time.Now().UTC()

	entry := result.NewResult(ctx, Schema)
	entry.Set("weekday", utcNow.Weekday().String())
	entry.Set("year", int32(utcNow.Year()))
	entry.Set("month", int32(utcNow.Month()))
	entry.Set("day", int32(utcNow.Day()))
	entry.Set("hour", int32(utcNow.Hour()))
	entry.Set("minutes", int32(utcNow.Minute()))
	entry.Set("seconds", int32(utcNow.Second()))
	entry.Set("timezone", "UTC")

	_, offset := time.Now().Zone()
	entry.Set("local_timezone", fmt.Sprintf("UTC%+d", offset/3600))

	entry.Set("unix_time", int32(utcNow.Unix()))
	entry.Set("timestamp", utcNow.Format("Mon Jan _2 15:04:05 2006 UTC"))
	entry.Set("datetime", utcNow.Format("2006-01-02 15:04:05"))
	entry.Set("iso_8601", utcNow.Format("2006-01-02T15:04:05Z"))
	// Convert Unix timestamp to Windows timestamp (100-nanosecond intervals since Jan 1, 1601)
	winTime := (utcNow.Unix() + SEC_TO_UNIX_EPOCH) * 1e7
	entry.Set("win_timestamp", int64(winTime))

	results := result.NewQueryResult()
	results.AppendResult(*entry)
	return results, nil
}
