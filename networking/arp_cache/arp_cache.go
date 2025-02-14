package arp_cache

import (
	"fmt"
	"strings"

	"github.com/StackExchange/wmi"
)

// ARPCache represents a single entry in the ARP cache.
type ARPCache struct {
	Address   string `json:"address"`
	MAC       string `json:"mac"`
	Interface string `json:"interface"`
	Permanent bool   `json:"permanent"`
}

// MSFT_NetNeighbor represents the WMI MSFT_NetNeighbor class structure.
type MSFT_NetNeighbor struct {
	AddressFamily    uint16
	InterfaceAlias   string
	InterfaceIndex   uint32
	IPAddress        string
	LinkLayerAddress string
	State            uint8
	Store            uint8
}

const (
	IPv4AddressFamily uint16 = 2
	IPv6AddressFamily uint16 = 23
	PermanentState    uint8  = 6
)

// GenARPCache retrieves the current ARP cache entries from the system.
// It returns a slice of ARPEntry and an error if the operation fails.
func GenARPCache() ([]ARPCache, error) {
	var neighbors []MSFT_NetNeighbor
	query := "SELECT * FROM MSFT_NetNeighbor"
	namespace := `ROOT\StandardCimv2`
	if err := wmi.QueryNamespace(query, &neighbors, namespace); err != nil {
		return nil, fmt.Errorf("failed to query WMI for ARP entries: %w", err)
	}

	entries := make([]ARPCache, 0, len(neighbors))
	for _, n := range neighbors {
		if n.AddressFamily != IPv4AddressFamily {
			continue
		}

		if n.LinkLayerAddress == "" || n.LinkLayerAddress == "00-00-00-00-00-00" {
			continue
		}

		entries = append(entries, ARPCache{
			Address:   n.IPAddress,
			MAC:       strings.ReplaceAll(n.LinkLayerAddress, "-", ":"),
			Interface: n.InterfaceAlias,
			Permanent: n.State == PermanentState,
		})
	}

	return entries, nil
}
