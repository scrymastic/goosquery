package main

import (
	"fmt"
	"log"
	"os"
	"syscall"
	"unsafe"

	"golang.org/x/sys/windows"
)

type Process struct {
	PID     uint32 `json:"pid"`
	Name    string `json:"name"`
	Path    string `json:"path"`
	CMDLine string `json:"cmdline"`
	State   string `json:"state"`
	CWD     string `json:"cwd"`
	Root    string `json:"root"`
	// UID
	// GID
	// EUID
	// EGID
	// SUID
	// SGID
	OnDisk       bool   `json:"on_disk"`
	WiredSize    uint64 `json:"wired_size"`
	ResidentSize uint64 `json:"resident_size"`
	TotalSize    uint64 `json:"total_size"`
	// UserTime
	// SystemTime
	DiskBytesRead    uint64 `json:"disk_bytes_read"`
	DiskBytesWritten uint64 `json:"disk_bytes_written"`
	// StartTime
	Parent uint32 `json:"parent"`
	// Pgroup
	Threads       uint32 `json:"threads"`
	Nice          int32  `json:"nice"`
	ElevatedToken bool   `json:"elevated_token"`
	// SecureProcess
	// ProtectionType
	// VirtualProcess
	ElapsedTime uint64 `json:"elapsed_time"`
	HandleCount uint32 `json:"handle_count"`
	// PercentProcessorTime
	// UPID
	// UPPID
	// CPUType
	// CPUSubtype
	// Translated
	// CgroupPath
}

type _PROCESS_MEMORY_COUNTERS_EX struct {
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

type _IO_COUNTERS struct {
	ReadOperationCount  uint64
	WriteOperationCount uint64
	OtherOperationCount uint64
	ReadTransferCount   uint64
	WriteTransferCount  uint64
	OtherTransferCount  uint64
}

type _TOKEN_ELEVATION struct {
	TokenIsElevated uint32
}

const (
	STILL_ACTIVE = 259
)

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
func getProcessRssInfo(proc windows.Handle) (*_PROCESS_MEMORY_COUNTERS_EX, error) {
	// https://docs.microsoft.com/en-us/windows/win32/api/psapi/nf-psapi-getprocessmemoryinfo
	var procGetProcessMemoryInfo uintptr

	kernel32, err := syscall.LoadLibrary("kernel32.dll")
	if err != nil {
		return nil, fmt.Errorf("failed to load kernel32.dll: %v", err)
	}
	// It's weird to unload kernel32.dll, just a convention
	defer syscall.FreeLibrary(kernel32)
	procGetProcessMemoryInfo, err = syscall.GetProcAddress(kernel32, "K32GetProcessMemoryInfo")
	if err != nil {
		return nil, fmt.Errorf("failed to get GetProcessMemoryInfo address: %v", err)
	}

	// If GetProcessMemoryInfo is not available, try with psapi.dll
	if procGetProcessMemoryInfo == 0 {
		psapi, err := syscall.LoadLibrary("psapi.dll")
		if err != nil {
			return nil, fmt.Errorf("failed to load psapi.dll: %v", err)
		}
		defer syscall.FreeLibrary(psapi)
		procGetProcessMemoryInfo, err = syscall.GetProcAddress(psapi, "GetProcessMemoryInfo")
		if err != nil {
			return nil, fmt.Errorf("failed to get GetProcessMemoryInfo address: %v", err)
		}
	}

	// If GetProcessMemoryInfo is not available, return an error
	if procGetProcessMemoryInfo == 0 {
		return nil, fmt.Errorf("failed to get GetProcessMemoryInfo address")
	}

	var pmcEx _PROCESS_MEMORY_COUNTERS_EX
	ret, _, _ := syscall.SyscallN(procGetProcessMemoryInfo,
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
	var elevation _TOKEN_ELEVATION
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

func getProcessDetails(procInfo Process) (Process, error) {
	// Try to open process with PROCESS_ALL_ACCESS first
	procHandle, err := windows.OpenProcess(windows.PROCESS_ALL_ACCESS, false, procInfo.PID)
	if err != nil {
		// If it fails, try with PROCESS_QUERY_LIMITED_INFORMATION
		procHandle, err = windows.OpenProcess(windows.PROCESS_QUERY_LIMITED_INFORMATION, false, procInfo.PID)
	}

	if err != nil {
		return procInfo, fmt.Errorf("failed to open process %d: %v", procInfo.PID, err)
	}

	// Get process path
	var pathBuf [windows.MAX_PATH]uint16
	pathLen := uint32(len(pathBuf))
	if err := windows.QueryFullProcessImageName(
		procHandle,
		0,
		&pathBuf[0],
		&pathLen); err == nil {
		procInfo.Path = windows.UTF16ToString(pathBuf[:pathLen])
	}

	// Get command line
	cmdline, err := getProcessCommandLine(procHandle)
	if err == nil {
		procInfo.CMDLine = cmdline
	}

	// Get process state
	var exitCode uint32
	if err := windows.GetExitCodeProcess(procHandle, &exitCode); err == nil {
		if exitCode == STILL_ACTIVE {
			procInfo.State = "STILL_ACTIVE"
		} else {
			procInfo.State = "EXITED"
		}
	}

	// Get working directory
	cwd, err := getProcessWorkingDirectory(procHandle)
	if err == nil {
		procInfo.CWD = cwd
	}

	// Get process root directory
	procInfo.Root = procInfo.CWD

	// Check if process is on disk
	if _, err := os.Stat(procInfo.Path); err == nil {
		procInfo.OnDisk = true
	}

	// Get process memory info
	pmc, err := getProcessRssInfo(procHandle)
	if err == nil {
		procInfo.WiredSize = pmc.QuotaNonPagedPoolUsage
		procInfo.ResidentSize = pmc.WorkingSetSize
		procInfo.TotalSize = pmc.PrivateUsage
	}

	// Get disk bytes read and written
	var ioCounters _IO_COUNTERS
	processIoCounters := kernel32.NewProc("GetProcessIoCounters")
	if ret, _, _ := processIoCounters.Call(
		uintptr(procHandle),
		uintptr(unsafe.Pointer(&ioCounters))); ret != 0 {
		procInfo.DiskBytesRead = ioCounters.ReadTransferCount
		procInfo.DiskBytesWritten = ioCounters.WriteTransferCount
	}

	// Get nice value
	nice, err := windows.GetPriorityClass(procHandle)
	if err == nil {
		procInfo.Nice = int32(nice)
	}

	// Get elevated token
	elevated, err := isProcessElevated(procHandle)
	if err == nil {
		procInfo.ElevatedToken = elevated
	}

	// Get elapsed time
	elapsedTime, err := getProcessElapsedTime(procHandle)
	if err == nil {
		procInfo.ElapsedTime = elapsedTime
	}

	// Get handle count
	var handleCount uint32
	procGetProcessHandleCount := kernel32.NewProc("GetProcessHandleCount")
	if ret, _, _ := procGetProcessHandleCount.Call(uintptr(procHandle), uintptr(unsafe.Pointer(&handleCount))); ret != 0 {
		procInfo.HandleCount = handleCount
	}

	windows.CloseHandle(procHandle)

	return procInfo, nil
}

func GenProcesses() ([]Process, error) {
	procs, err := windows.CreateToolhelp32Snapshot(windows.TH32CS_SNAPALL, 0)
	if err != nil {
		return nil, fmt.Errorf("failed to create snapshot: %v", err)
	}
	defer windows.CloseHandle(procs)

	var processes []Process
	var pe32 windows.ProcessEntry32
	pe32.Size = uint32(unsafe.Sizeof(pe32))

	// Get first process
	err = windows.Process32First(procs, &pe32)
	if err != nil {
		return nil, fmt.Errorf("failed to get first process: %v", err)
	}

	for {
		procInfo := Process{
			PID:     pe32.ProcessID,
			Parent:  pe32.ParentProcessID,
			Name:    windows.UTF16ToString(pe32.ExeFile[:]),
			Threads: pe32.Threads,
		}

		procInfo, err = getProcessDetails(procInfo)
		if err != nil {
			log.Printf("Failed to get process details: %v", err)
		}

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
