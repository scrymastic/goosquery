package uptime

import (
	"fmt"
	"syscall"

	"github.com/scrymastic/goosquery/sql/result"
	"github.com/scrymastic/goosquery/sql/sqlctx"
)

var (
	kernel32       = syscall.NewLazyDLL("kernel32.dll")
	getTickCount64 = kernel32.NewProc("GetTickCount64")
)

// GenUptime returns the system uptime information
func GenUptime(ctx *sqlctx.Context) (*result.Results, error) {
	// GetTickCount64 returns the number of milliseconds since system startup
	ret, _, err := getTickCount64.Call()
	if ret == 0 {
		return nil, fmt.Errorf("failed to get uptime: %v", err)
	}

	totalSeconds := uint64(ret) / 1000

	uptime := result.NewResult(ctx, Schema)

	uptime.Set("days", uint32(totalSeconds/86400))
	uptime.Set("hours", uint16(totalSeconds%86400/3600))
	uptime.Set("minutes", uint16(totalSeconds%3600/60))
	uptime.Set("seconds", uint16(totalSeconds%60))
	uptime.Set("total_seconds", totalSeconds)

	results := result.NewQueryResult()
	results.AppendResult(*uptime)
	return results, nil
}
