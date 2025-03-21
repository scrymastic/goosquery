package default_environment

import (
	"fmt"
	"log"

	"github.com/scrymastic/goosquery/sql/context"
	"github.com/scrymastic/goosquery/tables/specs"
	"golang.org/x/sys/windows/registry"
)

const (
	regKeyEnvironment = `SYSTEM\CurrentControlSet\Control\Session Manager\Environment`
)

func boolToInt32(b bool) int32 {
	if b {
		return 1
	}
	return 0
}

// GenDefaultEnvironments retrieves system environment variables from the Windows Registry
func GenDefaultEnvironments(ctx context.Context) ([]map[string]interface{}, error) {
	var results []map[string]interface{}

	// Open the Registry key
	key, err := registry.OpenKey(registry.LOCAL_MACHINE,
		regKeyEnvironment,
		registry.QUERY_VALUE,
	)
	if err != nil {
		return nil, fmt.Errorf("error opening registry key: %s, %v", regKeyEnvironment, err)
	}
	defer key.Close()

	// Get all value names under the key
	valueNames, err := key.ReadValueNames(-1)
	if err != nil {
		return nil, fmt.Errorf("error reading value names: %s, %v", regKeyEnvironment, err)
	}

	// Iterate through each value
	for _, name := range valueNames {
		value, valueType, err := key.GetStringValue(name)
		if err != nil {
			log.Printf("Error reading value for %s: %s, %v", name, regKeyEnvironment, err)
			continue
		}

		// Initialize all requested columns with default values
		envVar := specs.Init(ctx, Schema)

		if ctx.IsColumnUsed("variable") {
			envVar["variable"] = name
		}

		if ctx.IsColumnUsed("value") {
			envVar["value"] = value
		}

		if ctx.IsColumnUsed("expand") {
			envVar["expand"] = boolToInt32(valueType == registry.EXPAND_SZ)
		}

		results = append(results, envVar)
	}

	return results, nil
}
