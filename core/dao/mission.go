package dao

import (
	"context"
	"database/sql"
	"fmt"
	"strings"
	"time"

	"github.com/Masterminds/squirrel"
	"github.com/gnasnik/titan-quest/config"
	"github.com/gnasnik/titan-quest/core/generated/model"
	"github.com/jmoiron/sqlx"
)

func AddTwitterOAuth(ctx context.Context, oauth *model.TwitterOauth) error {
	query := `insert into twitter_oauth(username, request_token, redirect_uri, created_at) values(:username, :request_token, :redirect_uri, now())`

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
	query := `insert into discord_oauth(username, state, redirect_uri, created_at) values(:username, :state, :redirect_uri, now())`

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

func AddTelegramOAuth(ctx context.Context, oauth *model.TelegramOauth) error {
	query := `insert into telegram_oauth(username, code, redirect_uri, created_at) values(:username, :code, :redirect_uri, now())`

	_, err := DB.NamedExecContext(ctx, query, &oauth)
	return err
}

func GetTelegramOauthByUsername(ctx context.Context, username string) (*model.TelegramOauth, error) {
	query := `select * from telegram_oauth where username  = ? and telegram_user_id <> 0`

	var out model.TelegramOauth
	err := DB.GetContext(ctx, &out, query, username)
	if err != nil {
		return nil, err
	}

	return &out, nil
}

func AddTelegramUserInfo(ctx context.Context, to *model.TelegramOauth) error {
	query := `insert into telegram_oauth(code, username, telegram_user_id, telegram_username, created_at, updated_at) values 
		(:code, :username, :telegram_user_id, :telegram_username, now(), now())`
	_, err := DB.NamedExecContext(ctx, query, to)
	return err
}

func GetMissions(ctx context.Context) ([]*model.Mission, error) {
	query := `select * from mission where status = 1 order by sort_id`

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

func GetMissionById2(ctx context.Context, missionId int64) (*model.Mission, error) {
	query := `select * from mission where id = ? order by id`

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

func GetKOLCommissionCredits(ctx context.Context, kolUserId string) (float64, error) {
	var total float64

	countQuery := `select ifnull(sum(m.credit),0) * 0.2 from users u left join user_mission m on u.username = m.username where from_kol_user_id = ?  group by u.username;`
	countQueryIn, countQueryParams, err := sqlx.In(countQuery, kolUserId)
	if err != nil {
		return 0, err
	}

	err = DB.GetContext(ctx, &total, countQueryIn, countQueryParams...)
	if err != nil {
		return 0, err
	}

	return total, nil
}

func GetUserCreditsByKOLReferralCode(ctx context.Context, kolUserId string, option QueryOption) (int64, []*model.UserCredit, error) {
	limit := option.PageSize
	offset := option.Page
	if option.PageSize <= 0 {
		limit = 50
	}
	if option.Page > 0 {
		offset = limit * (option.Page - 1)
	}

	var total int64

	countQuery := `select count(1) from ( select u.username from users u left join user_mission m on u.username = m.username where from_kol_user_id = ?  group by u.username) d`
	countQueryIn, countQueryParams, err := sqlx.In(countQuery, kolUserId)
	if err != nil {
		return 0, nil, err
	}

	err = DB.GetContext(ctx, &total, countQueryIn, countQueryParams...)
	if err != nil {
		return 0, nil, err
	}

	query := `select * from (
		select u.username, u.from_kol_ref_code , IFNULL(sum(m.credit),0) as credits, count(1) as completed_mission_count, u.created_at from users u left join user_mission m on u.username = m.username where from_kol_user_id = ?  group by u.username
	) d order by created_at desc LIMIT ? OFFSET ?;`

	var out []*model.UserCredit
	err = DB.SelectContext(ctx, &out, query, kolUserId, limit, offset)
	if err != nil {
		return 0, nil, err
	}

	return total, out, nil
}

// AddUserMissionAndInviteLog 增加用户任务完成并且增加邀请的积分记录
func AddUserMissionAndInviteLog(ctx context.Context, um *model.UserMission) error {
	query, args, err := squirrel.Insert(um.TableName()).Columns("username, mission_id, sub_mission_id, type, credit, content, created_at").
		Values(um.Username, um.MissionID, um.SubMissionID, um.Type, um.Credit, um.Content, um.CreatedAt).ToSql()
	if err != nil {
		return fmt.Errorf("generate insert user_mission sql error:%w", err)
	}
	_, err = DB.ExecContext(ctx, query, args...)
	if err != nil {
		return err
	}

	// 查询该用户是否被邀请，有错误不返回
	userExt, err := GetUserExt(ctx, um.Username)
	switch err {
	case sql.ErrNoRows:
	case nil:
		// 存在的话则增加记录
		if strings.TrimSpace(userExt.InvitedCode) != "" {
			addInviteLog(ctx, userExt, um)
		}
	}

	return nil
}

func addInviteLog(ctx context.Context, ue *model.UsersExt, um *model.UserMission) error {
	// 根据邀请人的邀请码获取邀请人信息
	ui, err := GetUserExtByInviteCode(ctx, ue.InvitedCode)
	if err != nil {
		return err
	}

	// 增加记录
	credit := um.Credit * config.Cfg.InviteShareRate / 100
	err = CreateInviteLog(ctx, &model.InviteLog{
		Username:    ui.Username,
		InvitedName: ue.Username,
		MissionID:   um.MissionID,
		Credit:      credit,
		CreatedAt:   time.Now(),
	})

	return err
}

// SumInviteCredits 统计邀请积分的总和
func SumInviteCredits(ctx context.Context, username string) (int64, error) {
	var sum int64
	iv := model.InviteLog{}

	query, args, err := squirrel.Select("ifnull(sum(credit),0)").From(iv.TableName()).Where("username = ?", username).ToSql()
	if err != nil {
		return 0, fmt.Errorf("generate select sum of credits error:%w", err)
	}

	err = DB.QueryRowxContext(ctx, query, args...).Scan(&sum)
	return sum, err
}

// GetMissionLogs 获取任务完成记录
func GetMissionLogs(ctx context.Context, name string, option QueryOption) ([]*model.MissionLogResp, int64, error) {
	var (
		limit, offset int
		total         int64
		out           []*model.MissionLogResp
	)

	um := model.UserMission{}

	if option.PageSize <= 0 {
		limit = 50
	} else {
		limit = option.PageSize
	}
	if option.Page > 0 {
		offset = limit * (option.Page - 1)
	}

	// 获取总条数
	query, args, err := squirrel.Select("COUNT(id)").From(um.TableName()).Where("username = ?", name).ToSql()
	if err != nil {
		return nil, 0, fmt.Errorf("generate sql error:%w", err)
	}
	err = DB.GetContext(ctx, &total, query, args...)
	if err != nil {
		return nil, 0, fmt.Errorf("get total of invite_log error:%w", err)
	}
	// 获取详情
	query, args, err = squirrel.Select("title,title_cn,user_mission.created_at AS createdAt,user_mission.credit AS ucredit").
		From(um.TableName()).LeftJoin("mission ON user_mission.mission_id = mission.id").
		Where("username = ?", name).Limit(uint64(limit)).Offset(uint64(offset)).OrderBy("createdAt DESC").ToSql()
	if err != nil {
		return nil, 0, fmt.Errorf("generate sql error:%w", err)
	}
	err = DB.SelectContext(ctx, &out, query, args...)
	if err != nil {
		return nil, 0, fmt.Errorf("get list of invite_log error:%w", err)
	}

	return out, total, nil
}

func GetCreditsList(ctx context.Context, option QueryOption) (int64, []*model.UserCredit, error) {
	limit := option.PageSize
	offset := option.Page
	if option.PageSize <= 0 {
		limit = 50
	}
	if option.Page > 0 {
		offset = limit * (option.Page - 1)
	}

	var total int64

	countQuery := `select count(1) from ( select t.username from (select username, credit from user_mission union all select username, credit from invite_log) t group by t.username) d`
	countQueryIn, countQueryParams, err := sqlx.In(countQuery)
	if err != nil {
		return 0, nil, err
	}

	err = DB.GetContext(ctx, &total, countQueryIn, countQueryParams...)
	if err != nil {
		return 0, nil, err
	}

	query := `select * from (
		select t.username, IFNULL(sum(t.credit),0) as credits from (select username, credit from user_mission union all select username, credit from invite_log) t group by t.username
	) d order by d.credits desc LIMIT ? OFFSET ?;`

	var out []*model.UserCredit
	err = DB.SelectContext(ctx, &out, query, limit, offset)
	if err != nil {
		return 0, nil, err
	}

	return total, out, nil
}
