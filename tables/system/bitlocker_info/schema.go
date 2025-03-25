package bitlocker_info

import (
	"github.com/scrymastic/goosquery/sql/result"
)

var TableName = "bitlocker_info"
var Description = "Retrieve bitlocker status of the machine."
var Schema = result.Schema{
	result.Column{Name: "device_id", Type: "TEXT", Description: "ID of the encrypted drive."},
	result.Column{Name: "drive_letter", Type: "TEXT", Description: "Drive letter of the encrypted drive."},
	result.Column{Name: "persistent_volume_id", Type: "TEXT", Description: "Persistent ID of the drive."},
	result.Column{Name: "conversion_status", Type: "INTEGER", Description: "The bitlocker conversion status of the drive."},
	result.Column{Name: "protection_status", Type: "INTEGER", Description: "The bitlocker protection status of the drive."},
	result.Column{Name: "encryption_method", Type: "TEXT", Description: "The encryption type of the device."},
	result.Column{Name: "version", Type: "INTEGER", Description: "The FVE metadata version of the drive."},
	result.Column{Name: "percentage_encrypted", Type: "INTEGER", Description: "The percentage of the drive that is encrypted."},
	result.Column{Name: "lock_status", Type: "INTEGER", Description: "The accessibility status of the drive from Windows."},
}
