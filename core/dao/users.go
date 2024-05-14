package dao

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"github.com/Masterminds/squirrel"
	"github.com/gnasnik/titan-quest/core/generated/model"
)

func CreateUser(ctx context.Context, user *model.User) error {
	_, err := DB.NamedExecContext(ctx, fmt.Sprintf(
		`INSERT INTO users (username, pass_hash, user_email, wallet_address, role, referrer, referrer_user_id, referral_code, created_at)
			VALUES (:username, :pass_hash, :user_email, :wallet_address, :role, :referrer, :referrer_user_id, :referral_code, :created_at);`,
	), user)
	return err
}

func ResetPassword(ctx context.Context, passHash, username string) error {
	_, err := DB.DB.ExecContext(ctx, fmt.Sprintf(
		`UPDATE users SET pass_hash = '%s', updated_at = now() WHERE username = '%s'`, passHash, username))
	return err
}

func GetUserByUsername(ctx context.Context, username string) (*model.User, error) {
	var out model.User
	if err := DB.QueryRowxContext(ctx, fmt.Sprintf(
		`SELECT * FROM users WHERE username = ?`), username,
	).StructScan(&out); err != nil {
		return nil, err
	}

	return &out, nil
}

func GetUserIds(ctx context.Context) ([]string, error) {
	queryStatement := fmt.Sprintf(`SELECT username as user_id FROM users;`)
	var out []string
	err := DB.SelectContext(ctx, &out, queryStatement)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrNoRow
		}
		return nil, err
	}
	return out, nil
}

func UpdateUserKOLReferralCode(ctx context.Context, username, kolReferralCode, kolUserId string) error {
	_, err := DB.DB.ExecContext(ctx, `UPDATE users SET from_kol_ref_code = ?, from_kol_user_id = ?, updated_at = now() WHERE username = ?`, kolReferralCode, kolUserId, username)
	return err
}

// CreateUserInfo 创建用户信息
func CreateUserInfo(ctx context.Context, user *model.User, userExt *model.UsersExt) error {
	tx, err := DB.BeginTx(ctx, nil)
	if err != nil {
		return fmt.Errorf("begin transaction error:%w", err)
	}
	// 创建用户表
	query, args, err := squirrel.Insert(user.TableName()).Columns("username", "created_at", "updated_at", "referral_code").
		Values(user.Username, user.CreatedAt, user.UpdatedAt, user.ReferralCode).ToSql()
	if err != nil {
		return fmt.Errorf("generate insert user sql error:%w", err)
	}
	_, err = tx.ExecContext(ctx, query, args...)
	if err != nil {
		tx.Rollback()
		return err
	}
	// 创建用户附属表
	query, args, err = squirrel.Insert(userExt.TableName()).Columns("username", "invite_code", "invited_code").
		Values(userExt.Username, userExt.InviteCode, userExt.InvitedCode).ToSql()
	if err != nil {
		return fmt.Errorf("generate insert user_ext sql error:%w", err)
	}
	_, err = tx.ExecContext(ctx, query, args...)
	if err != nil {
		tx.Rollback()
		return err
	}

	err = tx.Commit()
	return err
}

// CreateUserExt 创建用户附属表
func CreateUserExt(ctx context.Context, userExt *model.UsersExt) error {
	query, args, err := squirrel.Insert(userExt.TableName()).Columns("username", "invite_code", "invited_code").
		Values(userExt.Username, userExt.InviteCode, userExt.InvitedCode).ToSql()
	if err != nil {
		return fmt.Errorf("generate insert user_ext sql error:%w", err)
	}

	_, err = DB.DB.ExecContext(ctx, query, args...)
	return err
}

// CreateInviteLog 增加邀请收益记录
func CreateInviteLog(ctx context.Context, inviteLog *model.InviteLog) error {
	query, args, err := squirrel.Insert(inviteLog.TableName()).Columns("username", "invited_name", "mission_id", "credit", "created_at").
		Values(inviteLog.Username, inviteLog.InvitedName, inviteLog.MissionID, inviteLog.Credit, inviteLog.CreatedAt).ToSql()
	if err != nil {
		return fmt.Errorf("generate insert invite_log sql error:%w", err)
	}

	_, err = DB.DB.ExecContext(ctx, query, args...)
	return err
}

// GetUserExt 获取用户附属信息
func GetUserExt(ctx context.Context, username string) (*model.UsersExt, error) {
	userExt := model.UsersExt{}

	query, args, err := squirrel.Select("*").From(userExt.TableName()).Where("username = ?", username).ToSql()
	if err != nil {
		return nil, fmt.Errorf("generate query user_ext sql error:%w", err)
	}

	err = DB.QueryRowxContext(ctx, query, args...).StructScan(&userExt)
	if err != nil {
		return nil, err
	}

	return &userExt, nil
}

// GetUserExtByInviteCode 根据邀请码获取用户附属信息
func GetUserExtByInviteCode(ctx context.Context, inviteCode string) (*model.UsersExt, error) {
	userExt := model.UsersExt{}

	query, args, err := squirrel.Select("*").From(userExt.TableName()).Where("invite_code = ?", inviteCode).ToSql()
	if err != nil {
		return nil, fmt.Errorf("generate query user_ext sql error:%w", err)
	}

	err = DB.QueryRowxContext(ctx, query, args...).StructScan(&userExt)
	if err != nil {
		return nil, err
	}

	return &userExt, nil
}

// GetUserResponse 获取用户响应
func GetUserResponse(ctx context.Context, username string) (*model.ResponseUser, error) {
	response := model.ResponseUser{}

	query, args, err := squirrel.Select("users.username AS un", "user_email", "wallet_address", "role", "created_at", "referral_code", "referrer", "from_kol_ref_code", "invite_code").
		From("users").LeftJoin("users_ext ON users.username = users_ext.username").
		Where("users.username = ?", username).ToSql()

	if err != nil {
		return nil, fmt.Errorf("generate query user_info sql error:%w", err)
	}

	err = DB.QueryRowxContext(ctx, query, args...).StructScan(&response)
	if err != nil {
		return nil, err
	}

	return &response, nil
}
