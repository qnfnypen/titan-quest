package dao

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
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
