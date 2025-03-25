package ec2_instance_tags

import (
	"github.com/scrymastic/goosquery/sql/result"
)

var TableName = "ec2_instance_tags"
var Description = "EC2 instance tag key value pairs."
var Schema = result.Schema{
	result.Column{Name: "instance_id", Type: "TEXT", Description: "EC2 instance ID"},
	result.Column{Name: "key", Type: "TEXT", Description: "Tag key"},
	result.Column{Name: "value", Type: "TEXT", Description: "Tag value"},
}
