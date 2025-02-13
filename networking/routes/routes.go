package routes

import (
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

type _in_addr struct {
	_ [4]byte
}

type SOCKADDR_IN struct {
	SinFamily uint16
	SinPort   uint16
	SinAddr   _in_addr
	SinZero   [8]byte
}

type _in6_addr struct {
	_ [16]byte
}

type SOCKADDR_IN6 struct {
	Sin6Family uint16
	Sin6Port   uint16
	Sin6Addr   _in6_addr
	Sin6Zero   [8]byte
}

type SOCKADDR_INET struct {
	Ipv4     SOCKADDR_IN
	Ipv6     SOCKADDR_IN6
	SiFamily uint16
}

type IP_ADDRESS_PREFIX struct {
	Prefix       SOCKADDR_INET
	PrefixLength uint8
}

type _MIB_IPFORWARD_ROW2 struct {
	InterfaceLuid        windows.LUID
	InterfaceIndex       uint32
	DestinationPrefix    IP_ADDRESS_PREFIX
	NextHop              SOCKADDR_INET
	SitePrefixLength     uint8
	ValidLifetime        uint32
	PreferredLifetime    uint32
	Metric               uint32
	Protocol             uint32
	Loopback             bool
	AutoconfigureAddress bool
	Publish              bool
	Immortal             bool
	Age                  uint32
	Origin               uint32
}

type MIB_IPFORWARD_TABLE2 struct {
	NumEntries uint32
	Table      [5]_MIB_IPFORWARD_ROW2
}

type adapterInfo struct {
	windows.IpAdapterInfo
	Next *adapterInfo
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
	var ipTable MIB_IPFORWARD_TABLE2

	// Get the IP forwarding table
	ret, _, _ := procGetIpForwardTable2.Call(
		uintptr(windows.AF_UNSPEC),
		uintptr(unsafe.Pointer(&ipTable)),
	)
	if windows.Errno(ret) != windows.NO_ERROR {
		return nil, fmt.Errorf("GetIpForwardTable2 failed: %v", ret)
	}
	defer procFreeMibTable.Call(uintptr(unsafe.Pointer(&ipTable)))

	// Get adapter information
	_, err := getAdapterAddressMapping()
	if err != nil {
		return nil, fmt.Errorf("failed to get adapter mapping: %v", err)
	}

	return routes, nil
}
