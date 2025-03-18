package python_packages

import (
	"bufio"
	"os"
	"path/filepath"
	"strings"

	"golang.org/x/sys/windows/registry"
)

// PythonPackage represents a Python package's metadata
type PythonPackage struct {
	Name      string `json:"name"`
	Version   string `json:"version"`
	Summary   string `json:"summary"`
	Author    string `json:"author"`
	License   string `json:"license"`
	Path      string `json:"path"`
	Directory string `json:"directory"`
}

// readPackageMetadata reads and parses package metadata from METADATA or PKG-INFO files
func readPackageMetadata(path string, directory string) (*PythonPackage, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	pkg := &PythonPackage{
		Path:      path,
		Directory: directory,
	}

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
			pkg.Name = value
		case "Version":
			pkg.Version = value
		case "Summary":
			pkg.Summary = value
		case "Author":
			pkg.Author = value
		case "License":
			pkg.License = value
		}
	}

	return pkg, scanner.Err()
}

// scanSitePackages scans a directory for Python packages
func scanSitePackages(siteDir string) ([]PythonPackage, error) {
	var packages []PythonPackage

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

		pkg, err := readPackageMetadata(metadataPath, siteDir)
		if err != nil {
			continue // Skip packages with unreadable metadata
		}

		packages = append(packages, *pkg)
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
func GenPythonPackages() ([]PythonPackage, error) {
	var allPackages []PythonPackage

	// Get Python installation paths from registry
	paths, err := getPythonInstallPaths()
	if err != nil {
		return nil, err
	}

	// Scan each site-packages directory
	for _, path := range paths {
		packages, err := scanSitePackages(path)
		if err != nil {
			continue // Skip directories we can't read
		}
		allPackages = append(allPackages, packages...)
	}

	return allPackages, nil
}
