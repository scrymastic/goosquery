package wmi_event_filters

import (
	"fmt"

	"github.com/StackExchange/wmi"
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

type WMIEventFilter struct {
	Name          string `json:"name"`
	Query         string `json:"query"`
	QueryLanguage string `json:"query_language"`
	Class         string `json:"class"`
	RelativePath  string `json:"relative_path"`
}

func GenWMIEventFilters() ([]WMIEventFilter, error) {

	var filters []_EventFilter
	query := "SELECT * FROM __EventFilter"
	namespace := `ROOT\Subscription`
	if err := wmi.QueryNamespace(query, &filters, namespace); err != nil {
		return nil, fmt.Errorf("failed to query WMI event filters: %w", err)
	}

	// Missing __CLASS and __RELPATH fields

	if len(filters) == 0 {
		return nil, fmt.Errorf("no WMI event filters found")
	}

	filterInfo := make([]WMIEventFilter, 0, len(filters))
	for _, filter := range filters {
		info := WMIEventFilter{
			Name:          filter.Name,
			Query:         filter.Query,
			QueryLanguage: filter.QueryLanguage,
			Class:         filter.__CLASS,
			RelativePath:  filter.__RELPATH,
		}
		filterInfo = append(filterInfo, info)
		fmt.Printf("Raw WMI Object: %+v\n", filter)
	}

	return filterInfo, nil
}
