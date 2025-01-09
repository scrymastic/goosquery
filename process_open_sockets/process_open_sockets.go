package process_open_sockets

import (
	"fmt"
	"syscall"
	"unsafe"

	"golang.org/x/sys/windows"
)

// ProcessOpenSocket represents the Windows process open sockets structure
type ProcessOpenSocket struct {
	PID           uint32 `json:"pid"`
	FD            uint32 `json:"fd"`
	Socket        uint32 `json:"socket"`
	Family        uint32 `json:"family"`
	Proto         uint32 `json:"proto"`
	LocalAddress  string `json:"local_address"`
	RemoteAddress string `json:"remote_address"`
	LocalPort     uint32 `json:"local_port"`
	RemotePort    uint32 `json:"remote_port"`
	Path          string `json:"path"`
	State         string `json:"state"`
	NetNamespace  string `json:"net_namespace"`
}

// type _MIB_TCPTABLE_OWNER_PID struct {
// 	DwNumEntries uint32
// 	Table        [0]_MIB_TCPROW_OWNER_PID
// }

type _MIB_TCPROW_OWNER_PID struct {
	DwState      uint32
	DwLocalAddr  uint32
	DwLocalPort  uint32
	DwRemoteAddr uint32
	DwRemotePort uint32
	DwOwningPid  uint32
}

// type _MIB_TCP6TABLE_OWNER_PID struct {
// 	DwNumEntries uint32
// 	Table        [0]_MIB_TCP6ROW_OWNER_PID
// }

type _MIB_TCP6ROW_OWNER_PID struct {
	UcLocalAddr     [16]byte
	DwLocalScopeId  uint32
	DwLocalPort     uint32
	UcRemoteAddr    [16]byte
	DwRemoteScopeId uint32
	DwRemotePort    uint32
	DwState         uint32
	DwOwningPid     uint32
}

// type _MIB_UDPTABLE_OWNER_PID struct {
// 	DwNumEntries uint32
// 	Table        [0]_MIB_UDPROW_OWNER_PID
// }

type _MIB_UDPROW_OWNER_PID struct {
	DwLocalAddr uint32
	DwLocalPort uint32
	DwOwningPid uint32
}

// type _MIB_UDP6TABLE_OWNER_PID struct {
// 	DwNumEntries uint32
// 	Table        [0]_MIB_UDP6ROW_OWNER_PID
// }

type _MIB_UDP6ROW_OWNER_PID struct {
	UcLocalAddr    [16]byte
	DwLocalScopeId uint32
	DwLocalPort    uint32
	DwOwningPid    uint32
}

var (
	modIphlpapi             windows.Handle
	procGetExtendedTcpTable uintptr
	procGetExtendedUdpTable uintptr
)

var (
	tcpStateMap = map[uint32]string{
		1:  "CLOSED",
		2:  "LISTEN",
		3:  "SYN_SENT",
		4:  "SYN_RCVD",
		5:  "ESTABLISHED",
		6:  "FIN_WAIT1",
		7:  "FIN_WAIT2",
		8:  "CLOSE_WAIT",
		9:  "CLOSING",
		10: "LAST_ACK",
		11: "TIME_WAIT",
		12: "DELETE_TCB",
	}
)

// TCP states mapping similar to osquery
func tcpStateToString(state uint32) string {
	if s, ok := tcpStateMap[state]; ok {
		return s
	}
	return fmt.Sprintf("UNKNOWN (%d)", state)
}

func allocateSocketTable(sockType string) ([]byte, error) {
	var err error
	// Load iphlpapi.dll
	modIphlpapi, err = windows.LoadLibrary("iphlpapi.dll")
	if err != nil {
		return nil, fmt.Errorf("error loading iphlpapi.dll: %v", err)
	}
	defer windows.FreeLibrary(modIphlpapi)

	// Get the GetExtendedTcpTable function
	procGetExtendedTcpTable, err = windows.GetProcAddress(modIphlpapi, "GetExtendedTcpTable")
	if err != nil {
		return nil, fmt.Errorf("error getting GetExtendedTcpTable function: %v", err)
	}

	// Get the GetExtendedUdpTable function
	procGetExtendedUdpTable, err = windows.GetProcAddress(modIphlpapi, "GetExtendedUdpTable")
	if err != nil {
		return nil, fmt.Errorf("error getting GetExtendedUdpTable function: %v", err)
	}

	var size uint32
	switch sockType {
	case "TCP":
		// TCP IPv4 table
		if ret, _, _ := syscall.SyscallN(procGetExtendedTcpTable,
			0,
			uintptr(unsafe.Pointer(&size)),
			1, // true for sorted
			syscall.AF_INET,
			5, // TCP_TABLE_OWNER_PID_ALL
			0,
		); ret != uintptr(windows.ERROR_INSUFFICIENT_BUFFER) {
			return nil, fmt.Errorf("error getting TCP table size: %v", ret)
		}

		// Allocate memory for the TCP table
		tcpTable := make([]byte, size)
		if ret, _, _ := syscall.SyscallN(procGetExtendedTcpTable,
			uintptr(unsafe.Pointer(&tcpTable[0])),
			uintptr(unsafe.Pointer(&size)),
			1, // true for sorted
			syscall.AF_INET,
			5, // TCP_TABLE_OWNER_PID_ALL
			0,
		); ret != uintptr(windows.ERROR_SUCCESS) {
			return nil, fmt.Errorf("error calling GetExtendedTcpTable: %v", ret)
		}
		return tcpTable, nil
	case "TCP6":
		// TCP IPv6 table
		if ret, _, _ := syscall.SyscallN(procGetExtendedTcpTable,
			0,
			uintptr(unsafe.Pointer(&size)),
			1, // true for sorted
			syscall.AF_INET6,
			5, // TCP_TABLE_OWNER_PID_ALL
			0,
		); ret != uintptr(windows.ERROR_INSUFFICIENT_BUFFER) {
			return nil, fmt.Errorf("error getting TCP6 table size: %v", ret)
		}

		// Allocate memory for the TCP6 table
		tcp6Table := make([]byte, size)
		if ret, _, _ := syscall.SyscallN(procGetExtendedTcpTable,
			uintptr(unsafe.Pointer(&tcp6Table[0])),
			uintptr(unsafe.Pointer(&size)),
			1, // true for sorted
			syscall.AF_INET6,
			5, // TCP_TABLE_OWNER_PID_ALL
			0,
		); ret != uintptr(windows.ERROR_SUCCESS) {
			return nil, fmt.Errorf("error calling GetExtendedTcpTable: %v", ret)
		}
		return tcp6Table, nil

	case "UDP":
		// UDP IPv4 table
		if ret, _, _ := syscall.SyscallN(procGetExtendedUdpTable,
			0,
			uintptr(unsafe.Pointer(&size)),
			1, // true for sorted
			syscall.AF_INET,
			1, // UDP_TABLE_OWNER_PID
			0,
		); ret != uintptr(windows.ERROR_INSUFFICIENT_BUFFER) {
			return nil, fmt.Errorf("error getting UDP table size: %v", ret)
		}

		// Allocate memory for the UDP table
		udpTable := make([]byte, size)
		if ret, _, _ := syscall.SyscallN(procGetExtendedUdpTable,
			uintptr(unsafe.Pointer(&udpTable[0])),
			uintptr(unsafe.Pointer(&size)),
			1, // true for sorted
			syscall.AF_INET,
			1, // UDP_TABLE_OWNER_PID
			0,
		); ret != uintptr(windows.ERROR_SUCCESS) {
			return nil, fmt.Errorf("error calling GetExtendedUdpTable: %v", ret)
		}
		return udpTable, nil

	case "UDP6":
		// UDP IPv6 table
		if ret, _, _ := syscall.SyscallN(procGetExtendedUdpTable,
			0,
			uintptr(unsafe.Pointer(&size)),
			1, // true for sorted
			syscall.AF_INET6,
			1, // UDP_TABLE_OWNER_PID
			0,
		); ret != uintptr(windows.ERROR_INSUFFICIENT_BUFFER) {
			return nil, fmt.Errorf("error getting UDP6 table size: %v", ret)
		}

		// Allocate memory for the UDP6 table
		udp6Table := make([]byte, size)
		if ret, _, _ := syscall.SyscallN(procGetExtendedUdpTable,
			uintptr(unsafe.Pointer(&udp6Table[0])),
			uintptr(unsafe.Pointer(&size)),
			1, // true for sorted
			syscall.AF_INET6,
			1, // UDP_TABLE_OWNER_PID
			0,
		); ret != uintptr(windows.ERROR_SUCCESS) {
			return nil, fmt.Errorf("error calling GetExtendedUdpTable: %v", ret)
		}
		return udp6Table, nil

	default:
		return nil, fmt.Errorf("unknown socket type: %s", sockType)
	}
}

// Convert network byte order (big-endian) to host byte order
func networkToHostPort(port uint32) uint32 {
	return ((port & 0xFF) << 8) | ((port & 0xFF00) >> 8)
}

// formatIPv6Address formats a 16-byte IPv6 address into proper string representation
func formatIPv6Address(addr [16]byte) string {
	return fmt.Sprintf("%02x%02x:%02x%02x:%02x%02x:%02x%02x:%02x%02x:%02x%02x:%02x%02x:%02x%02x",
		addr[0], addr[1], addr[2], addr[3],
		addr[4], addr[5], addr[6], addr[7],
		addr[8], addr[9], addr[10], addr[11],
		addr[12], addr[13], addr[14], addr[15])
}

func formatIPv4Address(addr uint32) string {
	return fmt.Sprintf("%d.%d.%d.%d", byte(addr), byte(addr>>8), byte(addr>>16), byte(addr>>24))
}

func parseSocketTable(sockType string, table []byte) ([]ProcessOpenSocket, error) {
	// Get the size of the TCP table
	DwNumEntries := *(*uint32)(unsafe.Pointer(&table[0]))

	switch sockType {
	case "TCP":
		// Get the first TCP row
		row := (*_MIB_TCPROW_OWNER_PID)(unsafe.Pointer(&table[4]))

		// Parse the TCP table
		sockets := make([]ProcessOpenSocket, 0, DwNumEntries)
		for i := uint32(0); i < DwNumEntries; i++ {
			sockets = append(sockets, ProcessOpenSocket{
				PID:           row.DwOwningPid,
				FD:            0,
				Socket:        0,
				Family:        syscall.AF_INET,
				Proto:         syscall.IPPROTO_TCP,
				LocalAddress:  formatIPv4Address(row.DwLocalAddr),
				RemoteAddress: formatIPv4Address(row.DwRemoteAddr),
				LocalPort:     networkToHostPort(row.DwLocalPort),
				RemotePort:    networkToHostPort(row.DwRemotePort),
				Path:          "",
				State:         tcpStateToString(row.DwState),
				NetNamespace:  "",
			})
			row = (*_MIB_TCPROW_OWNER_PID)(unsafe.Pointer(uintptr(unsafe.Pointer(row)) + unsafe.Sizeof(*row)))
		}
		return sockets, nil

	case "TCP6":
		// Get the first TCP6 row
		row := (*_MIB_TCP6ROW_OWNER_PID)(unsafe.Pointer(&table[4]))

		// Parse the TCP6 table
		sockets := make([]ProcessOpenSocket, 0, DwNumEntries)
		for i := uint32(0); i < DwNumEntries; i++ {
			sockets = append(sockets, ProcessOpenSocket{
				PID:           row.DwOwningPid,
				FD:            0,
				Socket:        0,
				Family:        syscall.AF_INET6,
				Proto:         syscall.IPPROTO_TCP,
				LocalAddress:  formatIPv6Address(row.UcLocalAddr),
				RemoteAddress: formatIPv6Address(row.UcRemoteAddr),
				LocalPort:     networkToHostPort(row.DwLocalPort),
				RemotePort:    networkToHostPort(row.DwRemotePort),
				Path:          "",
				State:         tcpStateToString(row.DwState),
				NetNamespace:  "",
			})
			row = (*_MIB_TCP6ROW_OWNER_PID)(unsafe.Pointer(uintptr(unsafe.Pointer(row)) + unsafe.Sizeof(*row)))
		}
		return sockets, nil

	case "UDP":
		// Get the first UDP row
		row := (*_MIB_UDPROW_OWNER_PID)(unsafe.Pointer(&table[4]))

		// Parse the UDP table
		sockets := make([]ProcessOpenSocket, 0, DwNumEntries)
		for i := uint32(0); i < DwNumEntries; i++ {
			sockets = append(sockets, ProcessOpenSocket{
				PID:           row.DwOwningPid,
				FD:            0,
				Socket:        0,
				Family:        syscall.AF_INET,
				Proto:         syscall.IPPROTO_UDP,
				LocalAddress:  formatIPv4Address(row.DwLocalAddr),
				RemoteAddress: "",
				LocalPort:     networkToHostPort(row.DwLocalPort),
				RemotePort:    0,
				Path:          "",
				State:         "",
				NetNamespace:  "",
			})
			row = (*_MIB_UDPROW_OWNER_PID)(unsafe.Pointer(uintptr(unsafe.Pointer(row)) + unsafe.Sizeof(*row)))
		}
		return sockets, nil

	case "UDP6":
		// Get the first UDP6 row
		row := (*_MIB_UDP6ROW_OWNER_PID)(unsafe.Pointer(&table[4]))

		// Parse the UDP6 table
		sockets := make([]ProcessOpenSocket, 0, DwNumEntries)
		for i := uint32(0); i < DwNumEntries; i++ {
			sockets = append(sockets, ProcessOpenSocket{
				PID:           row.DwOwningPid,
				FD:            0,
				Socket:        0,
				Family:        syscall.AF_INET6,
				Proto:         syscall.IPPROTO_UDP,
				LocalAddress:  formatIPv6Address(row.UcLocalAddr),
				RemoteAddress: "",
				LocalPort:     networkToHostPort(row.DwLocalPort),
				RemotePort:    0,
				Path:          "",
				State:         "",
				NetNamespace:  "",
			})
			row = (*_MIB_UDP6ROW_OWNER_PID)(unsafe.Pointer(uintptr(unsafe.Pointer(row)) + unsafe.Sizeof(*row)))
		}
		return sockets, nil

	default:
		return nil, fmt.Errorf("unknown socket type: %s", sockType)
	}
}

// GenProcessOpenSockets returns a list of open sockets for each process
func GenProcessOpenSockets() ([]ProcessOpenSocket, error) {

	// Allocate memory for the TCP table
	tcpTable, err := allocateSocketTable("TCP")
	if err != nil {
		return nil, err
	}

	// Allocate memory for the TCP6 table
	tcp6Table, err := allocateSocketTable("TCP6")
	if err != nil {
		return nil, err
	}

	// Allocate memory for the UDP table
	udpTable, err := allocateSocketTable("UDP")
	if err != nil {
		return nil, err
	}

	// Allocate memory for the UDP6 table
	udp6Table, err := allocateSocketTable("UDP6")
	if err != nil {
		return nil, err
	}

	// Parse the TCP table
	tcpSockets, err := parseSocketTable("TCP", tcpTable)
	if err != nil {
		return nil, err
	}

	// Parse the TCP6 table
	tcp6Sockets, err := parseSocketTable("TCP6", tcp6Table)
	if err != nil {
		return nil, err
	}

	// Parse the UDP table
	udpSockets, err := parseSocketTable("UDP", udpTable)
	if err != nil {
		return nil, err
	}

	// Parse the UDP6 table
	udp6Sockets, err := parseSocketTable("UDP6", udp6Table)
	if err != nil {
		return nil, err
	}

	// Combine all sockets
	sockets := append(tcpSockets, tcp6Sockets...)
	sockets = append(sockets, udpSockets...)
	sockets = append(sockets, udp6Sockets...)

	return sockets, nil
}
