package curl

import (
	"github.com/scrymastic/goosquery/sql/result"
)

var TableName = "curl"
var Description = "Perform an http request and return stats about it."
var Schema = result.Schema{
	result.Column{Name: "url", Type: "TEXT", Description: "The url for the request"},
	result.Column{Name: "method", Type: "TEXT", Description: "The HTTP method for the request"},
	result.Column{Name: "response_code", Type: "INTEGER", Description: "The HTTP status code for the response"},
	result.Column{Name: "round_trip_time", Type: "BIGINT", Description: "Time taken to complete the request"},
	result.Column{Name: "bytes", Type: "BIGINT", Description: "Number of bytes in the response"},
	result.Column{Name: "result", Type: "TEXT", Description: "The HTTP response body"},
}
