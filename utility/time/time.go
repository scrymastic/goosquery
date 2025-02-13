package time_info

import (
	"time"

	"golang.org/x/sys/windows"
)

// Time represents the system time information matching osquery's time table schema
type Time struct {
	Weekday       string `json:"weekday"`
	Year          int    `json:"year"`
	Month         int    `json:"month"`
	Day           int    `json:"day"`
	Hour          int    `json:"hour"`
	Minutes       int    `json:"minutes"`
	Seconds       int    `json:"seconds"`
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

// GetTime returns the current system time information
func GenTime() ([]Time, error) {
	utcNow := time.Now().UTC()

	result := Time{
		Weekday:   utcNow.Weekday().String(),
		Year:      utcNow.Year(),
		Month:     int(utcNow.Month()),
		Day:       utcNow.Day(),
		Hour:      utcNow.Hour(),
		Minutes:   utcNow.Minute(),
		Seconds:   utcNow.Second(),
		Timezone:  "UTC",
		UnixTime:  utcNow.Unix(),
		Timestamp: utcNow.Format("Mon Jan 2 15:04:05 2006 UTC"),
		Datetime:  utcNow.Format("2006-01-02T15:04:05Z"),
		ISO8601:   utcNow.Format("2006-01-02T15:04:05Z"),
	}

	// Windows timestamp (difference between Unix epoch and Windows epoch in 100ns intervals)
	winTimestamp := (utcNow.UnixNano() / WINDOWS_TICK) + SEC_TO_UNIX_EPOCH*10000000
	result.WinTimestamp = &winTimestamp

	// Get local timezone
	var timezoneInfo windows.Timezoneinformation
	_, err := windows.GetTimeZoneInformation(&timezoneInfo)
	if err != nil {
		return nil, err
	}
	result.LocalTimezone = windows.UTF16ToString(timezoneInfo.StandardName[:])

	return []Time{result}, nil
}
