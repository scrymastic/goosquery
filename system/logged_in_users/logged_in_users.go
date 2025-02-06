package windows

import (
	"fmt"
	"net"
	"syscall"
	"unsafe"

	"golang.org/x/sys/windows"
)

// Windows API constants
const (
	WTS_CURRENT_SERVER_HANDLE = 0
	WTSSessionInfo            = 5
	WTSClientInfo             = 8
	WTSClientName             = 10
	WTSTypeSessionInfoLevel1  = 1
	AF_INET                   = 2
	AF_INET6                  = 23
	AF_UNSPEC                 = 0
	CLIENTADDRESS_LENGTH      = 31
)

var (
	procWTSEnumerateSessionsExW     uintptr
	procWTSFreeMemoryEx             uintptr
	procWTSQuerySessionInformationW uintptr
	procWTSFreeMemory               uintptr
	procLookupAccountNameW          uintptr
)

// _WTS_SESSION_INFO_1W contains information about a Terminal Services session
// typedef struct _WTS_SESSION_INFO_1W {
// 	DWORD                  ExecEnvId;
// 	WTS_CONNECTSTATE_CLASS State;
// 	DWORD                  SessionId;
// 	LPWSTR                 pSessionName;
// 	LPWSTR                 pHostName;
// 	LPWSTR                 pUserName;
// 	LPWSTR                 pDomainName;
// 	LPWSTR                 pFarmName;
//   } WTS_SESSION_INFO_1W, *PWTS_SESSION_INFO_1W;

type _WTS_SESSION_INFO_1W struct {
	ExecEnvId    uint32
	State        uint32
	SessionId    uint32
	pSessionName *uint16
	pHostName    *uint16
	pUserName    *uint16
	pDomain      *uint16
	pFarmName    *uint16
}

// _WTSINFO contains information about a Terminal Services session
type _WTSINFOW struct {
	State              uint32
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
	ConnectTime        windows.Filetime
	DisconnectTime     windows.Filetime
	LastInputTime      windows.Filetime
	LogonTime          windows.Filetime
	CurrentTime        windows.Filetime
}

// _WTSCLIENT contains information about a Terminal Services client
type _WTSCLIENT struct {
	ClientName          [32]uint16
	Domain              [17]uint16
	UserName            [21]uint16
	WorkDirectory       [260]uint16
	InitialProgram      [260]uint16
	EncryptionLevel     byte
	ClientAddressFamily uint32
	ClientAddress       [31]uint16
	ClientHardwareId    uint32
	ClientProductId     uint16
	OutBufCountHost     uint16
	OutBufCountClient   uint16
	OutBufLength        uint16
	DeviceId            [32]uint16
}

// LoggedInUser represents a logged-in user on the system
type LoggedInUser struct {
	User         string
	Type         string
	Tty          string
	Host         string
	Time         int64
	Pid          int
	Sid          string
	RegistryHive string
}

// sessionStates maps WTS states to their string representations
var sessionStates = map[int]string{
	windows.WTSActive:       "active",
	windows.WTSDisconnected: "disconnected",
	windows.WTSConnected:    "connected",
	windows.WTSConnectQuery: "connectquery",
	windows.WTSShadow:       "shadow",
	windows.WTSIdle:         "idle",
	windows.WTSListen:       "listen",
	windows.WTSReset:        "reset",
	windows.WTSDown:         "down",
	windows.WTSInit:         "init",
}

// Convert Windows FILETIME to Unix timestamp
func fileTimeToUnixTime(ft windows.Filetime) int64 {
	// Convert to int64 for nanoseconds calculation
	nsec := int64(ft.HighDateTime)<<32 + int64(ft.LowDateTime)
	// Adjust from Windows epoch to Unix epoch
	nsec = (nsec - 116444736000000000) * 100
	return nsec / 1e9 // Convert to seconds
}

// GenLoggedInUsers returns information about logged-in users
func GenLoggedInUsers() ([]LoggedInUser, error) {
	modWtsapi32, err := windows.LoadLibrary("wtsapi32.dll")
	if err != nil {
		return nil, fmt.Errorf("failed to load wtsapi32.dll: %w", err)
	}
	defer windows.FreeLibrary(modWtsapi32)

	modAdvapi32, err := windows.LoadLibrary("advapi32.dll")
	if err != nil {
		return nil, fmt.Errorf("failed to load advapi32.dll: %w", err)
	}
	defer windows.FreeLibrary(modAdvapi32)

	procWTSEnumerateSessionsExW, err = windows.GetProcAddress(modWtsapi32, "WTSEnumerateSessionsExW")
	if err != nil {
		return nil, fmt.Errorf("failed to get WTSEnumerateSessionsExW: %w", err)
	}

	procWTSFreeMemory, err = windows.GetProcAddress(modWtsapi32, "WTSFreeMemory")
	if err != nil {
		return nil, fmt.Errorf("failed to get WTSFreeMemory: %w", err)
	}

	procWTSQuerySessionInformationW, err = windows.GetProcAddress(modWtsapi32, "WTSQuerySessionInformationW")
	if err != nil {
		return nil, fmt.Errorf("failed to get WTSQuerySessionInformationW: %w", err)
	}

	procLookupAccountNameW, err = windows.GetProcAddress(modAdvapi32, "LookupAccountNameW")
	if err != nil {
		return nil, fmt.Errorf("failed to get LookupAccountNameW: %w", err)
	}

	var users []LoggedInUser
	level := uint32(1)
	count := uint32(0)
	var sessions *_WTS_SESSION_INFO_1W

	ret, _, err := syscall.SyscallN(procWTSEnumerateSessionsExW,
		uintptr(WTS_CURRENT_SERVER_HANDLE),
		uintptr(unsafe.Pointer(&level)),
		0,
		uintptr(unsafe.Pointer(&sessions)),
		uintptr(unsafe.Pointer(&count)),
	)

	if ret == 0 {
		return nil, fmt.Errorf("failed to enumerate sessions: %w", err)
	}

	defer syscall.SyscallN(procWTSFreeMemoryEx,
		uintptr(WTSTypeSessionInfoLevel1),
		uintptr(unsafe.Pointer(sessions)),
		uintptr(count))

	// Convert sessions pointer to slice
	sessionSlice := unsafe.Slice(sessions, count)

	for i := uint32(0); i < count; i++ {
		session := sessionSlice[i]
		// Skip non-active sessions and session 0 (system session)
		if session.State != windows.WTSActive || session.SessionId == 0 {
			continue
		}

		var sessionInfo *_WTSINFOW
		var bytesRet uint32
		ret, _, err = syscall.SyscallN(procWTSQuerySessionInformationW,
			uintptr(WTS_CURRENT_SERVER_HANDLE),
			uintptr(session.SessionId),
			uintptr(WTSSessionInfo),
			uintptr(unsafe.Pointer(&sessionInfo)),
			uintptr(unsafe.Pointer(&bytesRet)),
		)
		if ret == 0 || sessionInfo == nil {
			continue
		}

		user := LoggedInUser{
			User: windows.UTF16ToString(sessionInfo.UserName[:]),
			Type: sessionStates[int(session.State)],
			Tty:  windows.UTF16PtrToString(session.pSessionName),
			Time: fileTimeToUnixTime(sessionInfo.ConnectTime),
			Pid:  -1,
		}

		// Get client information
		var clientInfo *_WTSCLIENT
		ret, _, _ = syscall.SyscallN(procWTSQuerySessionInformationW,
			uintptr(WTS_CURRENT_SERVER_HANDLE),
			uintptr(session.SessionId),
			uintptr(WTSClientInfo),
			uintptr(unsafe.Pointer(&clientInfo)),
			uintptr(unsafe.Pointer(&bytesRet)),
		)

		if ret != 0 && clientInfo != nil {
			switch clientInfo.ClientAddressFamily {
			case AF_INET:
				// Convert byte array to IPv4 address
				addr := clientInfo.ClientAddress[:4]
				user.Host = fmt.Sprintf("%d.%d.%d.%d", addr[0], addr[1], addr[2], addr[3])
			case AF_INET6:
				// Convert uint16 array to bytes for IPv6 address
				addr := make([]byte, 16)
				for i := 0; i < 16; i++ {
					addr[i] = byte(clientInfo.ClientAddress[i])
				}
				user.Host = net.IP(addr).String()
			case AF_UNSPEC:
				var clientName *uint16
				ret, _, _ = syscall.SyscallN(procWTSQuerySessionInformationW,
					uintptr(WTS_CURRENT_SERVER_HANDLE),
					uintptr(session.SessionId),
					uintptr(WTSClientName),
					uintptr(unsafe.Pointer(&clientName)),
					uintptr(unsafe.Pointer(&bytesRet)),
				)
				if ret != 0 && clientName != nil {
					user.Host = windows.UTF16PtrToString(clientName)
					syscall.SyscallN(procWTSFreeMemory, uintptr(unsafe.Pointer(clientName)))
				}
			}
			syscall.SyscallN(procWTSFreeMemory, uintptr(unsafe.Pointer(clientInfo)))
		}

		// Get user SID
		domainUser := windows.UTF16ToString(sessionInfo.Domain[:]) + "\\" + windows.UTF16ToString(sessionInfo.UserName[:])
		sid, err := getSidFromAccountName(domainUser)
		if err == nil {
			user.Sid = sid.String()
			user.RegistryHive = "HKEY_USERS\\" + user.Sid
		}

		syscall.SyscallN(procWTSFreeMemory, uintptr(unsafe.Pointer(sessionInfo)))
		users = append(users, user)
	}

	return users, nil
}

// getSidFromAccountName converts a username to a SID
func getSidFromAccountName(accountName string) (*windows.SID, error) {
	var sidSize uint32
	var domainSize uint32
	var sidUse uint32

	// First call to determine the buffer sizes
	ret, _, err := syscall.SyscallN(procLookupAccountNameW,
		0,
		uintptr(unsafe.Pointer(windows.StringToUTF16Ptr(accountName))),
		0,
		uintptr(unsafe.Pointer(&sidSize)),
		0,
		uintptr(unsafe.Pointer(&domainSize)),
		uintptr(unsafe.Pointer(&sidUse)),
	)

	if ret == 0 && err != windows.ERROR_INSUFFICIENT_BUFFER {
		return nil, fmt.Errorf("LookupAccountNameW failed: %w", err)
	}

	// Allocate buffers
	sid := make([]byte, sidSize)
	domain := make([]uint16, domainSize)

	// Second call to actually retrieve the SID
	ret, _, err = syscall.SyscallN(procLookupAccountNameW,
		0,
		uintptr(unsafe.Pointer(windows.StringToUTF16Ptr(accountName))),
		uintptr(unsafe.Pointer(&sid[0])),
		uintptr(unsafe.Pointer(&sidSize)),
		uintptr(unsafe.Pointer(&domain[0])),
		uintptr(unsafe.Pointer(&domainSize)),
		uintptr(unsafe.Pointer(&sidUse)),
	)

	if ret == 0 {
		return nil, err
	}

	return (*windows.SID)(unsafe.Pointer(&sid[0])), nil
}
