package logon_sessions

import (
	"fmt"
	"unsafe"

	"golang.org/x/sys/windows"
)

// LogonSession represents a Windows logon session
type LogonSession struct {
	LogonID               int32  `json:"logon_id"`
	User                  string `json:"user"`
	LogonDomain           string `json:"logon_domain"`
	AuthenticationPackage string `json:"authentication_package"`
	LogonType             string `json:"logon_type"`
	SessionID             int32  `json:"session_id"`
	LogonSID              string `json:"logon_sid"`
	LogonTime             int64  `json:"logon_time"`
	LogonServer           string `json:"logon_server"`
	DnsDomainName         string `json:"dns_domain_name"`
	UPN                   string `json:"upn"`
	LogonScript           string `json:"logon_script"`
	ProfilePath           string `json:"profile_path"`
	HomeDirectory         string `json:"home_directory"`
	HomeDirectoryDrive    string `json:"home_directory_drive"`
}

var logonTypes = map[uint32]string{
	0:  "Undefined Logon Type",
	2:  "Interactive",
	3:  "Network",
	4:  "Batch",
	5:  "Service",
	6:  "Proxy",
	7:  "Unlock",
	8:  "Network Cleartext",
	9:  "New Credentials",
	10: "Remote Interactive",
	11: "Cached Interactive",
	12: "Cached Remote Interactive",
	13: "Cached Unlock",
}

type SECURITY_LOGON_SESSION_DATA struct {
	Size                  uint32
	LogonID               windows.LUID
	UserName              windows.NTUnicodeString
	LogonDomain           windows.NTUnicodeString
	AuthenticationPackage windows.NTUnicodeString
	LogonType             uint32
	Session               uint32
	Sid                   *windows.SID
	LogonTime             windows.Filetime
	LogonServer           windows.NTUnicodeString
	DnsDomainName         windows.NTUnicodeString
	Upn                   windows.NTUnicodeString
}

func GenLogonSessions() ([]LogonSession, error) {
	var sessionCount uint32
	var sessionListPtr *windows.LUID

	// Load secur32.dll
	modSecur32 := windows.NewLazySystemDLL("secur32.dll")
	procLsaEnumerateLogonSessions := modSecur32.NewProc("LsaEnumerateLogonSessions")
	procLsaGetLogonSessionData := modSecur32.NewProc("LsaGetLogonSessionData")
	procLsaFreeReturnBuffer := modSecur32.NewProc("LsaFreeReturnBuffer")

	// Enumerate logon sessions
	ret, _, _ := procLsaEnumerateLogonSessions.Call(
		uintptr(unsafe.Pointer(&sessionCount)),
		uintptr(unsafe.Pointer(&sessionListPtr)),
	)
	if ret != 0 { // STATUS_SUCCESS = 0
		return nil, fmt.Errorf("failed to enumerate logon sessions: %x", ret)
	}
	// Free the session list
	defer procLsaFreeReturnBuffer.Call(uintptr(unsafe.Pointer(sessionListPtr)))

	sessionList := unsafe.Slice((*windows.LUID)(unsafe.Pointer(sessionListPtr)), sessionCount)

	var sessions []LogonSession

	for i := uint32(0); i < sessionCount; i++ {
		var sessionDataPtr *SECURITY_LOGON_SESSION_DATA
		ret, _, _ := procLsaGetLogonSessionData.Call(
			uintptr(unsafe.Pointer(&sessionList[i])),
			uintptr(unsafe.Pointer(&sessionDataPtr)),
		)
		if ret != 0 {
			fmt.Printf("LsaGetLogonSessionData failed: %x\n", ret)
			continue
		}

		sessionData := (*SECURITY_LOGON_SESSION_DATA)(unsafe.Pointer(sessionDataPtr))

		session := LogonSession{
			LogonID:               int32(sessionData.LogonID.LowPart),
			User:                  sessionData.UserName.String(),
			LogonDomain:           sessionData.LogonDomain.String(),
			AuthenticationPackage: sessionData.AuthenticationPackage.String(),
			LogonType:             logonTypes[sessionData.LogonType],
			SessionID:             int32(sessionData.Session),
			LogonSID:              sessionData.Sid.String(),
			LogonTime:             int64(sessionData.LogonTime.Nanoseconds() / 1e9),
			LogonServer:           sessionData.LogonServer.String(),
			DnsDomainName:         sessionData.DnsDomainName.String(),
			UPN:                   sessionData.Upn.String(),
		}

		sessions = append(sessions, session)
		procLsaFreeReturnBuffer.Call(uintptr(unsafe.Pointer(sessionDataPtr)))

	}

	return sessions, nil
}
