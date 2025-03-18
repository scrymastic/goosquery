package routes

import "golang.org/x/sys/windows"

// import "golang.org/x/sys/windows"

// type SOCKADDR_IN struct {
// 	SinFamily uint16
// 	SinPort   uint16
// 	SinAddr   [4]byte
// 	SinZero   [8]byte
// }

// type SOCKADDR_IN6 struct {
// 	Sin6Family uint16
// 	Sin6Port   uint16
// 	Sin6Addr   [16]byte
// 	Sin6Zero   [8]byte
// }

// type SOCKADDR_INET struct {
// 	Ipv4     SOCKADDR_IN
// 	Ipv6     SOCKADDR_IN6
// 	SiFamily uint16
// }

// type IP_ADDRESS_PREFIX struct {
// 	Prefix       SOCKADDR_INET
// 	PrefixLength uint8
// }

// type MIB_IPFORWARD_ROW2 struct {
// 	InterfaceLuid        windows.LUID
// 	InterfaceIndex       uint32
// 	DestinationPrefix    IP_ADDRESS_PREFIX
// 	NextHop              [16]byte
// 	SitePrefixLength     uint8
// 	_                    [3]byte
// 	ValidLifetime        uint32
// 	PreferredLifetime    uint32
// 	Metric               uint32
// 	Protocol             [4]byte
// 	Loopback             uint8
// 	_                    [3]byte
// 	AutoconfigureAddress uint8
// 	_                    [3]byte
// 	Publish              uint8
// 	_                    [3]byte
// 	Immortal             uint8
// 	_                    [3]byte
// 	Age                  uint32
// 	Origin               [4]byte
// }

// type MIB_IPFORWARD_TABLE2 struct {
// 	NumEntries uint32
// 	_          [4]byte // padding
// 	Table      [1]MIB_IPFORWARD_ROW2
// }

const (
	ScopeLevelInterface    uint32 = 1
	ScopeLevelLink         uint32 = 2
	ScopeLevelSubnet       uint32 = 3
	ScopeLevelAdmin        uint32 = 4
	ScopeLevelSite         uint32 = 5
	ScopeLevelOrganization uint32 = 8
	ScopeLevelGlobal       uint32 = 14
	ScopeLevelCount        uint32 = 16
)

// MIB_IPINTERFACE_ROW equivalent in Go
type MIB_IPINTERFACE_ROW struct {
	Family                               uint16
	InterfaceLuid                        windows.LUID
	InterfaceIndex                       uint32
	MaxReassemblySize                    uint32
	InterfaceIdentifier                  uint64
	MinRouterAdvertisementInterval       uint32
	MaxRouterAdvertisementInterval       uint32
	AdvertisingEnabled                   uint8
	ForwardingEnabled                    uint8
	WeakHostSend                         uint8
	WeakHostReceive                      uint8
	UseAutomaticMetric                   uint8
	UseNeighborUnreachabilityDetection   uint8
	ManagedAddressConfigurationSupported uint8
	OtherStatefulConfigurationSupported  uint8
	AdvertiseDefaultRoute                uint8
	RouterDiscoveryBehavior              int32
	DadTransmits                         uint32
	BaseReachableTime                    uint32
	RetransmitTime                       uint32
	PathMtuDiscoveryTimeout              uint32
	LinkLocalAddressBehavior             int32
	LinkLocalAddressTimeout              uint32
	ZoneIndices                          [ScopeLevelCount]uint32
	SitePrefixLength                     uint32
	Metric                               uint32
	NlMtu                                uint32
	Connected                            uint8
	SupportsWakeUpPatterns               uint8
	SupportsNeighborDiscovery            uint8
	SupportsRouterDiscovery              uint8
	ReachableTime                        uint32
	TransmitOffload                      [8]uint8
	ReceiveOffload                       [8]uint8
	DisableDefaultRoutes                 uint8
}
