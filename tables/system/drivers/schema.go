package drivers

import (
	"github.com/scrymastic/goosquery/tables/specs"
)

var TableName = "drivers"
var Description = "Details for in-use Windows device drivers. This does not display installed but unused drivers."
var Schema = specs.Schema{
	specs.Column{Name: "device_id", Type: "TEXT", Description: "Device ID"},
	specs.Column{Name: "device_name", Type: "TEXT", Description: "Device name"},
	specs.Column{Name: "image", Type: "TEXT", Description: "Path to driver image file"},
	specs.Column{Name: "description", Type: "TEXT", Description: "Driver description"},
	specs.Column{Name: "service", Type: "TEXT", Description: "Driver service name"},
	specs.Column{Name: "service_key", Type: "TEXT", Description: "Driver service registry key"},
	specs.Column{Name: "version", Type: "TEXT", Description: "Driver version"},
	specs.Column{Name: "inf", Type: "TEXT", Description: "Associated inf file"},
	specs.Column{Name: "class", Type: "TEXT", Description: "Device/driver class name"},
	specs.Column{Name: "provider", Type: "TEXT", Description: "Driver provider"},
	specs.Column{Name: "manufacturer", Type: "TEXT", Description: "Device manufacturer"},
	specs.Column{Name: "driver_key", Type: "TEXT", Description: "Driver key"},
	specs.Column{Name: "date", Type: "BIGINT", Description: "Driver date"},
	specs.Column{Name: "signed", Type: "INTEGER", Description: "Whether the driver is signed or not"},
}
