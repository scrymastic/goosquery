package default_environment

import (
	"fmt"
	"log"

	"golang.org/x/sys/windows/registry"
)

type DefaultEnvironment struct {
	Variable string
	Value    string
	Expand   int32
}

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
func GenDefaultEnvironments() ([]DefaultEnvironment, error) {
	var results []DefaultEnvironment

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

		// Create environment variable entry
		envVar := DefaultEnvironment{
			Variable: name,
			Value:    value,
			Expand:   boolToInt32(valueType == registry.EXPAND_SZ),
		}
		results = append(results, envVar)
	}

	return results, nil
}
