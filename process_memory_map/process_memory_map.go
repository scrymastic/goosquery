package process_memory_map

import (
	"fmt"
	"strings"
	"unsafe"

	"golang.org/x/sys/windows"
)

type ProcessMemoryMap struct {
	PID         uint32 `json:"pid"`
	Start       uint64 `json:"start"`
	End         uint64 `json:"end"`
	Permissions string `json:"permissions"`
	Offset      uint64 `json:"offset"`
	Device      string `json:"device"`
	Inode       uint64 `json:"inode"`
	Path        string `json:"path"`
	Pseudo      bool   `json:"pseudo"`
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
				PID:         pid,
				Start:       uint64(memInfo.BaseAddress),
				End:         uint64(memInfo.BaseAddress + memInfo.RegionSize),
				Permissions: formatPermissions(memInfo.Protect),
				Offset:      uint64(memInfo.AllocationBase),
				Device:      "",
				Inode:       0,
				Path:        windows.UTF16PtrToString(&me.ExePath[0]),
				Pseudo:      false,
			}
			memoryMaps = append(memoryMaps, memMap)
		}
		ret = windows.Module32Next(modSnap, &me)
	}

	return memoryMaps, nil
}
