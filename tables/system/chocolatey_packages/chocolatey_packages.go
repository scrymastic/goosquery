package chocolatey_packages

import (
	"encoding/xml"
	"fmt"
	"os"
	"path/filepath"

	"github.com/scrymastic/goosquery/sql/result"
	"github.com/scrymastic/goosquery/sql/sqlctx"
)

type NuspecXML struct {
	XMLName  xml.Name `xml:"package"`
	Metadata struct {
		ID         string `xml:"id"`
		Version    string `xml:"version"`
		Summary    string `xml:"summary"`
		Authors    string `xml:"authors"`
		LicenseUrl string `xml:"licenseUrl"`
	} `xml:"metadata"`
}

func parseNuspec(ctx *sqlctx.Context, nuspecPath string) (*result.Result, error) {
	content, err := os.ReadFile(nuspecPath)
	if err != nil {
		return nil, fmt.Errorf("failed to read nuspec file: %v", err)
	}

	var nuspec NuspecXML
	if err := xml.Unmarshal(content, &nuspec); err != nil {
		return nil, fmt.Errorf("failed to unmarshal nuspec file: %v", err)
	}

	pkg := result.NewResult(ctx, Schema)

	pkg.Set("name", nuspec.Metadata.ID)
	pkg.Set("version", nuspec.Metadata.Version)
	pkg.Set("summary", nuspec.Metadata.Summary)
	pkg.Set("author", nuspec.Metadata.Authors)
	pkg.Set("license", nuspec.Metadata.LicenseUrl)
	pkg.Set("path", nuspecPath)

	return pkg, nil
}

func GenChocolateyPackages(ctx *sqlctx.Context) (*result.Results, error) {
	chocoInstall := os.Getenv("ChocolateyInstall")
	if chocoInstall == "" {
		return nil, nil
	}

	pattern := filepath.Join(chocoInstall, "lib", "*", "*.nuspec")
	matches, err := filepath.Glob(pattern)
	if err != nil {
		return nil, fmt.Errorf("failed to glob nuspec files: %v", err)
	}

	packages := result.NewQueryResult()
	for _, nuspecPath := range matches {
		pkg, err := parseNuspec(ctx, nuspecPath)
		if err != nil {
			continue
		}
		packages.AppendResult(*pkg)
	}

	return packages, nil
}
