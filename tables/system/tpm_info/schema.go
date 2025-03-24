package tpm_info

import (
	"github.com/scrymastic/goosquery/tables/specs"
)

var TableName = "tpm_info"
var Description = "A table that lists the TPM related information."
var Schema = specs.Schema{
	specs.Column{Name: "activated", Type: "INTEGER", Description: "TPM is activated"},
	specs.Column{Name: "enabled", Type: "INTEGER", Description: "TPM is enabled"},
	specs.Column{Name: "owned", Type: "INTEGER", Description: "TPM is owned"},
	specs.Column{Name: "manufacturer_version", Type: "TEXT", Description: "TPM version"},
	specs.Column{Name: "manufacturer_id", Type: "INTEGER", Description: "TPM manufacturers ID"},
	specs.Column{Name: "manufacturer_name", Type: "TEXT", Description: "TPM manufacturers name"},
	specs.Column{Name: "product_name", Type: "TEXT", Description: "Product name of the TPM"},
	specs.Column{Name: "physical_presence_version", Type: "TEXT", Description: "Version of the Physical Presence Interface"},
	specs.Column{Name: "spec_version", Type: "TEXT", Description: "Trusted Computing Group specification that the TPM supports"},
}
