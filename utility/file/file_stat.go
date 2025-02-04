package file

import (
	"fmt"
	"os"
	"path/filepath"
	"slices"
	"strings"
	"syscall"
	"unsafe"

	"golang.org/x/sys/windows"
)

type FileStat struct {
	Path             string `json:"path"`
	Filename         string `json:"filename"`
	Symlink          int32  `json:"symlink"`
	FileID           string `json:"file_id"`
	Inode            int64  `json:"inode"`
	UID              int64  `json:"uid"`
	GID              int64  `json:"gid"`
	Mode             string `json:"mode"`
	Device           int64  `json:"device"`
	Size             int64  `json:"size"`
	BlockSize        int32  `json:"block_size"`
	Atime            int64  `json:"atime"`
	Mtime            int64  `json:"mtime"`
	Ctime            int64  `json:"ctime"`
	Btime            int64  `json:"btime"`
	HardLinks        int32  `json:"hard_links"`
	Type             string `json:"type"`
	Attributes       string `json:"attributes"`
	VolumeSerial     string `json:"volume_serial"`
	ProductVersion   string `json:"product_version"`
	FileVersion      string `json:"file_version"`
	OriginalFilename string `json:"original_filename"`
}

// typedef struct _FILE_BASIC_INFO {
// 	LARGE_INTEGER CreationTime;
// 	LARGE_INTEGER LastAccessTime;
// 	LARGE_INTEGER LastWriteTime;
// 	LARGE_INTEGER ChangeTime;
// 	DWORD         FileAttributes;
//   } FILE_BASIC_INFO, *PFILE_BASIC_INFO;

type _FILE_BASIC_INFO struct {
	CreationTime   windows.Filetime
	LastAccessTime windows.Filetime
	LastWriteTime  windows.Filetime
	ChangeTime     windows.Filetime
	FileAttributes int32
	_              [4]byte
}

var (
	modAdvapi32           = windows.NewLazySystemDLL("advapi32.dll")
	procGetSecurityInfo   = modAdvapi32.NewProc("GetSecurityInfo")
	modKernel32           = windows.NewLazySystemDLL("kernel32.dll")
	procGetFileType       = modKernel32.NewProc("GetFileType")
	procGetDiskFreeSpaceW = modKernel32.NewProc("GetDiskFreeSpaceW")
)

var driveLetters = []string{
	"A:", "B:", "C:", "D:", "E:", "F:", "G:", "H:", "I:", "J:", "K:", "L:", "M:", "N:", "O:", "P:", "Q:", "R:", "S:", "T:", "U:", "V:", "W:", "X:", "Y:", "Z:",
}

func HIWORD(value uint32) uint16 {
	return uint16(value >> 16)
}

func LOWORD(value uint32) uint16 {
	return uint16(value & 0xFFFF)
}

func getRidFromSid(sid *windows.SID) int64 {
	if sid == nil {
		return -1
	}
	return int64(sid.SubAuthority(uint32(sid.SubAuthorityCount() - 1)))
}

// Helper function to get file type string
func getFileTypeString(fileType uint32, attributes uint32) string {
	switch fileType {
	case windows.FILE_TYPE_CHAR:
		return "character"
	case windows.FILE_TYPE_DISK:
		if attributes&windows.FILE_ATTRIBUTE_DIRECTORY != 0 {
			return "directory"
		}
		if attributes&windows.FILE_ATTRIBUTE_REPARSE_POINT != 0 {
			return "symbolic"
		}
		if attributes&windows.FILE_ATTRIBUTE_ARCHIVE != 0 ||
			attributes&windows.FILE_ATTRIBUTE_NORMAL != 0 {
			return "regular"
		}
		return "disk"
	case windows.FILE_TYPE_PIPE:
		return "pipe or socket"
	default:
		return "unknown"
	}
}

// Helper function to get file attributes string
func getFileAttributesString(attrs uint32) string {
	var result strings.Builder
	if attrs&windows.FILE_ATTRIBUTE_ARCHIVE != 0 {
		result.WriteRune('A')
	}
	if attrs&windows.FILE_ATTRIBUTE_COMPRESSED != 0 {
		result.WriteRune('C')
	}
	if attrs&windows.FILE_ATTRIBUTE_ENCRYPTED != 0 {
		result.WriteRune('E')
	}
	if attrs&windows.FILE_ATTRIBUTE_REPARSE_POINT != 0 {
		result.WriteRune('L')
	}
	if attrs&windows.FILE_ATTRIBUTE_HIDDEN != 0 {
		result.WriteRune('H')
	}
	if attrs&windows.FILE_ATTRIBUTE_INTEGRITY_STREAM != 0 {
		result.WriteRune('V')
	}
	if attrs&windows.FILE_ATTRIBUTE_NORMAL != 0 {
		result.WriteRune('N')
	}
	if attrs&windows.FILE_ATTRIBUTE_NOT_CONTENT_INDEXED != 0 {
		result.WriteRune('I')
	}
	if attrs&windows.FILE_ATTRIBUTE_NO_SCRUB_DATA != 0 {
		result.WriteRune('X')
	}
	if attrs&windows.FILE_ATTRIBUTE_OFFLINE != 0 {
		result.WriteRune('O')
	}
	if attrs&windows.FILE_ATTRIBUTE_READONLY != 0 {
		result.WriteRune('R')
	}
	if attrs&windows.FILE_ATTRIBUTE_SYSTEM != 0 {
		result.WriteRune('S')
	}
	if attrs&windows.FILE_ATTRIBUTE_TEMPORARY != 0 {
		result.WriteRune('T')
	}
	return result.String()
}

// GetFileStat retrieves detailed file information
func GetFileStat(path string) (*FileStat, error) {
	fileStat := &FileStat{}

	fileInfo, err := os.Stat(path)
	if err != nil {
		return nil, fmt.Errorf("failed to get file info: %w", err)
	}

	FLAG_AND_ATTRIBUTES := windows.FILE_ATTRIBUTE_ARCHIVE |
		windows.FILE_ATTRIBUTE_ENCRYPTED | windows.FILE_ATTRIBUTE_HIDDEN |
		windows.FILE_ATTRIBUTE_NORMAL | windows.FILE_ATTRIBUTE_OFFLINE |
		windows.FILE_ATTRIBUTE_READONLY | windows.FILE_ATTRIBUTE_SYSTEM |
		windows.FILE_ATTRIBUTE_TEMPORARY

	if fileInfo.IsDir() {
		FLAG_AND_ATTRIBUTES |= windows.FILE_FLAG_BACKUP_SEMANTICS
	}

	fileStat.Path = path
	fileStat.Filename = fileInfo.Name()
	fileStat.Size = int64(fileInfo.Size())
	fileStat.Mode = fileInfo.Mode().String()
	winAttr := fileInfo.Sys().(*syscall.Win32FileAttributeData)
	fileStat.Atime = (*windows.Filetime)(&winAttr.LastAccessTime).Nanoseconds() / 1e9
	fileStat.Mtime = (*windows.Filetime)(&winAttr.LastWriteTime).Nanoseconds() / 1e9
	fileStat.Btime = (*windows.Filetime)(&winAttr.CreationTime).Nanoseconds() / 1e9

	hFile, err := windows.CreateFile(
		windows.StringToUTF16Ptr(path),
		windows.GENERIC_READ,
		windows.FILE_SHARE_READ|windows.FILE_SHARE_WRITE,
		nil,
		windows.OPEN_EXISTING,
		uint32(FLAG_AND_ATTRIBUTES),
		windows.Handle(0),
	)
	if hFile == windows.InvalidHandle {
		return nil, fmt.Errorf("failed to create file: %w", err)
	}
	defer windows.CloseHandle(hFile)

	var sidOwner, sidGroup *windows.SID
	var securityDescriptor *windows.SECURITY_DESCRIPTOR

	ret, _, err := procGetSecurityInfo.Call(
		uintptr(hFile),
		uintptr(windows.SE_FILE_OBJECT),
		uintptr(windows.OWNER_SECURITY_INFORMATION|windows.GROUP_SECURITY_INFORMATION),
		uintptr(unsafe.Pointer(&sidOwner)),
		uintptr(unsafe.Pointer(&sidGroup)),
		uintptr(0),
		uintptr(0),
		uintptr(unsafe.Pointer(&securityDescriptor)),
	)
	defer windows.LocalFree(windows.Handle(uintptr(unsafe.Pointer(securityDescriptor))))

	if windows.Errno(ret) != windows.ERROR_SUCCESS {
		return nil, fmt.Errorf("failed to get security info: %w", err)
	}

	byHandleFileInfo := &windows.ByHandleFileInformation{}

	if err := windows.GetFileInformationByHandle(hFile, byHandleFileInfo); err != nil {
		return nil, fmt.Errorf("failed to get file information: %w", err)
	}

	fileStat.FileID = fmt.Sprintf("0x%016X", uint64(byHandleFileInfo.FileIndexHigh)<<32|uint64(byHandleFileInfo.FileIndexLow))
	fileStat.Inode = int64(uint64(byHandleFileInfo.FileIndexHigh)<<32 | uint64(byHandleFileInfo.FileIndexLow))
	fileStat.UID = getRidFromSid(sidOwner)
	fileStat.GID = getRidFromSid(sidGroup)
	fileStat.Mode = "-1"
	fileStat.Symlink = 0
	fileStat.HardLinks = int32(byHandleFileInfo.NumberOfLinks)
	fileStat.Attributes = getFileAttributesString(byHandleFileInfo.FileAttributes)
	fileStat.Device = int64(byHandleFileInfo.VolumeSerialNumber)
	fileStat.VolumeSerial = fmt.Sprintf("%04X-%04X",
		HIWORD(byHandleFileInfo.VolumeSerialNumber),
		LOWORD(byHandleFileInfo.VolumeSerialNumber),
	)

	fileType, _, err := procGetFileType.Call(
		uintptr(hFile),
	)
	if err != windows.ERROR_SUCCESS {
		return nil, fmt.Errorf("failed to get file type: %w", err)
	}
	fileStat.Type = getFileTypeString(uint32(fileType), byHandleFileInfo.FileAttributes)
	if fileStat.Type == "pipe or socket" {
		if windows.GetNamedPipeInfo(hFile, nil, nil, nil, nil) != nil {
			fileStat.Type = "socket"
		} else {
			fileStat.Type = "pipe"
		}
	}
	// Extract drive letter from path
	driveLetter := filepath.VolumeName(path)
	// Check if drive letter is in driveLetters
	if slices.Contains(driveLetters, driveLetter) {
		var sectPerCluster, bytesPerSect, freeClusters, totalClusters uint32
		ret, _, err := procGetDiskFreeSpaceW.Call(
			uintptr(unsafe.Pointer(windows.StringToUTF16Ptr(driveLetter))),
			uintptr(unsafe.Pointer(&sectPerCluster)),
			uintptr(unsafe.Pointer(&bytesPerSect)),
			uintptr(unsafe.Pointer(&freeClusters)),
			uintptr(unsafe.Pointer(&totalClusters)),
		)
		if ret == 0 {
			return nil, fmt.Errorf("failed to get disk free space: %w", err)
		}
		fileStat.BlockSize = int32(bytesPerSect)
	}
	basicInfo := make([]byte, unsafe.Sizeof(_FILE_BASIC_INFO{}))
	err = windows.GetFileInformationByHandleEx(
		hFile,
		windows.FileBasicInfo,
		(*byte)(unsafe.Pointer(&basicInfo[0])),
		uint32(len(basicInfo)),
	)
	if err != nil {
		return nil, fmt.Errorf("failed to get file basic info: %w", err)
	}

	fileStat.Ctime = (*windows.Filetime)(unsafe.Pointer(&basicInfo[0])).Nanoseconds() / 1e9

	return fileStat, nil
}
