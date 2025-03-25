package physical_disk_performance

import (
	"github.com/scrymastic/goosquery/sql/result"
)

var TableName = "physical_disk_performance"
var Description = "Provides provides raw data from performance counters that monitor hard or fixed disk drives on the system."
var Schema = result.Schema{
	result.Column{Name: "name", Type: "TEXT", Description: "Name of the physical disk"},
	result.Column{Name: "avg_disk_bytes_per_read", Type: "BIGINT", Description: "Average number of bytes transferred from the disk during read operations"},
	result.Column{Name: "avg_disk_bytes_per_write", Type: "BIGINT", Description: "Average number of bytes transferred to the disk during write operations"},
	result.Column{Name: "avg_disk_read_queue_length", Type: "BIGINT", Description: "Average number of read requests that were queued for the selected disk during the sample interval"},
	result.Column{Name: "avg_disk_write_queue_length", Type: "BIGINT", Description: "Average number of write requests that were queued for the selected disk during the sample interval"},
	result.Column{Name: "avg_disk_sec_per_read", Type: "INTEGER", Description: "Average time"},
	result.Column{Name: "avg_disk_sec_per_write", Type: "INTEGER", Description: "Average time"},
	result.Column{Name: "current_disk_queue_length", Type: "INTEGER", Description: "Number of requests outstanding on the disk at the time the performance data is collected"},
	result.Column{Name: "percent_disk_read_time", Type: "BIGINT", Description: "Percentage of elapsed time that the selected disk drive is busy servicing read requests"},
	result.Column{Name: "percent_disk_write_time", Type: "BIGINT", Description: "Percentage of elapsed time that the selected disk drive is busy servicing write requests"},
	result.Column{Name: "percent_disk_time", Type: "BIGINT", Description: "Percentage of elapsed time that the selected disk drive is busy servicing read or write requests"},
	result.Column{Name: "percent_idle_time", Type: "BIGINT", Description: "Percentage of time during the sample interval that the disk was idle"},
}
