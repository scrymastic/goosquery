package curl

import (
	"github.com/scrymastic/goosquery/tables/specs"
)

var TableName = "curl"
var Description = "Perform an http request and return stats about it."
var Schema = specs.Schema{
	specs.Column{Name: "url", Type: "string", Description: "The url for the request"},
	specs.Column{Name: "method", Type: "string", Description: "The HTTP method for the request"},
	specs.Column{Name: "user_agent", Type: "string", Description: "The user-agent string to use for the request"},
	specs.Column{Name: "response_code", Type: "int32", Description: "The HTTP status code for the response"},
	specs.Column{Name: "round_trip_time", Type: "int64", Description: "Time taken to complete the request"},
	specs.Column{Name: "bytes", Type: "int64", Description: "Number of bytes in the response"},
	specs.Column{Name: "result", Type: "string", Description: "The HTTP response body"},
}
