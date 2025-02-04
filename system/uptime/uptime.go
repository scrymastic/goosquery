package uptime

import (
	"fmt"
	"syscall"
)

var (
	kernel32       = syscall.NewLazyDLL("kernel32.dll")
	getTickCount64 = kernel32.NewProc("GetTickCount64")
)

type Uptime struct {
	Days         uint32 `json:"days"`
	Hours        uint16 `json:"hours"`
	Minutes      uint16 `json:"minutes"`
	Seconds      uint16 `json:"seconds"`
	TotalSeconds uint64 `json:"total_seconds"`
}

func GenUptime() ([]Uptime, error) {
	// GetTickCount64 returns the number of milliseconds since system startup
	ret, _, err := getTickCount64.Call()
	if ret == 0 {
		return nil, fmt.Errorf("failed to get uptime: %v", err)
	}

	totalSeconds := uint64(ret) / 1000

	return []Uptime{
		{
			Days:         uint32(totalSeconds / 86400),
			Hours:        uint16(totalSeconds % 86400 / 3600),
			Minutes:      uint16(totalSeconds % 3600 / 60),
			Seconds:      uint16(totalSeconds % 60),
			TotalSeconds: totalSeconds,
		},
	}, nil
}
