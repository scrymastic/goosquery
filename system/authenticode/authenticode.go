package authenticode

import (
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
	procCryptDecodeObject                   uintptr
	procCryptMsgClose                       uintptr
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

// error invalid pointer attr.pszObjId
func getOriginalProgramName(signerInfo *_CMSG_SIGNER_INFO) (string, error) {
	var publisherInfoPtr *_CRYPT_ATTRIBUTE

	// Search through auth attributes for SPC_SP_OPUS_INFO_OBJID
	for i := uint32(0); i < signerInfo.AuthAttrs.cAttr; i++ {
		attr := (*_CRYPT_ATTRIBUTE)(
			unsafe.Pointer(
				uintptr(unsafe.Pointer(signerInfo.AuthAttrs.rgAttr)) +
					uintptr(i)*unsafe.Sizeof(_CRYPT_ATTRIBUTE{}),
			),
		)

		if windows.BytePtrToString(attr.pszObjId) == "1.3.6.1.4.1.311.2.1.12" { // SPC_SP_OPUS_INFO_OBJID
			publisherInfoPtr = attr
			break
		}
	}

	if publisherInfoPtr == nil {
		return "", fmt.Errorf("publisher information could not be found")
	}

	var publisherInfoSize uint32
	ret, _, _ := syscall.SyscallN(procCryptDecodeObject,
		uintptr(windows.X509_ASN_ENCODING|windows.PKCS_7_ASN_ENCODING),
		uintptr(unsafe.Pointer(publisherInfoPtr.pszObjId)),
		uintptr(unsafe.Pointer(publisherInfoPtr.rgValue.Data)),
		uintptr(publisherInfoPtr.rgValue.Size),
		0,
		0,
		uintptr(unsafe.Pointer(&publisherInfoSize)))

	if ret == 0 {
		return "", fmt.Errorf("failed to access the publisher information: %v", windows.GetLastError())
	}

	publisherInfoBlob := make([]byte, publisherInfoSize)
	ret, _, _ = syscall.SyscallN(procCryptDecodeObject,
		uintptr(windows.X509_ASN_ENCODING|windows.PKCS_7_ASN_ENCODING),
		uintptr(unsafe.Pointer(publisherInfoPtr.pszObjId)),
		uintptr(unsafe.Pointer(publisherInfoPtr.rgValue.Data)),
		uintptr(publisherInfoPtr.rgValue.Size),
		0,
		uintptr(unsafe.Pointer(&publisherInfoBlob[0])),
		uintptr(unsafe.Pointer(&publisherInfoSize)))

	if ret == 0 {
		return "", fmt.Errorf("failed to decode the publisher information: %v", windows.GetLastError())
	}

	// Extract program name from decoded blob
	// The SPC_SP_OPUS_INFO structure has a pwszProgramName field that's a wide string
	if len(publisherInfoBlob) < 2 {
		return "", nil
	}

	programNamePtr := (*uint16)(unsafe.Pointer(&publisherInfoBlob[0]))
	if programNamePtr == nil {
		return "", nil
	}

	programName := windows.UTF16PtrToString(programNamePtr)

	return programName, nil
}

func getSignatureInformation(path string, authenticode *Authenticode) error {
	// Convert path to UTF16
	uft16Path, err := windows.UTF16PtrFromString(path)
	if err != nil {
		return fmt.Errorf("failed to convert path: %v", err)
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
		return fmt.Errorf("failed to query object: %v", windows.GetLastError())
	}
	defer windows.CertCloseStore(certStore, 0)
	// Use syscall for CryptMsgClose since it's not in windows package
	defer syscall.SyscallN(procCryptMsgClose, uintptr(message))

	var signerInfoSize uint32
	ret, _, _ = syscall.SyscallN(procCryptMsgGetParam,
		uintptr(message),
		uintptr(_CMSG_SIGNER_INFO_PARAM),
		0,
		0,
		uintptr(unsafe.Pointer(&signerInfoSize)),
	)
	if ret == 0 {
		return fmt.Errorf("failed to get signer info size: %v", windows.GetLastError())
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
		return fmt.Errorf("failed to get signer info: %v", windows.GetLastError())
	}

	// // Get original program name
	// originalProgramName, err := getOriginalProgramName(signerInfoPtr)
	// if err != nil {
	// 	// Just log the error and continue, as this is not critical
	// 	fmt.Printf("Warning: failed to get original program name: %v\n", err)
	// }
	// authenticode.OriginalProgramName = originalProgramName

	// Find certificate in store
	certInfo := windows.CertInfo{
		Issuer:       signerInfoPtr.Issuer,
		SerialNumber: signerInfoPtr.SerialNumber,
	}

	certContext, err := windows.CertFindCertificateInStore(
		certStore,
		windows.X509_ASN_ENCODING|windows.PKCS_7_ASN_ENCODING,
		0,
		windows.CERT_FIND_SUBJECT_CERT,
		unsafe.Pointer(&certInfo),
		nil,
	)
	if err != nil {
		return fmt.Errorf("failed to find certificate in store: %v", err)
	}
	defer windows.CertFreeCertificateContext(certContext)

	// Get serial number
	serialNumber := make([]byte, certContext.CertInfo.SerialNumber.Size)
	copy(serialNumber, unsafe.Slice(certContext.CertInfo.SerialNumber.Data, certContext.CertInfo.SerialNumber.Size))

	// Format serial number as hex string
	var serialNumberStr string
	for i := len(serialNumber) - 1; i >= 0; i-- {
		serialNumberStr += fmt.Sprintf("%02x", serialNumber[i])
	}
	authenticode.SerialNumber = serialNumberStr

	// Get issuer name
	var issuerNameSize uint32
	// First call to get size
	issuerNameSize = uint32(windows.CertGetNameString(
		certContext,
		windows.CERT_NAME_SIMPLE_DISPLAY_TYPE,
		windows.CERT_NAME_ISSUER_FLAG,
		nil,
		nil,
		0,
	))
	if issuerNameSize == 0 {
		return fmt.Errorf("failed to get issuer name size: %v", windows.GetLastError())
	}

	issuerName := make([]uint16, issuerNameSize)
	// Second call to get actual name
	ret2 := windows.CertGetNameString(
		certContext,
		windows.CERT_NAME_SIMPLE_DISPLAY_TYPE,
		windows.CERT_NAME_ISSUER_FLAG,
		nil,
		&issuerName[0],
		issuerNameSize,
	)
	if ret2 == 0 {
		return fmt.Errorf("failed to get issuer name: %v", windows.GetLastError())
	}
	authenticode.IssuerName = windows.UTF16ToString(issuerName)

	// Get subject name
	var subjectNameSize uint32
	// First call to get size
	subjectNameSize = uint32(windows.CertGetNameString(
		certContext,
		windows.CERT_NAME_SIMPLE_DISPLAY_TYPE,
		0,
		nil,
		nil,
		0,
	))
	if subjectNameSize == 0 {
		return fmt.Errorf("failed to get subject name size: %v", windows.GetLastError())
	}

	subjectName := make([]uint16, subjectNameSize)
	// Second call to get actual name
	ret2 = windows.CertGetNameString(
		certContext,
		windows.CERT_NAME_SIMPLE_DISPLAY_TYPE,
		0,
		nil,
		&subjectName[0],
		subjectNameSize,
	)
	if ret2 == 0 {
		return fmt.Errorf("failed to get subject name: %v", windows.GetLastError())
	}
	authenticode.SubjectName = windows.UTF16ToString(subjectName)

	return nil
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
		return "trusted", nil
	case int(windows.TRUST_E_NOSIGNATURE):
		return "missing", nil
	case int(windows.CRYPT_E_SECURITY_SETTINGS):
		return "valid", nil
	case int(windows.TRUST_E_SUBJECT_NOT_TRUSTED):
		return "untrusted", nil
	default:
		return "unknown", fmt.Errorf("unknown verification status: %v", verificationStatus)
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

	procCryptDecodeObject, err = windows.GetProcAddress(modCrypt32, "CryptDecodeObject")
	if err != nil {
		return nil, fmt.Errorf("failed to get CryptDecodeObject address: %v", err)
	}

	procCryptMsgClose, err = windows.GetProcAddress(modCrypt32, "CryptMsgClose")
	if err != nil {
		return nil, fmt.Errorf("failed to get CryptMsgClose address: %v", err)
	}

	authenticode := Authenticode{}

	catalogFile, err := getCatalogPathForFilePath(path)
	if err != nil {
		catalogFile = path
	}

	fmt.Printf("Catalog file: %s\n", catalogFile)

	verificationResult, err := verifySignature(catalogFile)
	if err != nil {
		return nil, fmt.Errorf("failed to verify signature: %v", err)
	}

	authenticode.Result = verificationResult

	fmt.Printf("Authenticode: %s\n", authenticode)

	err = getSignatureInformation(catalogFile, &authenticode)
	if err != nil {
		return nil, fmt.Errorf("failed to get signature information: %v", err)
	}

	return []Authenticode{authenticode}, nil
}
