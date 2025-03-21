package file

import (
	"fmt"
	"os"
	"path/filepath"
	"slices"
	"strings"
	"syscall"
	"unsafe"

	"github.com/scrymastic/goosquery/sql/context"

	"golang.org/x/sys/windows"
)

type FILE_BASIC_INFO struct {
	CreationTime   windows.Filetime
	LastAccessTime windows.Filetime
	LastWriteTime  windows.Filetime
	ChangeTime     windows.Filetime
	FileAttributes int32
	_              [4]byte
}

type LANGANDCODEPAGE struct {
	WLanguage uint16
	WCodePage uint16
}

var (
	modAdvapi32           = windows.NewLazySystemDLL("advapi32.dll")
	procGetSecurityInfo   = modAdvapi32.NewProc("GetSecurityInfo")
	modKernel32           = windows.NewLazySystemDLL("kernel32.dll")
	procGetFileType       = modKernel32.NewProc("GetFileType")
	procGetDiskFreeSpaceW = modKernel32.NewProc("GetDiskFreeSpaceW")
	// Api-ms-win-core-version-l1-1-0.dll
	modApiMsWinCoreVersionL110    = windows.NewLazySystemDLL("Api-ms-win-core-version-l1-1-0.dll")
	procGetFileVersionInfoSizeExW = modApiMsWinCoreVersionL110.NewProc("GetFileVersionInfoSizeExW")
	procGetFileVersionInfoExW     = modApiMsWinCoreVersionL110.NewProc("GetFileVersionInfoExW")
)

const (
	FILE_VER_GET_NEUTRAL = 0x02
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
func getFileTypeString(fileType uint32, attributes uint32, hFile windows.Handle) string {
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
		if windows.GetNamedPipeInfo(hFile, nil, nil, nil, nil) != nil {
			return "socket"
		}
		return "pipe"
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

func getVersionInfo(path string) (string, string, error) {
	var handle windows.Handle
	verSize, err := windows.GetFileVersionInfoSize(path, &handle)
	if err != nil {
		return "", "", fmt.Errorf("failed to get file version info size: %v", err)
	}

	verInfo := make([]byte, verSize)
	err = windows.GetFileVersionInfo(path, uint32(handle), verSize, unsafe.Pointer(&verInfo[0]))
	if err != nil {
		return "", "", fmt.Errorf("failed to get file version info: %v", err)
	}

	var verInfoPtr *windows.VS_FIXEDFILEINFO
	verInfoSize := unsafe.Sizeof(verInfoPtr)
	err = windows.VerQueryValue(unsafe.Pointer(&verInfo[0]), "\\", unsafe.Pointer(&verInfoPtr), (*uint32)(unsafe.Pointer(&verInfoSize)))
	if err != nil {
		return "", "", fmt.Errorf("failed to get file version info: %v", err)
	}

	productVersion := fmt.Sprintf("%d.%d.%d.%d",
		HIWORD(verInfoPtr.ProductVersionMS),
		LOWORD(verInfoPtr.ProductVersionMS),
		HIWORD(verInfoPtr.ProductVersionLS),
		LOWORD(verInfoPtr.ProductVersionLS),
	)

	fileVersion := fmt.Sprintf("%d.%d.%d.%d",
		HIWORD(verInfoPtr.FileVersionMS),
		LOWORD(verInfoPtr.FileVersionMS),
		HIWORD(verInfoPtr.FileVersionLS),
		LOWORD(verInfoPtr.FileVersionLS),
	)

	return productVersion, fileVersion, nil
}

func getLanguagesAndCodepages(versionInfo *byte) ([]LANGANDCODEPAGE, error) {
	var langAndCodePagePtr *LANGANDCODEPAGE

	var verInfoSize uint32

	err := windows.VerQueryValue(unsafe.Pointer(versionInfo), "\\VarFileInfo\\Translation", unsafe.Pointer(&langAndCodePagePtr), &verInfoSize)
	if err != nil {
		return nil, fmt.Errorf("failed to get file version info: %v", err)
	}

	// Cast langAndCodePagePtr to LANGANDCODEPAGE Array
	langAndCodePage := unsafe.Slice(langAndCodePagePtr, int(verInfoSize)/int(unsafe.Sizeof(LANGANDCODEPAGE{})))

	return langAndCodePage, nil
}

func getOriginalFilenameForCodepage(versionInfo *byte, langAndCodePage *LANGANDCODEPAGE) (string, error) {
	// Construct wstring L"\\StringFileInfo\\%04x%04x\\OriginalFilename",
	wstring := fmt.Sprintf("\\StringFileInfo\\%04X%04X\\OriginalFilename", langAndCodePage.WLanguage, langAndCodePage.WCodePage)

	buffer := make([]uint16, 50)
	bufferSize := uint32(len(buffer) * 2)

	windows.VerQueryValue(unsafe.Pointer(versionInfo), wstring, unsafe.Pointer(&buffer), &bufferSize)
	return windows.UTF16PtrToString(&buffer[0]), nil
}

func getOriginalFilename(path string) (string, error) {
	ret, _, err := procGetFileVersionInfoSizeExW.Call(
		uintptr(FILE_VER_GET_NEUTRAL),
		uintptr(unsafe.Pointer(windows.StringToUTF16Ptr(path))),
		uintptr(0),
	)
	if ret == 0 {
		return "", fmt.Errorf("failed to get file version info size: %v", err)
	}
	verSize := uint32(ret)

	verInfo := make([]byte, verSize)
	ret, _, err = procGetFileVersionInfoExW.Call(
		uintptr(FILE_VER_GET_NEUTRAL),
		uintptr(unsafe.Pointer(windows.StringToUTF16Ptr(path))),
		uintptr(0),
		uintptr(verSize),
		uintptr(unsafe.Pointer(&verInfo[0])),
	)
	if ret == 0 {
		return "", fmt.Errorf("failed to get file version info: %v", err)
	}

	langAndCodePage, err := getLanguagesAndCodepages(&verInfo[0])
	if err != nil {
		return "", fmt.Errorf("failed to get languages and code pages: %v", err)
	}

	for _, langAndCodePage := range langAndCodePage {
		// Get original filename for each language and code page
		// Stop on first successful read
		originalFilename, err := getOriginalFilenameForCodepage(&verInfo[0], &langAndCodePage)
		if err == nil {
			return originalFilename, nil
		}
	}

	return "", fmt.Errorf("failed to get original filename")
}

// GetFileStat retrieves detailed file information
func GetFileStat(ctx context.Context, path string, fileStat *map[string]interface{}) error {

	fileInfo, err := os.Stat(path)
	if err != nil {
		return fmt.Errorf("failed to get file info: %w", err)
	}

	FLAGS_AND_ATTRIBUTES := windows.FILE_ATTRIBUTE_ARCHIVE |
		windows.FILE_ATTRIBUTE_ENCRYPTED | windows.FILE_ATTRIBUTE_HIDDEN |
		windows.FILE_ATTRIBUTE_NORMAL | windows.FILE_ATTRIBUTE_OFFLINE |
		windows.FILE_ATTRIBUTE_READONLY | windows.FILE_ATTRIBUTE_SYSTEM |
		windows.FILE_ATTRIBUTE_TEMPORARY

	if fileInfo.IsDir() {
		FLAGS_AND_ATTRIBUTES |= windows.FILE_FLAG_BACKUP_SEMANTICS
	}
	if ctx.IsColumnUsed("atime") {
		winAttr := fileInfo.Sys().(*syscall.Win32FileAttributeData)
		(*fileStat)["atime"] = (*windows.Filetime)(&winAttr.LastAccessTime).Nanoseconds() / 1e9
	}
	if ctx.IsColumnUsed("btime") {
		winAttr := fileInfo.Sys().(*syscall.Win32FileAttributeData)
		(*fileStat)["btime"] = (*windows.Filetime)(&winAttr.CreationTime).Nanoseconds() / 1e9
	}

	hFile, err := windows.CreateFile(
		windows.StringToUTF16Ptr(path),
		windows.GENERIC_READ,
		windows.FILE_SHARE_READ|windows.FILE_SHARE_WRITE,
		nil,
		windows.OPEN_EXISTING,
		uint32(FLAGS_AND_ATTRIBUTES),
		windows.Handle(0),
	)
	if hFile == windows.InvalidHandle {
		return fmt.Errorf("failed to create file: %w", err)
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
		return fmt.Errorf("failed to get security info: %w", err)
	}

	byHandleFileInfo := &windows.ByHandleFileInformation{}

	if err := windows.GetFileInformationByHandle(hFile, byHandleFileInfo); err != nil {
		return fmt.Errorf("failed to get file information: %w", err)
	}

	if ctx.IsColumnUsed("file_id") {
		(*fileStat)["file_id"] = fmt.Sprintf("0x%016X", uint64(byHandleFileInfo.FileIndexHigh)<<32|uint64(byHandleFileInfo.FileIndexLow))
	}
	if ctx.IsColumnUsed("inode") {
		(*fileStat)["inode"] = int64(uint64(byHandleFileInfo.FileIndexHigh)<<32 | uint64(byHandleFileInfo.FileIndexLow))
	}
	if ctx.IsColumnUsed("uid") {
		(*fileStat)["uid"] = getRidFromSid(sidOwner)
	}
	if ctx.IsColumnUsed("gid") {
		(*fileStat)["gid"] = getRidFromSid(sidGroup)
	}
	if ctx.IsColumnUsed("mode") {
		(*fileStat)["mode"] = "-1"
	}
	if ctx.IsColumnUsed("symlink") {
		(*fileStat)["symlink"] = 0
	}
	if ctx.IsColumnUsed("hard_links") {
		(*fileStat)["hard_links"] = int32(byHandleFileInfo.NumberOfLinks)
	}
	if ctx.IsColumnUsed("attributes") {
		(*fileStat)["attributes"] = getFileAttributesString(byHandleFileInfo.FileAttributes)
	}
	if ctx.IsColumnUsed("device") {
		(*fileStat)["device"] = int64(byHandleFileInfo.VolumeSerialNumber)
	}
	if ctx.IsColumnUsed("volume_serial") {
		(*fileStat)["volume_serial"] = fmt.Sprintf("%04X-%04X",
			HIWORD(byHandleFileInfo.VolumeSerialNumber),
			LOWORD(byHandleFileInfo.VolumeSerialNumber),
		)
	}

	fileType, _, err := procGetFileType.Call(
		uintptr(hFile),
	)
	if err != windows.ERROR_SUCCESS {
		return fmt.Errorf("failed to get file type: %w", err)
	}
	if ctx.IsColumnUsed("type") {
		(*fileStat)["type"] = getFileTypeString(uint32(fileType), byHandleFileInfo.FileAttributes, hFile)
	}

	if ctx.IsColumnUsed("block_size") {
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
				return fmt.Errorf("failed to get disk free space: %w", err)
			}

			(*fileStat)["block_size"] = int32(bytesPerSect)
		}
	}

	if ctx.IsColumnUsed("ctime") {
		basicInfo := make([]byte, unsafe.Sizeof(FILE_BASIC_INFO{}))
		err = windows.GetFileInformationByHandleEx(
			hFile,
			windows.FileBasicInfo,
			(*byte)(unsafe.Pointer(&basicInfo[0])),
			uint32(len(basicInfo)),
		)
		if err != nil {
			return fmt.Errorf("failed to get file basic info: %w", err)
		}
		// Cast basicInfo to FILE_BASIC_INFO
		basicInfoPtr := (*FILE_BASIC_INFO)(unsafe.Pointer(&basicInfo[0]))

		(*fileStat)["ctime"] = basicInfoPtr.ChangeTime.Nanoseconds() / 1e9
	}

	if ctx.IsAnyOfColumnsUsed([]string{"product_version", "file_version", "original_filename"}) {
		productVersion, fileVersion, err := getVersionInfo(path)
		if err != nil {
			// Log error
			fmt.Printf("failed to get version info: %v", err)
		}
		if ctx.IsColumnUsed("product_version") {
			(*fileStat)["product_version"] = productVersion
		}
		if ctx.IsColumnUsed("file_version") {
			(*fileStat)["file_version"] = fileVersion
		}
	}

	if ctx.IsColumnUsed("original_filename") {
		originalFilename, err := getOriginalFilename(path)
		if err != nil {
			// Log error
			fmt.Printf("failed to get original filename: %v", err)
		}
		(*fileStat)["original_filename"] = originalFilename
	}

	return nil
}
