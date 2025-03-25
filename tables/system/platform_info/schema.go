package platform_info

import (
	"github.com/scrymastic/goosquery/sql/result"
)

var TableName = "platform_info"
var Description = "Information about EFI/UEFI/ROM and platform/boot."
var Schema = result.Schema{
	result.Column{Name: "vendor", Type: "TEXT", Description: "Platform code vendor"},
	result.Column{Name: "version", Type: "TEXT", Description: "Platform code version"},
	result.Column{Name: "date", Type: "TEXT", Description: "Self-reported platform code update date"},
	result.Column{Name: "revision", Type: "TEXT", Description: "BIOS major and minor revision"},
	result.Column{Name: "extra", Type: "TEXT", Description: "Platform-specific additional information"},
	result.Column{Name: "firmware_type", Type: "TEXT", Description: "The type of firmware (uefi, bios, iboot, openfirmware, unknown)."},
}
