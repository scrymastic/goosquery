package arp_cache

import (
	"github.com/scrymastic/goosquery/tables/specs"
)

var TableName = "arp_cache"
var Description = "Address resolution cache, both static and dynamic (from ARP, NDP)."
var Schema = specs.Schema{
	specs.Column{Name: "address", Type: "string", Description: "IPv4 address target"},
	specs.Column{Name: "mac", Type: "string", Description: "MAC address of broadcasted address"},
	specs.Column{Name: "interface", Type: "string", Description: "Interface of the network for the MAC"},
	specs.Column{Name: "permanent", Type: "string", Description: "1 for true, 0 for false"},
}
