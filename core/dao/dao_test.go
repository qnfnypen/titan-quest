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
