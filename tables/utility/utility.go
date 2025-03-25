package utility

import (
	"github.com/scrymastic/goosquery/sql/result"
	"github.com/scrymastic/goosquery/sql/sqlctx"
	"github.com/scrymastic/goosquery/tables/utility/file"
	time_info "github.com/scrymastic/goosquery/tables/utility/time"
)

// GenFile generates file information for the specified path and directory
func GenFile(ctx *sqlctx.Context) (*result.Results, error) {
	return file.GenFiles(ctx)
}

// GenTime generates current date and time information
func GenTime(ctx *sqlctx.Context) (*result.Results, error) {
	return time_info.GenTime(ctx)
}
