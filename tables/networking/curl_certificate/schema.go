package curl_certificate

import (
	"github.com/scrymastic/goosquery/tables/specs"
)

var TableName = "curl_certificate"
var Description = "Inspect TLS certificates by connecting to input hostnames."
var Schema = specs.Schema{
	specs.Column{Name: "hostname", Type: "TEXT", Description: "Hostname to CURL"},
	specs.Column{Name: "common_name", Type: "TEXT", Description: "Common name of company issued to"},
	specs.Column{Name: "organization", Type: "TEXT", Description: "Organization issued to"},
	specs.Column{Name: "organization_unit", Type: "TEXT", Description: "Organization unit issued to"},
	specs.Column{Name: "serial_number", Type: "TEXT", Description: "Certificate serial number"},
	specs.Column{Name: "issuer_common_name", Type: "TEXT", Description: "Issuer common name"},
	specs.Column{Name: "issuer_organization", Type: "TEXT", Description: "Issuer organization"},
	specs.Column{Name: "issuer_organization_unit", Type: "TEXT", Description: "Issuer organization unit"},
	specs.Column{Name: "valid_from", Type: "TEXT", Description: "Period of validity start date"},
	specs.Column{Name: "valid_to", Type: "TEXT", Description: "Period of validity end date"},
	specs.Column{Name: "sha256_fingerprint", Type: "TEXT", Description: "SHA-256 fingerprint"},
	specs.Column{Name: "sha1_fingerprint", Type: "TEXT", Description: "SHA1 fingerprint"},
	specs.Column{Name: "version", Type: "INTEGER", Description: "Version Number"},
	specs.Column{Name: "signature_algorithm", Type: "TEXT", Description: "Signature Algorithm"},
	specs.Column{Name: "signature", Type: "TEXT", Description: "Signature"},
	specs.Column{Name: "subject_key_identifier", Type: "TEXT", Description: "Subject Key Identifier"},
	specs.Column{Name: "authority_key_identifier", Type: "TEXT", Description: "Authority Key Identifier"},
	specs.Column{Name: "key_usage", Type: "TEXT", Description: "Usage of key in certificate"},
	specs.Column{Name: "extended_key_usage", Type: "TEXT", Description: "Extended usage of key in certificate"},
	specs.Column{Name: "policies", Type: "TEXT", Description: "Certificate Policies"},
	specs.Column{Name: "subject_alternative_names", Type: "TEXT", Description: "Subject Alternative Name"},
	specs.Column{Name: "issuer_alternative_names", Type: "TEXT", Description: "Issuer Alternative Name"},
	specs.Column{Name: "info_access", Type: "TEXT", Description: "Authority Information Access"},
	specs.Column{Name: "subject_info_access", Type: "TEXT", Description: "Subject Information Access"},
	specs.Column{Name: "policy_mappings", Type: "TEXT", Description: "Policy Mappings"},
	specs.Column{Name: "has_expired", Type: "INTEGER", Description: "1 if the certificate has expired"},
	specs.Column{Name: "basic_constraint", Type: "TEXT", Description: "Basic Constraints"},
	specs.Column{Name: "name_constraints", Type: "TEXT", Description: "Name Constraints"},
	specs.Column{Name: "policy_constraints", Type: "TEXT", Description: "Policy Constraints"},
	specs.Column{Name: "dump_certificate", Type: "INTEGER", Description: "Set this value to 1 to dump certificate"},
	specs.Column{Name: "timeout", Type: "INTEGER", Description: "Set this value to the timeout in seconds to complete the TLS handshake"},
	specs.Column{Name: "pem", Type: "TEXT", Description: "Certificate PEM format"},
}
