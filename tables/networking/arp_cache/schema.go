package arp_cache

import (
	"github.com/scrymastic/goosquery/sql/result"
)

var TableName = "arp_cache"
var Description = "Address resolution cache, both static and dynamic (from ARP, NDP)."
var Schema = result.Schema{
	result.Column{Name: "address", Type: "TEXT", Description: "IPv4 address target"},
	result.Column{Name: "mac", Type: "TEXT", Description: "MAC address of broadcasted address"},
	result.Column{Name: "interface", Type: "TEXT", Description: "Interface of the network for the MAC"},
	result.Column{Name: "permanent", Type: "TEXT", Description: "1 for true"},
}
