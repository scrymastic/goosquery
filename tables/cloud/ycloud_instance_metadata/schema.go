package ycloud_instance_metadata

import (
	"github.com/scrymastic/goosquery/tables/specs"
)

var TableName = "ycloud_instance_metadata"
var Description = "Yandex.Cloud instance metadata."
var Schema = specs.Schema{
	specs.Column{Name: "instance_id", Type: "TEXT", Description: "Unique identifier for the VM"},
	specs.Column{Name: "folder_id", Type: "TEXT", Description: "Folder identifier for the VM"},
	specs.Column{Name: "cloud_id", Type: "TEXT", Description: "Cloud identifier for the VM"},
	specs.Column{Name: "name", Type: "TEXT", Description: "Name of the VM"},
	specs.Column{Name: "description", Type: "TEXT", Description: "Description of the VM"},
	specs.Column{Name: "hostname", Type: "TEXT", Description: "Hostname of the VM"},
	specs.Column{Name: "zone", Type: "TEXT", Description: "Availability zone of the VM"},
	specs.Column{Name: "ssh_public_key", Type: "TEXT", Description: "SSH public key. Only available if supplied at instance launch time"},
	specs.Column{Name: "serial_port_enabled", Type: "TEXT", Description: "Indicates if serial port is enabled for the VM"},
	specs.Column{Name: "metadata_endpoint", Type: "TEXT", Description: "Endpoint used to fetch VM metadata"},
}
