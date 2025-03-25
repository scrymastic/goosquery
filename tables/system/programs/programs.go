package programs

import (
	"fmt"
	"regexp"
	"strings"

	"github.com/scrymastic/goosquery/sql/result"
	"github.com/scrymastic/goosquery/sql/sqlctx"
	"golang.org/x/sys/windows/registry"
)

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

func keyEnumProgram(hive registry.Key, key string, ctx *sqlctx.Context) (*result.Results, error) {
	programs := result.NewQueryResult()

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

		// Initialize program
		prog := result.NewResult(ctx, Schema)

		// Initialize identifyingNumber
		var val string

		if ctx.IsColumnUsed("identifying_number") {
			// Check if subkey is a GUID
			if guidRegex.MatchString(subkey) {
				val = subkey
			} else {
				val, _, _ = subk.GetStringValue("BundleIdentifier")
			}
			prog.Set("identifying_number", val)
		}

		if ctx.IsColumnUsed("name") {
			val, _, _ = subk.GetStringValue("DisplayName")
			prog.Set("name", val)
		}
		if ctx.IsColumnUsed("version") {
			val, _, _ = subk.GetStringValue("DisplayVersion")
			prog.Set("version", val)
		}
		if ctx.IsColumnUsed("install_location") {
			val, _, _ = subk.GetStringValue("InstallLocation")
			prog.Set("install_location", val)
		}
		if ctx.IsColumnUsed("install_source") {
			val, _, _ = subk.GetStringValue("InstallSource")
			prog.Set("install_source", val)
		}
		if ctx.IsColumnUsed("language") {
			val, _, _ = subk.GetStringValue("Language")
			prog.Set("language", val)
		}
		if ctx.IsColumnUsed("publisher") {
			val, _, _ = subk.GetStringValue("Publisher")
			prog.Set("publisher", val)
		}
		if ctx.IsColumnUsed("uninstall_string") {
			val, _, _ = subk.GetStringValue("UninstallString")
			prog.Set("uninstall_string", val)
		}
		if ctx.IsColumnUsed("install_date") {
			val, _, _ = subk.GetStringValue("InstallDate")
			prog.Set("install_date", val)
		}

		// prog.RegKey = key + `\` + subkey

		// Check if the program does not empty
		if prog.Get("name") == "" &&
			prog.Get("version") == "" &&
			prog.Get("install_location") == "" &&
			prog.Get("install_source") == "" &&
			prog.Get("language") == "" &&
			prog.Get("publisher") == "" &&
			prog.Get("uninstall_string") == "" &&
			prog.Get("install_date") == "" &&
			prog.Get("identifying_number") == "" {
			continue
		}

		programs.AppendResult(*prog)
		subk.Close()
	}

	return programs, nil
}

func GenPrograms(ctx *sqlctx.Context) (*result.Results, error) {
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

	programs := result.NewQueryResult()
	for _, key := range localProgramKeys {
		prog, err := keyEnumProgram(registry.LOCAL_MACHINE, key, ctx)
		if err != nil {
			return nil, fmt.Errorf("failed to get programs from %s: %v", key, err)
		}
		programs.AppendResults(*prog)
	}

	for _, key := range userProgramKeys {
		prog, err := keyEnumProgram(registry.USERS, key, ctx)
		if err != nil {
			return nil, fmt.Errorf("failed to get programs from %s: %v", key, err)
		}
		programs.AppendResults(*prog)
	}

	return programs, nil
}
