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
	Netmask     int32  `json:"netmask"`
	Gateway     string `json:"gateway"`
	Source      string `json:"source"`
	Flags       int32  `json:"flags"`
	Interface   string `json:"interface"`
	MTU         int32  `json:"mtu"`
	Metric      int32  `json:"metric"`
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
	Prefix       SOCKADDR_INET
	PrefixLength uint8
	_            [3]byte
}

type MIB_IPFORWARD_ROW2 struct {
	InterfaceLuid        uint64
	InterfaceIndex       uint32
	DestinationPrefix    IP_ADDRESS_PREFIX
	NextHop              SOCKADDR_INET
	SitePrefixLength     uint8
	ValidLifetime        uint32
	PreferredLifetime    uint32
	Metric               uint32
	Protocol             uint32
	Loopback             uint8
	AutoconfigureAddress uint8
	Publish              uint8
	Immortal             uint8
	Age                  uint32
	Origin               uint32
}

type MIB_IPFORWARD_TABLE2 struct {
	NumEntries uint32
	_          [4]byte // padding
	Table      [1]MIB_IPFORWARD_ROW2
}

const (
	INET_ADDRSTRLEN  = 22
	INET6_ADDRSTRLEN = 65
)

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

// cStringToGoString converts a null-terminated C string to a Go string
func cStringToGoString(cString []byte) string {
	// Find the null terminator
	end := 0
	for i, b := range cString {
		if b == 0 {
			end = i
			break
		}
	}
	// If no null terminator found, use the whole string
	if end == 0 && len(cString) > 0 && cString[0] != 0 {
		end = len(cString)
	}
	return string(cString[:end])
}

// GenRoutes returns all routes on the system
func GenRoutes() ([]Route, error) {
	modIphlpapi := windows.NewLazySystemDLL("iphlpapi.dll")
	if modIphlpapi.Load() != nil {
		return nil, fmt.Errorf("failed to load iphlpapi.dll")
	}
	procGetIpForwardTable2 := modIphlpapi.NewProc("GetIpForwardTable2")
	procFreeMibTable := modIphlpapi.NewProc("FreeMibTable")
	procGetIpInterfaceEntry := modIphlpapi.NewProc("GetIpInterfaceEntry")
	modWs2_32 := windows.NewLazySystemDLL("ws2_32.dll")
	if modWs2_32.Load() != nil {
		return nil, fmt.Errorf("failed to load ws2_32.dll")
	}
	procInetNtopW := modWs2_32.NewProc("InetNtopW")

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

	// cast ipTablePtr to MIB_IPFORWARD_TABLE2 struct
	numEntries := ipTablePtr.NumEntries
	table := unsafe.Slice((*MIB_IPFORWARD_ROW2)(unsafe.Pointer(&ipTablePtr.Table[0])), numEntries)

	for _, currRow := range table {
		// fmt.Printf("Index: %v\n", row.InterfaceIndex)
		var route Route
		var actualInterface windows.MibIpInterfaceRow
		var interfaceIpAddress string
		actualInterface.Family = currRow.DestinationPrefix.Prefix.GetFamily()
		actualInterface.InterfaceLuid = currRow.InterfaceLuid
		actualInterface.InterfaceIndex = currRow.InterfaceIndex

		ret, _, _ := procGetIpInterfaceEntry.Call(
			uintptr(unsafe.Pointer(&actualInterface)),
		)
		if windows.Errno(ret) != windows.NO_ERROR {
			// return nil, fmt.Errorf("GetIpInterfaceEntry failed: %v", ret)
			route.Metric = -1
			route.MTU = -1
		} else {
			route.Metric = int32(actualInterface.Metric + currRow.Metric)
			route.MTU = int32(actualInterface.NlMtu)
		}

		addrFamily := actualInterface.Family

		if addrFamily == windows.AF_INET6 {
			route.Type = "local"
			ipAddress := currRow.DestinationPrefix.Prefix.GetIpv6().Sin6Addr
			gateway := currRow.NextHop.GetIpv6().Sin6Addr
			buffer := make([]uint16, INET6_ADDRSTRLEN)

			ret, _, _ := procInetNtopW.Call(
				uintptr(addrFamily),
				uintptr(unsafe.Pointer(&ipAddress)),
				uintptr(unsafe.Pointer(&buffer[0])),
				uintptr(INET6_ADDRSTRLEN),
			)

			if ret != 0 {
				route.Destination = windows.UTF16ToString(buffer)
			}

			ret, _, _ = procInetNtopW.Call(
				uintptr(addrFamily),
				uintptr(unsafe.Pointer(&gateway)),
				uintptr(unsafe.Pointer(&buffer[0])),
				uintptr(INET6_ADDRSTRLEN),
			)

			if ret != 0 {
				route.Gateway = windows.UTF16ToString(buffer)
			}

		} else if addrFamily == windows.AF_INET {
			ipAddress := currRow.DestinationPrefix.Prefix.GetIpv4().SinAddr
			gateway := currRow.NextHop.GetIpv4().SinAddr
			buffer := make([]uint16, INET_ADDRSTRLEN)

			ret, _, _ = procInetNtopW.Call(
				uintptr(addrFamily),
				uintptr(unsafe.Pointer(&ipAddress)),
				uintptr(unsafe.Pointer(&buffer[0])),
				uintptr(INET_ADDRSTRLEN),
			)

			if ret != 0 {
				route.Destination = windows.UTF16ToString(buffer)
			}

			adapters, err := getAdapterAddressMapping()
			if err != nil {
				return nil, fmt.Errorf("failed to get adapter mapping: %v", err)
			}
			if actualAdapter, ok := adapters[currRow.InterfaceIndex]; ok {
				interfaceIpAddress = cStringToGoString(actualAdapter.IpAddressList.IpAddress.String[:])
				route.Gateway = cStringToGoString(actualAdapter.GatewayList.IpAddress.String[:])
			} else {
				interfaceIpAddress = "127.0.0.1"
				ret, _, _ = procInetNtopW.Call(
					uintptr(addrFamily),
					uintptr(unsafe.Pointer(&gateway)),
					uintptr(unsafe.Pointer(&buffer[0])),
					uintptr(INET_ADDRSTRLEN),
				)
				if ret != 0 {
					route.Gateway = windows.UTF16ToString(buffer)
				}
			}
			if currRow.Loopback == 1 {
				route.Type = "local"
			} else {
				route.Type = "remote"
			}
		}

		route.Interface = interfaceIpAddress
		route.Netmask = int32(currRow.DestinationPrefix.PrefixLength)
		route.Flags = -1
		routes = append(routes, route)
	}
	return routes, nil
}

func HexDump(pointer unsafe.Pointer, size uint32) {
	data := unsafe.Slice((*byte)(pointer), size)
	hexDump := hex.Dump(data)
	fmt.Println(hexDump)
}
