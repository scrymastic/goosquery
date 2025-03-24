package video_info

import (
	"github.com/scrymastic/goosquery/tables/specs"
)

var TableName = "video_info"
var Description = "Retrieve video card information of the machine."
var Schema = specs.Schema{
	specs.Column{Name: "color_depth", Type: "INTEGER", Description: "The amount of bits per pixel to represent color."},
	specs.Column{Name: "driver", Type: "TEXT", Description: "The driver of the device."},
	specs.Column{Name: "driver_date", Type: "BIGINT", Description: "The date listed on the installed driver."},
	specs.Column{Name: "driver_version", Type: "TEXT", Description: "The version of the installed driver."},
	specs.Column{Name: "manufacturer", Type: "TEXT", Description: "The manufacturer of the gpu."},
	specs.Column{Name: "model", Type: "TEXT", Description: "The model of the gpu."},
	specs.Column{Name: "series", Type: "TEXT", Description: "The series of the gpu."},
	specs.Column{Name: "video_mode", Type: "TEXT", Description: "The current resolution of the display."},
}
