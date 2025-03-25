package process_memory_map

import (
	"fmt"
	"strconv"
	"strings"
	"unsafe"

	"github.com/scrymastic/goosquery/sql/result"
	"github.com/scrymastic/goosquery/sql/sqlctx"
	"golang.org/x/sys/windows"
)

type ProcessMemoryMap struct {
	PID         int64  `json:"pid"`
	Start       int64  `json:"start"`
	End         int64  `json:"end"`
	Permissions string `json:"permissions"`
	Offset      int64  `json:"offset"`
	Device      string `json:"device"`
	Inode       int64  `json:"inode"`
	Path        string `json:"path"`
	Pseudo      int32  `json:"pseudo"`
}

var memoryConstants = map[uint32]string{
	windows.PAGE_EXECUTE:           "PAGE_EXECUTE",
	windows.PAGE_EXECUTE_READ:      "PAGE_EXECUTE_READ",
	windows.PAGE_EXECUTE_READWRITE: "PAGE_EXECUTE_READWRITE",
	windows.PAGE_EXECUTE_WRITECOPY: "PAGE_EXECUTE_WRITECOPY",
	windows.PAGE_NOACCESS:          "PAGE_NOACCESS",
	windows.PAGE_READONLY:          "PAGE_READONLY",
	windows.PAGE_READWRITE:         "PAGE_READWRITE",
	windows.PAGE_WRITECOPY:         "PAGE_WRITECOPY",
	windows.PAGE_GUARD:             "PAGE_GUARD",
	windows.PAGE_NOCACHE:           "PAGE_NOCACHE",
	windows.PAGE_WRITECOMBINE:      "PAGE_WRITECOMBINE",
}

func formatPermissions(perm uint32) string {
	var perms []string
	for k, v := range memoryConstants {
		if perm&k != 0 {
			perms = append(perms, v)
		}
	}
	return strings.Join(perms, "|")
}

func GenProcessMemoryMap(ctx *sqlctx.Context) (*result.Results, error) {
	pids := ctx.GetConstants("pid")
	if len(pids) == 0 {
		return nil, fmt.Errorf("pid is not set")
	}
	pid64, err := strconv.ParseUint(pids[0], 10, 32)
	if err != nil {
		return nil, fmt.Errorf("invalid pid: %v", err)
	}
	pid := uint32(pid64)
	proc, err := windows.OpenProcess(windows.PROCESS_QUERY_INFORMATION, false, pid)
	if err != nil {
		return nil, fmt.Errorf("failed to open process: %v", err)
	}
	defer windows.CloseHandle(proc)

	memoryMaps := result.NewQueryResult()

	modSnap, err := windows.CreateToolhelp32Snapshot(windows.TH32CS_SNAPMODULE|windows.TH32CS_SNAPMODULE32, pid)
	if err != nil {
		return nil, fmt.Errorf("failed to create module snapshot: %v", err)
	}

	defer windows.CloseHandle(modSnap)

	var me windows.ModuleEntry32
	var memInfo windows.MemoryBasicInformation
	me.Size = uint32(unsafe.Sizeof(me))
	ret := windows.Module32First(modSnap, &me)
	for ret == nil {
		// Get the memory map for the module
		for addr := uintptr(me.ModBaseAddr); windows.VirtualQueryEx(proc, addr, &memInfo, unsafe.Sizeof(memInfo)) == nil &&
			addr < me.ModBaseAddr+uintptr(me.ModBaseSize); addr += uintptr(memInfo.RegionSize) {
			// Get the path for the module
			memMap := result.NewResult(ctx, Schema)
			memMap.Set("pid", int64(pid))
			memMap.Set("start", int64(memInfo.BaseAddress))
			memMap.Set("end", int64(memInfo.BaseAddress+memInfo.RegionSize))
			memMap.Set("permissions", formatPermissions(memInfo.Protect))
			memMap.Set("offset", int64(memInfo.AllocationBase))
			memMap.Set("device", "")
			memMap.Set("inode", 0)
			memMap.Set("path", windows.UTF16PtrToString(&me.ExePath[0]))
			memMap.Set("pseudo", 0)
			memoryMaps.AppendResult(*memMap)
		}
		ret = windows.Module32Next(modSnap, &me)
	}

	return memoryMaps, nil
}
