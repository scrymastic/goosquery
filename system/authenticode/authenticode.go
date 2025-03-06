package authenticode

import (
	"fmt"
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
	modWintrust                             = windows.NewLazyDLL("wintrust.dll")
	procCryptCATAdminAcquireContext         = modWintrust.NewProc("CryptCATAdminAcquireContext")
	procCryptCATAdminCalcHashFromFileHandle = modWintrust.NewProc("CryptCATAdminCalcHashFromFileHandle")
	procCryptCATAdminEnumCatalogFromHash    = modWintrust.NewProc("CryptCATAdminEnumCatalogFromHash")
	procCryptCATCatalogInfoFromContext      = modWintrust.NewProc("CryptCATCatalogInfoFromContext")
	procCryptCATAdminReleaseCatalogContext  = modWintrust.NewProc("CryptCATAdminReleaseCatalogContext")
	procCryptCATAdminReleaseContext         = modWintrust.NewProc("CryptCATAdminReleaseContext")
	procWinVerifyTrust                      = modWintrust.NewProc("WinVerifyTrust")
	modCrypt32                              = windows.NewLazyDLL("crypt32.dll")
	procCryptQueryObject                    = modCrypt32.NewProc("CryptQueryObject")
	procCryptMsgGetParam                    = modCrypt32.NewProc("CryptMsgGetParam")
	procCryptDecodeObject                   = modCrypt32.NewProc("CryptDecodeObject")
	procCryptMsgClose                       = modCrypt32.NewProc("CryptMsgClose")
)

var (
	DRIVER_ACTION_VERIFY              = windows.GUID(*ole.NewGUID("{f750e6c3-38ee-11d1-85e5-00c04fc295ee}"))
	WINTRUST_ACTION_GENERIC_VERIFY_V2 = windows.GUID(*ole.NewGUID("{00aac56b-cd44-11d0-8cc2-00c04fc295ee}"))
)

const (
	SHA512_HASH_SIZE       = 64
	CMSG_SIGNER_INFO_PARAM = 6
	SPC_SP_OPUS_INFO_OBJID = "1.3.6.1.4.1.311.2.1.12"
)

type CATALOG_INFO struct {
	cbStruct       uint32
	wszCatalogFile [windows.MAX_PATH]uint16
}

type CRYPT_ATTRIBUTE struct {
	pszObjId *byte
	cValue   uint32
	rgValue  *windows.CryptDataBlob
}

type CRYPT_ATTRIBUTES struct {
	cAttr  uint32
	rgAttr *CRYPT_ATTRIBUTE
}

type CMSG_SIGNER_INFO struct {
	dwVersion               uint32
	Issuer                  windows.CertNameBlob
	SerialNumber            windows.CryptIntegerBlob
	HashAlgorithm           windows.CryptAlgorithmIdentifier
	HashEncryptionAlgorithm windows.CryptAlgorithmIdentifier
	EncryptedHash           windows.CryptDataBlob
	AuthAttrs               CRYPT_ATTRIBUTES
	UnauthAttrs             CRYPT_ATTRIBUTES
}

type SPC_SP_OPUS_INFO struct {
	pwszProgramName *uint16
	pMoreInfo       *byte
	pPublisherInfo  *byte
}

func getOriginalProgramName(signerInfo *CMSG_SIGNER_INFO) (string, error) {
	var publisherInfoPtr *CRYPT_ATTRIBUTE

	// Search through auth attributes for SPC_SP_OPUS_INFO_OBJID
	for i := uint32(0); i < signerInfo.AuthAttrs.cAttr; i++ {
		attr := (*CRYPT_ATTRIBUTE)(
			unsafe.Pointer(
				uintptr(unsafe.Pointer(signerInfo.AuthAttrs.rgAttr)) +
					uintptr(i)*unsafe.Sizeof(CRYPT_ATTRIBUTE{}),
			),
		)
		objId := windows.BytePtrToString(attr.pszObjId)
		if objId == SPC_SP_OPUS_INFO_OBJID {
			publisherInfoPtr = attr
			break
		}
	}

	if publisherInfoPtr == nil {
		return "", fmt.Errorf("failed to find publisher information")
	}

	var publisherInfoSize uint32
	ret, _, err := procCryptDecodeObject.Call(
		uintptr(windows.X509_ASN_ENCODING|windows.PKCS_7_ASN_ENCODING),
		uintptr(unsafe.Pointer(publisherInfoPtr.pszObjId)),
		uintptr(unsafe.Pointer(publisherInfoPtr.rgValue.Data)),
		uintptr(publisherInfoPtr.rgValue.Size),
		0,
		0,
		uintptr(unsafe.Pointer(&publisherInfoSize)))

	if ret == 0 {
		return "", fmt.Errorf("failed to access the publisher information: %v", err)
	}

	publisherInfoBlob := make([]byte, publisherInfoSize)
	ret, _, err = procCryptDecodeObject.Call(
		uintptr(windows.X509_ASN_ENCODING|windows.PKCS_7_ASN_ENCODING),
		uintptr(unsafe.Pointer(publisherInfoPtr.pszObjId)),
		uintptr(unsafe.Pointer(publisherInfoPtr.rgValue.Data)),
		uintptr(publisherInfoPtr.rgValue.Size),
		0,
		uintptr(unsafe.Pointer(&publisherInfoBlob[0])),
		uintptr(unsafe.Pointer(&publisherInfoSize)))

	if ret == 0 {
		return "", fmt.Errorf("failed to decode the publisher information: %v", err)
	}

	// Cast to SPC_SP_OPUS_INFO
	publisherInfo := (*SPC_SP_OPUS_INFO)(unsafe.Pointer(&publisherInfoBlob[0]))

	programName := windows.UTF16PtrToString(publisherInfo.pwszProgramName)

	return programName, nil
}

func getSignatureInformation(path string, authenticode *Authenticode) error {
	// Convert path to UTF16
	uft16Path := windows.StringToUTF16Ptr(path)

	var certStore windows.Handle
	var message windows.Handle

	ret, _, err := procCryptQueryObject.Call(
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
		return fmt.Errorf("failed to query signature: %v", err)
	}
	defer windows.CertCloseStore(certStore, 0)
	// Use syscall for CryptMsgClose since it's not in windows package
	defer procCryptMsgClose.Call(uintptr(message))

	var signerInfoSize uint32
	ret, _, err = procCryptMsgGetParam.Call(
		uintptr(message),
		uintptr(CMSG_SIGNER_INFO_PARAM),
		0,
		0,
		uintptr(unsafe.Pointer(&signerInfoSize)),
	)
	if ret == 0 {
		return fmt.Errorf("failed to get signer info size: %v", err)
	}

	signerInfo := make([]byte, signerInfoSize)
	signerInfoPtr := (*CMSG_SIGNER_INFO)(unsafe.Pointer(&signerInfo[0]))
	ret, _, err = procCryptMsgGetParam.Call(
		uintptr(message),
		uintptr(CMSG_SIGNER_INFO_PARAM),
		0,
		uintptr(unsafe.Pointer(signerInfoPtr)),
		uintptr(unsafe.Pointer(&signerInfoSize)),
	)
	if ret == 0 {
		return fmt.Errorf("failed to get signer info: %v", err)
	}

	// Get original program name
	originalProgramName, err := getOriginalProgramName(signerInfoPtr)
	if err != nil {
		// Just log the error and continue, as this is not critical
		fmt.Printf("failed to get original program name: %v\n", err)
	}
	authenticode.OriginalProgramName = originalProgramName

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
		return fmt.Errorf("cert not found: %v", err)
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
	// First call to get size
	issuerNameSize := uint32(windows.CertGetNameString(
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

	// First call to get size
	subjectNameSize := uint32(windows.CertGetNameString(
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
	uft16Path := windows.StringToUTF16Ptr(path)

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
		return "", fmt.Errorf("CreateFile: failed to open file: %v", err)
	}
	defer windows.CloseHandle(handle)

	// Acquire catalog admin context
	var context windows.Handle
	ret, _, _ := procCryptCATAdminAcquireContext.Call(
		uintptr(unsafe.Pointer(&context)),
		uintptr(unsafe.Pointer(&DRIVER_ACTION_VERIFY)),
		0,
	)
	if ret == 0 {
		return "", fmt.Errorf("failed to get context: %v", err)
	}
	defer procCryptCATAdminReleaseContext.Call(
		uintptr(context),
		0,
	)

	// Calculate file hash
	hash := make([]byte, SHA512_HASH_SIZE)
	hashSize := uint32(len(hash))

	ret, _, err = procCryptCATAdminCalcHashFromFileHandle.Call(
		uintptr(handle),
		uintptr(unsafe.Pointer(&hashSize)),
		uintptr(unsafe.Pointer(&hash[0])),
		0,
	)
	if ret == 0 {
		return "", fmt.Errorf("failed to calc hash: %v", err)
	}

	// Find matching catalog
	var catalog windows.Handle
	ret, _, err = procCryptCATAdminEnumCatalogFromHash.Call(
		uintptr(context),
		uintptr(unsafe.Pointer(&hash[0])),
		uintptr(hashSize),
		0,
		0,
	)
	if ret == 0 {
		return "", fmt.Errorf("no matching catalog: %v", err)
	}
	catalog = windows.Handle(ret)

	defer procCryptCATAdminReleaseCatalogContext.Call(
		uintptr(context),
		uintptr(catalog),
		0,
	)

	// Get catalog info
	var info CATALOG_INFO
	info.cbStruct = uint32(unsafe.Sizeof(info))

	ret, _, err = procCryptCATCatalogInfoFromContext.Call(
		uintptr(catalog),
		uintptr(unsafe.Pointer(&info)),
		0,
	)
	if ret == 0 {
		return "", fmt.Errorf("failed to get info: %v", err)
	}

	// Convert catalog path from UTF16
	catalogFile := windows.UTF16ToString(info.wszCatalogFile[:])
	return catalogFile, nil
}

func verifySignature(path string) (string, error) {
	uft16Path := windows.StringToUTF16Ptr(path)

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

	authenticodePolicyProvider := WINTRUST_ACTION_GENERIC_VERIFY_V2
	ret, _, _ := procWinVerifyTrust.Call(
		uintptr(0),
		uintptr(unsafe.Pointer(&authenticodePolicyProvider)),
		uintptr(unsafe.Pointer(&trustProviderSettings)),
	)

	trustProviderSettings.StateAction = windows.WTD_STATEACTION_CLOSE
	defer procWinVerifyTrust.Call(
		uintptr(0),
		uintptr(unsafe.Pointer(&authenticodePolicyProvider)),
		uintptr(unsafe.Pointer(&trustProviderSettings)),
	)

	verificationStatus := int(ret)

	switch verificationStatus {
	case int(windows.TRUST_E_EXPLICIT_DISTRUST):
		return "distrusted", nil
	case int(windows.TRUST_E_NOSIGNATURE):
		return "missing", nil
	case int(windows.ERROR_SUCCESS):
		return "trusted", nil
	case int(windows.CRYPT_E_SECURITY_SETTINGS):
		return "valid", nil
	case int(windows.TRUST_E_SUBJECT_NOT_TRUSTED):
		return "untrusted", nil
	default:
		return "unknown", fmt.Errorf("unknown verification status: %v", verificationStatus)
	}
}

func GenAuthenticode(path string) ([]Authenticode, error) {
	authenticode := Authenticode{}

	authenticode.Path = path

	catalogFile, err := getCatalogPathForFilePath(path)
	if err != nil {
		catalogFile = path
	}

	verificationResult, err := verifySignature(catalogFile)
	if err != nil {
		return nil, fmt.Errorf("failed to verify signature: %v", err)
	}

	authenticode.Result = verificationResult

	if verificationResult == "missing" {
		return []Authenticode{authenticode}, nil
	}

	err = getSignatureInformation(catalogFile, &authenticode)
	if err != nil {
		return nil, fmt.Errorf("failed to get cert details: %v", err)
	}

	return []Authenticode{authenticode}, nil
}
