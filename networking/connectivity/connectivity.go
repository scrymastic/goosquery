package connectivity

import (
	"fmt"

	"github.com/go-ole/go-ole"
)

// Connectivity flag constants (from netlistmgr.h)
const (
	NLM_CONNECTIVITY_DISCONNECTED      = 0
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

// Connectivity represents the network connectivity state
type Connectivity struct {
	Disconnected     bool `json:"disconnected"`
	IPv4NoTraffic    bool `json:"ipv4_no_traffic"`
	IPv6NoTraffic    bool `json:"ipv6_no_traffic"`
	IPv4Subnet       bool `json:"ipv4_subnet"`
	IPv4LocalNetwork bool `json:"ipv4_local_network"`
	IPv4Internet     bool `json:"ipv4_internet"`
	IPv6Subnet       bool `json:"ipv6_subnet"`
	IPv6LocalNetwork bool `json:"ipv6_local_network"`
	IPv6Internet     bool `json:"ipv6_internet"`
}

// GenConnectivity initializes COM, creates the NetworkListManager instance,
// retrieves connectivity flags, and returns the result in a slice.
func GenConnectivity() ([]Connectivity, error) {
	// Initialize COM.
	err := ole.CoInitializeEx(0, ole.COINIT_MULTITHREADED)
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

	connectivity := Connectivity{
		Disconnected:     (flag & NLM_CONNECTIVITY_DISCONNECTED) != 0,
		IPv4NoTraffic:    (flag & NLM_CONNECTIVITY_IPV4_NOTRAFFIC) != 0,
		IPv6NoTraffic:    (flag & NLM_CONNECTIVITY_IPV6_NOTRAFFIC) != 0,
		IPv4Subnet:       (flag & NLM_CONNECTIVITY_IPV4_SUBNET) != 0,
		IPv4LocalNetwork: (flag & NLM_CONNECTIVITY_IPV4_LOCALNETWORK) != 0,
		IPv4Internet:     (flag & NLM_CONNECTIVITY_IPV4_INTERNET) != 0,
		IPv6Subnet:       (flag & NLM_CONNECTIVITY_IPV6_SUBNET) != 0,
		IPv6LocalNetwork: (flag & NLM_CONNECTIVITY_IPV6_LOCALNETWORK) != 0,
		IPv6Internet:     (flag & NLM_CONNECTIVITY_IPV6_INTERNET) != 0,
	}

	return []Connectivity{connectivity}, nil
}
