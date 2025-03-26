package user_groups

import (
	"github.com/scrymastic/goosquery/sql/result"
	"github.com/scrymastic/goosquery/sql/sqlctx"
	"github.com/scrymastic/goosquery/tables/system/users"
)

func GenUserGroups(ctx *sqlctx.Context) (*result.Results, error) {
	userGroups := result.NewQueryResult()

	users, err := users.GenUsers(ctx)
	if err != nil {
		return nil, err
	}

	for i := 0; i < users.Size(); i++ {
		user := users.GetRow(i)
		if user.Get("username") == "LOCAL SERVICE" || user.Get("username") == "SYSTEM" || user.Get("username") == "NETWORK SERVICE" {
			continue
		}

		userGroup := result.NewResult(ctx, Schema)
		userGroup.Set("uid", user.Get("uid"))
		userGroup.Set("gid", user.Get("gid"))

		userGroups.AppendResult(*userGroup)
	}

	return userGroups, nil
}
