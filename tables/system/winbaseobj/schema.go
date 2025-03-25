package winbaseobj

import (
	"github.com/scrymastic/goosquery/sql/result"
)

var TableName = "winbaseobj"
var Description = "Lists named Windows objects in the default object directories, across all terminal services sessions.  Example Windows ojbect types include Mutexes, Events, Jobs and Semaphors."
var Schema = result.Schema{
	result.Column{Name: "session_id", Type: "INTEGER", Description: "Terminal Services Session Id"},
	result.Column{Name: "object_name", Type: "TEXT", Description: "Object Name"},
	result.Column{Name: "object_type", Type: "TEXT", Description: "Object Type"},
}
