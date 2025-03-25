package azure_instance_metadata

import (
	"github.com/scrymastic/goosquery/sql/result"
)

var TableName = "azure_instance_metadata"
var Description = "Azure instance metadata."
var Schema = result.Schema{
	result.Column{Name: "location", Type: "TEXT", Description: "Azure Region the VM is running in"},
	result.Column{Name: "name", Type: "TEXT", Description: "Name of the VM"},
	result.Column{Name: "offer", Type: "TEXT", Description: "Offer information for the VM image"},
	result.Column{Name: "publisher", Type: "TEXT", Description: "Publisher of the VM image"},
	result.Column{Name: "sku", Type: "TEXT", Description: "SKU for the VM image"},
	result.Column{Name: "version", Type: "TEXT", Description: "Version of the VM image"},
	result.Column{Name: "os_type", Type: "TEXT", Description: "Linux or Windows"},
	result.Column{Name: "platform_update_domain", Type: "TEXT", Description: "Update domain the VM is running in"},
	result.Column{Name: "platform_fault_domain", Type: "TEXT", Description: "Fault domain the VM is running in"},
	result.Column{Name: "vm_id", Type: "TEXT", Description: "Unique identifier for the VM"},
	result.Column{Name: "vm_size", Type: "TEXT", Description: "VM size"},
	result.Column{Name: "subscription_id", Type: "TEXT", Description: "Azure subscription for the VM"},
	result.Column{Name: "resource_group_name", Type: "TEXT", Description: "Resource group for the VM"},
	result.Column{Name: "placement_group_id", Type: "TEXT", Description: "Placement group for the VM scale set"},
	result.Column{Name: "vm_scale_set_name", Type: "TEXT", Description: "VM scale set name"},
	result.Column{Name: "zone", Type: "TEXT", Description: "Availability zone of the VM"},
}
