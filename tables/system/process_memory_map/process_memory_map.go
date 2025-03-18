package process_memory_map

import (
	"fmt"
	"strings"
	"unsafe"

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

func GenProcessMemoryMap(pid uint32) ([]ProcessMemoryMap, error) {
	proc, err := windows.OpenProcess(windows.PROCESS_QUERY_INFORMATION, false, pid)
	if err != nil {
		return nil, fmt.Errorf("failed to open process: %v", err)
	}
	defer windows.CloseHandle(proc)

	var memoryMaps []ProcessMemoryMap

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
			memMap := ProcessMemoryMap{
				PID:         int64(pid),
				Start:       int64(memInfo.BaseAddress),
				End:         int64(memInfo.BaseAddress + memInfo.RegionSize),
				Permissions: formatPermissions(memInfo.Protect),
				Offset:      int64(memInfo.AllocationBase),
				Device:      "",
				Inode:       0,
				Path:        windows.UTF16PtrToString(&me.ExePath[0]),
				Pseudo:      0,
			}
			memoryMaps = append(memoryMaps, memMap)
		}
		ret = windows.Module32Next(modSnap, &me)
	}

	return memoryMaps, nil
}
