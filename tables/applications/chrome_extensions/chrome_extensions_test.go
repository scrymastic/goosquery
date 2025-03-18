package chrome_extensions

import (
	"encoding/json"
	"testing"
)

func TestGetUserInformationList(t *testing.T) {
	// Test case 1: Basic functionality
	userInfoList, err := getUserInformationList()
	if err != nil {
		t.Errorf("getUserInformationList() returned error: %v", err)
	}

	// We can't predict the exact number of users, but we can check if the list is populated
	if len(userInfoList) == 0 {
		t.Error("getUserInformationList() returned empty list")
	}

	// Check that each user has valid UID and path
	for _, user := range userInfoList {
		if user.Uid == 0 {
			t.Error("getUserInformationList() returned user with invalid UID")
		}
		if user.Path == "" {
			t.Error("getUserInformationList() returned user with empty path")
		}
	}
}

func TestGetChromeProfilePathList(t *testing.T) {
	// Test case 1: Basic functionality
	profilePaths, err := getChromeProfilePathList()
	if err != nil {
		t.Errorf("getChromeProfilePathList() returned error: %v", err)
	}

	// We expect multiple paths (number of users * number of browser types)
	if len(profilePaths) == 0 {
		t.Error("getChromeProfilePathList() returned empty list")
	}

	// Pretty print the profile paths
	json, err := json.MarshalIndent(profilePaths, "", "  ")
	if err != nil {
		t.Errorf("failed to marshal profile paths: %v", err)
	}
	t.Logf("Profile paths: %s", string(json))
}
