package curl

import (
	"github.com/scrymastic/goosquery/tables/specs"
)

var TableName = "curl"
var Description = "Perform an http request and return stats about it."
var Schema = specs.Schema{
	specs.Column{Name: "url", Type: "TEXT", Description: "The url for the request"},
	specs.Column{Name: "method", Type: "TEXT", Description: "The HTTP method for the request"},
	specs.Column{Name: "response_code", Type: "INTEGER", Description: "The HTTP status code for the response"},
	specs.Column{Name: "round_trip_time", Type: "BIGINT", Description: "Time taken to complete the request"},
	specs.Column{Name: "bytes", Type: "BIGINT", Description: "Number of bytes in the response"},
	specs.Column{Name: "result", Type: "TEXT", Description: "The HTTP response body"},
}
