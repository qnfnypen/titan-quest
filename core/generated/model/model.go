package model

import (
	"time"
)

type Language string

const (
	LanguageEN = "en"
	LanguageCN = "cn"
)

type errorCode int

type Location struct {
	ID          int64     `db:"id" json:"id"`
	Ip          string    `db:"ip" json:"ip"`
	Continent   string    `db:"continent" json:"continent"`
	Province    string    `db:"province" json:"province"`
	City        string    `db:"city" json:"city"`
	Country     string    `db:"country" json:"country"`
	Latitude    string    `db:"latitude" json:"latitude"`
	Longitude   string    `db:"longitude" json:"longitude"`
	AreaCode    string    `db:"area_code" json:"area_code"`
	Isp         string    `db:"isp" json:"isp"`
	ZipCode     string    `db:"zip_code" json:"zip_code"`
	Elevation   string    `db:"elevation" json:"elevation"`
	CreatedAt   time.Time `db:"created_at" json:"created_at"`
	UpdatedAt   time.Time `db:"updated_at" json:"updated_at"`
	CountryCode string    `db:"-" json:"country_code"`
}
