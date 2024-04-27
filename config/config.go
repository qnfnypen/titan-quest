package config

var Cfg Config

type Config struct {
	Mode                    string
	ApiListen               string
	DatabaseURL             string
	SecretKey               string
	RedisAddr               string
	RedisPassword           string
	Emails                  []EmailConfig
	UToolAPIKeys            []string
	TwitterAPIKey           string
	TwitterAPIKeySecret     string
	DiscordClientId         string
	DiscordClientSecret     string
	OfficialTwitterUserId   int64
	OfficialTelegramGroupId int64
	DiscordBotToken         string
	TelegramBotToken        string
	TelegramBotID           string
	TelegramCallback        string
	RedirectURI             string
	DisableDiscordBot       bool
}

type EmailConfig struct {
	From     string
	Nickname string
	SMTPHost string
	SMTPPort string
	Username string
	Password string
}
