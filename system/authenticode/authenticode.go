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
	procWinVerifyTrust                      uintptr
	procCryptQueryObject                    uintptr
	procCryptMsgGetParam                    uintptr
)

var (
	_DRIVER_ACTION_VERIFY              = windows.GUID(*ole.NewGUID("{f750e6c3-38ee-11d1-85e5-00c04fc295ee}"))
	_WINTRUST_ACTION_GENERIC_VERIFY_V2 = windows.GUID(*ole.NewGUID("{00aac56b-cd44-11d0-8cc2-00c04fc295ee}"))
)

const (
	_SHA512_HASH_SIZE       = 64
	_CMSG_SIGNER_INFO_PARAM = 2
)

type _CATALOG_INFO struct {
	cbStruct       uint32
	wszCatalogFile [windows.MAX_PATH]uint16
}

type _CRYPT_ATTRIBUTE struct {
	pszObjId *byte
	cValue   uint32
	rgValue  *windows.CryptDataBlob
}

type _CRYPT_ATTRIBUTES struct {
	cAttr  uint32
	rgAttr *_CRYPT_ATTRIBUTE
}

type _CMSG_SIGNER_INFO struct {
	dwVersion               uint32
	Issuer                  windows.CertNameBlob
	SerialNumber            windows.CryptIntegerBlob
	HashAlgorithm           windows.CryptAlgorithmIdentifier
	HashEncryptionAlgorithm windows.CryptAlgorithmIdentifier
	EncryptedHash           windows.CryptDataBlob
	AuthAttrs               _CRYPT_ATTRIBUTES
	UnauthAttrs             _CRYPT_ATTRIBUTES
}

func getSignatureInformation(path string) (string, error) {
	// Convert path to UTF16
	uft16Path, err := windows.UTF16PtrFromString(path)
	if err != nil {
		return "", fmt.Errorf("failed to convert path: %v", err)
	}

	var certStore windows.Handle
	var message windows.Handle

	ret, _, _ := syscall.SyscallN(procCryptQueryObject,
		uintptr(windows.CERT_QUERY_OBJECT_FILE),
		uintptr(unsafe.Pointer(uft16Path)),
		uintptr(windows.CERT_QUERY_CONTENT_FLAG_PKCS7_SIGNED_EMBED),
		uintptr(windows.CERT_QUERY_FORMAT_FLAG_BINARY),
		0,
		0,
		0,
		0,
		uintptr(unsafe.Pointer(&certStore)),
		uintptr(unsafe.Pointer(&message)),
		0,
	)
	if ret == 0 {
		return "", fmt.Errorf("failed to query object: %v", windows.GetLastError())
	}

	var signerInfoSize uint32
	ret, _, _ = syscall.SyscallN(procCryptMsgGetParam,
		uintptr(message),
		uintptr(_CMSG_SIGNER_INFO_PARAM),
		0,
		0,
		uintptr(unsafe.Pointer(&signerInfoSize)),
	)
	if ret == 0 {
		return "", fmt.Errorf("failed to get signer info size: %v", windows.GetLastError())
	}

	signerInfo := make([]byte, signerInfoSize)
	signerInfoPtr := (*_CMSG_SIGNER_INFO)(unsafe.Pointer(&signerInfo[0]))
	ret, _, _ = syscall.SyscallN(procCryptMsgGetParam,
		uintptr(message),
		uintptr(_CMSG_SIGNER_INFO_PARAM),
		0,
		uintptr(unsafe.Pointer(signerInfoPtr)),
		uintptr(unsafe.Pointer(&signerInfoSize)),
	)
	if ret == 0 {
		return "", fmt.Errorf("failed to get signer info: %v", windows.GetLastError())
	}

	return "", nil
}

// getCatalogPathForFilePath retrieves the catalog file path for a given file path
func getCatalogPathForFilePath(path string) (string, error) {
	// Convert path to UTF16
	uft16Path, err := windows.UTF16PtrFromString(path)
	if err != nil {
		return "", fmt.Errorf("failed to convert path: %v", err)
	}

	// Open the file
	handle, err := windows.CreateFile(
		uft16Path,
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
		uintptr(unsafe.Pointer(&_DRIVER_ACTION_VERIFY)),
		0,
	)
	if ret == 0 {
		return "", fmt.Errorf("failed to acquire catalog context: %v", err)
	}
	defer syscall.SyscallN(procCryptCATAdminReleaseContext,
		uintptr(context),
		0,
	)

	// Calculate file hash
	hash := make([]byte, _SHA512_HASH_SIZE)
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
	info.cbStruct = uint32(unsafe.Sizeof(info))

	ret, _, err = syscall.SyscallN(procCryptCATCatalogInfoFromContext,
		uintptr(catalog),
		uintptr(unsafe.Pointer(&info)),
		0,
	)
	if ret == 0 {
		return "", fmt.Errorf("no catalog info found: %v", windows.GetLastError())
	}

	// Convert catalog path from UTF16
	catalogFile := windows.UTF16ToString(info.wszCatalogFile[:])
	return catalogFile, nil
}

func verifySignature(path string) (string, error) {
	uft16Path, err := windows.UTF16PtrFromString(path)
	if err != nil {
		return "", fmt.Errorf("failed to convert path: %v", err)
	}

	trustProviderSettings := windows.WinTrustData{}
	trustProviderSettings.Size = uint32(unsafe.Sizeof(trustProviderSettings))

	fileInfo := windows.WinTrustFileInfo{
		Size:     uint32(unsafe.Sizeof(windows.WinTrustFileInfo{})),
		FilePath: uft16Path,
	}

	trustProviderSettings.FileOrCatalogOrBlobOrSgnrOrCert = unsafe.Pointer(&fileInfo)
	trustProviderSettings.RevocationChecks = windows.WTD_REVOKE_WHOLECHAIN
	trustProviderSettings.UIChoice = windows.WTD_UI_NONE
	trustProviderSettings.StateAction = windows.WTD_STATEACTION_VERIFY
	trustProviderSettings.UnionChoice = windows.WTD_CHOICE_FILE

	authenticodePolicyProvider := _WINTRUST_ACTION_GENERIC_VERIFY_V2
	ret, _, _ := syscall.SyscallN(procWinVerifyTrust,
		uintptr(0),
		uintptr(unsafe.Pointer(&authenticodePolicyProvider)),
		uintptr(unsafe.Pointer(&trustProviderSettings)),
	)

	verificationStatus := int(ret)

	switch verificationStatus {
	case int(windows.ERROR_SUCCESS):
		return "Trusted", nil
	case int(windows.TRUST_E_NOSIGNATURE):
		return "Missing", nil
	case int(windows.CRYPT_E_SECURITY_SETTINGS):
		return "Valid", nil
	case int(windows.TRUST_E_SUBJECT_NOT_TRUSTED):
		return "Untrusted", nil
	default:
		return "Unknown", fmt.Errorf("unknown verification status: %v", verificationStatus)
	}
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

	procWinVerifyTrust, err = windows.GetProcAddress(modWintrust, "WinVerifyTrust")
	if err != nil {
		return nil, fmt.Errorf("failed to get WinVerifyTrust address: %v", err)
	}

	modCrypt32, err := windows.LoadLibrary("crypt32.dll")
	if err != nil {
		return nil, fmt.Errorf("failed to load crypt32.dll: %v", err)
	}
	defer windows.FreeLibrary(modCrypt32)

	procCryptQueryObject, err = windows.GetProcAddress(modCrypt32, "CryptQueryObject")
	if err != nil {
		return nil, fmt.Errorf("failed to get CryptQueryObject address: %v", err)
	}

	procCryptMsgGetParam, err = windows.GetProcAddress(modCrypt32, "CryptMsgGetParam")
	if err != nil {
		return nil, fmt.Errorf("failed to get CryptMsgGetParam address: %v", err)
	}

	catalogFile, err := getCatalogPathForFilePath(path)
	if err != nil {
		catalogFile = path
	}

	fmt.Printf("Catalog file: %s\n", catalogFile)

	authenticode, err := verifySignature(catalogFile)
	if err != nil {
		return nil, fmt.Errorf("failed to verify signature: %v", err)
	}

	fmt.Printf("Authenticode: %s\n", authenticode)

	return nil, errors.New("not implemented")
}
