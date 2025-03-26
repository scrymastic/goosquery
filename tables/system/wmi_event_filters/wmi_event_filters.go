package wmi_event_filters

import (
	"fmt"

	"github.com/StackExchange/wmi"
	"github.com/scrymastic/goosquery/sql/result"
	"github.com/scrymastic/goosquery/sql/sqlctx"
)

type _EventFilter struct {
	__CLASS        string
	__RELPATH      string
	CreatorSID     []uint8
	EventAccess    string
	EventNamespace string
	Name           string
	Query          string
	QueryLanguage  string
}

func GenWmiEventFilters(ctx *sqlctx.Context) (*result.Results, error) {

	var filters []_EventFilter
	query := "SELECT * FROM __EventFilter"
	namespace := `ROOT\Subscription`
	if err := wmi.QueryNamespace(query, &filters, namespace); err != nil {
		return nil, fmt.Errorf("failed to query WMI event filters: %w", err)
	}

	// Missing __CLASS and __RELPATH fields

	filterInfo := result.NewQueryResult()
	for _, filter := range filters {
		info := result.NewResult(ctx, Schema)
		info.Set("name", filter.Name)
		info.Set("query", filter.Query)
		info.Set("query_language", filter.QueryLanguage)
		info.Set("class", filter.__CLASS)
		info.Set("relative_path", filter.__RELPATH)
		filterInfo.AppendResult(*info)
		// fmt.Printf("Raw WMI Object: %+v\n", filter)
	}

	return filterInfo, nil
}
