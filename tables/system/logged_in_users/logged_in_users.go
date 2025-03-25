package logged_in_users

import (
	"fmt"
	"log"
	"unsafe"

	"golang.org/x/sys/windows"

	"github.com/scrymastic/goosquery/sql/result"
	"github.com/scrymastic/goosquery/sql/sqlctx"
)

var (
	modWtsapi32              = windows.NewLazySystemDLL("wtsapi32.dll")
	procWtsEnumerateSessions = modWtsapi32.NewProc("WTSEnumerateSessionsExW")
	procWtsQuerySessionInfo  = modWtsapi32.NewProc("WTSQuerySessionInformationW")
	procWtsFreeMemory        = modWtsapi32.NewProc("WTSFreeMemory")
	procWtsFreeMemoryEx      = modWtsapi32.NewProc("WTSFreeMemoryExW")
	modAdvapi32              = windows.NewLazySystemDLL("advapi32.dll")
	procGetSidByName         = modAdvapi32.NewProc("LookupAccountNameW")
)

const (
	WTS_CURRENT_SERVER_HANDLE = 0
	WTSActive                 = 0

	WTSSessionInfo = 0x18
	WTSClientInfo  = 0x17
)

var sessionStates = map[int32]string{
	0: "active",
	1: "connected",
	2: "connectquery",
	3: "shadow",
	4: "disconnected",
	5: "idle",
	6: "listen",
	7: "reset",
	8: "down",
	9: "init",
}

// WTS_SESSION_INFO_1W structure
type WTS_SESSION_INFO_1W struct {
	ExecEnvId   uint32
	State       int32
	SessionId   uint32
	SessionName *uint16
	HostName    *uint16
	UserName    *uint16
	DomainName  *uint16
	FarmName    *uint16
}

// WTSINFOW structure
type WTSINFOW struct {
	State              int32
	SessionId          uint32
	IncomingBytes      uint32
	OutgoingBytes      uint32
	IncomingFrames     uint32
	OutgoingFrames     uint32
	IncomingCompressed uint32
	OutgoingCompressed uint32
	WinStationName     [32]uint16
	Domain             [17]uint16
	UserName           [21]uint16
	ConnectTime        int64
	DisconnectTime     int64
	LastInputTime      int64
	LogonTime          int64
	CurrentTime        int64
}

// WTSCLIENTW structure
type WTSCLIENTW struct {
	ClientName          [21]uint16
	Domain              [18]uint16
	UserName            [21]uint16
	WorkDirectory       [261]uint16
	InitialProgram      [261]uint16
	EncryptionLevel     byte
	ClientAddressFamily uint32
	ClientAddress       [31]uint16
	HRes                uint16
	VRes                uint16
	ColorDepth          uint16
	ClientDirectory     [261]uint16
	ClientBuildNumber   uint32
	ClientHardwareId    uint32
	ClientProductId     uint16
	OutBufCountHost     uint16
	OutBufCountClient   uint16
	OutBufLength        uint16
	DeviceId            [261]uint16
}

func filetimeToUnix(filetime int64) int64 {
	// Convert Windows FILETIME (100-nanosecond intervals since January 1, 1601 UTC)
	// to Unix time (seconds since January 1, 1970 UTC)
	return (filetime - 116444736000000000) / 10000000
}

func getSid(domain, username string) (string, error) {
	var sidSize uint32
	var domainSize uint32
	var sidUse uint32

	// First call to get required buffer sizes
	// domainPtr, _ := syscall.UTF16PtrFromString(domain)
	usernamePtr, _ := windows.UTF16PtrFromString(username)
	_, _, err := procGetSidByName.Call(
		0, // lpSystemName (NULL for local system)
		uintptr(unsafe.Pointer(usernamePtr)),
		0, // Sid buffer (NULL for size query)
		uintptr(unsafe.Pointer(&sidSize)),
		0, // ReferencedDomainName (NULL for size query)
		uintptr(unsafe.Pointer(&domainSize)),
		uintptr(unsafe.Pointer(&sidUse)),
	)

	if sidSize == 0 {
		return "", fmt.Errorf("failed to get SID size: %v", err)
	}

	// Allocate buffers
	sid := make([]byte, sidSize)
	referencedDomain := make([]uint16, domainSize)

	// Second call to get the actual SID
	ret, _, err := procGetSidByName.Call(
		0,
		uintptr(unsafe.Pointer(usernamePtr)),
		uintptr(unsafe.Pointer(&sid[0])),
		uintptr(unsafe.Pointer(&sidSize)),
		uintptr(unsafe.Pointer(&referencedDomain[0])),
		uintptr(unsafe.Pointer(&domainSize)),
		uintptr(unsafe.Pointer(&sidUse)),
	)
	if ret == 0 {
		return "", fmt.Errorf("failed to get SID: %v", err)
	}

	// Convert SID to string
	var sidString *uint16
	modAdvapi32.NewProc("ConvertSidToStringSidW").Call(
		uintptr(unsafe.Pointer(&sid[0])),
		uintptr(unsafe.Pointer(&sidString)),
	)
	defer windows.LocalFree(windows.Handle(unsafe.Pointer(sidString)))
	return windows.UTF16PtrToString(sidString), nil
}

func GenLoggedInUsers(ctx *sqlctx.Context) (*result.Results, error) {
	results := result.NewQueryResult()
	var count uint32
	var pSessionInfo *WTS_SESSION_INFO_1W
	level := uint32(1)

	// Enumerate sessions
	ret, _, err := procWtsEnumerateSessions.Call(
		uintptr(WTS_CURRENT_SERVER_HANDLE),
		uintptr(unsafe.Pointer(&level)),
		0,
		uintptr(unsafe.Pointer(&pSessionInfo)),
		uintptr(unsafe.Pointer(&count)),
	)
	if ret == 0 {
		return nil, fmt.Errorf("failed to enumerate sessions: %v", err)
	}
	defer procWtsFreeMemoryEx.Call(1, uintptr(unsafe.Pointer(pSessionInfo)), uintptr(count))

	sessions := unsafe.Slice((*WTS_SESSION_INFO_1W)(unsafe.Pointer(pSessionInfo)), count)

	for _, session := range sessions {
		if session.State != WTSActive || session.SessionId == 0 {
			continue
		}

		user := result.NewResult(ctx, Schema)

		// Get session info
		var sessionInfo *WTSINFOW
		var bytesRet uint32
		ret, _, _ := procWtsQuerySessionInfo.Call(
			uintptr(WTS_CURRENT_SERVER_HANDLE),
			uintptr(session.SessionId),
			uintptr(WTSSessionInfo),
			uintptr(unsafe.Pointer(&sessionInfo)),
			uintptr(unsafe.Pointer(&bytesRet)),
		)
		if ret != 0 {
			wts := (*WTSINFOW)(unsafe.Pointer(sessionInfo))
			user.Set("user", windows.UTF16PtrToString(session.UserName))
			user.Set("type", sessionStates[session.State])
			user.Set("tty", windows.UTF16PtrToString(session.SessionName))
			if wts.ConnectTime != 0 {
				user.Set("time", filetimeToUnix(wts.ConnectTime))
			}
			procWtsFreeMemory.Call(uintptr(unsafe.Pointer(sessionInfo)))
		}

		if ctx.IsAnyOfColumnsUsed([]string{"host", "sid", "registry_hive"}) {
			// Get client info
			var clientInfo *WTSCLIENTW
			ret, _, _ = procWtsQuerySessionInfo.Call(
				uintptr(WTS_CURRENT_SERVER_HANDLE),
				uintptr(session.SessionId),
				uintptr(WTSClientInfo),
				uintptr(unsafe.Pointer(&clientInfo)),
				uintptr(unsafe.Pointer(&bytesRet)),
			)
			if ret != 0 {
				client := (*WTSCLIENTW)(unsafe.Pointer(clientInfo))
				user.Set("host", windows.UTF16ToString(client.ClientName[:]))
				procWtsFreeMemory.Call(uintptr(unsafe.Pointer(clientInfo)))
			}

			// Get SID
			if user.Get("user") != "" {
				domain := windows.UTF16PtrToString(session.DomainName)
				if sid, err := getSid(domain, user.Get("user").(string)); err == nil {
					user.Set("sid", sid)
					user.Set("registry_hive", "HKEY_USERS\\"+sid)
				} else {
					log.Printf("Failed to get SID for %s: %v", user.Get("user"), err)
				}
			}

			results.AppendResult(*user)
		}
	}

	return results, nil
}
