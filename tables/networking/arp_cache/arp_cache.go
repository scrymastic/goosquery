package arp_cache

import (
	"fmt"
	"strings"

	"github.com/StackExchange/wmi"
	"github.com/scrymastic/goosquery/sql/context"
	"github.com/scrymastic/goosquery/util"
)

// Column definitions for the listening_ports table
var columnDefs = map[string]string{
	"address":   "string",
	"mac":       "string",
	"interface": "string",
	"permanent": "string",
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
// It returns a slice of map[string]interface{} and an error if the operation fails.
func GenARPCache(ctx context.Context) ([]map[string]interface{}, error) {
	var neighbors []MSFT_NetNeighbor
	query := "SELECT * FROM MSFT_NetNeighbor"
	namespace := `ROOT\StandardCimv2`
	if err := wmi.QueryNamespace(query, &neighbors, namespace); err != nil {
		return nil, fmt.Errorf("failed to query MSFT_NetNeighbor: %w", err)
	}

	entries := make([]map[string]interface{}, 0, len(neighbors))
	for _, n := range neighbors {
		if n.AddressFamily != IPv4AddressFamily {
			continue
		}

		if n.LinkLayerAddress == "" || n.LinkLayerAddress == "00-00-00-00-00-00" {
			continue
		}

		entry := util.InitColumns(ctx, columnDefs)

		if ctx.IsColumnUsed("address") {
			entry["address"] = n.IPAddress
		}

		if ctx.IsColumnUsed("mac") {
			entry["mac"] = strings.ReplaceAll(n.LinkLayerAddress, "-", ":")
		}

		if ctx.IsColumnUsed("interface") {
			entry["interface"] = n.InterfaceAlias
		}

		if ctx.IsColumnUsed("permanent") {
			entry["permanent"] = map[bool]string{true: "1", false: "0"}[n.State == PermanentState]
		}

		entries = append(entries, entry)
	}

	return entries, nil
}
