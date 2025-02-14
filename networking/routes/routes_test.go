package routes

import (
	"encoding/json"
	"fmt"
	"testing"
	"unsafe"
)

func TestGetRoutes(t *testing.T) {
	routes, err := GenRoutes()
	if err != nil {
		t.Fatalf("Failed to get routes: %v", err)
	}

	jsonData, err := json.MarshalIndent(routes, "", "  ")
	if err != nil {
		t.Fatalf("Failed to marshal routes to JSON: %v", err)
	}
	fmt.Printf("Routes Results:\n%s\n", string(jsonData))
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
