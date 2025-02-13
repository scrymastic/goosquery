package interface_addresses

import (
	"net"
	"syscall"
	"unsafe"

	"golang.org/x/sys/windows"
)

// InterfaceAddress represents a network interface address
type InterfaceAddress struct {
	Interface    uint32 `json:"interface"`
	Address      string `json:"address"`
	Mask         string `json:"mask"`
	Broadcast    string `json:"broadcast"`
	PointToPoint string `json:"point_to_point"`
	Type         string `json:"type"`
	FriendlyName string `json:"friendly_name"`
}

// Windows-specific constants for IP address suffix origin
const (
	IpSuffixOriginOther            = 0
	IpSuffixOriginManual           = 1
	IpSuffixOriginWellKnown        = 2
	IpSuffixOriginDhcp             = 3
	IpSuffixOriginLinkLayerAddress = 4
	IpSuffixOriginRandom           = 5
)

// Helper function to convert IP mask to string representation
func ipNetMaskToString(prefixLength uint8, isIPv6 bool) string {
	// Create a mask from the prefix length, using appropriate size for IPv4/IPv6
	if isIPv6 {
		mask := net.CIDRMask(int(prefixLength), 128)
		return net.IP(mask).String()
	}
	mask := net.CIDRMask(int(prefixLength), 32)
	return net.IP(mask).String()
}

// Helper function to calculate broadcast address for IPv4
func calculateBroadcast(ip net.IP, prefixLength uint8) string {
	if ip.To4() == nil {
		return ""
	}

	mask := net.CIDRMask(int(prefixLength), 32)
	broadcast := make(net.IP, len(ip.To4()))

	for i := 0; i < len(ip.To4()); i++ {
		broadcast[i] = ip[i] | ^mask[i]
	}

	return broadcast.String()
}

// processUnicastAddress handles a single unicast address and returns an InterfaceAddress
func processUnicastAddress(addr *windows.IpAdapterAddresses, unicastAddr *windows.IpAdapterUnicastAddress) (ifaceAddr InterfaceAddress, isOk bool) {

	// Set interface index and friendly name
	ifaceAddr.Interface = addr.IfIndex
	ifaceAddr.FriendlyName = windows.UTF16PtrToString(addr.FriendlyName)

	// Get the IP address from the unicast address
	sockAddr := (*syscall.RawSockaddrAny)(unsafe.Pointer(unicastAddr.Address.Sockaddr))
	ip, isIPv6, ok := getIPFromSockAddr(sockAddr)
	if !ok {
		return ifaceAddr, false
	}

	// Set IP address and mask
	ifaceAddr.Address = ip.String()
	ifaceAddr.Mask = ipNetMaskToString(unicastAddr.OnLinkPrefixLength, isIPv6)

	// Set address type
	ifaceAddr.Type = getAddressType(addr.IfType, unicastAddr.SuffixOrigin)

	// Set broadcast for IPv4
	if ip.To4() != nil {
		ifaceAddr.Broadcast = calculateBroadcast(ip, unicastAddr.OnLinkPrefixLength)
	}

	// Set point-to-point if applicable
	if addr.IfType == windows.IF_TYPE_PPP {
		ifaceAddr.PointToPoint = ""
	}

	return ifaceAddr, true
}

// getIPFromSockAddr extracts IP address from a socket address
func getIPFromSockAddr(sockAddr *syscall.RawSockaddrAny) (ip net.IP, isIPv6 bool, isOk bool) {
	switch sockAddr.Addr.Family {
	case syscall.AF_INET:
		sa := (*syscall.RawSockaddrInet4)(unsafe.Pointer(sockAddr))
		ip = net.IP(sa.Addr[:])
		isIPv6 = false
		isOk = true
	case syscall.AF_INET6:
		sa := (*syscall.RawSockaddrInet6)(unsafe.Pointer(sockAddr))
		ip = net.IP(sa.Addr[:])
		isIPv6 = true
		isOk = true
	default:
		ip = nil
		isIPv6 = false
		isOk = false
	}
	return
}

// getAddressType determines the address type based on interface and suffix origin
func getAddressType(ifType uint32, suffixOrigin int32) string {
	if ifType == windows.IF_TYPE_SOFTWARE_LOOPBACK {
		return "other"
	}

	switch suffixOrigin {
	case IpSuffixOriginManual:
		return "manual"
	case IpSuffixOriginDhcp:
		return "dhcp"
	case IpSuffixOriginLinkLayerAddress, IpSuffixOriginRandom:
		return "auto"
	default:
		return "unknown"
	}
}

// GenInterfaceAddresses returns a list of all interface addresses on the system
// It returns a slice of InterfaceAddress and an error if the operation fails.
func GenInterfaceAddresses() ([]InterfaceAddress, error) {
	var results []InterfaceAddress

	// Get required buffer size
	var size uint32
	err := windows.GetAdaptersAddresses(
		syscall.AF_UNSPEC,
		windows.GAA_FLAG_INCLUDE_PREFIX|windows.GAA_FLAG_INCLUDE_GATEWAYS,
		0,
		nil,
		&size,
	)
	if err != windows.ERROR_BUFFER_OVERFLOW {
		return nil, err
	}

	// Allocate buffer and make the actual call
	buffer := make([]byte, size)
	addr := (*windows.IpAdapterAddresses)(unsafe.Pointer(&buffer[0]))
	err = windows.GetAdaptersAddresses(
		syscall.AF_UNSPEC,
		windows.GAA_FLAG_INCLUDE_PREFIX|windows.GAA_FLAG_INCLUDE_GATEWAYS,
		0,
		addr,
		&size,
	)
	if err != nil {
		return nil, err
	}

	// Iterate through all adapters
	for ; addr != nil; addr = addr.Next {
		for unicastAddr := addr.FirstUnicastAddress; unicastAddr != nil; unicastAddr = unicastAddr.Next {
			if ifaceAddr, ok := processUnicastAddress(addr, unicastAddr); ok {
				results = append(results, ifaceAddr)
			}
		}
	}

	return results, nil
}
