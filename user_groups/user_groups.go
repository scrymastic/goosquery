package user_groups

import "osquery/users"

type UserGroup struct {
	UID int64 `json:"uid"`
	GID int64 `json:"gid"`
}

func GenUserGroups() ([]UserGroup, error) {
	userGroups := []UserGroup{}

	users, err := users.GenUsers()
	if err != nil {
		return nil, err
	}

	for _, user := range users {
		if user.Username == "LOCAL SERVICE" || user.Username == "SYSTEM" || user.Username == "NETWORK SERVICE" {
			continue
		}

		userGroup := UserGroup{
			UID: user.UID,
			GID: user.GID,
		}

		userGroups = append(userGroups, userGroup)
	}

	return userGroups, nil
}
