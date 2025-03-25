package tpm_info

import (
	"github.com/scrymastic/goosquery/sql/result"
)

var TableName = "tpm_info"
var Description = "A table that lists the TPM related information."
var Schema = result.Schema{
	result.Column{Name: "activated", Type: "INTEGER", Description: "TPM is activated"},
	result.Column{Name: "enabled", Type: "INTEGER", Description: "TPM is enabled"},
	result.Column{Name: "owned", Type: "INTEGER", Description: "TPM is owned"},
	result.Column{Name: "manufacturer_version", Type: "TEXT", Description: "TPM version"},
	result.Column{Name: "manufacturer_id", Type: "INTEGER", Description: "TPM manufacturers ID"},
	result.Column{Name: "manufacturer_name", Type: "TEXT", Description: "TPM manufacturers name"},
	result.Column{Name: "product_name", Type: "TEXT", Description: "Product name of the TPM"},
	result.Column{Name: "physical_presence_version", Type: "TEXT", Description: "Version of the Physical Presence Interface"},
	result.Column{Name: "spec_version", Type: "TEXT", Description: "Trusted Computing Group specification that the TPM supports"},
}
