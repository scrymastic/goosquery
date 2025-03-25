package registry

import (
	"fmt"
	"strings"

	"github.com/scrymastic/goosquery/sql/result"
	"github.com/scrymastic/goosquery/sql/sqlctx"
	"golang.org/x/sys/windows"
	"golang.org/x/sys/windows/registry"
)

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
	default:
		return 0, "", fmt.Errorf("invalid registry key: %s", parts[0])
	}

	if len(parts) == 1 {
		return rootKey, "", nil
	}

	return rootKey, parts[1], nil
}

func getRegistryValues(rootKey registry.Key, keyPath string, searchKey string, ctx *sqlctx.Context) (*result.Results, error) {
	// Open the registry key with permissions to enumerate subkeys
	key, err := registry.OpenKey(rootKey, keyPath, registry.QUERY_VALUE|registry.ENUMERATE_SUB_KEYS)
	if err != nil {
		return nil, fmt.Errorf("failed to open registry key: %v", err)
	}
	defer key.Close()

	info, err := key.Stat()
	if err != nil {
		return nil, fmt.Errorf("failed to get registry key info: %v", err)
	}

	modTime := uint64(info.ModTime().Unix())

	regkeys := result.NewQueryResult()
	var reg registry.Key

	// Get and process subkeys
	subkeys, err := key.ReadSubKeyNames(-1)
	if err != nil {
		return nil, fmt.Errorf("failed to read subkeys: %v", err)
	}

	// Process each subkey
	for _, subkeyName := range subkeys {
		subReg := result.NewResult(ctx, Schema)
		subReg.Set("key", searchKey)
		subReg.Set("path", searchKey+"\\"+subkeyName)
		subReg.Set("name", subkeyName)
		subReg.Set("type", "subkey")
		subReg.Set("data", "")

		if keyPath == "" {
			reg, err = registry.OpenKey(rootKey, subkeyName, registry.QUERY_VALUE)
		} else {
			reg, err = registry.OpenKey(rootKey, keyPath+"\\"+subkeyName, registry.QUERY_VALUE)
		}
		defer reg.Close()

		if err == nil {
			subkeyInfo, err := reg.Stat()
			if err == nil {
				subReg.Set("mtime", uint64(subkeyInfo.ModTime().Unix()))
			}
		}

		regkeys.AppendResult(*subReg)
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

		if valueName == "" {
			valueName = "(Default)"
		}

		regType, ok := registryTypeMap[valType]
		if !ok {
			regType = fmt.Sprintf("Unknown Type: %d", valType)
		}

		entry := result.NewResult(ctx, Schema)
		entry.Set("key", searchKey)
		entry.Set("path", searchKey+"\\"+valueName)
		entry.Set("name", valueName)
		entry.Set("type", regType)
		entry.Set("data", val)
		entry.Set("mtime", modTime)

		regkeys.AppendResult(*entry)
	}

	return regkeys, nil
}

func GenRegistry(ctx *sqlctx.Context) (*result.Results, error) {
	searchKey := ctx.GetConstants("key")
	results := result.NewQueryResult()
	rootKey, keyPath, err := parseSearchKey(searchKey[0])
	if err != nil {
		return nil, fmt.Errorf("failed to parse search key: %v", err)
	}
	regkeys, err := getRegistryValues(rootKey, keyPath, searchKey[0], ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to get registry values: %v", err)
	}
	results.AppendResults(*regkeys)

	return results, nil
}
