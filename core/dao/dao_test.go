package dao

import (
	"context"
	"database/sql"
	"testing"
	"time"

	"github.com/gnasnik/titan-quest/core/generated/model"
	"github.com/jmoiron/sqlx"
)

func TestMain(m *testing.M) {
	db, err := sqlx.Connect("mysql", "root:123456@tcp(127.0.0.1:3306)/titan_quest?charset=utf8mb4&parseTime=True&loc=Local")
	if err != nil {
		log.Fatal(err)
	}

	db.SetMaxOpenConns(maxOpenConnections)
	db.SetConnMaxLifetime(connMaxLifetime * time.Second)
	db.SetMaxIdleConns(maxIdleConnections)
	db.SetConnMaxIdleTime(connMaxIdleTime * time.Second)

	DB = db

	m.Run()
}

func TestCreateUserExt(t *testing.T) {
	ctx := context.Background()
	// err := CreateUserExt(ctx, &model.UsersExt{
	// 	Username:    "1",
	// 	InviteCode:  "2",
	// 	InvitedCode: "3",
	// })
	// if err != nil {
	// 	t.Fatal(err)
	// }

	// 判断用户附属表是否存在
	_, err := GetUserExt(ctx, "2")
	switch err {
	case sql.ErrNoRows:
		err = CreateUserExt(ctx, &model.UsersExt{
			Username:    "2",
			InviteCode:  "3",
			InvitedCode: "3",
		})
		if err != nil {
			t.Fatal(err)
		}
	case nil:
	default:
		t.Fatal(err)
	}
}

func TestGetMissionLogs(t *testing.T) {
	list, total, err := GetMissionLogs(context.Background(), "0x30c16b1c6e07b5f685ee668b9e69a28512f74159", QueryOption{})
	if err != nil {
		t.Fatal(err)
	}

	t.Log(list)
	t.Log(total)
}

func TestGetUserInviteLogs(t *testing.T) {
	list, total, err := GetUserInviteLogs(context.Background(), "dingktester3@163.com", QueryOption{Page: 1, PageSize: 10})
	if err != nil {
		t.Fatal(err)
	}

	for _, v := range list {
		t.Log(*v)
	}
	t.Log(total)
}

func TestComplete(t *testing.T) {
	ctx := context.Background()
	var missionID int64 = 1011
	username := "1661628099@qq.com"
	mission, err := GetMissionById(ctx, missionID)
	if err != nil {
		t.Fatal(err)
	}

	ums, err := GetUserMissionByMissionId(ctx, username, mission.ID, QueryOption{})
	if err != nil {
		t.Fatal(err)
	}

	if len(ums) == 0 {
		err := AddUserMissionAndInviteLog(ctx, &model.UserMission{
			Username:  username,
			MissionID: mission.ID,
			Type:      mission.Type,
			Credit:    mission.Credit,
			Content:   username,
			CreatedAt: time.Now(),
		})

		if err != nil {
			t.Fatal(err)
		}
	}
}

func TestAddUserMissionAndInviteLog(t *testing.T) {
	err := AddUserMissionAndInviteLog(context.Background(), &model.UserMission{
		Username:  "2",
		MissionID: 1003,
		Type:      4,
		Credit:    100,
		CreatedAt: time.Now(),
	})
	if err != nil {
		t.Fatal(err)
	}
}
