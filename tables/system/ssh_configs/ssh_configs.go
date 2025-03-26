package ssh_configs

import (
	"bufio"
	"os"
	"path/filepath"
	"strings"

	"github.com/scrymastic/goosquery/sql/result"
	"github.com/scrymastic/goosquery/sql/sqlctx"
	"github.com/scrymastic/goosquery/tables/system/users"
)

const (
	WindowsSystemwideSSHConfig = "C:\\ProgramData\\ssh\\ssh_config"
)

func genSshConfig(ctx *sqlctx.Context, uid int64, configFilePath string) (*result.Results, error) {
	sshConfigs := result.NewQueryResult()

	file, err := os.Open(configFilePath)
	if err != nil {
		return sshConfigs, err
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	block := ""

	for scanner.Scan() {
		sshConfig := result.NewResult(ctx, Schema)
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
			sshConfig.Set("uid", uid)
			sshConfig.Set("block", block)
			sshConfig.Set("option", line)
			sshConfig.Set("ssh_config_file", configFilePath)
			sshConfigs.AppendResult(*sshConfig)
		}
	}

	if err := scanner.Err(); err != nil {
		return sshConfigs, err
	}

	return sshConfigs, nil
}

func genSshConfigForUser(ctx *sqlctx.Context, uid int64, directory string) (*result.Results, error) {
	sshConfigFile := filepath.Join(directory, ".ssh", "config")

	// Check if the file exists
	if _, err := os.Stat(sshConfigFile); os.IsNotExist(err) {
		return result.NewQueryResult(), nil
	}

	return genSshConfig(ctx, uid, sshConfigFile)
}

func GenSshConfigs(ctx *sqlctx.Context) (*result.Results, error) {
	results := result.NewQueryResult()

	// Get all users
	users, err := users.GenUsers(ctx)
	if err != nil {
		return nil, err
	}

	// Process each user's SSH config
	for i := 0; i < users.Size(); i++ {
		user := users.GetRow(i)
		if user.Get("directory") != "" {
			configs, err := genSshConfigForUser(ctx, user.Get("uid").(int64), user.Get("directory").(string))
			if err == nil {
				results.AppendResults(*configs)
			}
			// We don't return on error for individual users, just continue
		}
	}

	// Check for system-wide SSH config
	if _, err := os.Stat(WindowsSystemwideSSHConfig); err == nil {
		configs, err := genSshConfig(ctx, 0, WindowsSystemwideSSHConfig)
		if err == nil {
			results.AppendResults(*configs)
		}
	}

	return results, nil
}
