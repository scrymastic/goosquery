package interface_details

import (
	"fmt"
	"strconv"
	"strings"
	"unsafe"

	"github.com/StackExchange/wmi"
	"github.com/scrymastic/goosquery/sql/context"
	"github.com/scrymastic/goosquery/util"
	"golang.org/x/sys/windows"
)

// Column definitions for the interface_details table
var columnDefs = map[string]string{
	"interface":                      "string",
	"mac":                            "string",
	"type":                           "int32",
	"mtu":                            "int32",
	"metric":                         "int32",
	"flags":                          "int32",
	"ipackets":                       "int64",
	"opackets":                       "int64",
	"ibytes":                         "int64",
	"obytes":                         "int64",
	"ierrors":                        "int64",
	"oerrors":                        "int64",
	"idrops":                         "int64",
	"odrops":                         "int64",
	"collisions":                     "int64",
	"last_change":                    "int64",
	"friendly_name":                  "string",
	"description":                    "string",
	"manufacturer":                   "string",
	"connection_id":                  "string",
	"connection_status":              "string",
	"enabled":                        "int32",
	"physical_adapter":               "int32",
	"service":                        "string",
	"speed":                          "int32",
	"dhcp_enabled":                   "int32",
	"dhcp_lease_expires":             "string",
	"dhcp_lease_obtained":            "string",
	"dhcp_server":                    "string",
	"dns_domain":                     "string",
	"dns_domain_suffix_search_order": "string",
	"dns_host_name":                  "string",
	"dns_server_search_order":        "string",
}

func boolToInt32(b bool) int32 {
	if b {
		return 1
	}
	return 0
}

func getInterfaceStats(ifDetail map[string]interface{}, ctx context.Context) error {
	ifDesc, ok := ifDetail["description"].(string)
	if !ok || !ctx.IsAnyOfColumnsUsed([]string{"ipackets", "opackets", "ibytes", "obytes", "ierrors", "oerrors", "idrops", "odrops"}) {
		return nil
	}

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

	query := fmt.Sprintf("SELECT * FROM Win32_PerfRawData_Tcpip_NetworkInterface WHERE Name = %q", ifDesc)
	err := wmi.Query(query, &dst)
	if err != nil {
		return fmt.Errorf("failed to query interface stats: %v", err)
	}

	if len(dst) > 0 {
		if ctx.IsColumnUsed("ipackets") {
			ifDetail["ipackets"], _ = strconv.ParseInt(dst[0].PacketsReceivedPerSec, 10, 64)
		}
		if ctx.IsColumnUsed("opackets") {
			ifDetail["opackets"], _ = strconv.ParseInt(dst[0].PacketsSentPerSec, 10, 64)
		}
		if ctx.IsColumnUsed("ibytes") {
			ifDetail["ibytes"], _ = strconv.ParseInt(dst[0].BytesReceivedPerSec, 10, 64)
		}
		if ctx.IsColumnUsed("obytes") {
			ifDetail["obytes"], _ = strconv.ParseInt(dst[0].BytesSentPerSec, 10, 64)
		}
		if ctx.IsColumnUsed("ierrors") {
			ifDetail["ierrors"], _ = strconv.ParseInt(dst[0].PacketsReceivedErrors, 10, 64)
		}
		if ctx.IsColumnUsed("oerrors") {
			ifDetail["oerrors"], _ = strconv.ParseInt(dst[0].PacketsOutboundErrors, 10, 64)
		}
		if ctx.IsColumnUsed("idrops") {
			ifDetail["idrops"], _ = strconv.ParseInt(dst[0].PacketsReceivedDiscarded, 10, 64)
		}
		if ctx.IsColumnUsed("odrops") {
			ifDetail["odrops"], _ = strconv.ParseInt(dst[0].PacketsOutboundDiscarded, 10, 64)
		}
	}

	return nil
}

func getAdapterDetails(ifDetail map[string]interface{}, ctx context.Context) error {
	if !ctx.IsAnyOfColumnsUsed([]string{
		"manufacturer", "connection_id", "connection_status", "enabled",
		"physical_adapter", "service", "speed"}) {
		return nil
	}

	ifIndex, err := strconv.ParseInt(ifDetail["interface"].(string), 10, 64)
	if err != nil {
		return fmt.Errorf("failed to parse interface index: %v", err)
	}

	var dst []struct {
		Manufacturer        string
		NetConnectionID     string
		NetConnectionStatus uint32
		NetEnabled          bool
		PhysicalAdapter     bool
		ServiceName         string
		Speed               uint64
	}

	query := fmt.Sprintf("SELECT * FROM Win32_NetworkAdapter WHERE InterfaceIndex = %d", ifIndex)
	err = wmi.Query(query, &dst)
	if err != nil {
		return fmt.Errorf("failed to query adapter details: %v", err)
	}

	if len(dst) > 0 {
		if ctx.IsColumnUsed("manufacturer") {
			ifDetail["manufacturer"] = dst[0].Manufacturer
		}
		if ctx.IsColumnUsed("connection_id") {
			ifDetail["connection_id"] = dst[0].NetConnectionID
		}
		if ctx.IsColumnUsed("connection_status") {
			ifDetail["connection_status"] = strconv.FormatUint(uint64(dst[0].NetConnectionStatus), 10)
		}
		if ctx.IsColumnUsed("enabled") {
			ifDetail["enabled"] = boolToInt32(dst[0].NetEnabled)
		}
		if ctx.IsColumnUsed("physical_adapter") {
			ifDetail["physical_adapter"] = boolToInt32(dst[0].PhysicalAdapter)
		}
		if ctx.IsColumnUsed("service") {
			ifDetail["service"] = dst[0].ServiceName
		}
		if ctx.IsColumnUsed("speed") {
			ifDetail["speed"] = int32(dst[0].Speed)
		}
	}

	return nil
}

func getDHCPAndDNSInfo(ifDetail map[string]interface{}, ctx context.Context) error {
	if !ctx.IsAnyOfColumnsUsed([]string{
		"dhcp_enabled", "dhcp_lease_expires", "dhcp_lease_obtained", "dhcp_server",
		"dns_domain", "dns_domain_suffix_search_order", "dns_host_name", "dns_server_search_order"}) {
		return nil
	}

	ifIndex, err := strconv.ParseInt(ifDetail["interface"].(string), 10, 64)
	if err != nil {
		return fmt.Errorf("failed to parse interface index: %v", err)
	}

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

	query := fmt.Sprintf("SELECT * FROM Win32_NetworkAdapterConfiguration WHERE InterfaceIndex = %d", ifIndex)
	err = wmi.Query(query, &dst)
	if err != nil {
		return fmt.Errorf("failed to query DHCP and DNS info: %v", err)
	}

	if len(dst) > 0 {
		if ctx.IsColumnUsed("dhcp_enabled") {
			ifDetail["dhcp_enabled"] = boolToInt32(dst[0].DHCPEnabled)
		}
		if ctx.IsColumnUsed("dhcp_lease_expires") {
			ifDetail["dhcp_lease_expires"] = dst[0].DHCPLeaseExpires
		}
		if ctx.IsColumnUsed("dhcp_lease_obtained") {
			ifDetail["dhcp_lease_obtained"] = dst[0].DHCPLeaseObtained
		}
		if ctx.IsColumnUsed("dhcp_server") {
			ifDetail["dhcp_server"] = dst[0].DHCPServer
		}
		if ctx.IsColumnUsed("dns_domain") {
			ifDetail["dns_domain"] = dst[0].DNSDomain
		}
		if ctx.IsColumnUsed("dns_domain_suffix_search_order") {
			ifDetail["dns_domain_suffix_search_order"] = strings.Join(dst[0].DNSDomainSuffixSearchOrder, ", ")
		}
		if ctx.IsColumnUsed("dns_host_name") {
			ifDetail["dns_host_name"] = dst[0].DNSHostName
		}
		if ctx.IsColumnUsed("dns_server_search_order") {
			ifDetail["dns_server_search_order"] = strings.Join(dst[0].DNSServerSearchOrder, ", ")
		}
	}

	return nil
}

func GenInterfaceDetails(ctx context.Context) ([]map[string]interface{}, error) {
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

	var interfaces []map[string]interface{}
	current := (*windows.IpAdapterAddresses)(unsafe.Pointer(&result[0]))

	for current != nil {
		ifDetail := util.InitColumns(ctx, columnDefs)

		// Basic interface details
		if ctx.IsColumnUsed("interface") {
			ifDetail["interface"] = strconv.FormatInt(int64(current.IfIndex), 10)
		}
		if ctx.IsColumnUsed("mtu") {
			ifDetail["mtu"] = int32(current.Mtu)
		}
		if ctx.IsColumnUsed("type") {
			ifDetail["type"] = int32(current.IfType)
		}
		if ctx.IsColumnUsed("description") {
			ifDetail["description"] = windows.UTF16PtrToString(current.Description)
		}
		if ctx.IsColumnUsed("flags") {
			ifDetail["flags"] = int32(current.Flags)
		}
		if ctx.IsColumnUsed("metric") {
			ifDetail["metric"] = int32(current.Ipv4Metric)
		}
		if ctx.IsColumnUsed("last_change") {
			ifDetail["last_change"] = int64(-1)
		}
		if ctx.IsColumnUsed("collisions") {
			ifDetail["collisions"] = int64(-1)
		}

		// Convert physical address (MAC) to string
		if ctx.IsColumnUsed("mac") {
			macBytes := make([]string, current.PhysicalAddressLength)
			for i := uint32(0); i < current.PhysicalAddressLength; i++ {
				macBytes[i] = fmt.Sprintf("%02x", current.PhysicalAddress[i])
			}
			ifDetail["mac"] = strings.Join(macBytes, ":")
		}

		// Only get additional details if we have interface ID and there are columns that need them
		if ifDetail["interface"] != nil {
			// Get network interface statistics using WMI if needed
			_ = getInterfaceStats(ifDetail, ctx)

			// Get physical adapter details using WMI if needed
			_ = getAdapterDetails(ifDetail, ctx)

			// Get DHCP and DNS information using WMI if needed
			_ = getDHCPAndDNSInfo(ifDetail, ctx)
		}

		interfaces = append(interfaces, ifDetail)
		current = current.Next
	}

	return interfaces, nil
}
