package process_open_sockets

import (
	"fmt"
	"net"
	"syscall"
	"unsafe"

	"github.com/scrymastic/goosquery/sql/context"
	"github.com/scrymastic/goosquery/util"
	"golang.org/x/sys/windows"
)

// Column definitions for the process_open_sockets table
var columnDefs = map[string]string{
	"pid":            "int32",
	"fd":             "int64",
	"socket":         "int64",
	"family":         "int32",
	"proto":          "int32",
	"local_address":  "string",
	"remote_address": "string",
	"local_port":     "int32",
	"remote_port":    "int32",
	"path":           "string",
	"state":          "string",
}

type MIB_TCPROW_OWNER_PID struct {
	DwState      uint32
	DwLocalAddr  uint32
	DwLocalPort  uint32
	DwRemoteAddr uint32
	DwRemotePort uint32
	DwOwningPid  uint32
}

type MIB_TCP6ROW_OWNER_PID struct {
	UcLocalAddr     [16]byte
	DwLocalScopeId  uint32
	DwLocalPort     uint32
	UcRemoteAddr    [16]byte
	DwRemoteScopeId uint32
	DwRemotePort    uint32
	DwState         uint32
	DwOwningPid     uint32
}

type MIB_UDPROW_OWNER_PID struct {
	DwLocalAddr uint32
	DwLocalPort uint32
	DwOwningPid uint32
}

type MIB_UDP6ROW_OWNER_PID struct {
	UcLocalAddr    [16]byte
	DwLocalScopeId uint32
	DwLocalPort    uint32
	DwOwningPid    uint32
}

var (
	procGetExtendedTcpTable *windows.LazyProc
	procGetExtendedUdpTable *windows.LazyProc
)

func init() {
	modIphlpapi := windows.NewLazySystemDLL("iphlpapi.dll")
	if modIphlpapi.Load() != nil {
		return
	}
	procGetExtendedTcpTable = modIphlpapi.NewProc("GetExtendedTcpTable")
	procGetExtendedUdpTable = modIphlpapi.NewProc("GetExtendedUdpTable")
}

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

const (
	UDP_TABLE_BASIC        = 0
	UDP_TABLE_OWNER_PID    = 1
	UDP_TABLE_OWNER_MODULE = 2
)

const (
	TCP_TABLE_BASIC_LISTENER           = 0
	TCP_TABLE_BASIC_CONNECTIONS        = 1
	TCP_TABLE_BASIC_ALL                = 2
	TCP_TABLE_OWNER_PID_LISTENER       = 3
	TCP_TABLE_OWNER_PID_CONNECTIONS    = 4
	TCP_TABLE_OWNER_PID_ALL            = 5
	TCP_TABLE_OWNER_MODULE_LISTENER    = 6
	TCP_TABLE_OWNER_MODULE_CONNECTIONS = 7
	TCP_TABLE_OWNER_MODULE_ALL         = 8
)

// TCP states mapping similar to osquery
func tcpStateToString(state uint32) string {
	if s, ok := tcpStateMap[state]; ok {
		return s
	}
	return fmt.Sprintf("UNKNOWN (%d)", state)
}

// Helper function to handle table allocation
func allocateTable(proc *windows.LazyProc, family uint32, class uint32) ([]byte, error) {
	var size uint32
	if ret, _, _ := proc.Call(
		0,
		uintptr(unsafe.Pointer(&size)),
		1, // true for sorted
		uintptr(family),
		uintptr(class),
		0,
	); syscall.Errno(ret) != windows.ERROR_INSUFFICIENT_BUFFER {
		return nil, fmt.Errorf("error getting table size: %v", ret)
	}

	table := make([]byte, size)
	if ret, _, _ := proc.Call(
		uintptr(unsafe.Pointer(&table[0])),
		uintptr(unsafe.Pointer(&size)),
		1, // true for sorted
		uintptr(family),
		uintptr(class),
		0,
	); syscall.Errno(ret) != windows.ERROR_SUCCESS {
		return nil, fmt.Errorf("error calling GetExtendedTable: %v", ret)
	}
	return table, nil
}

func allocateSocketTable(sockType string) ([]byte, error) {
	switch sockType {
	case "TCP":
		return allocateTable(procGetExtendedTcpTable, syscall.AF_INET, TCP_TABLE_OWNER_PID_ALL)
	case "TCP6":
		return allocateTable(procGetExtendedTcpTable, syscall.AF_INET6, TCP_TABLE_OWNER_PID_ALL)
	case "UDP":
		return allocateTable(procGetExtendedUdpTable, syscall.AF_INET, UDP_TABLE_OWNER_PID)
	case "UDP6":
		return allocateTable(procGetExtendedUdpTable, syscall.AF_INET6, UDP_TABLE_OWNER_PID)
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
	ip := net.IP(addr[:])
	return ip.String()
}

func formatIPv4Address(addr uint32) string {
	ip := make(net.IP, 4)
	ip[0] = byte(addr)
	ip[1] = byte(addr >> 8)
	ip[2] = byte(addr >> 16)
	ip[3] = byte(addr >> 24)
	return ip.String()
}

func parseSocketTable(sockType string, table []byte, ctx context.Context) ([]map[string]interface{}, error) {
	// Get the size of the TCP table
	DwNumEntries := *(*uint32)(unsafe.Pointer(&table[0]))

	switch sockType {
	case "TCP":
		// Get the first TCP row
		row := (*MIB_TCPROW_OWNER_PID)(unsafe.Pointer(&table[4]))

		// Parse the TCP table
		sockets := make([]map[string]interface{}, 0, DwNumEntries)
		for i := uint32(0); i < DwNumEntries; i++ {
			socket := util.InitColumns(ctx, columnDefs)

			if ctx.IsColumnUsed("pid") {
				socket["pid"] = int32(row.DwOwningPid)
			}
			if ctx.IsColumnUsed("fd") {
				socket["fd"] = int32(0)
			}
			if ctx.IsColumnUsed("socket") {
				socket["socket"] = int32(0)
			}
			if ctx.IsColumnUsed("family") {
				socket["family"] = int32(syscall.AF_INET)
			}
			if ctx.IsColumnUsed("proto") {
				socket["proto"] = int32(syscall.IPPROTO_TCP)
			}
			if ctx.IsColumnUsed("local_address") {
				socket["local_address"] = formatIPv4Address(row.DwLocalAddr)
			}
			if ctx.IsColumnUsed("remote_address") {
				socket["remote_address"] = formatIPv4Address(row.DwRemoteAddr)
			}
			if ctx.IsColumnUsed("local_port") {
				socket["local_port"] = int32(networkToHostPort(row.DwLocalPort))
			}
			if ctx.IsColumnUsed("remote_port") {
				socket["remote_port"] = int32(networkToHostPort(row.DwRemotePort))
			}
			if ctx.IsColumnUsed("path") {
				socket["path"] = ""
			}
			if ctx.IsColumnUsed("state") {
				socket["state"] = tcpStateToString(row.DwState)
			}

			sockets = append(sockets, socket)
			row = (*MIB_TCPROW_OWNER_PID)(unsafe.Pointer(uintptr(unsafe.Pointer(row)) + unsafe.Sizeof(*row)))
		}
		return sockets, nil

	case "TCP6":
		// Get the first TCP6 row
		row := (*MIB_TCP6ROW_OWNER_PID)(unsafe.Pointer(&table[4]))

		// Parse the TCP6 table
		sockets := make([]map[string]interface{}, 0, DwNumEntries)
		for i := uint32(0); i < DwNumEntries; i++ {
			socket := make(map[string]interface{})

			if ctx.IsColumnUsed("pid") {
				socket["pid"] = int32(row.DwOwningPid)
			}
			if ctx.IsColumnUsed("fd") {
				socket["fd"] = int32(0)
			}
			if ctx.IsColumnUsed("socket") {
				socket["socket"] = int32(0)
			}
			if ctx.IsColumnUsed("family") {
				socket["family"] = int32(syscall.AF_INET6)
			}
			if ctx.IsColumnUsed("proto") {
				socket["proto"] = int32(syscall.IPPROTO_TCP)
			}
			if ctx.IsColumnUsed("local_address") {
				socket["local_address"] = formatIPv6Address(row.UcLocalAddr)
			}
			if ctx.IsColumnUsed("remote_address") {
				socket["remote_address"] = formatIPv6Address(row.UcRemoteAddr)
			}
			if ctx.IsColumnUsed("local_port") {
				socket["local_port"] = int32(networkToHostPort(row.DwLocalPort))
			}
			if ctx.IsColumnUsed("remote_port") {
				socket["remote_port"] = int32(networkToHostPort(row.DwRemotePort))
			}
			if ctx.IsColumnUsed("path") {
				socket["path"] = ""
			}
			if ctx.IsColumnUsed("state") {
				socket["state"] = tcpStateToString(row.DwState)
			}
			if ctx.IsColumnUsed("net_namespace") {
				socket["net_namespace"] = ""
			}

			sockets = append(sockets, socket)
			row = (*MIB_TCP6ROW_OWNER_PID)(unsafe.Pointer(uintptr(unsafe.Pointer(row)) + unsafe.Sizeof(*row)))
		}
		return sockets, nil

	case "UDP":
		// Get the first UDP row
		row := (*MIB_UDPROW_OWNER_PID)(unsafe.Pointer(&table[4]))

		// Parse the UDP table
		sockets := make([]map[string]interface{}, 0, DwNumEntries)
		for i := uint32(0); i < DwNumEntries; i++ {
			socket := make(map[string]interface{})

			if ctx.IsColumnUsed("pid") {
				socket["pid"] = int32(row.DwOwningPid)
			}
			if ctx.IsColumnUsed("fd") {
				socket["fd"] = int32(0)
			}
			if ctx.IsColumnUsed("socket") {
				socket["socket"] = int32(0)
			}
			if ctx.IsColumnUsed("family") {
				socket["family"] = int32(syscall.AF_INET)
			}
			if ctx.IsColumnUsed("proto") {
				socket["proto"] = int32(syscall.IPPROTO_UDP)
			}
			if ctx.IsColumnUsed("local_address") {
				socket["local_address"] = formatIPv4Address(row.DwLocalAddr)
			}
			if ctx.IsColumnUsed("remote_address") {
				socket["remote_address"] = ""
			}
			if ctx.IsColumnUsed("local_port") {
				socket["local_port"] = int32(networkToHostPort(row.DwLocalPort))
			}
			if ctx.IsColumnUsed("remote_port") {
				socket["remote_port"] = int32(0)
			}
			if ctx.IsColumnUsed("path") {
				socket["path"] = ""
			}
			if ctx.IsColumnUsed("state") {
				socket["state"] = ""
			}
			if ctx.IsColumnUsed("net_namespace") {
				socket["net_namespace"] = ""
			}

			sockets = append(sockets, socket)
			row = (*MIB_UDPROW_OWNER_PID)(unsafe.Pointer(uintptr(unsafe.Pointer(row)) + unsafe.Sizeof(*row)))
		}
		return sockets, nil

	case "UDP6":
		// Get the first UDP6 row
		row := (*MIB_UDP6ROW_OWNER_PID)(unsafe.Pointer(&table[4]))

		// Parse the UDP6 table
		sockets := make([]map[string]interface{}, 0, DwNumEntries)
		for i := uint32(0); i < DwNumEntries; i++ {
			socket := make(map[string]interface{})

			if ctx.IsColumnUsed("pid") {
				socket["pid"] = int32(row.DwOwningPid)
			}
			if ctx.IsColumnUsed("fd") {
				socket["fd"] = int32(0)
			}
			if ctx.IsColumnUsed("socket") {
				socket["socket"] = int32(0)
			}
			if ctx.IsColumnUsed("family") {
				socket["family"] = int32(syscall.AF_INET6)
			}
			if ctx.IsColumnUsed("proto") {
				socket["proto"] = int32(syscall.IPPROTO_UDP)
			}
			if ctx.IsColumnUsed("local_address") {
				socket["local_address"] = formatIPv6Address(row.UcLocalAddr)
			}
			if ctx.IsColumnUsed("remote_address") {
				socket["remote_address"] = ""
			}
			if ctx.IsColumnUsed("local_port") {
				socket["local_port"] = int32(networkToHostPort(row.DwLocalPort))
			}
			if ctx.IsColumnUsed("remote_port") {
				socket["remote_port"] = int32(0)
			}
			if ctx.IsColumnUsed("path") {
				socket["path"] = ""
			}
			if ctx.IsColumnUsed("state") {
				socket["state"] = ""
			}
			if ctx.IsColumnUsed("net_namespace") {
				socket["net_namespace"] = ""
			}

			sockets = append(sockets, socket)
			row = (*MIB_UDP6ROW_OWNER_PID)(unsafe.Pointer(uintptr(unsafe.Pointer(row)) + unsafe.Sizeof(*row)))
		}
		return sockets, nil

	default:
		return nil, fmt.Errorf("unknown socket type: %s", sockType)
	}
}

// GenProcessOpenSockets returns a list of open sockets for each process
func GenProcessOpenSockets(ctx context.Context) ([]map[string]interface{}, error) {
	if procGetExtendedTcpTable == nil || procGetExtendedUdpTable == nil {
		return nil, fmt.Errorf("failed to initialize iphlpapi.dll")
	}

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
	tcpSockets, err := parseSocketTable("TCP", tcpTable, ctx)
	if err != nil {
		return nil, err
	}

	// Parse the TCP6 table
	tcp6Sockets, err := parseSocketTable("TCP6", tcp6Table, ctx)
	if err != nil {
		return nil, err
	}

	// Parse the UDP table
	udpSockets, err := parseSocketTable("UDP", udpTable, ctx)
	if err != nil {
		return nil, err
	}

	// Parse the UDP6 table
	udp6Sockets, err := parseSocketTable("UDP6", udp6Table, ctx)
	if err != nil {
		return nil, err
	}

	// Combine all sockets
	sockets := append(tcpSockets, tcp6Sockets...)
	sockets = append(sockets, udpSockets...)
	sockets = append(sockets, udp6Sockets...)

	return sockets, nil
}
