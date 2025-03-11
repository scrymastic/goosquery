// Package collector provides a common implementation for all collector functions
package collector

import (
	"fmt"
	"reflect"
	"runtime"
	"strings"
	"time"
)

// Wrapper provides a common implementation for all collector functions
type Wrapper struct {
	Name        string
	Description string
	Category    string
	Function    interface{}
	RequiresArg bool   // Indicates if the function requires arguments
	ArgType     string // Type of argument required (e.g., "string", "uint32")
	ArgDesc     string // Description of the argument
}

// Registry stores all registered collectors
var Registry = make(map[string]Wrapper)

// Register registers a collector function in the registry
func Register(name, description, category string, fn interface{}, requiresArg bool, argType, argDesc string) {
	Registry[name] = Wrapper{
		Name:        name,
		Description: description,
		Category:    category,
		Function:    fn,
		RequiresArg: requiresArg,
		ArgType:     argType,
		ArgDesc:     argDesc,
	}
}

// RegisterSimple is a backward compatibility function for the old Register function
func RegisterSimple(name, description, category string, fn interface{}) {
	Register(name, description, category, fn, false, "", "")
}

// GetByCategory returns all collectors in a specific category
func GetByCategory(category string) []Wrapper {
	var collectors []Wrapper
	for _, collector := range Registry {
		if collector.Category == category {
			collectors = append(collectors, collector)
		}
	}
	return collectors
}

// GetAll returns all registered collectors
func GetAll() []Wrapper {
	var collectors []Wrapper
	for _, collector := range Registry {
		collectors = append(collectors, collector)
	}
	return collectors
}

// GetFunctionName returns the name of a function
func GetFunctionName(fn interface{}) string {
	if fn == nil {
		return "nil"
	}

	// Get the function name using reflection
	funcValue := reflect.ValueOf(fn)
	if funcValue.Kind() != reflect.Func {
		return "not a function"
	}

	funcName := runtime.FuncForPC(funcValue.Pointer()).Name()

	// Extract just the function name without the package path
	parts := strings.Split(funcName, ".")
	return parts[len(parts)-1]
}

// WrapGenFunction creates a wrapper function that logs execution time and handles errors
func WrapGenFunction(fn interface{}) interface{} {
	// Get the function type
	fnType := reflect.TypeOf(fn)

	// Check if it's a function
	if fnType.Kind() != reflect.Func {
		panic("Not a function")
	}

	// Get the function name
	funcName := GetFunctionName(fn)

	// Create a wrapper function with the same signature
	wrapper := reflect.MakeFunc(fnType, func(args []reflect.Value) []reflect.Value {
		// Log start time
		startTime := time.Now()
		fmt.Printf("Starting %s...\n", funcName)

		// Call the original function
		results := reflect.ValueOf(fn).Call(args)

		// Log end time and duration
		duration := time.Since(startTime)

		// Check if there was an error (assuming the last return value is an error)
		var errMsg string
		if len(results) > 1 && !results[len(results)-1].IsNil() {
			err := results[len(results)-1].Interface().(error)
			errMsg = fmt.Sprintf(" (Error: %v)", err)
		}

		fmt.Printf("Completed %s in %v%s\n", funcName, duration, errMsg)

		return results
	})

	return wrapper.Interface()
}

// Initialize initializes all collectors from the registry
func Initialize() {
	// This function will be called at startup to register all collectors
	// Each package should have an init() function that calls Register
	fmt.Println("Initializing collectors...")
	fmt.Printf("Registered %d collectors\n", len(Registry))
}

// ExecuteWithArg executes a collector function that requires an argument
func ExecuteWithArg(name string, arg interface{}) (interface{}, error) {
	collector, exists := Registry[name]
	if !exists {
		return nil, fmt.Errorf("collector %s not found", name)
	}

	if !collector.RequiresArg {
		return nil, fmt.Errorf("collector %s does not require an argument", name)
	}

	// Get the function type
	fnType := reflect.TypeOf(collector.Function)
	if fnType.Kind() != reflect.Func {
		return nil, fmt.Errorf("collector %s is not a function", name)
	}

	// Check if the function has at least one input parameter
	if fnType.NumIn() < 1 {
		return nil, fmt.Errorf("collector %s does not accept arguments", name)
	}

	// Check if we have multiple argument types
	argTypes := strings.Split(collector.ArgType, ",")
	argValues := []reflect.Value{}

	// Handle single argument case
	if len(argTypes) == 1 {
		// Convert the argument to the expected type
		var argValue reflect.Value
		switch collector.ArgType {
		case "string":
			if strArg, ok := arg.(string); ok {
				argValue = reflect.ValueOf(strArg)
			} else {
				return nil, fmt.Errorf("collector %s requires a string argument", name)
			}
		case "uint32":
			switch v := arg.(type) {
			case uint32:
				argValue = reflect.ValueOf(v)
			case int:
				argValue = reflect.ValueOf(uint32(v))
			case float64:
				argValue = reflect.ValueOf(uint32(v))
			default:
				return nil, fmt.Errorf("collector %s requires a uint32 argument", name)
			}
		default:
			return nil, fmt.Errorf("unsupported argument type %s for collector %s", collector.ArgType, name)
		}
		argValues = append(argValues, argValue)
	} else {
		// Handle multiple arguments case
		// For now, we only support string arguments in the multiple arguments case
		// and expect the arg to be a slice or array
		argSlice, ok := arg.([]interface{})
		if !ok {
			return nil, fmt.Errorf("collector %s requires multiple arguments as a slice", name)
		}

		if len(argSlice) != len(argTypes) {
			return nil, fmt.Errorf("collector %s requires %d arguments, got %d", name, len(argTypes), len(argSlice))
		}

		for i, argType := range argTypes {
			switch argType {
			case "string":
				if strArg, ok := argSlice[i].(string); ok {
					argValues = append(argValues, reflect.ValueOf(strArg))
				} else {
					return nil, fmt.Errorf("collector %s requires a string argument at position %d", name, i)
				}
			case "uint32":
				switch v := argSlice[i].(type) {
				case uint32:
					argValues = append(argValues, reflect.ValueOf(v))
				case int:
					argValues = append(argValues, reflect.ValueOf(uint32(v)))
				case float64:
					argValues = append(argValues, reflect.ValueOf(uint32(v)))
				default:
					return nil, fmt.Errorf("collector %s requires a uint32 argument at position %d", name, i)
				}
			default:
				return nil, fmt.Errorf("unsupported argument type %s for collector %s at position %d", argType, name, i)
			}
		}
	}

	// Call the function with the arguments
	results := reflect.ValueOf(collector.Function).Call(argValues)

	// Check if there was an error
	if len(results) > 1 && !results[1].IsNil() {
		err := results[1].Interface().(error)
		return nil, err
	}

	// Return the result
	if len(results) > 0 {
		return results[0].Interface(), nil
	}

	return nil, nil
}

// Execute executes a collector function that doesn't require an argument
func Execute(name string) (interface{}, error) {
	collector, exists := Registry[name]
	if !exists {
		return nil, fmt.Errorf("collector %s not found", name)
	}

	if collector.RequiresArg {
		return nil, fmt.Errorf("collector %s requires an argument, use ExecuteWithArg instead", name)
	}

	// Get the function type
	fnType := reflect.TypeOf(collector.Function)
	if fnType.Kind() != reflect.Func {
		return nil, fmt.Errorf("collector %s is not a function", name)
	}

	// Call the function
	results := reflect.ValueOf(collector.Function).Call([]reflect.Value{})

	// Check if there was an error
	if len(results) > 1 && !results[1].IsNil() {
		err := results[1].Interface().(error)
		return nil, err
	}

	// Return the result
	if len(results) > 0 {
		return results[0].Interface(), nil
	}

	return nil, nil
}

// RequiresArgument checks if a collector requires an argument
func RequiresArgument(name string) (bool, string, string, error) {
	collector, exists := Registry[name]
	if !exists {
		return false, "", "", fmt.Errorf("collector %s not found", name)
	}

	// If we have multiple argument types, split them
	argTypes := strings.Split(collector.ArgType, ",")
	argDescs := strings.Split(collector.ArgDesc, ",")

	// Make sure we have the same number of descriptions as types
	if len(argTypes) > 1 && len(argTypes) != len(argDescs) {
		return collector.RequiresArg, collector.ArgType, collector.ArgDesc,
			fmt.Errorf("mismatched argument types and descriptions for %s", name)
	}

	return collector.RequiresArg, collector.ArgType, collector.ArgDesc, nil
}

// ExecuteWithMultipleArgs executes a collector function that requires multiple arguments
func ExecuteWithMultipleArgs(name string, args []interface{}) (interface{}, error) {
	return ExecuteWithArg(name, args)
}
