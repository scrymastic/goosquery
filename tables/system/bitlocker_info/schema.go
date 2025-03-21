package bitlocker_info

import (
	"github.com/scrymastic/goosquery/tables/specs"
)

var TableName = "bitlocker_info"
var Description = "Retrieve bitlocker status of the machine."
var Schema = specs.Schema{
	specs.Column{Name: "device_id", Type: "string", Description: "ID of the encrypted drive."},
	specs.Column{Name: "drive_letter", Type: "string", Description: "Drive letter of the encrypted drive."},
	specs.Column{Name: "persistent_volume_id", Type: "string", Description: "Persistent ID of the drive."},
	specs.Column{Name: "conversion_status", Type: "int32", Description: "The bitlocker conversion status of the drive."},
	specs.Column{Name: "protection_status", Type: "int32", Description: "The bitlocker protection status of the drive."},
	specs.Column{Name: "encryption_method", Type: "string", Description: "The encryption type of the device."},
	specs.Column{Name: "version", Type: "int32", Description: "The FVE metadata version of the drive."},
	specs.Column{Name: "percentage_encrypted", Type: "int32", Description: "The percentage of the drive that is encrypted."},
	specs.Column{Name: "lock_status", Type: "int32", Description: "The accessibility status of the drive from Windows."},
}
