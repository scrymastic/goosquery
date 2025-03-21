package connectivity

import (
	"fmt"

	"github.com/go-ole/go-ole"
	"github.com/scrymastic/goosquery/sql/context"
	"github.com/scrymastic/goosquery/tables/specs"
)

// Connectivity flag constants (from netlistmgr.h)
const (
	NLM_CONNECTIVITY_DISCONNECTED      = 0x0
	NLM_CONNECTIVITY_IPV4_NOTRAFFIC    = 0x1
	NLM_CONNECTIVITY_IPV6_NOTRAFFIC    = 0x2
	NLM_CONNECTIVITY_IPV4_SUBNET       = 0x10
	NLM_CONNECTIVITY_IPV4_LOCALNETWORK = 0x20
	NLM_CONNECTIVITY_IPV4_INTERNET     = 0x40
	NLM_CONNECTIVITY_IPV6_SUBNET       = 0x100
	NLM_CONNECTIVITY_IPV6_LOCALNETWORK = 0x200
	NLM_CONNECTIVITY_IPV6_INTERNET     = 0x400
)

// CLSID and IID for the NetworkListManager COM interface.
var (
	CLSID_NetworkListManager = ole.NewGUID("{DCB00C01-570F-4A9B-8D69-199FDBA5723B}")
	IID_INetworkListManager  = ole.NewGUID("{DCB00001-570F-4A9B-8D69-199FDBA5723B}")
)

func boolToInt32(b bool) int32 {
	if b {
		return 1
	}
	return 0
}

// GenConnectivity initializes COM, creates the NetworkListManager instance,
// retrieves connectivity flags, and returns the result in a slice of maps.
func GenConnectivity(ctx context.Context) ([]map[string]interface{}, error) {
	// Initialize COM.
	err := ole.CoInitializeEx(0, ole.COINIT_APARTMENTTHREADED)
	if err != nil {
		return nil, fmt.Errorf("failed to initialize COM: %v", err)
	}
	defer ole.CoUninitialize()

	// Create an instance of INetworkListManager.
	unknown, err := ole.CreateInstance(CLSID_NetworkListManager, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create INetworkListManager instance: %v", err)
	}
	defer unknown.Release()

	mgr, err := unknown.QueryInterface(ole.IID_IDispatch)
	if err != nil {
		return nil, fmt.Errorf("failed to query interface: %v", err)
	}
	defer mgr.Release()

	connectivityVariant, err := mgr.CallMethod("GetConnectivity")
	if err != nil {
		return nil, fmt.Errorf("failed to call GetConnectivity: %v", err)
	}

	// Get the connectivity value directly as an integer
	flag := uint32(connectivityVariant.Val)

	// Create a map instead of a struct
	connectivityMap := specs.Init(ctx, Schema)

	// Populate map based on requested columns
	if ctx.IsColumnUsed("disconnected") {
		connectivityMap["disconnected"] = boolToInt32(flag&NLM_CONNECTIVITY_DISCONNECTED != 0)
	}
	if ctx.IsColumnUsed("ipv4_no_traffic") {
		connectivityMap["ipv4_no_traffic"] = boolToInt32(flag&NLM_CONNECTIVITY_IPV4_NOTRAFFIC != 0)
	}
	if ctx.IsColumnUsed("ipv6_no_traffic") {
		connectivityMap["ipv6_no_traffic"] = boolToInt32(flag&NLM_CONNECTIVITY_IPV6_NOTRAFFIC != 0)
	}
	if ctx.IsColumnUsed("ipv4_subnet") {
		connectivityMap["ipv4_subnet"] = boolToInt32(flag&NLM_CONNECTIVITY_IPV4_SUBNET != 0)
	}
	if ctx.IsColumnUsed("ipv4_local_network") {
		connectivityMap["ipv4_local_network"] = boolToInt32(flag&NLM_CONNECTIVITY_IPV4_LOCALNETWORK != 0)
	}
	if ctx.IsColumnUsed("ipv4_internet") {
		connectivityMap["ipv4_internet"] = boolToInt32(flag&NLM_CONNECTIVITY_IPV4_INTERNET != 0)
	}
	if ctx.IsColumnUsed("ipv6_subnet") {
		connectivityMap["ipv6_subnet"] = boolToInt32(flag&NLM_CONNECTIVITY_IPV6_SUBNET != 0)
	}
	if ctx.IsColumnUsed("ipv6_local_network") {
		connectivityMap["ipv6_local_network"] = boolToInt32(flag&NLM_CONNECTIVITY_IPV6_LOCALNETWORK != 0)
	}
	if ctx.IsColumnUsed("ipv6_internet") {
		connectivityMap["ipv6_internet"] = boolToInt32(flag&NLM_CONNECTIVITY_IPV6_INTERNET != 0)
	}

	return []map[string]interface{}{connectivityMap}, nil
}
