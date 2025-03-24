package platform_info

import (
	"github.com/scrymastic/goosquery/tables/specs"
)

var TableName = "platform_info"
var Description = "Information about EFI/UEFI/ROM and platform/boot."
var Schema = specs.Schema{
	specs.Column{Name: "vendor", Type: "TEXT", Description: "Platform code vendor"},
	specs.Column{Name: "version", Type: "TEXT", Description: "Platform code version"},
	specs.Column{Name: "date", Type: "TEXT", Description: "Self-reported platform code update date"},
	specs.Column{Name: "revision", Type: "TEXT", Description: "BIOS major and minor revision"},
	specs.Column{Name: "extra", Type: "TEXT", Description: "Platform-specific additional information"},
	specs.Column{Name: "firmware_type", Type: "TEXT", Description: "The type of firmware (uefi, bios, iboot, openfirmware, unknown)."},
}
