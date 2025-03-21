package utility

import (
	ctxPkg "github.com/scrymastic/goosquery/sql/context"
	"github.com/scrymastic/goosquery/tables/utility/file"
	time_info "github.com/scrymastic/goosquery/tables/utility/time"
)

// GenFile generates file information for the specified path and directory
func GenFile(ctx ctxPkg.Context) ([]map[string]interface{}, error) {
	return file.GenFiles(ctx)
}

// GenTime generates current date and time information
func GenTime(ctx ctxPkg.Context) ([]map[string]interface{}, error) {
	return time_info.GenTime(ctx)
}
