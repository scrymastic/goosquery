package certificates

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"syscall"

	"osquery/system/users"

	"golang.org/x/sys/windows"
)

type Certificate struct {
	CommonName       string `json:"common_name"`
	Subject          string `json:"subject"`
	Issuer           string `json:"issuer"`
	CA               bool   `json:"ca"`
	SelfSigned       bool   `json:"self_signed"`
	NotValidBefore   int64  `json:"not_valid_before"`
	NotValidAfter    int64  `json:"not_valid_after"`
	SigningAlgorithm string `json:"signing_algorithm"`
	KeyAlgorithm     string `json:"key_algorithm"`
	KeyStrength      uint32 `json:"key_strength"`
	KeyUsage         string `json:"key_usage"`
	SubjectKeyID     string `json:"subject_key_id"`
	AuthorityKeyID   string `json:"authority_key_id"`
	SHA1             string `json:"sha1"`
	Path             string `json:"path"`
	Serial           string `json:"serial"`
	SID              string `json:"sid"`
	StoreLocation    string `json:"store_location"`
	Store            string `json:"store"`
	StoreID          string `json:"store_id"`
	Username         string `json:"username"`
}

func GenPersonalCertsFromDisk() ([]Certificate, error) {
	var results []Certificate

	// Get all users from the system
	users, err := users.GenUsers()
	if err != nil {
		return nil, fmt.Errorf("failed to get users: %w", err)
	}

	for _, user := range users {
		// Get user's home directory
		homeDir := user.Directory
		if homeDir == "" {
			continue
		}

		// Construct certificates path
		certsPath := filepath.Join(homeDir,
			"AppData",
			"Roaming",
			"Microsoft",
			"SystemCertificates",
			"My",
			"Certificates")

		// Read all files in the certificates directory
		files, err := os.ReadDir(certsPath)
		if err != nil {
			log.Printf("Error reading directory %s: %v", certsPath, err)
			continue
		}

		for _, file := range files {
			if file.IsDir() {
				continue
			}
		}
	}

	return results, nil
}

func GenNonPersonalCerts() ([]Certificate, error) {
	var results []Certificate
	// Open the system store
	store, err := windows.CertOpenSystemStore(0, windows.StringToUTF16Ptr("ROOT"))
	if err != nil {
		return nil, fmt.Errorf("failed to open ROOT system store: %w", err)
	}
	defer windows.CertCloseStore(store, 0)

	// Get certificates from ROOT store
	rootCerts, err := getCertsFromStore(store, "ROOT", "LocalMachine")
	if err != nil {
		return nil, fmt.Errorf("failed to get ROOT certificates: %w", err)
	}
	results = append(results, rootCerts...)

	// Open and get certificates from CA store
	store, err = windows.CertOpenSystemStore(0, windows.StringToUTF16Ptr("CA"))
	if err != nil {
		return nil, fmt.Errorf("failed to open CA system store: %w", err)
	}
	defer windows.CertCloseStore(store, 0)

	caCerts, err := getCertsFromStore(store, "CA", "LocalMachine")
	if err != nil {
		return nil, fmt.Errorf("failed to get CA certificates: %w", err)
	}
	results = append(results, caCerts...)

	// Open and get certificates from MY store
	store, err = windows.CertOpenSystemStore(0, windows.StringToUTF16Ptr("MY"))
	if err != nil {
		return nil, fmt.Errorf("failed to open MY system store: %w", err)
	}
	defer windows.CertCloseStore(store, 0)

	myCerts, err := getCertsFromStore(store, "MY", "LocalMachine")
	if err != nil {
		return nil, fmt.Errorf("failed to get MY certificates: %w", err)
	}
	results = append(results, myCerts...)

	return results, nil
}

func getCertsFromStore(store windows.Handle, storeName string, storeLocation string) ([]Certificate, error) {
	var results []Certificate
	var cert *windows.CertContext
	var err error

	for {
		cert, err = windows.CertEnumCertificatesInStore(store, cert)
		if err != nil {
			if errno, ok := err.(syscall.Errno); ok {
				if errno == syscall.Errno(0x80092004) { // CRYPT_E_NOT_FOUND
					break
				}
			}
			return nil, fmt.Errorf("failed to enumerate certificates: %w", err)
		}

		if cert == nil {
			break
		}

		// Get certificate info
		certInfo := cert.CertInfo

		// Create Certificate struct
		certificate := Certificate{
			CommonName:       getCertName(cert),
			Subject:          getCertSubject(cert),
			Issuer:           getCertIssuer(cert),
			CA:               isCACert(cert),
			SelfSigned:       isSelfSigned(cert),
			NotValidBefore:   getNotBefore(certInfo),
			NotValidAfter:    getNotAfter(certInfo),
			SigningAlgorithm: getSigningAlgorithm(certInfo),
			KeyAlgorithm:     getKeyAlgorithm(certInfo),
			KeyStrength:      getKeyStrength(certInfo),
			KeyUsage:         getKeyUsage(certInfo),
			SubjectKeyID:     getSubjectKeyID(cert),
			AuthorityKeyID:   getAuthorityKeyID(cert),
			SHA1:             getSHA1(cert),
			Path:             fmt.Sprintf("%s\\%s", storeLocation, storeName),
			Serial:           getSerial(certInfo),
			StoreLocation:    storeLocation,
			Store:            storeName,
		}

		results = append(results, certificate)
	}

	return results, nil
}

// Helper functions to extract certificate information
func getCertName(cert *windows.CertContext) string {
	// Implementation needed
	return ""
}

func getCertSubject(cert *windows.CertContext) string {
	// Implementation needed
	return ""
}

func getCertIssuer(cert *windows.CertContext) string {
	// Implementation needed
	return ""
}

func isCACert(cert *windows.CertContext) bool {
	// Implementation needed
	return false
}

func isSelfSigned(cert *windows.CertContext) bool {
	// Implementation needed
	return false
}

func getNotBefore(certInfo *windows.CertInfo) int64 {
	// Implementation needed
	return 0
}

func getNotAfter(certInfo *windows.CertInfo) int64 {
	// Implementation needed
	return 0
}

func getSigningAlgorithm(certInfo *windows.CertInfo) string {
	// Implementation needed
	return ""
}

func getKeyAlgorithm(certInfo *windows.CertInfo) string {
	// Implementation needed
	return ""
}

func getKeyStrength(certInfo *windows.CertInfo) uint32 {
	// Implementation needed
	return 0
}

func getKeyUsage(certInfo *windows.CertInfo) string {
	// Implementation needed
	return ""
}

func getSubjectKeyID(cert *windows.CertContext) string {
	// Implementation needed
	return ""
}

func getAuthorityKeyID(cert *windows.CertContext) string {
	// Implementation needed
	return ""
}

func getSHA1(cert *windows.CertContext) string {
	// Implementation needed
	return ""
}

func getSerial(certInfo *windows.CertInfo) string {
	// Implementation needed
	return ""
}
