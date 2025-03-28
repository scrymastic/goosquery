package connectivity

import (
	"fmt"

	"github.com/go-ole/go-ole"
	"github.com/scrymastic/goosquery/sql/result"
	"github.com/scrymastic/goosquery/sql/sqlctx"
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
func GenConnectivity(ctx *sqlctx.Context) (*result.Results, error) {
	// Initialize COM.
	ole.CoInitializeEx(0, ole.COINIT_APARTMENTTHREADED)
	// if err != nil {
	// 	return nil, fmt.Errorf("failed to initialize COM: %v", err)
	// }
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
	connectivityMap := result.NewResult(ctx, Schema)

	// Populate map based on requested columns
	connectivityMap.Set("disconnected", boolToInt32(flag&NLM_CONNECTIVITY_DISCONNECTED != 0))
	connectivityMap.Set("ipv4_no_traffic", boolToInt32(flag&NLM_CONNECTIVITY_IPV4_NOTRAFFIC != 0))
	connectivityMap.Set("ipv6_no_traffic", boolToInt32(flag&NLM_CONNECTIVITY_IPV6_NOTRAFFIC != 0))
	connectivityMap.Set("ipv6_no_traffic", boolToInt32(flag&NLM_CONNECTIVITY_IPV6_NOTRAFFIC != 0))
	connectivityMap.Set("ipv4_subnet", boolToInt32(flag&NLM_CONNECTIVITY_IPV4_SUBNET != 0))
	connectivityMap.Set("ipv4_local_network", boolToInt32(flag&NLM_CONNECTIVITY_IPV4_LOCALNETWORK != 0))
	connectivityMap.Set("ipv4_internet", boolToInt32(flag&NLM_CONNECTIVITY_IPV4_INTERNET != 0))
	connectivityMap.Set("ipv6_subnet", boolToInt32(flag&NLM_CONNECTIVITY_IPV6_SUBNET != 0))
	connectivityMap.Set("ipv6_local_network", boolToInt32(flag&NLM_CONNECTIVITY_IPV6_LOCALNETWORK != 0))
	connectivityMap.Set("ipv6_internet", boolToInt32(flag&NLM_CONNECTIVITY_IPV6_INTERNET != 0))

	entries := result.NewQueryResult()
	entries.AppendResult(*connectivityMap)
	return entries, nil
}
