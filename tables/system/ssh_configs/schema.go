package ssh_configs

import (
	"github.com/scrymastic/goosquery/sql/result"
)

var TableName = "ssh_configs"
var Description = "A table of parsed ssh_configs."
var Schema = result.Schema{
	result.Column{Name: "uid", Type: "BIGINT", Description: "The local owner of the ssh_config file"},
	result.Column{Name: "block", Type: "TEXT", Description: "The host or match block"},
	result.Column{Name: "option", Type: "TEXT", Description: "The option and value"},
	result.Column{Name: "ssh_config_file", Type: "TEXT", Description: "Path to the ssh_config file"},
}
