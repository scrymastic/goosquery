package logon_sessions

import (
	"fmt"
	"unsafe"

	"golang.org/x/sys/windows"

	"github.com/scrymastic/goosquery/sql/result"
	"github.com/scrymastic/goosquery/sql/sqlctx"
)

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

func GenLogonSessions(ctx *sqlctx.Context) (*result.Results, error) {
	results := result.NewQueryResult()
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

		session := result.NewResult(ctx, Schema)
		session.Set("logon_id", int32(sessionData.LogonID.LowPart))
		session.Set("user", sessionData.UserName.String())
		session.Set("logon_domain", sessionData.LogonDomain.String())
		session.Set("authentication_package", sessionData.AuthenticationPackage.String())
		session.Set("logon_type", logonTypes[sessionData.LogonType])
		session.Set("session_id", int32(sessionData.Session))
		session.Set("logon_sid", sessionData.Sid.String())
		session.Set("logon_time", int64(sessionData.LogonTime.Nanoseconds()/1e9))
		session.Set("logon_server", sessionData.LogonServer.String())
		session.Set("dns_domain_name", sessionData.DnsDomainName.String())
		session.Set("upn", sessionData.Upn.String())

		results.AppendResult(*session)
		procLsaFreeReturnBuffer.Call(uintptr(unsafe.Pointer(sessionDataPtr)))

	}

	return results, nil
}
