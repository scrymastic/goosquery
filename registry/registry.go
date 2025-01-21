package registry

import (
	"fmt"
	"strings"

	"golang.org/x/sys/windows"
	"golang.org/x/sys/windows/registry"
)

type Registry struct {
	Key   string `json:"key"`
	Path  string `json:"path"`
	Name  string `json:"name"`
	Type  string `json:"type"`
	Data  string `json:"data"`
	MTime uint64 `json:"mtime"`
}

var registryTypeMap = map[uint32]string{
	windows.REG_BINARY:                   "REG_BINARY",
	windows.REG_DWORD:                    "REG_DWORD",
	windows.REG_DWORD_BIG_ENDIAN:         "REG_DWORD_BIG_ENDIAN",
	windows.REG_EXPAND_SZ:                "REG_EXPAND_SZ",
	windows.REG_LINK:                     "REG_LINK",
	windows.REG_MULTI_SZ:                 "REG_MULTI_SZ",
	windows.REG_NONE:                     "REG_NONE",
	windows.REG_QWORD:                    "REG_QWORD",
	windows.REG_SZ:                       "REG_SZ",
	windows.REG_FULL_RESOURCE_DESCRIPTOR: "REG_FULL_RESOURCE_DESCRIPTOR",
	windows.REG_RESOURCE_LIST:            "REG_RESOURCE_LIST",
}

func parseSearchKey(searchKey string) (registry.Key, string, error) {
	parts := strings.SplitN(searchKey, "\\", 2)

	var rootKey registry.Key
	switch strings.ToUpper(parts[0]) {
	case "HKEY_LOCAL_MACHINE", "HKLM":
		rootKey = registry.LOCAL_MACHINE
	case "HKEY_CURRENT_USER", "HKCU":
		rootKey = registry.CURRENT_USER
	case "HKEY_USERS", "HKU":
		rootKey = registry.USERS
	case "HKEY_CLASSES_ROOT", "HKCR":
		rootKey = registry.CLASSES_ROOT
	}

	if len(parts) == 1 {
		return rootKey, "", nil
	}

	return rootKey, parts[1], nil
}

func getRegistryValues(rootKey registry.Key, keyPath string, searchKey string) ([]Registry, error) {
	// Open the registry key with permissions to enumerate subkeys
	key, err := registry.OpenKey(rootKey, keyPath, registry.QUERY_VALUE|registry.ENUMERATE_SUB_KEYS)
	if err != nil {
		return nil, fmt.Errorf("failed to open registry key: %v", err)
	}
	defer key.Close()

	var results []Registry
	var reg registry.Key

	// Get and process subkeys
	subkeys, err := key.ReadSubKeyNames(-1)
	if err != nil {
		return nil, fmt.Errorf("failed to read subkeys: %v", err)
	}

	// Process each subkey
	for _, subkeyName := range subkeys {
		subReg := Registry{
			Key:  searchKey,
			Path: searchKey + "\\" + subkeyName,
			Name: subkeyName,
			Type: "subkey",
			Data: "",
		}

		if keyPath == "" {
			reg, err = registry.OpenKey(rootKey, subkeyName, registry.QUERY_VALUE)
		} else {
			reg, err = registry.OpenKey(rootKey, keyPath+"\\"+subkeyName, registry.QUERY_VALUE)
		}
		defer reg.Close()

		if err == nil {
			subkeyInfo, err := reg.Stat()
			if err == nil {
				subReg.MTime = uint64(subkeyInfo.ModTime().Unix())
			}
		}

		results = append(results, subReg)
	}

	// Get and process values
	valueNames, err := key.ReadValueNames(-1)
	if err != nil {
		return nil, fmt.Errorf("failed to read registry values: %v", err)
	}

	// Process each value
	for _, valueName := range valueNames {
		val, valType, err := key.GetStringValue(valueName)
		if err != nil {
			continue // Skip values that can't be read
		}

		entry := Registry{
			Key:  searchKey,
			Path: searchKey + "\\" + valueName,
			Name: valueName,
			Type: registryTypeMap[valType],
			Data: val,
		}
		// Need to be fixed
		if keyPath == "" {
			reg, err = registry.OpenKey(rootKey, valueName, registry.QUERY_VALUE)
		} else {
			reg, err = registry.OpenKey(rootKey, keyPath+"\\"+valueName, registry.QUERY_VALUE)
		}
		defer reg.Close()

		if err == nil {
			info, err := reg.Stat()
			if err == nil {
				entry.MTime = uint64(info.ModTime().Unix())
			}
		}

		results = append(results, entry)
	}

	return results, nil
}

func GenRegistry(searchKey string) ([]Registry, error) {
	rootKey, keyPath, err := parseSearchKey(searchKey)
	if err != nil {
		return nil, fmt.Errorf("failed to parse search key: %v", err)
	}

	return getRegistryValues(rootKey, keyPath, searchKey)
}
