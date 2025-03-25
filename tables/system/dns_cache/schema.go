package dns_cache

import (
	"github.com/scrymastic/goosquery/sql/result"
)

var TableName = "dns_cache"
var Description = "Enumerate the DNS cache using the undocumented DnsGetCacheDataTable function in dnsapi.dll."
var Schema = result.Schema{
	result.Column{Name: "name", Type: "TEXT", Description: "DNS record name"},
	result.Column{Name: "type", Type: "TEXT", Description: "DNS record type"},
	result.Column{Name: "flags", Type: "INTEGER", Description: "DNS record flags"},
}
