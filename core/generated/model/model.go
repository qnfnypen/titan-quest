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
	Username       string    `db:"un" json:"username"`
	UserEmail      string    `db:"user_email" json:"user_email"`
	WalletAddress  string    `db:"wallet_address" json:"wallet_address"`
	Role           int32     `db:"role" json:"role"`
	CreatedAt      time.Time `db:"created_at" json:"created_at"`
	ReferralCode   string    `db:"referral_code" json:"referral_code"`
	Referrer       string    `db:"referrer" json:"referrer"`
	Credits        int64     `db:"credits" json:"credits"`
	FromKolRefCode string    `db:"from_kol_ref_code" json:"from_kol_ref_code"`
	InviteCode     string    `db:"invite_code" json:"code"` // 邀请码
}

type (
	// InviteLogResp 邀请记录
	InviteLogResp struct {
		InvitedName string    `db:"invited_name" json:"email"`
		Tasks       int64     `db:"tasks" json:"task_nums"`
		CreatedAt   time.Time `db:"created_at" json:"created_at"`
		Credits     int64     `db:"credits" json:"credits"`
	}
	// MissionLogResp 任务完成记录
	MissionLogResp struct {
		Title     string    `db:"title" json:"title"`
		TitleCn   string    `db:"title_cn" json:"title_cn"`
		CreatedAt time.Time `db:"createdAt" json:"created_at"`
		Credit    int64     `db:"ucredit" json:"credit"`
	}
)

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

type UserCredit struct {
	Id                    int64   `json:"id" db:"-"`
	Username              string  `json:"username" db:"username"`
	Credits               float64 `json:"credits" db:"credits"`
	CompletedMissionCount int     `json:"completed_mission_count" db:"completed_mission_count"`
	FromKOLRefCode        string  `json:"from_kol_ref_code" db:"from_kol_ref_code"`
	CreatedAt             string  `json:"created_at" db:"created_at"`
}

// TableName 表名映射
func (InviteLog) TableName() string {
	return "invite_log"
}

// TableName 表名映射
func (UsersExt) TableName() string {
	return "users_ext"
}

// TableName 表明映射
func (User) TableName() string {
	return "users"
}

// TableName 表明映射
func (UserMission) TableName() string {
	return "user_mission"
}

// TableName 表明映射
func (Mission) TableName() string {
	return "mission"
}
