package arp_cache

import (
	"fmt"
	"strings"

	"github.com/StackExchange/wmi"
)

// msftNetNeighbor represents the WMI MSFT_NetNeighbor class structure.
type msft_NetNeighbor struct {
	AddressFamily    uint16
	InterfaceAlias   string
	InterfaceIndex   uint32
	IPAddress        string
	LinkLayerAddress string
	State            uint8
	Store            uint8
}

// ARPCache represents a single entry in the ARP cache.
type ARPCache struct {
	Address   string `json:"address"`
	MAC       string `json:"mac"`
	Interface string `json:"interface"`
	Permanent bool   `json:"permanent"`
}

const (
	ipv4AddressFamily uint16 = 2
	ipv6AddressFamily uint16 = 23
	permanentState    uint8  = 6
)

// var (
// 	// addressFamilyMap maps address family numbers to their string representations.
// 	addressFamilyMap = map[uint16]string{
// 		ipv4AddressFamily: "IPv4",
// 		ipv6AddressFamily: "IPv6",
// 	}

// 	// stateMap maps neighbor state numbers to their string representations.
// 	stateMap = map[uint8]string{
// 		0: "Unreachable",
// 		1: "Incomplete",
// 		2: "Probe",
// 		3: "Delay",
// 		4: "Stale",
// 		5: "Reachable",
// 		6: "Permanent",
// 		7: "TBD",
// 	}
// )

// GenARPCache retrieves the current ARP cache entries from the system.
// It returns a slice of ARPEntry and an error if the operation fails.
func GenARPCache() ([]ARPCache, error) {
	var neighbors []msft_NetNeighbor
	query := "SELECT * FROM MSFT_NetNeighbor"
	namespace := `ROOT\StandardCimv2`
	if err := wmi.QueryNamespace(query, &neighbors, namespace); err != nil {
		return nil, fmt.Errorf("failed to query WMI for ARP entries: %w", err)
	}

	if len(neighbors) == 0 {
		return nil, fmt.Errorf("no ARP cache entries found")
	}

	entries := make([]ARPCache, 0, len(neighbors))
	for _, n := range neighbors {
		if n.AddressFamily != ipv4AddressFamily {
			continue
		}

		if n.LinkLayerAddress == "" || n.LinkLayerAddress == "00-00-00-00-00-00" {
			continue
		}

		entries = append(entries, ARPCache{
			Address:   n.IPAddress,
			MAC:       strings.ReplaceAll(n.LinkLayerAddress, "-", ":"),
			Interface: n.InterfaceAlias,
			Permanent: n.State == permanentState,
		})
	}

	return entries, nil
}
