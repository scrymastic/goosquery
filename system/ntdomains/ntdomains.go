package ntdomains

import (
	"fmt"

	"github.com/StackExchange/wmi"
)

type NTDomain struct {
	Name                    string `json:"name"`
	ClientSiteName          string `json:"client_site_name"`
	DCSiteName              string `json:"dc_site_name"`
	DnsForestName           string `json:"dns_forest_name"`
	DomainControllerAddress string `json:"domain_controller_address"`
	DomainControllerName    string `json:"domain_controller_name"`
	DomainName              string `json:"domain_name"`
	Status                  string `json:"status"`
}

type Win32_NTDomain struct {
	Name                    string
	ClientSiteName          string
	DCSiteName              string
	DnsForestName           string
	DomainControllerAddress string
	DomainControllerName    string
	DomainName              string
	Status                  string
}

func GenNTDomains() ([]NTDomain, error) {
	query := `select * from Win32_NtDomain`
	var dst []Win32_NTDomain
	err := wmi.Query(query, &dst)
	if err != nil {
		return nil, fmt.Errorf("failed to query Win32_NtDomain: %w", err)
	}

	results := make([]NTDomain, 0, len(dst))
	for _, d := range dst {
		results = append(results, NTDomain(d))
	}
	return results, nil
}
