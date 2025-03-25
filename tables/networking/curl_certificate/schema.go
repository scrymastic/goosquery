package curl_certificate

import (
	"github.com/scrymastic/goosquery/sql/result"
)

var TableName = "curl_certificate"
var Description = "Inspect TLS certificates by connecting to input hostnames."
var Schema = result.Schema{
	result.Column{Name: "hostname", Type: "TEXT", Description: "Hostname to CURL"},
	result.Column{Name: "common_name", Type: "TEXT", Description: "Common name of company issued to"},
	result.Column{Name: "organization", Type: "TEXT", Description: "Organization issued to"},
	result.Column{Name: "organization_unit", Type: "TEXT", Description: "Organization unit issued to"},
	result.Column{Name: "serial_number", Type: "TEXT", Description: "Certificate serial number"},
	result.Column{Name: "issuer_common_name", Type: "TEXT", Description: "Issuer common name"},
	result.Column{Name: "issuer_organization", Type: "TEXT", Description: "Issuer organization"},
	result.Column{Name: "issuer_organization_unit", Type: "TEXT", Description: "Issuer organization unit"},
	result.Column{Name: "valid_from", Type: "TEXT", Description: "Period of validity start date"},
	result.Column{Name: "valid_to", Type: "TEXT", Description: "Period of validity end date"},
	result.Column{Name: "sha256_fingerprint", Type: "TEXT", Description: "SHA-256 fingerprint"},
	result.Column{Name: "sha1_fingerprint", Type: "TEXT", Description: "SHA1 fingerprint"},
	result.Column{Name: "version", Type: "INTEGER", Description: "Version Number"},
	result.Column{Name: "signature_algorithm", Type: "TEXT", Description: "Signature Algorithm"},
	result.Column{Name: "signature", Type: "TEXT", Description: "Signature"},
	result.Column{Name: "subject_key_identifier", Type: "TEXT", Description: "Subject Key Identifier"},
	result.Column{Name: "authority_key_identifier", Type: "TEXT", Description: "Authority Key Identifier"},
	result.Column{Name: "key_usage", Type: "TEXT", Description: "Usage of key in certificate"},
	result.Column{Name: "extended_key_usage", Type: "TEXT", Description: "Extended usage of key in certificate"},
	result.Column{Name: "policies", Type: "TEXT", Description: "Certificate Policies"},
	result.Column{Name: "subject_alternative_names", Type: "TEXT", Description: "Subject Alternative Name"},
	result.Column{Name: "issuer_alternative_names", Type: "TEXT", Description: "Issuer Alternative Name"},
	result.Column{Name: "info_access", Type: "TEXT", Description: "Authority Information Access"},
	result.Column{Name: "subject_info_access", Type: "TEXT", Description: "Subject Information Access"},
	result.Column{Name: "policy_mappings", Type: "TEXT", Description: "Policy Mappings"},
	result.Column{Name: "has_expired", Type: "INTEGER", Description: "1 if the certificate has expired"},
	result.Column{Name: "basic_constraint", Type: "TEXT", Description: "Basic Constraints"},
	result.Column{Name: "name_constraints", Type: "TEXT", Description: "Name Constraints"},
	result.Column{Name: "policy_constraints", Type: "TEXT", Description: "Policy Constraints"},
	result.Column{Name: "dump_certificate", Type: "INTEGER", Description: "Set this value to 1 to dump certificate"},
	result.Column{Name: "timeout", Type: "INTEGER", Description: "Set this value to the timeout in seconds to complete the TLS handshake"},
	result.Column{Name: "pem", Type: "TEXT", Description: "Certificate PEM format"},
}
