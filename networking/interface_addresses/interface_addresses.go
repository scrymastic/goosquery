package interface_addresses

import (
	"net"
	"syscall"
	"unsafe"

	"golang.org/x/sys/windows"
)

// InterfaceAddress represents a network interface address
type InterfaceAddress struct {
	Interface    uint32 `json:"interface"`      // Interface name
	Address      string `json:"address"`        // Specific address for interface
	Mask         string `json:"mask"`           // Interface netmask
	Broadcast    string `json:"broadcast"`      // Broadcast address for the interface
	PointToPoint string `json:"point_to_point"` // PtP address for the interface
	Type         string `json:"type"`           // Type of address (dhcp, manual, auto, other, unknown)
	FriendlyName string `json:"friendly_name"`  // Windows only: friendly display name
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

// GenInterfaceAddresses returns a list of all interface addresses on the system
func GenInterfaceAddresses() ([]InterfaceAddress, error) {
	var results []InterfaceAddress

	// First call to get required buffer size
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
		// Get unicast addresses for this adapter
		for unicastAddr := addr.FirstUnicastAddress; unicastAddr != nil; unicastAddr = unicastAddr.Next {
			var ifaceAddr InterfaceAddress

			// Set interface index
			ifaceAddr.Interface = addr.IfIndex

			// Get the IP address from the unicast address
			sockAddr := (*syscall.RawSockaddrAny)(unsafe.Pointer(unicastAddr.Address.Sockaddr))
			var ip net.IP
			isIPv6 := false

			switch sockAddr.Addr.Family {
			case syscall.AF_INET:
				sa := (*syscall.RawSockaddrInet4)(unsafe.Pointer(sockAddr))
				ip = net.IP(sa.Addr[:])
			case syscall.AF_INET6:
				sa := (*syscall.RawSockaddrInet6)(unsafe.Pointer(sockAddr))
				ip = net.IP(sa.Addr[:])
				isIPv6 = true
			default:
				continue
			}

			// Set IP address
			ifaceAddr.Address = ip.String()

			// Set netmask from prefix length
			ifaceAddr.Mask = ipNetMaskToString(unicastAddr.OnLinkPrefixLength, isIPv6)

			// Set address type based on suffix origin
			if addr.IfType == windows.IF_TYPE_SOFTWARE_LOOPBACK {
				ifaceAddr.Type = "other"
			} else {
				switch unicastAddr.SuffixOrigin {
				case IpSuffixOriginManual:
					ifaceAddr.Type = "manual"
				case IpSuffixOriginDhcp:
					ifaceAddr.Type = "dhcp"
				case IpSuffixOriginLinkLayerAddress, IpSuffixOriginRandom:
					ifaceAddr.Type = "auto"
				default:
					ifaceAddr.Type = "unknown"
				}
			}

			// Set broadcast address for IPv4
			if ip.To4() != nil {
				broadcast := calculateBroadcast(ip, unicastAddr.OnLinkPrefixLength)
				if broadcast != "" {
					ifaceAddr.Broadcast = broadcast
				}
			}

			// Set point-to-point address if applicable
			if addr.IfType == windows.IF_TYPE_PPP {
				// For PPP interfaces, we could get the remote address from addr.FirstPrefix
				// but for now we'll leave it empty as it requires more complex processing
				ifaceAddr.PointToPoint = ""
			}

			// Set friendly name
			ifaceAddr.FriendlyName = windows.UTF16PtrToString(addr.FriendlyName)

			results = append(results, ifaceAddr)
		}
	}

	return results, nil
}
