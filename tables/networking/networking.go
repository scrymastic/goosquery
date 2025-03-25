package networking

import (
	"github.com/scrymastic/goosquery/sql/result"
	"github.com/scrymastic/goosquery/sql/sqlctx"
	"github.com/scrymastic/goosquery/tables/networking/arp_cache"
	"github.com/scrymastic/goosquery/tables/networking/connectivity"
	"github.com/scrymastic/goosquery/tables/networking/curl"
	"github.com/scrymastic/goosquery/tables/networking/curl_certificate"
	"github.com/scrymastic/goosquery/tables/networking/etc_hosts"
	"github.com/scrymastic/goosquery/tables/networking/etc_protocols"
	"github.com/scrymastic/goosquery/tables/networking/etc_services"
	"github.com/scrymastic/goosquery/tables/networking/interface_addresses"
	"github.com/scrymastic/goosquery/tables/networking/interface_details"
	"github.com/scrymastic/goosquery/tables/networking/listening_ports"
	"github.com/scrymastic/goosquery/tables/networking/process_open_sockets"
	"github.com/scrymastic/goosquery/tables/networking/routes"
	"github.com/scrymastic/goosquery/tables/networking/windows_firewall_rules"
)

// GenARPCache generates ARP cache entries
func GenARPCache(ctx *sqlctx.Context) (*result.Results, error) {
	return arp_cache.GenARPCache(ctx)
}

// GenConnectivity generates connectivity information
func GenConnectivity(ctx *sqlctx.Context) (*result.Results, error) {
	return connectivity.GenConnectivity(ctx)
}

// GenCurl generates results from a curl request
func GenCurl(ctx *sqlctx.Context) (*result.Results, error) {
	return curl.GenCurl(ctx)
}

// GenCurlCertificate generates information about the curl certificate
func GenCurlCertificate(ctx *sqlctx.Context) (*result.Results, error) {
	return curl_certificate.GenCurlCertificate(ctx)
}

// GenEtcHosts generates entries from the hosts file
func GenEtcHosts(ctx *sqlctx.Context) (*result.Results, error) {
	return etc_hosts.GenEtcHosts(ctx)
}

// GenEtcProtocols generates entries from the protocols file
func GenEtcProtocols(ctx *sqlctx.Context) (*result.Results, error) {
	return etc_protocols.GenEtcProtocols(ctx)
}

// GenEtcServices generates entries from the services file
func GenEtcServices(ctx *sqlctx.Context) (*result.Results, error) {
	return etc_services.GenEtcServices(ctx)
}

// GenInterfaceAddresses generates entries from the interface addresses file
func GenInterfaceAddresses(ctx *sqlctx.Context) (*result.Results, error) {
	return interface_addresses.GenInterfaceAddresses(ctx)
}

// GenInterfaceDetails generates entries from the interface details file
func GenInterfaceDetails(ctx *sqlctx.Context) (*result.Results, error) {
	return interface_details.GenInterfaceDetails(ctx)
}

// GenListeningPorts generates information about listening ports
func GenListeningPorts(ctx *sqlctx.Context) (*result.Results, error) {
	return listening_ports.GenListeningPorts(ctx)
}

// GenProcessOpenSockets generates information about process open sockets
func GenProcessOpenSockets(ctx *sqlctx.Context) (*result.Results, error) {
	return process_open_sockets.GenProcessOpenSockets(ctx)
}

// GenRoutes generates network routing information
func GenRoutes(ctx *sqlctx.Context) (*result.Results, error) {
	return routes.GenRoutes(ctx)
}

// GenWindowsFirewallRules generates Windows firewall rules
func GenWindowsFirewallRules(ctx *sqlctx.Context) (*result.Results, error) {
	return windows_firewall_rules.GenWindowsFirewallRules(ctx)
}
