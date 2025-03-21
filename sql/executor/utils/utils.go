package utils

import (
	"fmt"
	"regexp"
	"strconv"
	"strings"
)

// ToFloat64 attempts to convert a value to float64
func ToFloat64(v interface{}) (float64, bool) {
	switch val := v.(type) {
	case int:
		return float64(val), true
	case int8:
		return float64(val), true
	case int16:
		return float64(val), true
	case int32:
		return float64(val), true
	case int64:
		return float64(val), true
	case uint:
		return float64(val), true
	case uint8:
		return float64(val), true
	case uint16:
		return float64(val), true
	case uint32:
		return float64(val), true
	case uint64:
		return float64(val), true
	case float32:
		return float64(val), true
	case float64:
		return val, true
	case string:
		// Try to convert string to number
		if f, err := strconv.ParseFloat(val, 64); err == nil {
			return f, true
		}
	}
	return 0, false
}

// Compare compares two values, returns:
//
//	 1 if a > b
//	 0 if a == b
//	-1 if a < b
func Compare(a, b interface{}) int {
	// Handle nil values
	if a == nil && b == nil {
		return 0
	}
	if a == nil {
		return -1
	}
	if b == nil {
		return 1
	}

	// Handle numeric comparison
	aFloat, aIsNum := ToFloat64(a)
	bFloat, bIsNum := ToFloat64(b)
	if aIsNum && bIsNum {
		if aFloat < bFloat {
			return -1
		}
		if aFloat > bFloat {
			return 1
		}
		return 0
	}

	// String comparison for non-numeric values
	aStr := fmt.Sprintf("%v", a)
	bStr := fmt.Sprintf("%v", b)
	return strings.Compare(aStr, bStr)
}

// MatchesLikePattern checks if a value matches a SQL LIKE pattern
func MatchesLikePattern(a, pattern interface{}) bool {
	if a == nil || pattern == nil {
		return false
	}

	aStr := fmt.Sprintf("%v", a)
	patternStr := fmt.Sprintf("%v", pattern)

	// Convert SQL LIKE pattern to Go regex
	// % becomes .*, _ becomes .
	regexPattern := strings.ReplaceAll(patternStr, "%", ".*")
	regexPattern = strings.ReplaceAll(regexPattern, "_", ".")

	// Escape regex metacharacters in the non-wildcard parts
	regexPattern = strings.ReplaceAll(regexPattern, "[", "\\[")
	regexPattern = strings.ReplaceAll(regexPattern, "]", "\\]")
	regexPattern = strings.ReplaceAll(regexPattern, "(", "\\(")
	regexPattern = strings.ReplaceAll(regexPattern, ")", "\\)")

	matched, err := regexp.MatchString("^"+regexPattern+"$", aStr)
	if err != nil {
		return false
	}
	return matched
}

// Contains checks if a string is in a slice of strings
func Contains(slice []string, s string) bool {
	for _, item := range slice {
		if item == s {
			return true
		}
	}
	return false
}

// CaseInsensitiveGet gets a value from a map, trying case-insensitive match if exact match fails
func CaseInsensitiveGet(m map[string]interface{}, key string) (interface{}, bool) {
	// Try exact match first
	if val, exists := m[key]; exists {
		return val, true
	}

	// Try case-insensitive matching
	for k, v := range m {
		if strings.EqualFold(k, key) {
			return v, true
		}
	}

	return nil, false
}
