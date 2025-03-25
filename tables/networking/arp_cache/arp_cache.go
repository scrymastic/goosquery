package arp_cache

import (
	"fmt"
	"strings"

	"github.com/StackExchange/wmi"
	"github.com/scrymastic/goosquery/sql/result"
	"github.com/scrymastic/goosquery/sql/sqlctx"
)

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
// It returns a slice of map[string]interface{} and an error if the operation fails.
func GenARPCache(ctx *sqlctx.Context) (*result.Results, error) {
	var neighbors []MSFT_NetNeighbor
	query := "SELECT * FROM MSFT_NetNeighbor"
	namespace := `ROOT\StandardCimv2`
	if err := wmi.QueryNamespace(query, &neighbors, namespace); err != nil {
		return nil, fmt.Errorf("failed to query MSFT_NetNeighbor: %w", err)
	}

	entries := result.NewQueryResult()
	for _, n := range neighbors {
		if n.AddressFamily != IPv4AddressFamily {
			continue
		}

		if n.LinkLayerAddress == "" || n.LinkLayerAddress == "00-00-00-00-00-00" {
			continue
		}

		entry := result.NewResult(ctx, Schema)

		entry.Set("address", n.IPAddress)
		entry.Set("mac", strings.ReplaceAll(n.LinkLayerAddress, "-", ":"))
		entry.Set("interface", n.InterfaceAlias)
		entry.Set("permanent", map[bool]string{true: "1", false: "0"}[n.State == PermanentState])

		entries.AppendResult(*entry)
	}

	return entries, nil
}
