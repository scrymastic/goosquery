package dns_cache

import (
	"github.com/scrymastic/goosquery/tables/specs"
)

var TableName = "dns_cache"
var Description = "Enumerate the DNS cache using the undocumented DnsGetCacheDataTable function in dnsapi.dll."
var Schema = specs.Schema{
	specs.Column{Name: "name", Type: "TEXT", Description: "DNS record name"},
	specs.Column{Name: "type", Type: "TEXT", Description: "DNS record type"},
	specs.Column{Name: "flags", Type: "INTEGER", Description: "DNS record flags"},
}
