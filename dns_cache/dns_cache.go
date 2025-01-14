package dns_cache

import (
	"fmt"
	"syscall"
	"unsafe"

	"golang.org/x/sys/windows"
)

// _DNS_CACHE_ENTRY represents the Windows DNS cache entry structure
// https://github.com/malcomvetter/DnsCache/blob/master/DnsCache/DnsCache.cpp
type _DNS_CACHE_ENTRY struct {
	pNext       *_DNS_CACHE_ENTRY
	pszName     *uint16
	wType       uint16
	wDataLength uint16
	dwFlags     uint32
}

type DNSCache struct {
	Name  string `json:"name"`
	Type  string `json:"type"`
	Flags uint32 `json:"flags"`
}

var (
	procDnsGetCacheDataTable uintptr
)

var (
	dnsTypeMap = map[uint16]string{
		1:     "A",
		2:     "NS",
		5:     "CNAME",
		6:     "SOA",
		12:    "PTR",
		13:    "HINFO",
		15:    "MX",
		16:    "TXT",
		17:    "RP",
		18:    "AFSDB",
		24:    "SIG",
		25:    "KEY",
		28:    "AAAA",
		29:    "LOC",
		33:    "SRV",
		35:    "NAPTR",
		36:    "KX",
		37:    "CERT",
		39:    "DNAME",
		41:    "OPT",
		42:    "APL",
		43:    "DS",
		44:    "SSHFP",
		45:    "IPSECKEY",
		46:    "RRSIG",
		47:    "NSEC",
		48:    "DNSKEY",
		49:    "DHCID",
		50:    "NSEC3",
		51:    "NSEC3PARAM",
		52:    "TLSA",
		53:    "SMIMEA",
		55:    "HIP",
		59:    "CDS",
		60:    "CDNSKEY",
		61:    "OPENPGPKEY",
		62:    "CSYNC",
		63:    "ZONEMD",
		249:   "TKEY",
		250:   "TSIG",
		251:   "IXFR",
		252:   "AXFR",
		255:   "*",
		256:   "URI",
		257:   "CAA",
		32768: "TA",
		32769: "DLV",
	}
)

func GenDNSCache() ([]DNSCache, error) {
	// Load dnsapi.dll
	modDnsapi, err := syscall.LoadLibrary("dnsapi.dll")
	if err != nil {
		return nil, fmt.Errorf("error loading dnsapi.dll: %v", err)
	}
	defer syscall.FreeLibrary(modDnsapi)

	// Get the DnsGetCacheDataTable function
	procDnsGetCacheDataTable, err = syscall.GetProcAddress(modDnsapi, "DnsGetCacheDataTable")
	if err != nil {
		return nil, fmt.Errorf("error getting DnsGetCacheDataTable function: %v", err)
	}

	// Allocate memory for the first entry
	entry := &_DNS_CACHE_ENTRY{}

	// Call DnsGetCacheDataTable using SyscallN
	if ret, _, err := syscall.SyscallN(procDnsGetCacheDataTable,
		uintptr(unsafe.Pointer(&entry)),
	); err != syscall.Errno(0) {
		return nil, fmt.Errorf("error calling DnsGetCacheDataTable: %v", ret)
	}

	var dnsCache []DNSCache

	dnsType := dnsTypeMap[entry.wType]
	if dnsType == "" {
		dnsType = fmt.Sprintf("Unknown (%d)", entry.wType)
	}

	// Iterate through the linked list of DNS_CACHE_ENTRY
	for entry != nil {
		dnsCache = append(dnsCache, DNSCache{
			Name:  windows.UTF16PtrToString(entry.pszName),
			Type:  dnsType,
			Flags: entry.dwFlags,
		})

		entry = entry.pNext
	}

	return dnsCache, nil
}
