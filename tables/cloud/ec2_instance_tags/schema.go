package ec2_instance_tags

import (
	"github.com/scrymastic/goosquery/tables/specs"
)

var TableName = "ec2_instance_tags"
var Description = "EC2 instance tag key value pairs."
var Schema = specs.Schema{
	specs.Column{Name: "instance_id", Type: "TEXT", Description: "EC2 instance ID"},
	specs.Column{Name: "key", Type: "TEXT", Description: "Tag key"},
	specs.Column{Name: "value", Type: "TEXT", Description: "Tag value"},
}
