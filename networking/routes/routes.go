package routes

import (
	"encoding/hex"
	"fmt"
	"unsafe"

	"golang.org/x/sys/windows"
)

// Route represents a single route entry
type Route struct {
	Destination string `json:"destination"`
	Netmask     int    `json:"netmask"`
	Gateway     string `json:"gateway"`
	Source      string `json:"source"`
	Flags       int    `json:"flags"`
	Interface   string `json:"interface"`
	MTU         int    `json:"mtu"`
	Metric      int    `json:"metric"`
	Type        string `json:"type"`
}

type SOCKADDR_IN struct {
	SinFamily uint16
	SinPort   uint16
	SinAddr   [4]byte
	SinZero   [8]byte
}

type SOCKADDR_IN6 struct {
	Sin6Family uint16
	Sin6Port   uint16
	Sin6Addr   [16]byte
	Sin6Zero   [8]byte
}

// type SOCKADDR_INET struct {
// 	Ipv4     SOCKADDR_IN
// 	Ipv6     SOCKADDR_IN6
// 	SiFamily uint16
// }

type SOCKADDR_INET [28]byte

func (s *SOCKADDR_INET) GetFamily() uint16 {
	return *(*uint16)(unsafe.Pointer(&s[0]))
}

func (s *SOCKADDR_INET) GetIpv4() SOCKADDR_IN {
	return *(*SOCKADDR_IN)(unsafe.Pointer(&s[0]))
}

func (s *SOCKADDR_INET) GetIpv6() SOCKADDR_IN6 {
	return *(*SOCKADDR_IN6)(unsafe.Pointer(&s[0]))
}

type IP_ADDRESS_PREFIX struct {
	Prefix       [28]byte
	PrefixLength uint8
}

type MIB_IPFORWARD_ROW2 struct {
	InterfaceLuid        windows.LUID
	InterfaceIndex       uint32
	DestinationPrefix    IP_ADDRESS_PREFIX
	NextHop              [16]byte
	SitePrefixLength     uint8
	_                    [3]byte
	ValidLifetime        uint32
	PreferredLifetime    uint32
	Metric               uint32
	Protocol             [4]byte
	Loopback             uint8
	_                    [3]byte
	AutoconfigureAddress uint8
	_                    [3]byte
	Publish              uint8
	_                    [3]byte
	Immortal             uint8
	_                    [3]byte
	Age                  uint32
	Origin               [4]byte
}

type MIB_IPFORWARD_TABLE2 struct {
	NumEntries uint32
	_          [4]byte // padding
	Table      [1]MIB_IPFORWARD_ROW2
}

// getAdapterAddressMapping returns a map of adapter indices to their information
func getAdapterAddressMapping() (map[uint32]*windows.IpAdapterInfo, error) {
	var size uint32
	result := make(map[uint32]*windows.IpAdapterInfo)

	// First call to get buffer size
	err := windows.GetAdaptersInfo(nil, &size)
	if err != nil && err != windows.ERROR_BUFFER_OVERFLOW {
		return nil, fmt.Errorf("GetAdaptersInfo failed: %v", err)
	}

	// Allocate buffer
	buffer := make([]byte, size)
	adapterInfo := (*windows.IpAdapterInfo)(unsafe.Pointer(&buffer[0]))

	// Second call to get actual data
	err = windows.GetAdaptersInfo(adapterInfo, &size)
	if err != nil {
		return nil, fmt.Errorf("GetAdaptersInfo failed: %v", err)
	}

	// Build the map
	for curr := adapterInfo; curr != nil; curr = curr.Next {
		result[curr.Index] = curr
	}

	return result, nil
}

// GenRoutes returns all routes on the system
func GenRoutes() ([]Route, error) {
	modIphlpapi := windows.NewLazyDLL("iphlpapi.dll")
	procGetIpForwardTable2 := modIphlpapi.NewProc("GetIpForwardTable2")
	procFreeMibTable := modIphlpapi.NewProc("FreeMibTable")
	// procGetIpInterfaceEntry := modIphlpapi.NewProc("GetIpInterfaceEntry")
	// procGetAdaptersInfo := modIphlpapi.NewProc("GetAdaptersInfo")

	var routes []Route
	var ipTablePtr *MIB_IPFORWARD_TABLE2

	// Get the IP forwarding table
	ret, _, _ := procGetIpForwardTable2.Call(
		uintptr(windows.AF_UNSPEC),
		uintptr(unsafe.Pointer(&ipTablePtr)),
	)
	if windows.Errno(ret) != windows.NO_ERROR {
		return nil, fmt.Errorf("GetIpForwardTable2 failed: %v", ret)
	}
	defer procFreeMibTable.Call(uintptr(unsafe.Pointer(ipTablePtr)))

	// Get adapter information
	_, err := getAdapterAddressMapping()
	if err != nil {
		return nil, fmt.Errorf("failed to get adapter mapping: %v", err)
	}

	// cast ipTablePtr to MIB_IPFORWARD_TABLE2 struct
	numEntries := ipTablePtr.NumEntries
	table := unsafe.Slice((*MIB_IPFORWARD_ROW2)(unsafe.Pointer(&ipTablePtr.Table[0])), numEntries)

	for _, row := range table {
		fmt.Printf("Index: %v\n", row.InterfaceIndex)
		var actualInterface MIB_IPINTERFACE_ROW
		// actualInterface.Family = row.DestinationPrefix.Prefix.SiFamily
		actualInterface.InterfaceLuid = row.InterfaceLuid
		actualInterface.InterfaceIndex = row.InterfaceIndex

		// result := GetIpInterfaceEntry(&actualInterface)
	}
	return routes, nil
}

func HexDump(pointer unsafe.Pointer, size uint32) {
	data := unsafe.Slice((*byte)(pointer), size)
	hexDump := hex.Dump(data)
	fmt.Println(hexDump)
}
