package kernel_info

import (
	"github.com/scrymastic/goosquery/tables/specs"
)

var TableName = "kernel_info"
var Description = "Basic active kernel information."
var Schema = specs.Schema{
	specs.Column{Name: "version", Type: "TEXT", Description: "Kernel version"},
	specs.Column{Name: "arguments", Type: "TEXT", Description: "Kernel arguments"},
	specs.Column{Name: "path", Type: "TEXT", Description: "Kernel path"},
	specs.Column{Name: "device", Type: "TEXT", Description: "Kernel device identifier"},
}
