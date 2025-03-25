package kernel_info

import (
	"github.com/scrymastic/goosquery/sql/result"
)

var TableName = "kernel_info"
var Description = "Basic active kernel information."
var Schema = result.Schema{
	result.Column{Name: "version", Type: "TEXT", Description: "Kernel version"},
	result.Column{Name: "arguments", Type: "TEXT", Description: "Kernel arguments"},
	result.Column{Name: "path", Type: "TEXT", Description: "Kernel path"},
	result.Column{Name: "device", Type: "TEXT", Description: "Kernel device identifier"},
}
