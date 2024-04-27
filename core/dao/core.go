package dao

import (
	"context"
	"fmt"
	"github.com/gnasnik/titan-quest/config"
	"github.com/gnasnik/titan-quest/core/generated/model"
	"github.com/go-redis/redis/v9"
	_ "github.com/go-sql-driver/mysql"
	"github.com/jmoiron/sqlx"
	"time"
)

var (
	// DB reference to database

	DB *sqlx.DB

	// RedisCache  redis caching instance
	RedisCache *redis.Client
)

const (
	maxOpenConnections = 60
	connMaxLifetime    = 120
	maxIdleConnections = 30
	connMaxIdleTime    = 20
)

var ErrNoRow = fmt.Errorf("no matching row found")

func Init(cfg *config.Config) error {
	if cfg.DatabaseURL == "" {
		return fmt.Errorf("database url not setup")
	}

	db, err := sqlx.Connect("mysql", cfg.DatabaseURL)
	if err != nil {
		return err
	}

	db.SetMaxOpenConns(maxOpenConnections)
	db.SetConnMaxLifetime(connMaxLifetime * time.Second)
	db.SetMaxIdleConns(maxIdleConnections)
	db.SetConnMaxIdleTime(connMaxIdleTime * time.Second)

	client := redis.NewClient(&redis.Options{
		Addr:     cfg.RedisAddr,
		Password: cfg.RedisPassword,
		PoolSize: 100,
	})
	_, err = client.Ping(context.Background()).Result()
	if err != nil {
		return err
	}

	DB = db
	RedisCache = client
	return nil
}

type QueryOption struct {
	Page       int            `json:"page"`
	PageSize   int            `json:"page_size"`
	Order      string         `json:"order"`
	OrderField string         `json:"order_field"`
	StartTime  string         `json:"start_time"`
	EndTime    string         `json:"end_time" `
	UserID     string         `json:"user_id"`
	Content    string         `json:"content"`
	Lang       model.Language `json:"lang"`
}
