package dao

import (
	"context"
	"fmt"

	"github.com/Masterminds/squirrel"
	"github.com/gnasnik/titan-quest/core/generated/model"
)

// GetUserInviteLogs 获取
func GetUserInviteLogs(ctx context.Context, name string, option QueryOption) ([]*model.InviteLogResp, int64, error) {
	var (
		limit, offset int
		total         int64
		out           []*model.InviteLogResp
	)
	inviteLog := model.InviteLog{}

	if option.PageSize <= 0 {
		limit = 50
	} else {
		limit = option.PageSize
	}
	if option.Page > 0 {
		offset = limit * (option.Page - 1)
	}

	// 获取总条数
	query, args, err := squirrel.Select("COUNT(DISTINCT(invited_name))").From(inviteLog.TableName()).Where("username = ?", name).ToSql()
	if err != nil {
		return nil, 0, fmt.Errorf("generate sql error:%w", err)
	}
	err = DB.GetContext(ctx, &total, query, args...)
	if err != nil {
		return nil, 0, fmt.Errorf("get total of invite_log error:%w", err)
	}
	// 获取详情
	query, args, err = squirrel.Select("invited_name,MIN(created_at) AS created_at,COUNT(id) AS tasks,IFNULL(SUM(credit),0) AS credits").
		From(inviteLog.TableName()).Where("username = ?", name).Limit(uint64(limit)).Offset(uint64(offset)).GroupBy("invited_name").OrderBy("created_at DESC").ToSql()
	if err != nil {
		return nil, 0, fmt.Errorf("generate sql error:%w", err)
	}

	err = DB.SelectContext(ctx, &out, query, args...)
	if err != nil {
		return nil, 0, fmt.Errorf("get list of invite_log error:%w", err)
	}

	return out, total, nil
}
