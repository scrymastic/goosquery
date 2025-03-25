package routes

import (
	"encoding/json"
	"fmt"
	"testing"
	"unsafe"

	"github.com/scrymastic/goosquery/sql/sqlctx"
)

func TestGetRoutes(t *testing.T) {
	// Create context with all columns used
	ctx := sqlctx.NewContext()
	// Add all possible columns to ensure they're all included in test
	ctx.Columns = []string{
		"destination", "netmask", "gateway", "source",
		"flags", "interface", "mtu", "metric", "type",
	}

	routes, err := GenRoutes(ctx)
	if err != nil {
		t.Fatalf("Failed to get routes: %v", err)
	}

	// Print results as JSON
	jsonData, err := json.MarshalIndent(routes, "", "  ")
	if err != nil {
		t.Fatalf("Failed to marshal routes to JSON: %v", err)
	}
	fmt.Printf("Routes Results:\n%s\n", string(jsonData))
	fmt.Printf("Total routes: %d\n", routes.Size())
}

// func TestMIB_IPINTERFACE_ROW(t *testing.T) {
// 	fmt.Printf("Size of MIB_IPINTERFACE_ROW: %d bytes\n", unsafe.Sizeof(MIB_IPINTERFACE_ROW{}))
// }

func TestSOCKADDR_INET(t *testing.T) {
	fmt.Printf("Size of SOCKADDR_INET: %d bytes\n", unsafe.Sizeof(SOCKADDR_INET{}))
	// SOCKADDR_IN
	fmt.Printf("Size of SOCKADDR_IN: %d bytes\n", unsafe.Sizeof(SOCKADDR_IN{}))
	// SOCKADDR_IN6
	fmt.Printf("Size of SOCKADDR_IN6: %d bytes\n", unsafe.Sizeof(SOCKADDR_IN6{}))
}

func TestIP_ADDRESS_PREFIX(t *testing.T) {
	fmt.Printf("Size of IP_ADDRESS_PREFIX: %d bytes\n", unsafe.Sizeof(IP_ADDRESS_PREFIX{}))
}
