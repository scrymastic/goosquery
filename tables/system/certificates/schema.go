package certificates

import (
	"github.com/scrymastic/goosquery/sql/result"
)

var TableName = "certificates"
var Description = "Certificate Authorities installed in Keychains/ca-bundles. NOTE: osquery limits frequent access to keychain files on macOS. This limit is controlled by keychain_access_interval flag."
var Schema = result.Schema{
	result.Column{Name: "common_name", Type: "TEXT", Description: "Certificate CommonName"},
	result.Column{Name: "subject", Type: "TEXT", Description: "Certificate distinguished name (deprecated, use subject2)"},
	result.Column{Name: "issuer", Type: "TEXT", Description: "Certificate issuer distinguished name (deprecated, use issuer2)"},
	result.Column{Name: "ca", Type: "INTEGER", Description: "1 if CA: true (certificate is an authority) else 0"},
	result.Column{Name: "self_signed", Type: "INTEGER", Description: "1 if self-signed, else 0"},
	result.Column{Name: "not_valid_before", Type: "certificates", Description: "Lower bound of valid date"},
	result.Column{Name: "not_valid_after", Type: "DATETIME", Description: "Certificate expiration data"},
	result.Column{Name: "signing_algorithm", Type: "TEXT", Description: "Signing algorithm used"},
	result.Column{Name: "key_algorithm", Type: "TEXT", Description: "Key algorithm used"},
	result.Column{Name: "key_strength", Type: "TEXT", Description: "Key size used for RSA/DSA, or curve name"},
	result.Column{Name: "key_usage", Type: "TEXT", Description: "Certificate key usage and extended key usage"},
	result.Column{Name: "subject_key_id", Type: "TEXT", Description: "SKID an optionally included SHA1"},
	result.Column{Name: "authority_key_id", Type: "TEXT", Description: "AKID an optionally included SHA1"},
	result.Column{Name: "sha1", Type: "TEXT", Description: "SHA1 hash of the raw certificate contents"},
	result.Column{Name: "path", Type: "TEXT", Description: "Path to Keychain or PEM bundle"},
	result.Column{Name: "serial", Type: "TEXT", Description: "Certificate serial number"},

	result.Column{Name: "sid", Type: "TEXT", Description: "SID"},
	result.Column{Name: "store_location", Type: "TEXT", Description: "Certificate system store location"},
	result.Column{Name: "store", Type: "TEXT", Description: "Certificate system store"},
	result.Column{Name: "username", Type: "TEXT", Description: "Username"},
	result.Column{Name: "store_id", Type: "TEXT", Description: "Exists for service/user stores. Contains raw store id provided by WinAPI."},
}
