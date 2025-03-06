package chocolatey_packages

import (
	"encoding/xml"
	"fmt"
	"os"
	"path/filepath"
)

type ChocolateyPackage struct {
	Name    string `json:"name"`
	Version string `json:"version"`
	Summary string `json:"summary"`
	Author  string `json:"author"`
	License string `json:"license"`
	Path    string `json:"path"`
}

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

func parseNuspec(nuspecPath string) (*ChocolateyPackage, error) {
	content, err := os.ReadFile(nuspecPath)
	if err != nil {
		return nil, fmt.Errorf("failed to read nuspec file: %v", err)
	}

	var nuspec NuspecXML
	if err := xml.Unmarshal(content, &nuspec); err != nil {
		return nil, fmt.Errorf("failed to unmarshal nuspec file: %v", err)
	}

	return &ChocolateyPackage{
		Name:    nuspec.Metadata.ID,
		Version: nuspec.Metadata.Version,
		Summary: nuspec.Metadata.Summary,
		Author:  nuspec.Metadata.Authors,
		License: nuspec.Metadata.LicenseUrl,
		Path:    nuspecPath,
	}, nil
}

func GenChocolateyPackages() ([]ChocolateyPackage, error) {
	chocoInstall := os.Getenv("ChocolateyInstall")
	if chocoInstall == "" {
		return nil, nil
	}

	pattern := filepath.Join(chocoInstall, "lib", "*", "*.nuspec")
	matches, err := filepath.Glob(pattern)
	if err != nil {
		return nil, fmt.Errorf("failed to glob nuspec files: %v", err)
	}

	var packages []ChocolateyPackage
	for _, nuspecPath := range matches {
		pkg, err := parseNuspec(nuspecPath)
		if err != nil {
			continue
		}
		packages = append(packages, *pkg)
	}

	return packages, nil
}
