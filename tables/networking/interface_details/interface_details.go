package interface_details

import (
	"fmt"
	"strconv"
	"strings"
	"unsafe"

	"github.com/StackExchange/wmi"
	"golang.org/x/sys/windows"
)

type InterfaceDetail struct {
	Interface                  string `json:"interface"`
	MAC                        string `json:"mac"`
	Type                       int32  `json:"type"`
	MTU                        int32  `json:"mtu"`
	Metric                     int32  `json:"metric"`
	Flags                      int32  `json:"flags"`
	IPackets                   int64  `json:"ipackets"`
	OPackets                   int64  `json:"opackets"`
	IBytes                     int64  `json:"ibytes"`
	OBytes                     int64  `json:"obytes"`
	IErrors                    int64  `json:"ierrors"`
	OErrors                    int64  `json:"oerrors"`
	IDrops                     int64  `json:"idrops"`
	ODrops                     int64  `json:"odrops"`
	Collisions                 int64  `json:"collisions"`
	LastChange                 int64  `json:"last_change"`
	FriendlyName               string `json:"friendly_name"`
	Description                string `json:"description"`
	Manufacturer               string `json:"manufacturer"`
	ConnectionID               string `json:"connection_id"`
	ConnectionStatus           string `json:"connection_status"`
	Enabled                    int32  `json:"enabled"`
	PhysicalAdapter            int32  `json:"physical_adapter"`
	Speed                      int32  `json:"speed"`
	Service                    string `json:"service"`
	DHCPEnabled                int32  `json:"dhcp_enabled"`
	DHCPLeaseExpires           string `json:"dhcp_lease_expires"`
	DHCPLeaseObtained          string `json:"dhcp_lease_obtained"`
	DHCPServer                 string `json:"dhcp_server"`
	DNSDomain                  string `json:"dns_domain"`
	DNSDomainSuffixSearchOrder string `json:"dns_domain_suffix_search_order"`
	DNSHostName                string `json:"dns_host_name"`
	DNSServerSearchOrder       string `json:"dns_server_search_order"`
}

func boolToInt32(b bool) int32 {
	if b {
		return 1
	}
	return 0
}

func getInterfaceStats(detail *InterfaceDetail) error {
	var dst []struct {
		PacketsReceivedPerSec    string
		PacketsSentPerSec        string
		BytesReceivedPerSec      string
		BytesSentPerSec          string
		PacketsReceivedErrors    string
		PacketsOutboundErrors    string
		PacketsReceivedDiscarded string
		PacketsOutboundDiscarded string
	}

	query := fmt.Sprintf("SELECT * FROM Win32_PerfRawData_Tcpip_NetworkInterface WHERE Name = %q", detail.Description)
	err := wmi.Query(query, &dst)
	if err != nil {
		return fmt.Errorf("failed to query interface stats: %v", err)
	}

	if len(dst) > 0 {
		detail.IPackets, _ = strconv.ParseInt(dst[0].PacketsReceivedPerSec, 10, 64)
		detail.OPackets, _ = strconv.ParseInt(dst[0].PacketsSentPerSec, 10, 64)
		detail.IBytes, _ = strconv.ParseInt(dst[0].BytesReceivedPerSec, 10, 64)
		detail.OBytes, _ = strconv.ParseInt(dst[0].BytesSentPerSec, 10, 64)
		detail.IErrors, _ = strconv.ParseInt(dst[0].PacketsReceivedErrors, 10, 64)
		detail.OErrors, _ = strconv.ParseInt(dst[0].PacketsOutboundErrors, 10, 64)
		detail.IDrops, _ = strconv.ParseInt(dst[0].PacketsReceivedDiscarded, 10, 64)
		detail.ODrops, _ = strconv.ParseInt(dst[0].PacketsOutboundDiscarded, 10, 64)
	}

	return nil
}

func getAdapterDetails(detail *InterfaceDetail) error {
	var dst []struct {
		Manufacturer        string
		NetConnectionID     string
		NetConnectionStatus uint32
		NetEnabled          bool
		PhysicalAdapter     bool
		ServiceName         string
		Speed               uint64
	}

	ifIndex, _ := strconv.ParseInt(detail.Interface, 10, 64)
	query := fmt.Sprintf("SELECT * FROM Win32_NetworkAdapter WHERE InterfaceIndex = %d", ifIndex)
	err := wmi.Query(query, &dst)
	if err != nil {
		return fmt.Errorf("failed to query adapter details: %v", err)
	}

	if len(dst) > 0 {
		detail.Manufacturer = dst[0].Manufacturer
		detail.ConnectionID = dst[0].NetConnectionID
		detail.ConnectionStatus = strconv.FormatUint(uint64(dst[0].NetConnectionStatus), 10)
		detail.Enabled = boolToInt32(dst[0].NetEnabled)
		detail.PhysicalAdapter = boolToInt32(dst[0].PhysicalAdapter)
		detail.Service = dst[0].ServiceName
		detail.Speed = int32(dst[0].Speed)
	}

	return nil
}

func getDHCPAndDNSInfo(detail *InterfaceDetail) error {
	var dst []struct {
		DHCPEnabled                bool
		DHCPLeaseExpires           string
		DHCPLeaseObtained          string
		DHCPServer                 string
		DNSDomain                  string
		DNSDomainSuffixSearchOrder []string
		DNSHostName                string
		DNSServerSearchOrder       []string
	}

	ifIndex, _ := strconv.ParseInt(detail.Interface, 10, 64)
	query := fmt.Sprintf("SELECT * FROM Win32_NetworkAdapterConfiguration WHERE InterfaceIndex = %d", ifIndex)
	err := wmi.Query(query, &dst)
	if err != nil {
		return fmt.Errorf("failed to query DHCP and DNS info: %v", err)
	}

	if len(dst) > 0 {
		detail.DHCPEnabled = boolToInt32(dst[0].DHCPEnabled)
		detail.DHCPLeaseExpires = dst[0].DHCPLeaseExpires
		detail.DHCPLeaseObtained = dst[0].DHCPLeaseObtained
		detail.DHCPServer = dst[0].DHCPServer
		detail.DNSDomain = dst[0].DNSDomain
		detail.DNSDomainSuffixSearchOrder = strings.Join(dst[0].DNSDomainSuffixSearchOrder, ", ")
		detail.DNSHostName = dst[0].DNSHostName
		detail.DNSServerSearchOrder = strings.Join(dst[0].DNSServerSearchOrder, ", ")
	}

	return nil
}

func GenInterfaceDetails() ([]InterfaceDetail, error) {
	const (
		maxBufferAllocRetries = 3
		initialBufferSize     = 15000
	)

	var bufLen uint32 = initialBufferSize
	var result []byte
	var err error

	// Try to get the adapter addresses with potentially multiple attempts
	for i := 0; i < maxBufferAllocRetries; i++ {
		result = make([]byte, bufLen)
		err = windows.GetAdaptersAddresses(
			windows.AF_UNSPEC,
			windows.GAA_FLAG_INCLUDE_PREFIX|windows.GAA_FLAG_SKIP_ANYCAST|windows.GAA_FLAG_SKIP_MULTICAST,
			0,
			(*windows.IpAdapterAddresses)(unsafe.Pointer(&result[0])),
			&bufLen,
		)
		if err == nil {
			break
		}
		if err != windows.ERROR_BUFFER_OVERFLOW {
			return nil, fmt.Errorf("GetAdaptersAddresses failed: %v", err)
		}
	}
	if err != nil {
		return nil, fmt.Errorf("GetAdaptersAddresses failed after retries: %v", err)
	}

	var interfaces []InterfaceDetail
	current := (*windows.IpAdapterAddresses)(unsafe.Pointer(&result[0]))

	for current != nil {
		detail := InterfaceDetail{
			Interface:   strconv.FormatInt(int64(current.IfIndex), 10),
			MTU:         int32(current.Mtu),
			Type:        int32(current.IfType),
			Description: windows.UTF16PtrToString(current.Description),
			Flags:       int32(current.Flags),
			Metric:      int32(current.Ipv4Metric),
			LastChange:  int64(-1),
			Collisions:  int64(-1),
		}

		// Convert physical address (MAC) to string
		macBytes := make([]string, current.PhysicalAddressLength)
		for i := uint32(0); i < current.PhysicalAddressLength; i++ {
			macBytes[i] = fmt.Sprintf("%02x", current.PhysicalAddress[i])
		}
		detail.MAC = strings.Join(macBytes, ":")

		// Get network interface statistics using WMI
		_ = getInterfaceStats(&detail)

		// Get physical adapter details using WMI
		_ = getAdapterDetails(&detail)

		// Get DHCP and DNS information using WMI
		_ = getDHCPAndDNSInfo(&detail)

		interfaces = append(interfaces, detail)
		current = current.Next
	}

	return interfaces, nil
}
