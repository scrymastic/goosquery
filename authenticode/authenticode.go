package authenticode

import (
	"crypto/x509"
	"errors"

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

func verifySignature(result uint32, path string) uint32 {
	var trustProviderSettings windows.WinTrustData
	trustProviderSettings.StructSize = uint32(unsafe.Sizeof(trustProviderSettings))
}

func GenAuthenticode(path string) ([]Authenticode, error) {
	return nil, errors.New("not implemented")
}
