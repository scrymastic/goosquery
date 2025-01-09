package authenticode

import (
	"errors"
	"fmt"
	"syscall"
	"unsafe"

	"github.com/go-ole/go-ole"
	"golang.org/x/sys/windows"
)

type Authenticode struct {
	Path                string `json:"path"`
	OriginalProgramName string `json:"original_program_name"`
	SerialNumber        string `json:"serial_number"`
	IssuerName          string `json:"issuer_name"`
	SubjectName         string `json:"subject_name"`
	Result              string `json:"result"`
}

var (
	procCryptCATAdminAcquireContext         uintptr
	procCryptCATAdminCalcHashFromFileHandle uintptr
	procCryptCATAdminEnumCatalogFromHash    uintptr
	procCryptCATCatalogInfoFromContext      uintptr
	procCryptCATAdminReleaseCatalogContext  uintptr
	procCryptCATAdminReleaseContext         uintptr
)

// // GUID for driver verification
// var DRIVER_ACTION_VERIFY = windows.GUID{
// 	Data1: 0xf750e6c3,
// 	Data2: 0x38ee,
// 	Data3: 0x11d1,
// 	Data4: [8]byte{0x85, 0xe5, 0x00, 0xc0, 0x4f, 0xc2, 0x95, 0xee},
// }

var DRIVER_ACTION_VERIFY = ole.NewGUID("{f750e6c3-38ee-11d1-85e5-00c04fc295ee}")

const SHA512_HASH_SIZE = 64

// _CATALOG_INFO structure
type _CATALOG_INFO struct {
	CbStruct       uint32
	WSzCatalogFile [windows.MAX_PATH]uint16
}

// GetCatalogPathForFilePath retrieves the catalog file path for a given file path
func GetCatalogPathForFilePath(path string) (catalogFile string, err error) {
	// Convert path to UTF16
	pathPtr, err := windows.UTF16PtrFromString(path)
	if err != nil {
		return "", fmt.Errorf("failed to convert path: %v", err)
	}

	// Open the file
	handle, err := windows.CreateFile(
		pathPtr,
		windows.GENERIC_READ,
		windows.FILE_SHARE_READ,
		nil,
		windows.OPEN_EXISTING,
		windows.FILE_ATTRIBUTE_NORMAL,
		0,
	)
	if err != nil || handle == windows.InvalidHandle {
		return "", fmt.Errorf("failed to open file: %v", err)
	}
	defer windows.CloseHandle(handle)

	// Acquire catalog admin context
	var context windows.Handle
	ret, _, _ := syscall.SyscallN(procCryptCATAdminAcquireContext,
		uintptr(unsafe.Pointer(&context)),
		uintptr(unsafe.Pointer(&DRIVER_ACTION_VERIFY)),
		0)
	if ret == 0 {
		return "", fmt.Errorf("failed to acquire catalog context: %v", err)
	}
	defer syscall.SyscallN(procCryptCATAdminReleaseContext,
		uintptr(context),
		0,
	)

	// Calculate file hash
	hash := make([]byte, SHA512_HASH_SIZE)
	hashSize := uint32(len(hash))

	ret, _, _ = syscall.SyscallN(procCryptCATAdminCalcHashFromFileHandle,
		uintptr(handle),
		uintptr(unsafe.Pointer(&hashSize)),
		uintptr(unsafe.Pointer(&hash[0])),
		0,
	)
	if ret == 0 {
		return "", fmt.Errorf("failed to calculate hash: %v", windows.GetLastError())
	}

	// Find matching catalog
	var catalog windows.Handle
	ret, _, err = syscall.SyscallN(procCryptCATAdminEnumCatalogFromHash,
		uintptr(context),
		uintptr(unsafe.Pointer(&hash[0])),
		uintptr(hashSize),
		0,
		0,
	)
	if ret == 0 {
		return "", fmt.Errorf("failed to enumerate catalog: %v", windows.GetLastError())
	}
	catalog = windows.Handle(ret)

	defer syscall.SyscallN(procCryptCATAdminReleaseCatalogContext,
		uintptr(context),
		uintptr(catalog),
		0,
	)

	// Get catalog info
	var info _CATALOG_INFO
	info.CbStruct = uint32(unsafe.Sizeof(info))

	ret, _, err = syscall.SyscallN(procCryptCATCatalogInfoFromContext,
		uintptr(catalog),
		uintptr(unsafe.Pointer(&info)),
		0,
	)
	if ret == 0 {
		return "", fmt.Errorf("no catalog info found: %v", windows.GetLastError())
	}

	// Convert catalog path from UTF16
	catalogFile = windows.UTF16ToString(info.WSzCatalogFile[:])
	return catalogFile, nil
}

func GenAuthenticode(path string) ([]Authenticode, error) {
	// Load Windows API functions
	modWintrust, err := windows.LoadLibrary("wintrust.dll")
	if err != nil {
		return nil, fmt.Errorf("failed to load wintrust.dll: %v", err)
	}
	defer windows.FreeLibrary(modWintrust)

	procCryptCATAdminAcquireContext, err = windows.GetProcAddress(modWintrust, "CryptCATAdminAcquireContext")
	if err != nil {
		return nil, fmt.Errorf("failed to get CryptCATAdminAcquireContext address: %v", err)
	}

	procCryptCATAdminCalcHashFromFileHandle, err = windows.GetProcAddress(modWintrust, "CryptCATAdminCalcHashFromFileHandle")
	if err != nil {
		return nil, fmt.Errorf("failed to get CryptCATAdminCalcHashFromFileHandle address: %v", err)
	}

	procCryptCATAdminEnumCatalogFromHash, err = windows.GetProcAddress(modWintrust, "CryptCATAdminEnumCatalogFromHash")
	if err != nil {
		return nil, fmt.Errorf("failed to get CryptCATAdminEnumCatalogFromHash address: %v", err)
	}

	procCryptCATCatalogInfoFromContext, err = windows.GetProcAddress(modWintrust, "CryptCATCatalogInfoFromContext")
	if err != nil {
		return nil, fmt.Errorf("failed to get CryptCATCatalogInfoFromContext address: %v", err)
	}

	procCryptCATAdminReleaseCatalogContext, err = windows.GetProcAddress(modWintrust, "CryptCATAdminReleaseCatalogContext")
	if err != nil {
		return nil, fmt.Errorf("failed to get CryptCATAdminReleaseCatalogContext address: %v", err)
	}

	procCryptCATAdminReleaseContext, err = windows.GetProcAddress(modWintrust, "CryptCATAdminReleaseContext")
	if err != nil {
		return nil, fmt.Errorf("failed to get CryptCATAdminReleaseContext address: %v", err)
	}

	catalogFile, err := GetCatalogPathForFilePath(path)
	if err != nil {
		return nil, fmt.Errorf("failed to get catalog file: %v", err)
	}

	fmt.Printf("Catalog file: %s\n", catalogFile)

	return nil, errors.New("not implemented")
}
