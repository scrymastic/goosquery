package routes

import (
	"encoding/hex"
	"fmt"
	"math"
	"unsafe"

	"golang.org/x/sys/windows"
)

// Route represents a single route entry
type Route struct {
	Destination string `json:"destination"`
	Netmask     uint32 `json:"netmask"`
	Gateway     string `json:"gateway"`
	Source      string `json:"source"`
	Flags       uint32 `json:"flags"`
	Interface   string `json:"interface"`
	MTU         uint32 `json:"mtu"`
	Metric      uint32 `json:"metric"`
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
	Prefix       SOCKADDR_INET
	PrefixLength uint8
	_            [3]byte
}

type MIB_IPFORWARD_ROW2 struct {
	InterfaceLuid        windows.LUID
	InterfaceIndex       uint32
	DestinationPrefix    IP_ADDRESS_PREFIX
	NextHop              SOCKADDR_INET
	SitePrefixLength     uint8
	_                    [3]byte
	ValidLifetime        uint32
	PreferredLifetime    uint32
	Metric               uint32
	Protocol             [4]byte
	Loopback             uint8
	AutoconfigureAddress uint8
	Publish              uint8
	Immortal             uint8
	Age                  uint32
	Origin               [4]byte
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

// GenRoutes returns all routes on the system
func GenRoutes() ([]Route, error) {
	modIphlpapi := windows.NewLazyDLL("iphlpapi.dll")
	procGetIpForwardTable2 := modIphlpapi.NewProc("GetIpForwardTable2")
	procFreeMibTable := modIphlpapi.NewProc("FreeMibTable")
	procGetIpInterfaceEntry := modIphlpapi.NewProc("GetIpInterfaceEntry")
	modWs2_32 := windows.NewLazyDLL("ws2_32.dll")
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
		var actualInterface MIB_IPINTERFACE_ROW
		var interfaceIpAddress string
		actualInterface.Family = currRow.DestinationPrefix.Prefix.GetFamily()
		actualInterface.InterfaceLuid = currRow.InterfaceLuid
		actualInterface.InterfaceIndex = currRow.InterfaceIndex

		ret, _, _ := procGetIpInterfaceEntry.Call(
			uintptr(unsafe.Pointer(&actualInterface)),
		)
		if windows.Errno(ret) != windows.NO_ERROR {
			// return nil, fmt.Errorf("GetIpInterfaceEntry failed: %v", ret)
			route.Metric = uint32(math.MaxUint32)
			route.MTU = uint32(math.MaxUint32)
		} else {
			route.Metric = actualInterface.Metric
			route.MTU = actualInterface.NlMtu
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
			route.Type = "global"
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
				interfaceIpAddress = string(actualAdapter.IpAddressList.IpAddress.String[:])
				route.Gateway = string(actualAdapter.GatewayList.IpAddress.String[:])
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
		}

		route.Interface = interfaceIpAddress
		route.Netmask = uint32(currRow.DestinationPrefix.PrefixLength)
		// TODO: add flags
		route.Flags = 0

		// print route
		fmt.Printf("Route: %+v\n", route)
		routes = append(routes, route)

	}
	return routes, nil
}

func HexDump(pointer unsafe.Pointer, size uint32) {
	data := unsafe.Slice((*byte)(pointer), size)
	hexDump := hex.Dump(data)
	fmt.Println(hexDump)
}
