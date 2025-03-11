// Package networking provides access to various networking-related information and operations.
package networking

import (
	"github.com/scrymastic/goosquery/collector"
	"github.com/scrymastic/goosquery/networking/arp_cache"
	"github.com/scrymastic/goosquery/networking/connectivity"
	"github.com/scrymastic/goosquery/networking/curl"
	"github.com/scrymastic/goosquery/networking/etc_hosts"
	"github.com/scrymastic/goosquery/networking/etc_protocols"
	"github.com/scrymastic/goosquery/networking/etc_services"
	"github.com/scrymastic/goosquery/networking/interface_addresses"
	"github.com/scrymastic/goosquery/networking/interface_details"
	"github.com/scrymastic/goosquery/networking/listening_ports"
	"github.com/scrymastic/goosquery/networking/process_open_sockets"
	"github.com/scrymastic/goosquery/networking/routes"
	"github.com/scrymastic/goosquery/networking/windows_firewall_rules"
)

// Initialize all networking collectors
func init() {
	// Register all networking collectors
	collector.Register("arp_cache", "Retrieves the current ARP cache entries from the system", "networking", GenARPCache, false, "", "")
	collector.Register("connectivity", "Performs connectivity checks", "networking", GenConnectivity, false, "", "")
	collector.Register("curl", "Performs an HTTP request and returns details about the response", "networking", GenCurl, true, "string,string", "URL to request,User agent string")
	collector.Register("etc_hosts", "Retrieves entries from the hosts file", "networking", GenEtcHosts, false, "", "")
	collector.Register("etc_protocols", "Retrieves protocol entries", "networking", GenEtcProtocols, false, "", "")
	collector.Register("etc_services", "Retrieves service entries", "networking", GenEtcServices, false, "", "")
	collector.Register("interface_addresses", "Retrieves network interface addresses", "networking", GenInterfaceAddresses, false, "", "")
	collector.Register("interface_details", "Retrieves detailed information about network interfaces", "networking", GenInterfaceDetails, false, "", "")
	collector.Register("listening_ports", "Retrieves information about listening ports", "networking", GenListeningPorts, false, "", "")
	collector.Register("process_open_sockets", "Retrieves information about open sockets by processes", "networking", GenProcessOpenSockets, false, "", "")
	collector.Register("routes", "Retrieves network routing information", "networking", GenRoutes, false, "", "")
	collector.Register("windows_firewall_rules", "Retrieves Windows Firewall rules", "networking", GenWindowsFirewallRules, false, "", "")
}

// ARPCache represents a single entry in the ARP cache.
type ARPCache = arp_cache.ARPCache

// GenARPCache retrieves the current ARP cache entries from the system.
func GenARPCache() ([]ARPCache, error) {
	return arp_cache.GenARPCache()
}

// Connectivity represents the network connectivity state
type Connectivity = connectivity.Connectivity

// GenConnectivity performs connectivity checks
func GenConnectivity() ([]Connectivity, error) {
	return connectivity.GenConnectivity()
}

// Curl represents a curl request and response
type Curl = curl.Curl

// GenCurl performs an HTTP request and returns details about the response
func GenCurl(url string, userAgent string) (Curl, error) {
	return curl.GenCurl(url, userAgent)
}

// HostEntry represents a single hosts file entry
type HostEntry = etc_hosts.HostEntry

// GenEtcHosts retrieves entries from the hosts file
func GenEtcHosts() ([]HostEntry, error) {
	return etc_hosts.GenEtcHosts()
}

// EtcProtocol represents a protocol entry
type EtcProtocol = etc_protocols.EtcProtocol

// GenEtcProtocols retrieves protocol entries
func GenEtcProtocols() ([]EtcProtocol, error) {
	return etc_protocols.GenEtcProtocols()
}

// ServiceEntry represents a service entry
type ServiceEntry = etc_services.ServiceEntry

// GenEtcServices retrieves service entries
func GenEtcServices() ([]ServiceEntry, error) {
	return etc_services.GenEtcServices()
}

// InterfaceAddress represents a network interface address
type InterfaceAddress = interface_addresses.InterfaceAddress

// GenInterfaceAddresses retrieves network interface addresses
func GenInterfaceAddresses() ([]InterfaceAddress, error) {
	return interface_addresses.GenInterfaceAddresses()
}

// InterfaceDetail represents detailed information about a network interface
type InterfaceDetail = interface_details.InterfaceDetail

// GenInterfaceDetails retrieves detailed information about network interfaces
func GenInterfaceDetails() ([]InterfaceDetail, error) {
	return interface_details.GenInterfaceDetails()
}

// ListeningPort represents a listening port on the system
type ListeningPort = listening_ports.ListeningPort

// GenListeningPorts retrieves information about listening ports
func GenListeningPorts() ([]ListeningPort, error) {
	return listening_ports.GenListeningPorts()
}

// ProcessOpenSocket represents an open socket by a process
type ProcessOpenSocket = process_open_sockets.ProcessOpenSocket

// GenProcessOpenSockets retrieves information about open sockets by processes
func GenProcessOpenSockets() ([]ProcessOpenSocket, error) {
	return process_open_sockets.GenProcessOpenSockets()
}

// Route represents a network route
type Route = routes.Route

// GenRoutes retrieves network routing information
func GenRoutes() ([]Route, error) {
	return routes.GenRoutes()
}

// WindowsFirewallRules represents a Windows Firewall rule
type WindowsFirewallRules = windows_firewall_rules.WindowsFirewallRules

// GenWindowsFirewallRules retrieves Windows Firewall rules
func GenWindowsFirewallRules() ([]WindowsFirewallRules, error) {
	return windows_firewall_rules.GenWindowsFirewallRules()
}
