package networking

import (
	ctxPkg "github.com/scrymastic/goosquery/sql/context"
	"github.com/scrymastic/goosquery/tables/networking/arp_cache"
	"github.com/scrymastic/goosquery/tables/networking/connectivity"
	"github.com/scrymastic/goosquery/tables/networking/curl"
	"github.com/scrymastic/goosquery/tables/networking/etc_hosts"
	"github.com/scrymastic/goosquery/tables/networking/listening_ports"
	"github.com/scrymastic/goosquery/tables/networking/process_open_sockets"
	"github.com/scrymastic/goosquery/tables/networking/routes"
	"github.com/scrymastic/goosquery/tables/networking/windows_firewall_rules"
)

// GenARPCache generates ARP cache entries
func GenARPCache(ctx ctxPkg.Context) ([]map[string]interface{}, error) {
	return arp_cache.GenARPCache(ctx)
}

// GenConnectivity generates connectivity information
func GenConnectivity(ctx ctxPkg.Context) ([]map[string]interface{}, error) {
	return connectivity.GenConnectivity(ctx)
}

// GenCurl generates results from a curl request
func GenCurl(ctx ctxPkg.Context) ([]map[string]interface{}, error) {
	return curl.GenCurl(ctx)
}

// GenEtcHosts generates entries from the hosts file
func GenEtcHosts(ctx ctxPkg.Context) ([]map[string]interface{}, error) {
	return etc_hosts.GenEtcHosts(ctx)
}

// GenListeningPorts generates information about listening ports
func GenListeningPorts(ctx ctxPkg.Context) ([]map[string]interface{}, error) {
	return listening_ports.GenListeningPorts(ctx)
}

// GenProcessOpenSockets generates information about process open sockets
func GenProcessOpenSockets(ctx ctxPkg.Context) ([]map[string]interface{}, error) {
	return process_open_sockets.GenProcessOpenSockets(ctx)
}

// GenRoutes generates network routing information
func GenRoutes(ctx ctxPkg.Context) ([]map[string]interface{}, error) {
	return routes.GenRoutes(ctx)
}

// GenWindowsFirewallRules generates Windows firewall rules
func GenWindowsFirewallRules(ctx ctxPkg.Context) ([]map[string]interface{}, error) {
	return windows_firewall_rules.GenWindowsFirewallRules(ctx)
}
