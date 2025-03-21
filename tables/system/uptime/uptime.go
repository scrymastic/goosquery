package uptime

import (
	"fmt"
	"syscall"

	"github.com/scrymastic/goosquery/sql/context"
	"github.com/scrymastic/goosquery/tables/specs"
)

var (
	kernel32       = syscall.NewLazyDLL("kernel32.dll")
	getTickCount64 = kernel32.NewProc("GetTickCount64")
)

// GenUptime returns the system uptime information
func GenUptime(ctx context.Context) ([]map[string]interface{}, error) {
	// GetTickCount64 returns the number of milliseconds since system startup
	ret, _, err := getTickCount64.Call()
	if ret == 0 {
		return nil, fmt.Errorf("failed to get uptime: %v", err)
	}

	totalSeconds := uint64(ret) / 1000

	result := specs.Init(ctx, Schema)

	if ctx.IsColumnUsed("days") {
		result["days"] = uint32(totalSeconds / 86400)
	}

	if ctx.IsColumnUsed("hours") {
		result["hours"] = uint16(totalSeconds % 86400 / 3600)
	}

	if ctx.IsColumnUsed("minutes") {
		result["minutes"] = uint16(totalSeconds % 3600 / 60)
	}

	if ctx.IsColumnUsed("seconds") {
		result["seconds"] = uint16(totalSeconds % 60)
	}

	if ctx.IsColumnUsed("total_seconds") {
		result["total_seconds"] = totalSeconds
	}

	return []map[string]interface{}{result}, nil
}
