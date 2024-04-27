package dao

import (
	"context"
	"fmt"
	"github.com/gnasnik/titan-quest/core/generated/model"
	logging "github.com/ipfs/go-log/v2"
)

var log = logging.Logger("dao")

func GetLoginLocation(ctx context.Context, userId string) string {
	query := `select login_location from login_log where login_username = ? order by id desc limit 1`

	var areaId string
	if err := DB.GetContext(ctx, &areaId, query, userId); err != nil {
		log.Errorf("get user area id: %v", err)
		return ""
	}

	return areaId
}

func AddLoginLog(ctx context.Context, log *model.LoginLog) error {
	_, err := DB.NamedExecContext(ctx,
		`INSERT INTO login_log (login_username, ip_address, login_location, browser, os, status, msg, created_at) VALUES 
		(:login_username, :ip_address, :login_location, :browser, :os, :status, :msg, now());`, log)
	return err
}

func ListLoginLog(ctx context.Context, option QueryOption) ([]*model.LoginLog, int64, error) {
	var args []interface{}
	var total int64
	var out []*model.LoginLog

	limit := option.PageSize
	offset := option.Page
	if option.PageSize <= 0 {
		limit = 50
	}
	if option.Page > 0 {
		offset = limit * (option.Page - 1)
	}

	err := DB.GetContext(ctx, &total,
		`SELECT count(*) FROM login_log`, args)
	if err != nil {
		return nil, 0, err
	}

	err = DB.SelectContext(ctx, &out, fmt.Sprintf(
		`SELECT * FROM login_log LIMIT %d OFFSET %d`, limit, offset,
	), args...)
	if err != nil {
		return nil, 0, err
	}

	return out, total, err
}
