package file

import (
	"fmt"
	"os"
	"path/filepath"
	"unsafe"

	"github.com/go-ole/go-ole"
	"github.com/go-ole/go-ole/oleutil"
	"golang.org/x/sys/windows"
)

var (
	modShell32         = windows.NewLazySystemDLL("shell32.dll")
	procSHGetFileInfoW = modShell32.NewProc("SHGetFileInfoW")
)

const (
	ShellLinkHeaderSizeFieldSize     = 4
	ShellLinkHeaderSizeExpectedValue = 0x0000004C
	SHGFI_TYPENAME                   = 0x000000400
	SHGFI_USEFILEATTRIBUTES          = 0x000000010
)

// LnkData represents the data extracted from a shortcut file
type LnkData struct {
	TargetPath     string
	TargetType     string
	TargetLocation string
	StartIn        string
	Comment        string
	Run            string
}

type SHFILEINFOW struct {
	hIcon         windows.Handle
	iIcon         int32
	dwAttributes  uint32
	szDisplayName [windows.MAX_PATH]uint16
	szTypeName    [80]uint16
}

// isValidShellLinkHeader checks if the file has a valid shell link header
func isValidShellLinkHeader(path string) (bool, error) {
	// Open the file
	file, err := os.Open(path)
	if err != nil {
		return false, err
	}
	defer file.Close()

	// Read the first ShellLinkHeaderSizeFieldSize bytes
	data := make([]byte, ShellLinkHeaderSizeFieldSize)
	n, err := file.Read(data)
	if err != nil {
		return false, err
	}
	if n != ShellLinkHeaderSizeFieldSize {
		return false, fmt.Errorf("read %d bytes, expected %d", n, ShellLinkHeaderSizeFieldSize)
	}

	// Convert first 4 bytes to uint32 (little-endian)
	headerSize := uint32(data[0]) | uint32(data[1])<<8 |
		uint32(data[2])<<16 | uint32(data[3])<<24

	return headerSize == ShellLinkHeaderSizeExpectedValue, nil
}

// showCmdToString converts Windows show command to string
func showCmdToString(showCmd int) string {
	switch showCmd {
	case 1: // SW_SHOWNORMAL
		return "Normal window"
	case 3: // SW_SHOWMAXIMIZED
		return "Maximized"
	case 7: // SW_SHOWMINIMIZED
		return "Minimized"
	default:
		return "Unknown"
	}
}

// ParseLnkData parses a Windows shortcut file
func ParseLnkData(linkPath string) (*LnkData, error) {
	// Initialize COM
	err := ole.CoInitialize(0)
	if err != nil {
		return nil, fmt.Errorf("failed to initialize COM: %v", err)
	}
	defer ole.CoUninitialize()

	// Create ShellLink object
	unknown, err := oleutil.CreateObject("WScript.Shell")
	if err != nil {
		return nil, fmt.Errorf("failed to create WScript.Shell: %v", err)
	}
	defer unknown.Release()

	shell, err := unknown.QueryInterface(ole.IID_IDispatch)
	if err != nil {
		return nil, fmt.Errorf("failed to query interface: %v", err)
	}
	defer shell.Release()

	// Create shortcut object
	shortcut, err := oleutil.CallMethod(shell, "CreateShortcut", linkPath)
	if err != nil {
		return nil, fmt.Errorf("failed to create shortcut: %v", err)
	}
	defer shortcut.Clear()

	// Validate header
	valid, err := isValidShellLinkHeader(linkPath)
	if err != nil || !valid {
		return nil, fmt.Errorf("invalid shell link header")
	}

	// Get shortcut properties
	targetPath, _ := oleutil.GetProperty(shortcut.ToIDispatch(), "TargetPath")
	workingDir, _ := oleutil.GetProperty(shortcut.ToIDispatch(), "WorkingDirectory")
	description, _ := oleutil.GetProperty(shortcut.ToIDispatch(), "Description")
	windowStyle, _ := oleutil.GetProperty(shortcut.ToIDispatch(), "WindowStyle")

	// Create LnkData structure
	lnkData := &LnkData{
		TargetPath: targetPath.ToString(),
		StartIn:    workingDir.ToString(),
		Comment:    description.ToString(),
		Run:        showCmdToString(int(windowStyle.Val)),
	}

	// Get attributes of target path
	attributes, err := windows.GetFileAttributes(windows.StringToUTF16Ptr(lnkData.TargetPath))
	if err != nil {
		return nil, fmt.Errorf("failed to get file attributes: %v", err)
	}

	var fileInfo SHFILEINFOW

	ret, _, err := procSHGetFileInfoW.Call(
		uintptr(unsafe.Pointer(windows.StringToUTF16Ptr(lnkData.TargetPath))),
		uintptr(attributes),
		uintptr(unsafe.Pointer(&fileInfo)),
		uintptr(unsafe.Sizeof(SHFILEINFOW{})),
		uintptr(SHGFI_TYPENAME|SHGFI_USEFILEATTRIBUTES),
	)

	if ret == 0 {
		return nil, fmt.Errorf("failed to get file info: %v", err)
	}

	lnkData.TargetType = windows.UTF16ToString(fileInfo.szTypeName[:])

	// Target location is name of the folder that hold targettype
	lnkData.TargetLocation = filepath.Base(filepath.Dir(lnkData.TargetPath))

	return lnkData, nil
}
