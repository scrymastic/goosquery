package programs

import (
	"fmt"
	"regexp"
	"strings"

	"golang.org/x/sys/windows/registry"
)

type Program struct {
	Name              string `json:"name"`
	Version           string `json:"version"`
	InstallLocation   string `json:"install_location"`
	InstallSource     string `json:"install_source"`
	Language          string `json:"language"`
	Publisher         string `json:"publisher"`
	UninstallString   string `json:"uninstall_string"`
	InstallDate       string `json:"install_date"`
	IdentifyingNumber string `json:"identifying_number"`
	// RegKey            string `json:"reg_key"`
}

func genUserSIDs() ([]string, error) {
	k, err := registry.OpenKey(registry.USERS, "", registry.READ)
	if err != nil {
		return nil, fmt.Errorf("failed to open HKEY_USERS: %v", err)
	}
	defer k.Close()

	info, err := k.Stat()
	if err != nil {
		return nil, err
	}

	sids, err := k.ReadSubKeyNames(int(info.SubKeyCount))
	if err != nil {
		return nil, err
	}

	// Filter out special SIDs and classes
	var validSIDs []string
	for _, sid := range sids {
		// Skip .DEFAULT, _Classes, and other special keys
		if !strings.HasPrefix(sid, "S-1-5-") || strings.HasSuffix(sid, "_Classes") {
			continue
		}
		validSIDs = append(validSIDs, sid)
	}

	return validSIDs, nil
}

func keyEnumProgram(hive registry.Key, key string) ([]Program, error) {
	var programs []Program

	k, err := registry.OpenKey(hive, key, registry.READ)
	if err != nil {
		return nil, fmt.Errorf("failed to open registry key: %v", err)
	}
	defer k.Close()

	keyInfo, err := k.Stat()
	if err != nil {
		return nil, fmt.Errorf("failed to get key info: %v", err)
	}

	// Get subkey names
	subkeys, err := k.ReadSubKeyNames(int(keyInfo.SubKeyCount))
	if err != nil {
		return nil, fmt.Errorf("failed to read subkey names: %v", err)
	}

	guidRegex := regexp.MustCompile(`{[a-fA-F0-9]+-[a-fA-F0-9]+-[a-fA-F0-9]+-[a-fA-F0-9]+-[a-fA-F0-9]+}$`)
	for _, subkey := range subkeys {
		// Open the subkey
		subk, err := registry.OpenKey(hive, key+`\`+subkey, registry.READ)
		if err != nil {
			subk.Close()
			return nil, fmt.Errorf("failed to open subkey: %v", err)
		}

		// Initialize identifyingNumber
		var identifyingNumber string

		// Check if subkey is a GUID
		if guidRegex.MatchString(subkey) {
			identifyingNumber = subkey
		} else {
			identifyingNumber, _, _ = subk.GetStringValue("BundleIdentifier")
		}

		// Initialize program
		var prog Program
		prog.Name, _, _ = subk.GetStringValue("DisplayName")
		prog.Version, _, _ = subk.GetStringValue("DisplayVersion")
		prog.InstallLocation, _, _ = subk.GetStringValue("InstallLocation")
		prog.InstallSource, _, _ = subk.GetStringValue("InstallSource")
		prog.Language, _, _ = subk.GetStringValue("Language")
		prog.Publisher, _, _ = subk.GetStringValue("Publisher")
		prog.UninstallString, _, _ = subk.GetStringValue("UninstallString")
		prog.InstallDate, _, _ = subk.GetStringValue("InstallDate")
		prog.IdentifyingNumber = identifyingNumber
		// prog.RegKey = key + `\` + subkey

		// Check if the program does not empty
		if prog.Name == "" &&
			prog.Version == "" &&
			prog.InstallLocation == "" &&
			prog.InstallSource == "" &&
			prog.Language == "" &&
			prog.Publisher == "" &&
			prog.UninstallString == "" &&
			prog.InstallDate == "" &&
			prog.IdentifyingNumber == "" {
			continue
		}

		programs = append(programs, prog)
		subk.Close()
	}

	return programs, nil
}

func GenPrograms() ([]Program, error) {
	var localProgramKeys = []string{
		`SOFTWARE\Microsoft\Windows\CurrentVersion\Uninstall`,
		`SOFTWARE\WOW6432Node\Microsoft\Windows\CurrentVersion\Uninstall`,
	}

	userProgramKey := `%s\Software\Microsoft\Windows\CurrentVersion\Uninstall`

	userSIDs, err := genUserSIDs()
	if err != nil {
		return nil, fmt.Errorf("failed to get user SIDs: %v", err)
	}

	var userProgramKeys []string
	for _, userSID := range userSIDs {
		userProgramKeys = append(userProgramKeys, fmt.Sprintf(userProgramKey, userSID))
	}

	programs := make([]Program, 0, len(localProgramKeys)+len(userProgramKeys))
	for _, key := range localProgramKeys {
		prog, err := keyEnumProgram(registry.LOCAL_MACHINE, key)
		if err != nil {
			return nil, fmt.Errorf("failed to get programs from %s: %v", key, err)
		}
		programs = append(programs, prog...)
	}

	for _, key := range userProgramKeys {
		prog, err := keyEnumProgram(registry.USERS, key)
		if err != nil {
			return nil, fmt.Errorf("failed to get programs from %s: %v", key, err)
		}
		programs = append(programs, prog...)
	}

	return programs, nil
}
