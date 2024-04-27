package dao

import (
	"context"
	"github.com/gnasnik/titan-quest/core/generated/model"
)

func AddTwitterOAuth(ctx context.Context, oauth *model.TwitterOauth) error {
	query := `insert into twitter_oauth(username, request_token, redirect_uri) values(:username, :request_token, :redirect_uri)`

	_, err := DB.NamedExecContext(ctx, query, &oauth)
	return err
}

func GetTwitterOauthByUsername(ctx context.Context, username string) (*model.TwitterOauth, error) {
	query := `select * from twitter_oauth where username  = ? and twitter_user_id <> ''`

	var out model.TwitterOauth
	err := DB.GetContext(ctx, &out, query, username)
	if err != nil {

		return nil, err
	}

	return &out, nil
}

func GetTwitterOauth(ctx context.Context, twitterUserId string) (*model.TwitterOauth, error) {
	query := `select * from twitter_oauth where twitter_user_id  = ? and username <> ''`

	var out model.TwitterOauth
	err := DB.GetContext(ctx, &out, query, twitterUserId)
	if err != nil {

		return nil, err
	}

	return &out, nil
}

func AddDiscordOAuth(ctx context.Context, oauth *model.DiscordOauth) error {
	query := `insert into discord_oauth(username, state, redirect_uri) values(:username, :state, :redirect_uri)`

	_, err := DB.NamedExecContext(ctx, query, &oauth)
	return err
}

func GetDiscordOAuthByUsername(ctx context.Context, username string) (*model.DiscordOauth, error) {
	query := `select * from discord_oauth where username  = ? and discord_user_id <> ''`

	var out model.DiscordOauth
	err := DB.GetContext(ctx, &out, query, username)
	if err != nil {

		return nil, err
	}

	return &out, nil
}

func GetDiscordOAuthByState(ctx context.Context, state string) (*model.DiscordOauth, error) {
	query := `select * from discord_oauth where state  = ? order by created_at desc`

	var out model.DiscordOauth
	err := DB.GetContext(ctx, &out, query, state)
	if err != nil {

		return nil, err
	}

	return &out, nil
}

func GetDiscordOAuth(ctx context.Context, discordUserId string) (*model.DiscordOauth, error) {
	query := `select * from discord_oauth where discord_user_id  = ? and username <> ''`

	var out model.DiscordOauth
	err := DB.GetContext(ctx, &out, query, discordUserId)
	if err != nil {

		return nil, err
	}

	return &out, nil
}

func UpdateDiscordUserInfo(ctx context.Context, state string, discordUserId, email string) error {
	query := `update discord_oauth set discord_user_id = ?, email =? where state = ?`
	_, err := DB.ExecContext(ctx, query, discordUserId, email, state)
	return err
}

func UpdateTwitterUserInfo(ctx context.Context, token string, twitterUserId, twitterScreenName string) error {
	query := `update twitter_oauth set twitter_user_id = ?, twitter_screen_name =? where request_token = ?`
	_, err := DB.ExecContext(ctx, query, twitterUserId, twitterScreenName, token)
	return err
}

func GetMissions(ctx context.Context) ([]*model.Mission, error) {
	query := `select * from mission where status = 1 order by id`

	var out []*model.Mission
	err := DB.SelectContext(ctx, &out, query)
	if err != nil {
		return nil, err
	}

	return out, nil
}

func GetSubMissions(ctx context.Context, parentId int64) ([]*model.Mission, error) {
	query := `select * from sub_mission where status = 1 and parent_id = ? order by id`

	var out []*model.Mission
	err := DB.SelectContext(ctx, &out, query, parentId)
	if err != nil {
		return nil, err
	}

	return out, nil
}

func GetMissionById(ctx context.Context, missionId int64) (*model.Mission, error) {
	query := `select * from mission where status = 1 and id = ? order by id`

	var out model.Mission
	err := DB.GetContext(ctx, &out, query, missionId)
	if err != nil {
		return nil, err
	}

	return &out, nil
}

func GetUserMissionByMissionId(ctx context.Context, username string, missionId int64, opt QueryOption) ([]*model.UserMission, error) {
	args := []interface{}{username, missionId}
	var where = ` where username = ? and mission_id = ?`

	if opt.StartTime != "" {
		where += ` and created_at >= ?`
		args = append(args, opt.StartTime)
	}

	if opt.Content != "" {
		where += ` and content = ?`
		args = append(args, opt.Content)
	}

	query := `select * from user_mission` + where

	var out []*model.UserMission
	err := DB.SelectContext(ctx, &out, query, args...)
	if err != nil {
		return nil, err
	}

	return out, nil
}

func GetUserMissionByUser(ctx context.Context, username string, t int32, opt QueryOption) ([]*model.UserMission, error) {
	args := []interface{}{username, t}
	var where = ` where username = ? and type = ?`

	if opt.StartTime != "" {
		where += ` and created_at >= ?`
		args = append(args, opt.StartTime)
	}

	query := `select * from user_mission` + where

	var out []*model.UserMission
	err := DB.SelectContext(ctx, &out, query, args...)
	if err != nil {
		return nil, err
	}

	return out, nil
}

func SumUserCredits(ctx context.Context, username string) (int64, error) {
	query := `select ifnull(sum(credit),0) from user_mission where username = ?`

	var credits int64
	err := DB.GetContext(ctx, &credits, query, username)
	if err != nil {
		return 0, err
	}

	return credits, nil
}

func AddUserMission(ctx context.Context, um *model.UserMission) error {
	query := `insert into user_mission(username, mission_id, sub_mission_id, type, credit, content, created_at) values(:username, :mission_id, :sub_mission_id, :type, :credit, :content, :created_at)`

	_, err := DB.NamedExecContext(ctx, query, &um)
	return err
}

func AddUserTwitterLink(ctx context.Context, link *model.UserTwitterLink) error {
	query := `insert into user_twitter_link(username, mission_id, link, created_at) values(:username, :mission_id, :link, :created_at)`

	_, err := DB.NamedExecContext(ctx, query, link)
	return err
}

func GetUserTwitterLink(ctx context.Context, username string, missionId int64, startTime string) (*model.UserTwitterLink, error) {
	query := `select * from user_twitter_link where username = ? and mission_id = ? and created_at >= ? order by created_at desc;`

	var out model.UserTwitterLink
	err := DB.GetContext(ctx, &out, query, username, missionId, startTime)
	if err != nil {
		return nil, err
	}

	return &out, nil
}
