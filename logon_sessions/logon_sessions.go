package logon_sessions

import (
	"fmt"
	"syscall"
	"unsafe"

	"golang.org/x/sys/windows"
)

type LogonSession struct {
	LogonID               uint32 `json:"logon_id"`
	User                  string `json:"user"`
	LogonDomain           string `json:"logon_domain"`
	AuthenticationPackage string `json:"authentication_package"`
	LogonType             string `json:"logon_type"`
	SessionID             uint32 `json:"session_id"`
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

var (
	secur32                       windows.Handle
	procLsaEnumerateLogonSessions uintptr
	procLsaGetLogonSessionData    uintptr
	procLsaFreeReturnBuffer       uintptr
)

type _LUID struct {
	LowPart  uint32
	HighPart int32
}

type _LSA_LAST_INTER_LOGON_INFO struct {
	LastSuccessfulLogon                int64
	LastFailedLogon                    int64
	FailedAttemptCountSinceLastSuccess int32
}

type _SECURITY_LOGON_SESSION_DATA struct {
	Size                  uint32
	LogonId               _LUID
	UserName              windows.NTUnicodeString
	LogonDomain           windows.NTUnicodeString
	AuthenticationPackage windows.NTUnicodeString
	LogonType             uint32
	Session               uint32
	Sid                   *windows.SID
	LogonTime             int64
	LogonServer           windows.NTUnicodeString
	DnsDomainName         windows.NTUnicodeString
	Upn                   windows.NTUnicodeString
	UserFlags             uint32
	LastLogonInfo         _LSA_LAST_INTER_LOGON_INFO
	LogonScript           windows.NTUnicodeString
	ProfilePath           windows.NTUnicodeString
	HomeDirectory         windows.NTUnicodeString
	HomeDirectoryDrive    windows.NTUnicodeString
	LogoffTime            int64
	KickOffTime           int64
	PasswordLastSet       int64
	PasswordCanChange     int64
	PasswordMustChange    int64
}

var securityLogonTypes = map[uint32]string{
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

func fileTimeToUnixTime(fileTime uint64) int64 {
	// Windows FILETIME is in 100-nanosecond intervals since January 1, 1601 UTC
	// Need to convert to Unix timestamp (seconds since January 1, 1970 UTC)
	const (
		WINDOWS_TICK      = 100         // nanoseconds
		SEC_TO_UNIX_EPOCH = 11644473600 // seconds between 1601 and 1970
	)

	return int64((fileTime*WINDOWS_TICK)/1e9) - SEC_TO_UNIX_EPOCH
}

func GenLogonSessions() ([]LogonSession, error) {
	var sessionCount uint32
	var sessions *_LUID

	var err error
	// Load secur32.dll
	if secur32, err = windows.LoadLibrary("secur32.dll"); err != nil {
		return nil, fmt.Errorf("error loading secur32.dll: %v", err)
	}
	defer windows.FreeLibrary(secur32)

	// Get the LsaEnumerateLogonSessions function
	if procLsaEnumerateLogonSessions, err = windows.GetProcAddress(secur32, "LsaEnumerateLogonSessions"); err != nil {
		return nil, fmt.Errorf("error getting LsaEnumerateLogonSessions function: %v", err)
	}

	// Get the LsaGetLogonSessionData function
	if procLsaGetLogonSessionData, err = windows.GetProcAddress(secur32, "LsaGetLogonSessionData"); err != nil {
		return nil, fmt.Errorf("error getting LsaGetLogonSessionData function: %v", err)
	}

	// Get the LsaFreeReturnBuffer function
	if procLsaFreeReturnBuffer, err = windows.GetProcAddress(secur32, "LsaFreeReturnBuffer"); err != nil {
		return nil, fmt.Errorf("error getting LsaFreeReturnBuffer function: %v", err)
	}

	// Call LsaEnumerateLogonSessions to get the number of logon sessions
	if ret, _, _ := syscall.SyscallN(uintptr(procLsaEnumerateLogonSessions),
		uintptr(unsafe.Pointer(&sessionCount)),
		uintptr(unsafe.Pointer(&sessions)),
	); ret != uintptr(windows.ERROR_SUCCESS) {
		return nil, fmt.Errorf("error calling LsaEnumerateLogonSessions: %v", ret)
	}

	defer syscall.SyscallN(uintptr(procLsaFreeReturnBuffer),
		uintptr(unsafe.Pointer(sessions)),
	)

	sessionsSlice := unsafe.Slice(sessions, sessionCount)

	var logonSessions []LogonSession
	for i := uint32(0); i < sessionCount; i++ {
		var sessionData *_SECURITY_LOGON_SESSION_DATA
		if ret, _, _ := syscall.SyscallN(uintptr(procLsaGetLogonSessionData),
			uintptr(unsafe.Pointer(&sessionsSlice[i])),
			uintptr(unsafe.Pointer(&sessionData)),
		); ret != uintptr(windows.ERROR_SUCCESS) {
			return nil, fmt.Errorf("error calling LsaGetLogonSessionData: %v", ret)
		}

		securityLogonType := securityLogonTypes[sessionData.LogonType]
		if securityLogonType == "" {
			securityLogonType = fmt.Sprintf("Unknown (%d)", sessionData.LogonType)
		}

		logonSession := LogonSession{
			LogonID:               sessionData.LogonId.LowPart,
			User:                  windows.UTF16PtrToString(sessionData.UserName.Buffer),
			LogonDomain:           windows.UTF16PtrToString(sessionData.LogonDomain.Buffer),
			AuthenticationPackage: windows.UTF16PtrToString(sessionData.AuthenticationPackage.Buffer),
			LogonType:             securityLogonType,
			SessionID:             sessionData.Session,
			LogonSID:              sessionData.Sid.String(),
			LogonTime:             fileTimeToUnixTime(uint64(sessionData.LogonTime)),
			LogonServer:           windows.UTF16PtrToString(sessionData.LogonServer.Buffer),
			DnsDomainName:         windows.UTF16PtrToString(sessionData.DnsDomainName.Buffer),
			UPN:                   windows.UTF16PtrToString(sessionData.Upn.Buffer),
			LogonScript:           windows.UTF16PtrToString(sessionData.LogonScript.Buffer),
			ProfilePath:           windows.UTF16PtrToString(sessionData.ProfilePath.Buffer),
			HomeDirectory:         windows.UTF16PtrToString(sessionData.HomeDirectory.Buffer),
			HomeDirectoryDrive:    windows.UTF16PtrToString(sessionData.HomeDirectoryDrive.Buffer),
		}

		logonSessions = append(logonSessions, logonSession)
	}

	return logonSessions, nil
}
