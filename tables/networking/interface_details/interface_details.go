package interface_details

import (
	"fmt"
	"strconv"
	"strings"
	"unsafe"

	"github.com/StackExchange/wmi"
	"github.com/scrymastic/goosquery/sql/result"
	"github.com/scrymastic/goosquery/sql/sqlctx"
	"golang.org/x/sys/windows"
)

func boolToInt32(b bool) int32 {
	if b {
		return 1
	}
	return 0
}

func getInterfaceStats(ifDetail *result.Result, ctx *sqlctx.Context) error {
	ifDesc, ok := ifDetail.Get("description").(string)
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
		if ipackets, err := strconv.ParseInt(dst[0].PacketsReceivedPerSec, 10, 64); err == nil {
			ifDetail.Set("ipackets", ipackets)
		}
		if opackets, err := strconv.ParseInt(dst[0].PacketsSentPerSec, 10, 64); err == nil {
			ifDetail.Set("opackets", opackets)
		}
		if ibytes, err := strconv.ParseInt(dst[0].BytesReceivedPerSec, 10, 64); err == nil {
			ifDetail.Set("ibytes", ibytes)
		}
		if obytes, err := strconv.ParseInt(dst[0].BytesSentPerSec, 10, 64); err == nil {
			ifDetail.Set("obytes", obytes)
		}
		if ierrors, err := strconv.ParseInt(dst[0].PacketsReceivedErrors, 10, 64); err == nil {
			ifDetail.Set("ierrors", ierrors)
		}
		if oerrors, err := strconv.ParseInt(dst[0].PacketsOutboundErrors, 10, 64); err == nil {
			ifDetail.Set("oerrors", oerrors)
		}
		if idrops, err := strconv.ParseInt(dst[0].PacketsReceivedDiscarded, 10, 64); err == nil {
			ifDetail.Set("idrops", idrops)
		}
		if odrops, err := strconv.ParseInt(dst[0].PacketsOutboundDiscarded, 10, 64); err == nil {
			ifDetail.Set("odrops", odrops)
		}
	}

	return nil
}

func getAdapterDetails(ifDetail *result.Result, ctx *sqlctx.Context) error {
	if !ctx.IsAnyOfColumnsUsed([]string{
		"manufacturer", "connection_id", "connection_status", "enabled",
		"physical_adapter", "service", "speed"}) {
		return nil
	}

	ifIndex, err := strconv.ParseInt(ifDetail.Get("interface").(string), 10, 64)
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
		ifDetail.Set("manufacturer", dst[0].Manufacturer)
		ifDetail.Set("connection_id", dst[0].NetConnectionID)
		ifDetail.Set("connection_status", strconv.FormatUint(uint64(dst[0].NetConnectionStatus), 10))
		ifDetail.Set("enabled", boolToInt32(dst[0].NetEnabled))
		ifDetail.Set("physical_adapter", boolToInt32(dst[0].PhysicalAdapter))
		ifDetail.Set("service", dst[0].ServiceName)
		ifDetail.Set("speed", int32(dst[0].Speed))
	}

	return nil
}

func getDHCPAndDNSInfo(ifDetail *result.Result, ctx *sqlctx.Context) error {
	if !ctx.IsAnyOfColumnsUsed([]string{
		"dhcp_enabled", "dhcp_lease_expires", "dhcp_lease_obtained", "dhcp_server",
		"dns_domain", "dns_domain_suffix_search_order", "dns_host_name", "dns_server_search_order"}) {
		return nil
	}

	ifIndex, err := strconv.ParseInt(ifDetail.Get("interface").(string), 10, 64)
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
		ifDetail.Set("dhcp_enabled", boolToInt32(dst[0].DHCPEnabled))
		ifDetail.Set("dhcp_lease_expires", dst[0].DHCPLeaseExpires)
		ifDetail.Set("dhcp_lease_obtained", dst[0].DHCPLeaseObtained)
		ifDetail.Set("dhcp_server", dst[0].DHCPServer)
		ifDetail.Set("dns_domain", dst[0].DNSDomain)
		ifDetail.Set("dns_domain_suffix_search_order", strings.Join(dst[0].DNSDomainSuffixSearchOrder, ", "))
		ifDetail.Set("dns_host_name", dst[0].DNSHostName)
		ifDetail.Set("dns_server_search_order", strings.Join(dst[0].DNSServerSearchOrder, ", "))
	}

	return nil
}

func GenInterfaceDetails(ctx *sqlctx.Context) (*result.Results, error) {
	const (
		maxBufferAllocRetries = 3
		initialBufferSize     = 15000
	)

	var bufLen uint32 = initialBufferSize
	var buff []byte
	var err error

	// Try to get the adapter addresses with potentially multiple attempts
	for i := 0; i < maxBufferAllocRetries; i++ {
		buff = make([]byte, bufLen)
		err = windows.GetAdaptersAddresses(
			windows.AF_UNSPEC,
			windows.GAA_FLAG_INCLUDE_PREFIX|windows.GAA_FLAG_SKIP_ANYCAST|windows.GAA_FLAG_SKIP_MULTICAST,
			0,
			(*windows.IpAdapterAddresses)(unsafe.Pointer(&buff[0])),
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

	interfaces := result.NewQueryResult()
	current := (*windows.IpAdapterAddresses)(unsafe.Pointer(&buff[0]))

	for current != nil {
		ifDetail := result.NewResult(ctx, Schema)

		// Basic interface details
		ifDetail.Set("interface", strconv.FormatInt(int64(current.IfIndex), 10))
		ifDetail.Set("mtu", int32(current.Mtu))
		ifDetail.Set("type", int32(current.IfType))
		ifDetail.Set("description", windows.UTF16PtrToString(current.Description))
		ifDetail.Set("flags", int32(current.Flags))
		ifDetail.Set("metric", int32(current.Ipv4Metric))
		ifDetail.Set("last_change", int64(-1))
		ifDetail.Set("collisions", int64(-1))

		// Convert physical address (MAC) to string
		macBytes := make([]string, current.PhysicalAddressLength)
		for i := uint32(0); i < current.PhysicalAddressLength; i++ {
			macBytes[i] = fmt.Sprintf("%02x", current.PhysicalAddress[i])
		}
		ifDetail.Set("mac", strings.Join(macBytes, ":"))

		// Only get additional details if we have interface ID and there are columns that need them
		if ifDetail.Get("interface") != nil {
			// Get network interface statistics using WMI if needed
			_ = getInterfaceStats(ifDetail, ctx)

			// Get physical adapter details using WMI if needed
			_ = getAdapterDetails(ifDetail, ctx)

			// Get DHCP and DNS information using WMI if needed
			_ = getDHCPAndDNSInfo(ifDetail, ctx)
		}

		interfaces.AppendResult(*ifDetail)
		current = current.Next
	}

	return interfaces, nil
}
