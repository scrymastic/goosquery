package dns_cache

import (
	"fmt"
	"syscall"
	"unsafe"

	"golang.org/x/sys/windows"
)

type DNSCache struct {
	Name  string `json:"name"`
	Type  string `json:"type"`
	Flags int32  `json:"flags"`
}

// DNS_CACHE_ENTRY represents the Windows DNS cache entry structure
// https://github.com/malcomvetter/DnsCache/blob/master/DnsCache/DnsCache.cpp
type DNS_CACHE_ENTRY struct {
	pNext       *DNS_CACHE_ENTRY
	pszName     *uint16
	wType       uint16
	wDataLength uint16
	dwFlags     uint32
}

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
	modDnsapi := windows.NewLazySystemDLL("dnsapi.dll")
	if modDnsapi.Load() != nil {
		return nil, fmt.Errorf("failed to load dnsapi.dll")
	}
	procDnsGetCacheDataTable := modDnsapi.NewProc("DnsGetCacheDataTable")

	// Allocate memory for the first entry
	entry := &DNS_CACHE_ENTRY{}

	// Call DnsGetCacheDataTable using SyscallN
	if _, _, err := procDnsGetCacheDataTable.Call(
		uintptr(unsafe.Pointer(&entry)),
	); err != syscall.Errno(0) {
		return nil, fmt.Errorf("error calling DnsGetCacheDataTable: %v", err)
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
			Flags: int32(entry.dwFlags),
		})

		entry = entry.pNext
	}

	return dnsCache, nil
}
