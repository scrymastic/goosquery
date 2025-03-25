package ntdomains

import (
	"fmt"

	"github.com/StackExchange/wmi"
	"github.com/scrymastic/goosquery/sql/result"
	"github.com/scrymastic/goosquery/sql/sqlctx"
)

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

func GenNTDomains(ctx *sqlctx.Context) (*result.Results, error) {
	query := `select * from Win32_NtDomain`
	var dst []Win32_NTDomain
	err := wmi.Query(query, &dst)
	if err != nil {
		return nil, fmt.Errorf("failed to query Win32_NtDomain: %w", err)
	}

	domains := result.NewQueryResult()
	for _, d := range dst {
		domainInfo := result.NewResult(ctx, Schema)
		domainInfo.Set("name", d.Name)
		domainInfo.Set("client_site_name", d.ClientSiteName)
		domainInfo.Set("dc_site_name", d.DCSiteName)
		domainInfo.Set("dns_forest_name", d.DnsForestName)
		domainInfo.Set("domain_controller_address", d.DomainControllerAddress)
		domainInfo.Set("domain_controller_name", d.DomainControllerName)
		domainInfo.Set("domain_name", d.DomainName)
		domainInfo.Set("status", d.Status)
		domains.AppendResult(*domainInfo)
	}
	return domains, nil
}
