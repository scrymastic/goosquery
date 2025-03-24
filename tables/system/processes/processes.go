package processes

import (
	"fmt"
	"log"
	"os"
	"syscall"
	"unsafe"

	"github.com/scrymastic/goosquery/sql/context"
	"github.com/scrymastic/goosquery/tables/specs"
	"golang.org/x/sys/windows"
)

// Process-related constants
const (
	STILL_ACTIVE = 259
)

// Windows-specific structs needed for API calls
type PROCESS_MEMORY_COUNTERS_EX struct {
	cb                         uint32
	PageFaultCount             uint32
	PeakWorkingSetSize         uint64
	WorkingSetSize             uint64
	QuotaPeakPagedPoolUsage    uint64
	QuotaPagedPoolUsage        uint64
	QuotaPeakNonPagedPoolUsage uint64
	QuotaNonPagedPoolUsage     uint64
	PagefileUsage              uint64
	PeakPagefileUsage          uint64
	PrivateUsage               uint64
}

type IO_COUNTERS struct {
	ReadOperationCount  uint64
	WriteOperationCount uint64
	OtherOperationCount uint64
	ReadTransferCount   uint64
	WriteTransferCount  uint64
	OtherTransferCount  uint64
}

type TOKEN_ELEVATION struct {
	TokenIsElevated uint32
}

var (
	kernel32 = windows.NewLazySystemDLL("kernel32.dll")
)

func getProcessCommandLine(proc windows.Handle) (string, error) {
	var cmdlineLen uint32
	if err := windows.NtQueryInformationProcess(
		proc,
		windows.ProcessCommandLineInformation,
		nil,
		0,
		&cmdlineLen); err != windows.STATUS_INFO_LENGTH_MISMATCH &&
		err != windows.STATUS_BUFFER_TOO_SMALL &&
		err != windows.STATUS_BUFFER_OVERFLOW {
		return "", fmt.Errorf("failed to get command line length: %v", err)
	}

	cmdlineBuf := make([]byte, cmdlineLen)
	if err := windows.NtQueryInformationProcess(
		proc,
		windows.ProcessCommandLineInformation,
		unsafe.Pointer(&cmdlineBuf[0]),
		cmdlineLen,
		&cmdlineLen); err != nil {
		return "", fmt.Errorf("failed to get command line: %v", err)
	}

	ustr := (*windows.NTUnicodeString)(unsafe.Pointer(&cmdlineBuf[0]))
	return windows.UTF16PtrToString(ustr.Buffer), nil
}

func getUserProcessParameters(proc windows.Handle) (*windows.RTL_USER_PROCESS_PARAMETERS, error) {
	var pbi windows.PROCESS_BASIC_INFORMATION
	var returnLength uint32
	if err := windows.NtQueryInformationProcess(
		proc,
		windows.ProcessBasicInformation,
		unsafe.Pointer(&pbi),
		uint32(unsafe.Sizeof(pbi)),
		&returnLength); err != nil {
		return nil, fmt.Errorf("failed to query process information: %v", err)
	}

	var peb windows.PEB
	var bytesRead uintptr
	if err := windows.ReadProcessMemory(
		proc,
		uintptr(unsafe.Pointer(pbi.PebBaseAddress)),
		(*byte)(unsafe.Pointer(&peb)),
		uintptr(unsafe.Sizeof(peb)),
		&bytesRead); err != nil {
		return nil, fmt.Errorf("failed to read process memory: %v", err)
	}

	var params windows.RTL_USER_PROCESS_PARAMETERS
	if err := windows.ReadProcessMemory(
		proc,
		uintptr(unsafe.Pointer(peb.ProcessParameters)),
		(*byte)(unsafe.Pointer(&params)),
		uintptr(unsafe.Sizeof(params)),
		&bytesRead); err != nil {
		return nil, fmt.Errorf("failed to read process memory: %v", err)
	}

	return &params, nil
}

// Regardless of the Windows version, the CWD of a process is only possible to
// retrieve by reading it from the process's PEB structure.
func getProcessWorkingDirectory(proc windows.Handle) (string, error) {
	params, err := getUserProcessParameters(proc)
	if err != nil {
		return "", fmt.Errorf("failed to get process parameters: %v", err)
	}

	var bytesRead uintptr
	var cwdBuf [windows.MAX_PATH]uint16
	if err := windows.ReadProcessMemory(
		proc,
		uintptr(unsafe.Pointer(params.CurrentDirectory.DosPath.Buffer)),
		(*byte)(unsafe.Pointer(&cwdBuf[0])),
		uintptr(params.CurrentDirectory.DosPath.Length),
		&bytesRead); err != nil {
		return "", fmt.Errorf("failed to read process memory: %v", err)
	}

	return windows.UTF16PtrToString(&cwdBuf[0]), nil
}

// getProcessRssInfo retrieves memory-related information about the calling process.
func getProcessRssInfo(proc windows.Handle) (*PROCESS_MEMORY_COUNTERS_EX, error) {
	// https://docs.microsoft.com/en-us/windows/win32/api/psapi/nf-psapi-getprocessmemoryinfo
	modKernel32 := windows.NewLazySystemDLL("kernel32.dll")
	procGetProcessMemoryInfo := modKernel32.NewProc("K32GetProcessMemoryInfo")

	// If GetProcessMemoryInfo is not available, try with psapi.dll
	if procGetProcessMemoryInfo == nil {
		modPsapi := windows.NewLazySystemDLL("psapi.dll")
		procGetProcessMemoryInfo = modPsapi.NewProc("GetProcessMemoryInfo")
	}

	// If GetProcessMemoryInfo is not available, return an error
	if procGetProcessMemoryInfo == nil {
		return nil, fmt.Errorf("failed to get GetProcessMemoryInfo address")
	}

	var pmcEx PROCESS_MEMORY_COUNTERS_EX
	ret, _, _ := procGetProcessMemoryInfo.Call(
		uintptr(proc),
		uintptr(unsafe.Pointer(&pmcEx)), // _PROCESS_MEMORY_COUNTERS_EX
		uintptr(unsafe.Sizeof(pmcEx)),
	)
	if ret == 0 {
		return nil, fmt.Errorf("failed to get process memory info: %v", syscall.GetLastError())
	}

	return &pmcEx, nil
}

func isProcessElevated(proc windows.Handle) (bool, error) {
	var token windows.Token
	if err := windows.OpenProcessToken(
		proc,
		windows.TOKEN_READ,
		&token); err != nil {
		return false, fmt.Errorf("failed to open process token: %v", err)
	}
	defer token.Close()

	if token == 0 {
		return false, fmt.Errorf("failed to open process token")
	}

	// Check if token is elevated
	var elevation TOKEN_ELEVATION
	var returnLength uint32
	if err := windows.GetTokenInformation(
		token,
		windows.TokenElevation,
		(*byte)(unsafe.Pointer(&elevation)),
		uint32(unsafe.Sizeof(elevation)),
		&returnLength); err != nil || returnLength != uint32(unsafe.Sizeof(elevation)) {
		return false, fmt.Errorf("failed to get token elevation: %v", err)
	}

	return elevation.TokenIsElevated != 0, nil
}

func getProcessElapsedTime(proc windows.Handle) (uint64, error) {
	var creationTime, exitTime, kernelTime, userTime windows.Filetime
	if err := windows.GetProcessTimes(
		proc,
		&creationTime,
		&exitTime,
		&kernelTime,
		&userTime); err != nil {
		return 0, fmt.Errorf("failed to get process times: %v", err)
	}

	// Calculate elapsed time in seconds
	var currentTime windows.Filetime
	windows.GetSystemTimeAsFileTime(&currentTime)
	elapsedTime := (currentTime.Nanoseconds() - creationTime.Nanoseconds()) / 1e9

	return uint64(elapsedTime), nil
}

// Collect detailed information about a process and populate the procInfo map
func getProcessDetails(ctx context.Context, procInfo map[string]interface{}) error {
	processID := uint32(procInfo["pid"].(int64))

	// Try to open process with PROCESS_ALL_ACCESS first
	procHandle, err := windows.OpenProcess(windows.PROCESS_ALL_ACCESS, false, processID)
	if err != nil {
		// If it fails, try with PROCESS_QUERY_LIMITED_INFORMATION
		procHandle, err = windows.OpenProcess(windows.PROCESS_QUERY_LIMITED_INFORMATION, false, processID)
	}

	if err != nil && processID != 0 {
		return fmt.Errorf("failed to open process %d: %v", processID, err)
	}

	if ctx.IsColumnUsed("path") {
		// Get process path
		var pathBuf [windows.MAX_PATH]uint16
		pathLen := uint32(len(pathBuf))
		if err := windows.QueryFullProcessImageName(
			procHandle,
			0,
			&pathBuf[0],
			&pathLen); err == nil {
			procInfo["path"] = windows.UTF16ToString(pathBuf[:pathLen])
		}
	}

	if ctx.IsColumnUsed("cmdline") {
		// Get command line
		cmdline, err := getProcessCommandLine(procHandle)
		if err == nil {
			procInfo["cmdline"] = cmdline
		}
	}

	if ctx.IsColumnUsed("state") {
		// Get process state
		var exitCode uint32
		if err := windows.GetExitCodeProcess(procHandle, &exitCode); err == nil {
			if exitCode == STILL_ACTIVE {
				procInfo["state"] = "STILL_ACTIVE"
			} else {
				procInfo["state"] = "EXITED"
			}
		}
	}

	if ctx.IsColumnUsed("cwd") {
		// Get working directory
		cwd, err := getProcessWorkingDirectory(procHandle)
		if err == nil {
			procInfo["cwd"] = cwd
		}
	}

	if ctx.IsColumnUsed("root") {
		// Get process root directory
		if cwd, ok := procInfo["cwd"].(string); ok {
			procInfo["root"] = cwd
		}
	}

	if ctx.IsColumnUsed("on_disk") {
		// Check if process is on disk
		if path, ok := procInfo["path"].(string); ok {
			if _, err := os.Stat(path); err == nil {
				procInfo["on_disk"] = int32(1)
			} else {
				procInfo["on_disk"] = int32(0)
			}
		}
	}

	if ctx.IsAnyOfColumnsUsed([]string{"wired_size", "resident_size", "total_size"}) {
		// Get process memory info
		pmc, err := getProcessRssInfo(procHandle)
		if err == nil {
			if ctx.IsColumnUsed("wired_size") {
				procInfo["wired_size"] = int64(pmc.QuotaNonPagedPoolUsage)
			}
			if ctx.IsColumnUsed("resident_size") {
				procInfo["resident_size"] = int64(pmc.WorkingSetSize)
			}
			if ctx.IsColumnUsed("total_size") {
				procInfo["total_size"] = int64(pmc.PrivateUsage)
			}
		}
	}

	if ctx.IsAnyOfColumnsUsed([]string{"disk_bytes_read", "disk_bytes_written"}) {
		// Get disk bytes read and written
		var ioCounters IO_COUNTERS
		processIoCounters := kernel32.NewProc("GetProcessIoCounters")
		if ret, _, _ := processIoCounters.Call(
			uintptr(procHandle),
			uintptr(unsafe.Pointer(&ioCounters))); ret != 0 {
			if ctx.IsColumnUsed("disk_bytes_read") {
				procInfo["disk_bytes_read"] = int64(ioCounters.ReadTransferCount)
			}
			if ctx.IsColumnUsed("disk_bytes_written") {
				procInfo["disk_bytes_written"] = int64(ioCounters.WriteTransferCount)
			}
		}
	}

	if ctx.IsColumnUsed("nice") {
		// Get nice value
		nice, err := windows.GetPriorityClass(procHandle)
		if err == nil {
			procInfo["nice"] = int32(nice)
		}
	}

	if ctx.IsColumnUsed("elevated_token") {
		// Get elevated token
		elevated, err := isProcessElevated(procHandle)
		if err == nil {
			if elevated {
				procInfo["elevated_token"] = int32(1)
			} else {
				procInfo["elevated_token"] = int32(0)
			}
		}
	}

	if ctx.IsColumnUsed("elapsed_time") {
		// Get elapsed time
		elapsedTime, err := getProcessElapsedTime(procHandle)
		if err == nil {
			procInfo["elapsed_time"] = int64(elapsedTime)
		}
	}

	if ctx.IsColumnUsed("handle_count") {
		// Get handle count
		var handleCount uint32
		procGetProcessHandleCount := kernel32.NewProc("GetProcessHandleCount")
		if ret, _, _ := procGetProcessHandleCount.Call(uintptr(procHandle), uintptr(unsafe.Pointer(&handleCount))); ret != 0 {
			procInfo["handle_count"] = int64(handleCount)
		}
	}

	windows.CloseHandle(procHandle)

	return nil
}

// GenProcesses returns information about all processes as a slice of maps
func GenProcesses(ctx context.Context) ([]map[string]interface{}, error) {
	procs, err := windows.CreateToolhelp32Snapshot(windows.TH32CS_SNAPALL, 0)
	if err != nil {
		return nil, fmt.Errorf("failed to create snapshot: %v", err)
	}
	defer windows.CloseHandle(procs)

	var processes []map[string]interface{}
	var pe32 windows.ProcessEntry32
	pe32.Size = uint32(unsafe.Sizeof(pe32))

	// Get first process
	err = windows.Process32First(procs, &pe32)
	if err != nil {
		return nil, fmt.Errorf("failed to get first process: %v", err)
	}

	for {
		// Create a new map for the process
		procInfo := make(map[string]interface{})

		// Initialize all requested columns with default values
		procInfo = specs.Init(ctx, Schema)

		// Always set pid regardless of whether it was explicitly requested
		procInfo["pid"] = int64(pe32.ProcessID)

		if ctx.IsColumnUsed("parent") {
			procInfo["parent"] = int64(pe32.ParentProcessID)
		}
		if ctx.IsColumnUsed("name") {
			procInfo["name"] = windows.UTF16ToString(pe32.ExeFile[:])
		}
		if ctx.IsColumnUsed("threads") {
			procInfo["threads"] = int32(pe32.Threads)
		}

		// Get additional process details
		err = getProcessDetails(ctx, procInfo)
		if err != nil {
			log.Printf("Failed to get process details: %v", err)
		}

		// Add process to the results
		processes = append(processes, procInfo)

		// Get next process
		err = windows.Process32Next(procs, &pe32)
		if err != nil {
			if err == windows.ERROR_NO_MORE_FILES {
				break
			}
			return nil, fmt.Errorf("failed to get next process: %v", err)
		}
	}

	return processes, nil
}
