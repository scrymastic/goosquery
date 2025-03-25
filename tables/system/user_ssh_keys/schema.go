package user_ssh_keys

import (
	"github.com/scrymastic/goosquery/sql/result"
)

var TableName = "user_ssh_keys"
var Description = "Returns the private keys in the users ~/.ssh directory and whether or not they are encrypted."
var Schema = result.Schema{
	result.Column{Name: "uid", Type: "BIGINT", Description: "The local user that owns the key file"},
	result.Column{Name: "path", Type: "TEXT", Description: "Path to key file"},
	result.Column{Name: "encrypted", Type: "INTEGER", Description: "1 if key is encrypted, 0 otherwise"},
	result.Column{Name: "key_type", Type: "TEXT", Description: "The type of the private key. One of [rsa, dsa, dh, ec, hmac, cmac], or the empty string."},
	result.Column{Name: "key_group_name", Type: "TEXT", Description: "The group of the private key. Supported for a subset of key_types implemented by OpenSSL"},
	result.Column{Name: "key_length", Type: "INTEGER", Description: "The cryptographic length of the cryptosystem to which the private key belongs, in bits. Definition of cryptographic length is specific to cryptosystem. -1 if unavailable"},
	result.Column{Name: "key_security_bits", Type: "INTEGER", Description: "The number of security bits of the private key, bits of security as defined in NIST SP800-57. -1 if unavailable"},
}
