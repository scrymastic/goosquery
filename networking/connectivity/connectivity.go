package connectivity

import (
	"fmt"
	"syscall"
	"unsafe"

	"github.com/go-ole/go-ole"
	"golang.org/x/sys/windows"
)

// GUID definitions for COM interfaces
var (
	CLSID_NetworkListManager = ole.NewGUID("{DCB00C01-570F-4A9B-8D69-199FDBA5723B}")
	IID_INetworkListManager  = ole.NewGUID("{DCB00000-570F-4A9B-8D69-199FDBA5723B}")
)

// NLM_CONNECTIVITY flags
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

// GenConnectivity generates network connectivity information
func GenConnectivity() ([]Connectivity, error) {
	// Initialize COM
	if err := windows.CoInitializeEx(0, windows.COINIT_MULTITHREADED); err != nil {
		return nil, fmt.Errorf("failed to initialize COM: %w", err)
	}
	defer windows.CoUninitialize()

	// Create instance of NetworkListManager
	var unknown *ole.IUnknown
	unknown, err := ole.CreateInstance(CLSID_NetworkListManager, IID_INetworkListManager)
	if err != nil {
		return nil, fmt.Errorf("failed to create NetworkListManager instance: %v", err)
	}
	defer unknown.Release()

	// Get the NetworkListManager interface
	nlm := (*INetworkListManager)(unsafe.Pointer(unknown))

	// Get connectivity status
	var connectivity uint32
	hr := nlm.GetConnectivity(&connectivity)
	if hr != windows.S_OK {
		return nil, fmt.Errorf("GetConnectivity failed: %v", hr)
	}

	// Create result
	result := Connectivity{
		Disconnected:     connectivity&NLM_CONNECTIVITY_DISCONNECTED != 0,
		IPv4NoTraffic:    connectivity&NLM_CONNECTIVITY_IPV4_NOTRAFFIC != 0,
		IPv6NoTraffic:    connectivity&NLM_CONNECTIVITY_IPV6_NOTRAFFIC != 0,
		IPv4Subnet:       connectivity&NLM_CONNECTIVITY_IPV4_SUBNET != 0,
		IPv4LocalNetwork: connectivity&NLM_CONNECTIVITY_IPV4_LOCALNETWORK != 0,
		IPv4Internet:     connectivity&NLM_CONNECTIVITY_IPV4_INTERNET != 0,
		IPv6Subnet:       connectivity&NLM_CONNECTIVITY_IPV6_SUBNET != 0,
		IPv6LocalNetwork: connectivity&NLM_CONNECTIVITY_IPV6_LOCALNETWORK != 0,
		IPv6Internet:     connectivity&NLM_CONNECTIVITY_IPV6_INTERNET != 0,
	}

	return []Connectivity{result}, nil
}

// INetworkListManager COM interface
type INetworkListManager struct {
	ole.IUnknown
}

// GetConnectivity gets the current connectivity state
func (nlm *INetworkListManager) GetConnectivity(connectivity *uint32) windows.Handle {
	ret, _, _ := syscall.Syscall(
		nlm.VTable().GetConnectivity,
		2,
		uintptr(unsafe.Pointer(nlm)),
		uintptr(unsafe.Pointer(connectivity)),
		0)
	return windows.Handle(ret)
}

// INetworkListManagerVtbl represents the COM vtable for INetworkListManager
type INetworkListManagerVtbl struct {
	ole.IUnknownVtbl
	GetConnectivity uintptr
}

func (i *INetworkListManager) VTable() *INetworkListManagerVtbl {
	return (*INetworkListManagerVtbl)(unsafe.Pointer(i.RawVTable))
}
