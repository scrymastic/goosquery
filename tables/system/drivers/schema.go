package drivers

import (
	"github.com/scrymastic/goosquery/sql/result"
)

var TableName = "drivers"
var Description = "Details for in-use Windows device drivers. This does not display installed but unused drivers."
var Schema = result.Schema{
	result.Column{Name: "device_id", Type: "TEXT", Description: "Device ID"},
	result.Column{Name: "device_name", Type: "TEXT", Description: "Device name"},
	result.Column{Name: "image", Type: "TEXT", Description: "Path to driver image file"},
	result.Column{Name: "description", Type: "TEXT", Description: "Driver description"},
	result.Column{Name: "service", Type: "TEXT", Description: "Driver service name"},
	result.Column{Name: "service_key", Type: "TEXT", Description: "Driver service registry key"},
	result.Column{Name: "version", Type: "TEXT", Description: "Driver version"},
	result.Column{Name: "inf", Type: "TEXT", Description: "Associated inf file"},
	result.Column{Name: "class", Type: "TEXT", Description: "Device/driver class name"},
	result.Column{Name: "provider", Type: "TEXT", Description: "Driver provider"},
	result.Column{Name: "manufacturer", Type: "TEXT", Description: "Device manufacturer"},
	result.Column{Name: "driver_key", Type: "TEXT", Description: "Driver key"},
	result.Column{Name: "date", Type: "BIGINT", Description: "Driver date"},
	result.Column{Name: "signed", Type: "INTEGER", Description: "Whether the driver is signed or not"},
}
