package ec2_instance_metadata

import (
	"github.com/scrymastic/goosquery/sql/result"
)

var TableName = "ec2_instance_metadata"
var Description = "EC2 instance metadata."
var Schema = result.Schema{
	result.Column{Name: "instance_id", Type: "TEXT", Description: "EC2 instance ID"},
	result.Column{Name: "instance_type", Type: "TEXT", Description: "EC2 instance type"},
	result.Column{Name: "architecture", Type: "TEXT", Description: "Hardware architecture of this EC2 instance"},
	result.Column{Name: "region", Type: "TEXT", Description: "AWS region in which this instance launched"},
	result.Column{Name: "availability_zone", Type: "TEXT", Description: "Availability zone in which this instance launched"},
	result.Column{Name: "local_hostname", Type: "TEXT", Description: "Private IPv4 DNS hostname of the first interface of this instance"},
	result.Column{Name: "local_ipv4", Type: "TEXT", Description: "Private IPv4 address of the first interface of this instance"},
	result.Column{Name: "mac", Type: "TEXT", Description: "MAC address for the first network interface of this EC2 instance"},
	result.Column{Name: "security_groups", Type: "TEXT", Description: "Comma separated list of security group names"},
	result.Column{Name: "iam_arn", Type: "TEXT", Description: "If there is an IAM role associated with the instance"},
	result.Column{Name: "ami_id", Type: "TEXT", Description: "AMI ID used to launch this EC2 instance"},
	result.Column{Name: "reservation_id", Type: "TEXT", Description: "ID of the reservation"},
	result.Column{Name: "account_id", Type: "TEXT", Description: "AWS account ID which owns this EC2 instance"},
	result.Column{Name: "ssh_public_key", Type: "TEXT", Description: "SSH public key. Only available if supplied at instance launch time"},
}
