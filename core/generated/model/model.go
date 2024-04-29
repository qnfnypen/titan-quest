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

type ResponseUser struct {
	Username       string    `db:"username" json:"username"`
	UserEmail      string    `db:"user_email" json:"user_email"`
	WalletAddress  string    `db:"wallet_address" json:"wallet_address"`
	Role           int32     `db:"role" json:"role"`
	CreatedAt      time.Time `db:"created_at" json:"created_at"`
	ReferralCode   string    `db:"referral_code" json:"referral_code"`
	Referrer       string    `db:"referrer" json:"referrer"`
	Credits        int64     `db:"credits" json:"credits"`
	FromKolRefCode string    `db:"from_kol_ref_code" json:"from_kol_ref_code"`
}

func (u *User) ToResponseUser() *ResponseUser {
	return &ResponseUser{
		Username:       u.Username,
		UserEmail:      u.UserEmail,
		WalletAddress:  u.WalletAddress,
		Role:           u.Role,
		CreatedAt:      u.CreatedAt,
		ReferralCode:   u.ReferralCode,
		Referrer:       u.Referrer,
		Credits:        u.Credits,
		FromKolRefCode: u.FromKolRefCode,
	}
}

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
