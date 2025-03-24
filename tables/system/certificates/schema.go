package certificates

import (
	"github.com/scrymastic/goosquery/tables/specs"
)

var TableName = "certificates"
var Description = "Certificate Authorities installed in Keychains/ca-bundles. NOTE: osquery limits frequent access to keychain files on macOS. This limit is controlled by keychain_access_interval flag."
var Schema = specs.Schema{
	specs.Column{Name: "common_name", Type: "TEXT", Description: "Certificate CommonName"},
	specs.Column{Name: "subject", Type: "TEXT", Description: "Certificate distinguished name (deprecated, use subject2)"},
	specs.Column{Name: "issuer", Type: "TEXT", Description: "Certificate issuer distinguished name (deprecated, use issuer2)"},
	specs.Column{Name: "ca", Type: "INTEGER", Description: "1 if CA: true (certificate is an authority) else 0"},
	specs.Column{Name: "self_signed", Type: "INTEGER", Description: "1 if self-signed, else 0"},
	specs.Column{Name: "not_valid_before", Type: "certificates", Description: "Lower bound of valid date"},
	specs.Column{Name: "not_valid_after", Type: "DATETIME", Description: "Certificate expiration data"},
	specs.Column{Name: "signing_algorithm", Type: "TEXT", Description: "Signing algorithm used"},
	specs.Column{Name: "key_algorithm", Type: "TEXT", Description: "Key algorithm used"},
	specs.Column{Name: "key_strength", Type: "TEXT", Description: "Key size used for RSA/DSA, or curve name"},
	specs.Column{Name: "key_usage", Type: "TEXT", Description: "Certificate key usage and extended key usage"},
	specs.Column{Name: "subject_key_id", Type: "TEXT", Description: "SKID an optionally included SHA1"},
	specs.Column{Name: "authority_key_id", Type: "TEXT", Description: "AKID an optionally included SHA1"},
	specs.Column{Name: "sha1", Type: "TEXT", Description: "SHA1 hash of the raw certificate contents"},
	specs.Column{Name: "path", Type: "TEXT", Description: "Path to Keychain or PEM bundle"},
	specs.Column{Name: "serial", Type: "TEXT", Description: "Certificate serial number"},

	specs.Column{Name: "sid", Type: "TEXT", Description: "SID"},
	specs.Column{Name: "store_location", Type: "TEXT", Description: "Certificate system store location"},
	specs.Column{Name: "store", Type: "TEXT", Description: "Certificate system store"},
	specs.Column{Name: "username", Type: "TEXT", Description: "Username"},
	specs.Column{Name: "store_id", Type: "TEXT", Description: "Exists for service/user stores. Contains raw store id provided by WinAPI."},
}
