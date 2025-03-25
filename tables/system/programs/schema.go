package programs

import (
	"github.com/scrymastic/goosquery/sql/result"
)

var TableName = "programs"
var Description = "Represents products as they are installed by Windows Installer. A product generally correlates to one installation package on Windows. Some fields may be blank as Windows installation details are left to the discretion of the product author."
var Schema = result.Schema{
	result.Column{Name: "name", Type: "TEXT", Description: "Commonly used product name."},
	result.Column{Name: "version", Type: "TEXT", Description: "Product version information."},
	result.Column{Name: "install_location", Type: "TEXT", Description: "The installation location directory of the product."},
	result.Column{Name: "install_source", Type: "TEXT", Description: "The installation source of the product."},
	result.Column{Name: "language", Type: "TEXT", Description: "The language of the product."},
	result.Column{Name: "publisher", Type: "TEXT", Description: "Name of the product supplier."},
	result.Column{Name: "uninstall_string", Type: "TEXT", Description: "Path and filename of the uninstaller."},
	result.Column{Name: "install_date", Type: "TEXT", Description: "Date that this product was installed on the system. "},
	result.Column{Name: "identifying_number", Type: "TEXT", Description: "Product identification such as a serial number on software"},
}
