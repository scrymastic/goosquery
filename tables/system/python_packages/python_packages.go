package python_packages

import (
	"bufio"
	"os"
	"path/filepath"
	"strings"

	"github.com/scrymastic/goosquery/sql/result"
	"github.com/scrymastic/goosquery/sql/sqlctx"
	"golang.org/x/sys/windows/registry"
)

// readPackageMetadata reads and parses package metadata from METADATA or PKG-INFO files
func readPackageMetadata(path string, directory string, ctx *sqlctx.Context) (*result.Result, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	pkg := result.NewResult(ctx, Schema)
	pkg.Set("path", path)
	pkg.Set("directory", directory)

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()
		parts := strings.SplitN(line, ":", 2)
		if len(parts) != 2 {
			continue
		}

		key := strings.TrimSpace(parts[0])
		value := strings.TrimSpace(parts[1])

		switch key {
		case "Name":
			pkg.Set("name", value)
		case "Version":
			pkg.Set("version", value)
		case "Summary":
			pkg.Set("summary", value)
		case "Author":
			pkg.Set("author", value)
		case "License":
			pkg.Set("license", value)
		}
	}

	return pkg, scanner.Err()
}

// scanSitePackages scans a directory for Python packages
func scanSitePackages(siteDir string, ctx *sqlctx.Context) (*result.Results, error) {
	packages := result.NewQueryResult()

	entries, err := os.ReadDir(siteDir)
	if err != nil {
		return nil, err
	}

	for _, entry := range entries {
		if !entry.IsDir() {
			continue
		}

		dirName := entry.Name()
		fullPath := filepath.Join(siteDir, dirName)

		var metadataPath string
		if strings.HasSuffix(dirName, ".dist-info") {
			metadataPath = filepath.Join(fullPath, "METADATA")
		} else if strings.HasSuffix(dirName, ".egg-info") {
			metadataPath = filepath.Join(fullPath, "PKG-INFO")
		} else {
			continue
		}

		pkg, err := readPackageMetadata(metadataPath, siteDir, ctx)
		if err != nil {
			continue // Skip packages with unreadable metadata
		}

		packages.AppendResult(*pkg)
	}

	return packages, nil
}

// getPythonInstallPaths gets Python installation paths from Windows registry
func getPythonInstallPaths() ([]string, error) {
	var paths []string

	// Check HKEY_LOCAL_MACHINE
	k, err := registry.OpenKey(registry.LOCAL_MACHINE, `SOFTWARE\Python\PythonCore`, registry.READ)
	if err == nil {
		defer k.Close()
		versions, err := k.ReadSubKeyNames(-1)
		if err == nil {
			for _, ver := range versions {
				installPath := `SOFTWARE\Python\PythonCore\` + ver + `\InstallPath`
				k2, err := registry.OpenKey(registry.LOCAL_MACHINE, installPath, registry.READ)
				if err == nil {
					if path, _, err := k2.GetStringValue(""); err == nil {
						sitePkgs := filepath.Join(path, "Lib", "site-packages")
						paths = append(paths, sitePkgs)
					}
					k2.Close()
				}
			}
		}
	}

	// Check HKEY_CURRENT_USER (similar pattern for all users would use HKEY_USERS)
	k, err = registry.OpenKey(registry.CURRENT_USER, `SOFTWARE\Python\PythonCore`, registry.READ)
	if err == nil {
		defer k.Close()
		versions, err := k.ReadSubKeyNames(-1)
		if err == nil {
			for _, ver := range versions {
				installPath := `SOFTWARE\Python\PythonCore\` + ver + `\InstallPath`
				k2, err := registry.OpenKey(registry.CURRENT_USER, installPath, registry.READ)
				if err == nil {
					if path, _, err := k2.GetStringValue(""); err == nil {
						sitePkgs := filepath.Join(path, "Lib", "site-packages")
						paths = append(paths, sitePkgs)
					}
					k2.Close()
				}
			}
		}
	}

	return paths, nil
}

// GenPythonPackages returns all Python packages installed on the system
func GenPythonPackages(ctx *sqlctx.Context) (*result.Results, error) {
	allPackages := result.NewQueryResult()

	// Get Python installation paths from registry
	paths, err := getPythonInstallPaths()
	if err != nil {
		return nil, err
	}

	// Scan each site-packages directory
	for _, path := range paths {
		packages, err := scanSitePackages(path, ctx)
		if err != nil {
			continue // Skip directories we can't read
		}
		allPackages.AppendResults(*packages)
	}

	return allPackages, nil
}
