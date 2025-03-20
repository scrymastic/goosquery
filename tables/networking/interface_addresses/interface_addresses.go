package interface_addresses

import (
	"fmt"
	"net"
	"syscall"
	"unsafe"

	"github.com/scrymastic/goosquery/sql/context"
	"github.com/scrymastic/goosquery/util"
	"golang.org/x/sys/windows"
)

// Column definitions for the interface_addresses table
var columnDefs = map[string]string{
	"interface":      "string",
	"address":        "string",
	"mask":           "string",
	"type":           "string",
	"broadcast":      "string",
	"point_to_point": "string",
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

// processUnicastAddress handles a single unicast address and returns address information
func processUnicastAddress(addr *windows.IpAdapterAddresses, unicastAddr *windows.IpAdapterUnicastAddress, ctx context.Context) (map[string]interface{}, bool) {
	result := util.InitColumns(ctx, columnDefs)

	// Get the IP address from the unicast address
	sockAddr := (*syscall.RawSockaddrAny)(unsafe.Pointer(unicastAddr.Address.Sockaddr))
	ip, isIPv6, ok := getIPFromSockAddr(sockAddr)
	if !ok {
		return nil, false
	}

	// Set interface index and friendly name
	if ctx.IsColumnUsed("interface") {
		result["interface"] = fmt.Sprintf("%d", addr.IfIndex)
	}

	if ctx.IsColumnUsed("friendly_name") {
		result["friendly_name"] = windows.UTF16PtrToString(addr.FriendlyName)
	}

	// Set IP address
	if ctx.IsColumnUsed("address") {
		result["address"] = ip.String()
	}

	// Set mask
	if ctx.IsColumnUsed("mask") {
		result["mask"] = ipNetMaskToString(unicastAddr.OnLinkPrefixLength, isIPv6)
	}

	// Set address type
	if ctx.IsColumnUsed("type") {
		result["type"] = getAddressType(addr.IfType, unicastAddr.SuffixOrigin)
	}

	// Set broadcast for IPv4
	if ctx.IsColumnUsed("broadcast") && ip.To4() != nil {
		result["broadcast"] = calculateBroadcast(ip, unicastAddr.OnLinkPrefixLength)
	}

	// Set point-to-point if applicable
	if ctx.IsColumnUsed("point_to_point") {
		if addr.IfType == windows.IF_TYPE_PPP {
			result["point_to_point"] = "true"
		} else {
			result["point_to_point"] = "false"
		}
	}

	return result, true
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
		return "loopback"
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
// It returns a slice of map[string]interface{} and an error if the operation fails.
func GenInterfaceAddresses(ctx context.Context) ([]map[string]interface{}, error) {
	var results []map[string]interface{}

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
			if ifaceAddr, ok := processUnicastAddress(addr, unicastAddr, ctx); ok {
				results = append(results, ifaceAddr)
			}
		}
	}

	return results, nil
}
