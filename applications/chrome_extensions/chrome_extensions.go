package chrome_extensions

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/scrymastic/goosquery/system/users"
)

type ChromeExtension struct {
	BrowserType             string `json:"browser_type"`
	Uid                     int64  `json:"uid"`
	Name                    string `json:"name"`
	Profile                 string `json:"profile"`
	ProfilePath             string `json:"profile_path"`
	ReferencedIdentifier    string `json:"referenced_identifier"`
	Identifier              string `json:"identifier"`
	Version                 string `json:"version"`
	Description             string `json:"description"`
	DefaultLocale           string `json:"default_locale"`
	CurrentLocale           string `json:"current_locale"`
	UpdateUrl               string `json:"update_url"`
	Author                  string `json:"author"`
	Persistent              int32  `json:"persistent"`
	Path                    string `json:"path"`
	Permissions             string `json:"permissions"`
	PermissionsJson         string `json:"permissions_json"`
	OptionalPermissions     string `json:"optional_permissions"`
	OptionalPermissionsJson string `json:"optional_permissions_json"`
	ManifestHash            string `json:"manifest_hash"`
	Referenced              int64  `json:"referenced"`
	FromWebstore            string `json:"from_webstore"`
	State                   string `json:"state"`
	InstallTime             string `json:"install_time"`
	InstallTimestamp        int64  `json:"install_timestamp"`
	ManifestJson            string `json:"manifest_json"`
	Key                     string `json:"key"`
}

const (
	GoogleChrome int32 = iota
	GoogleChromeBeta
	GoogleChromeDev
	GoogleChromeCanary
	Brave
	Chromium
	Yandex
	Opera
	Edge
	EdgeBeta
	Vivaldi
	Arc
)

const (
	kLocalizedMessagePrefix = "__MSG_"
	kPreferencesFile        = "Preferences"
	kSecurePreferencesFile  = "Secure Preferences"
	kExtensionManifestFile  = "manifest.json"
	kExtensionsFolderName   = "Extensions"
)

var WindowsChromePathSuffixMap = map[int32]string{
	GoogleChrome:       "AppData\\Local\\Google\\Chrome\\User Data",
	GoogleChromeBeta:   "AppData\\Local\\Google\\Chrome Beta\\User Data",
	GoogleChromeDev:    "AppData\\Local\\Google\\Chrome Dev\\User Data",
	GoogleChromeCanary: "AppData\\Local\\Google\\Chrome SxS\\User Data",
	Brave:              "AppData\\Local\\BraveSoftware\\Brave-Browser\\User Data",
	Chromium:           "AppData\\Local\\Chromium\\User Data",
	Yandex:             "AppData\\Local\\Yandex\\YandexBrowser\\User Data",
	Edge:               "AppData\\Local\\Microsoft\\Edge\\User Data",
	EdgeBeta:           "AppData\\Local\\Microsoft\\Edge Beta\\User Data",
	Opera:              "AppData\\Roaming\\Opera Software\\Opera Stable",
	Vivaldi:            "AppData\\Local\\Vivaldi\\User Data",
}

type UserInformation struct {
	Uid  int64  `json:"uid"`
	Path string `json:"path"`
}

type ChromeProfilePath struct {
	Type  int32
	Value string
	Uid   int64
}

// Extension represents a single extension found inside the profile
type Extension struct {
	// The absolute path to the extension folder
	Path string

	// The contents of the manifest file
	Manifest string
}

// ChromeProfileSnapshot represents a snapshot of a Chrome profile
type ChromeProfileSnapshot struct {
	// Profile type
	Type int32

	// Absolute path to this profile
	Path string

	// The contents of the 'Preferences' file
	Preferences string

	// The contents of the 'Secure Preferences' file
	SecurePreferences string

	// The user id
	Uid int64

	// A map of all the extensions discovered in the preferences
	ReferencedExtensions map[string]Extension

	// A map of all the extensions that are not present in the preferences
	UnreferencedExtensions map[string]Extension
}

func getUserInformationList() ([]UserInformation, error) {
	userInfoList := []UserInformation{}

	users, err := users.GenUsers()
	if err != nil {
		return nil, fmt.Errorf("failed to generate users: %w", err)
	}

	for _, user := range users {
		if user.UID == 0 || user.Directory == "" {
			continue
		}

		userInfoList = append(userInfoList, UserInformation{
			Uid:  user.UID,
			Path: user.Directory,
		})
	}

	return userInfoList, nil
}

func isValidChromeProfile(path string) bool {
	// Valid chrome profile contains 'Preferences' and 'Secure Preferences'
	preferencesPath := filepath.Join(path, kPreferencesFile)
	securePreferencesPath := filepath.Join(path, kSecurePreferencesFile)

	// Check if both files exist
	if _, err := os.Stat(preferencesPath); os.IsNotExist(err) {
		return false
	}
	if _, err := os.Stat(securePreferencesPath); os.IsNotExist(err) {
		return false
	}

	return true
}

func getChromeProfilePathList() ([]ChromeProfilePath, error) {
	userInfoList, err := getUserInformationList()
	if err != nil {
		return nil, fmt.Errorf("failed to get user information list: %w", err)
	}

	chromeProfilePathList := []ChromeProfilePath{}

	for _, userInfo := range userInfoList {
		for browserType, pathSuffix := range WindowsChromePathSuffixMap {
			// Check if the user data path exists
			userDataPath := filepath.Join(userInfo.Path, pathSuffix)
			if _, err := os.Stat(userDataPath); os.IsNotExist(err) {
				continue
			}

			// List all subdirectories in the user data path and check if they are valid chrome profiles
			subdirectories, err := os.ReadDir(userDataPath)
			if err != nil {
				continue
			}

			for _, subdirectory := range subdirectories {
				if !subdirectory.IsDir() {
					continue
				}

				subdirectoryPath := filepath.Join(userDataPath, subdirectory.Name())
				if !isValidChromeProfile(subdirectoryPath) {
					continue
				}

				// Add the chrome profile path to the list
				chromeProfilePathList = append(chromeProfilePathList, ChromeProfilePath{
					Type:  browserType,
					Value: subdirectoryPath,
					Uid:   userInfo.Uid,
				})
			}
		}
	}

	return chromeProfilePathList, nil
}

func captureProfileSnapshotSettingsFromPath(chromeProfileSnapshot *ChromeProfileSnapshot, chromeProfilePath ChromeProfilePath) error {
	chromeProfileSnapshot.Type = chromeProfilePath.Type
	chromeProfileSnapshot.Path = chromeProfilePath.Value
	chromeProfileSnapshot.Uid = chromeProfilePath.Uid

	preferencesPath := filepath.Join(chromeProfilePath.Value, kPreferencesFile)
	securePreferencesPath := filepath.Join(chromeProfilePath.Value, kSecurePreferencesFile)

	preferences, err := os.ReadFile(preferencesPath)
	if err != nil {
		return fmt.Errorf("failed to read preferences: %w", err)
	}

	chromeProfileSnapshot.Preferences = string(preferences)

	securePreferences, err := os.ReadFile(securePreferencesPath)
	if err != nil {
		return fmt.Errorf("failed to read secure preferences: %w", err)
	}

	chromeProfileSnapshot.SecurePreferences = string(securePreferences)

	return nil
}

func captureProfileSnapshotExtensionsFromPath(chromeProfileSnapshot *ChromeProfileSnapshot, chromeProfilePath ChromeProfilePath) error {
	extensionsFolderPath := filepath.Join(chromeProfilePath.Value, kExtensionsFolderName)
	if _, err := os.Stat(extensionsFolderPath); os.IsNotExist(err) {
		return nil
	}

	extensions, err := os.ReadDir(extensionsFolderPath)
	if err != nil {
		return fmt.Errorf("failed to read extensions: %w", err)
	}

	for _, extension := range extensions {
		if !extension.IsDir() {
			continue
		}
	}

	return nil
}

func getChromeProfileSnapshotList() ([]ChromeProfileSnapshot, error) {
	chromeProfilePathList, err := getChromeProfilePathList()
	if err != nil {
		return nil, fmt.Errorf("failed to get chrome profile path list: %w", err)
	}

	chromeProfileSnapshotList := []ChromeProfileSnapshot{}

	for _, chromeProfilePath := range chromeProfilePathList {
		chromeProfileSnapshot := ChromeProfileSnapshot{}
		err := captureProfileSnapshotSettingsFromPath(&chromeProfileSnapshot, chromeProfilePath)
		if err != nil {
			return nil, fmt.Errorf("failed to capture profile snapshot settings from path: %w", err)
		}

		chromeProfileSnapshotList = append(chromeProfileSnapshotList, chromeProfileSnapshot)
	}

	return chromeProfileSnapshotList, nil
}
