package ssh_configs

import (
	"bufio"
	"os"
	"path/filepath"
	"strings"

	"github.com/scrymastic/goosquery/system/users"
)

type SSHConfig struct {
	UID           int64  `json:"uid"`
	Block         string `json:"block"`
	Option        string `json:"option"`
	SSHConfigFile string `json:"ssh_config_file"`
}

const (
	WindowsSystemwideSSHConfig = "C:\\ProgramData\\ssh\\ssh_config"
)

func genSSHConfig(uid int64, configFilePath string) ([]SSHConfig, error) {
	var results []SSHConfig

	file, err := os.Open(configFilePath)
	if err != nil {
		return results, err
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	block := ""

	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())

		// Skip empty lines and comments
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}

		// Check for host or match blocks
		lineToLower := strings.ToLower(line)
		if strings.HasPrefix(lineToLower, "host ") || strings.HasPrefix(lineToLower, "match ") {
			block = line
		} else {
			// Add the option to results
			results = append(results, SSHConfig{
				UID:           uid,
				Block:         block,
				Option:        line,
				SSHConfigFile: configFilePath,
			})
		}
	}

	if err := scanner.Err(); err != nil {
		return results, err
	}

	return results, nil
}

func genSSHConfigForUser(uid int64, directory string) ([]SSHConfig, error) {
	sshConfigFile := filepath.Join(directory, ".ssh", "config")

	// Check if the file exists
	if _, err := os.Stat(sshConfigFile); os.IsNotExist(err) {
		return []SSHConfig{}, nil
	}

	return genSSHConfig(uid, sshConfigFile)
}

func GenSSHConfigs() ([]SSHConfig, error) {
	var allConfigs []SSHConfig

	// Get all users
	users, err := users.GenUsers()
	if err != nil {
		return nil, err
	}

	// Process each user's SSH config
	for _, user := range users {
		if user.Directory != "" {
			configs, err := genSSHConfigForUser(user.UID, user.Directory)
			if err == nil {
				allConfigs = append(allConfigs, configs...)
			}
			// We don't return on error for individual users, just continue
		}
	}

	// Check for system-wide SSH config
	if _, err := os.Stat(WindowsSystemwideSSHConfig); err == nil {
		configs, err := genSSHConfig(0, WindowsSystemwideSSHConfig)
		if err == nil {
			allConfigs = append(allConfigs, configs...)
		}
	}

	return allConfigs, nil
}
