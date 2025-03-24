package ssh_configs

import (
	"github.com/scrymastic/goosquery/tables/specs"
)

var TableName = "ssh_configs"
var Description = "A table of parsed ssh_configs."
var Schema = specs.Schema{
	specs.Column{Name: "uid", Type: "BIGINT", Description: "The local owner of the ssh_config file"},
	specs.Column{Name: "block", Type: "TEXT", Description: "The host or match block"},
	specs.Column{Name: "option", Type: "TEXT", Description: "The option and value"},
	specs.Column{Name: "ssh_config_file", Type: "TEXT", Description: "Path to the ssh_config file"},
}
