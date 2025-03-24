package azure_instance_metadata

import (
	"github.com/scrymastic/goosquery/tables/specs"
)

var TableName = "azure_instance_metadata"
var Description = "Azure instance metadata."
var Schema = specs.Schema{
	specs.Column{Name: "location", Type: "TEXT", Description: "Azure Region the VM is running in"},
	specs.Column{Name: "name", Type: "TEXT", Description: "Name of the VM"},
	specs.Column{Name: "offer", Type: "TEXT", Description: "Offer information for the VM image"},
	specs.Column{Name: "publisher", Type: "TEXT", Description: "Publisher of the VM image"},
	specs.Column{Name: "sku", Type: "TEXT", Description: "SKU for the VM image"},
	specs.Column{Name: "version", Type: "TEXT", Description: "Version of the VM image"},
	specs.Column{Name: "os_type", Type: "TEXT", Description: "Linux or Windows"},
	specs.Column{Name: "platform_update_domain", Type: "TEXT", Description: "Update domain the VM is running in"},
	specs.Column{Name: "platform_fault_domain", Type: "TEXT", Description: "Fault domain the VM is running in"},
	specs.Column{Name: "vm_id", Type: "TEXT", Description: "Unique identifier for the VM"},
	specs.Column{Name: "vm_size", Type: "TEXT", Description: "VM size"},
	specs.Column{Name: "subscription_id", Type: "TEXT", Description: "Azure subscription for the VM"},
	specs.Column{Name: "resource_group_name", Type: "TEXT", Description: "Resource group for the VM"},
	specs.Column{Name: "placement_group_id", Type: "TEXT", Description: "Placement group for the VM scale set"},
	specs.Column{Name: "vm_scale_set_name", Type: "TEXT", Description: "VM scale set name"},
	specs.Column{Name: "zone", Type: "TEXT", Description: "Availability zone of the VM"},
}
