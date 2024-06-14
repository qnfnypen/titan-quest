package api

import (
	"github.com/bwmarrin/discordgo"
	"github.com/gnasnik/titan-quest/core/bot/discord"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/gnasnik/titan-quest/config"
	logging "github.com/ipfs/go-log/v2"
	tele "gopkg.in/telebot.v3"
)

var (
	TeleBot *tele.Bot
	DCBot   *discordgo.Session
)

var log = logging.Logger("api")

func InitBot() {
	pref := tele.Settings{
		Token:  config.Cfg.TelegramBotToken,
		Poller: &tele.LongPoller{Timeout: 10 * time.Second},
	}

	b, err := tele.NewBot(pref)
	if err != nil {
		log.Errorf("create telegram cmd: %v", err)
		return
	}

	dbot, err := discord.NewBot(config.Cfg.DiscordBotToken)
	if err != nil {
		log.Errorf("create discord bot: %v", err)
	}

	if !config.Cfg.DisableDiscordBot {
		go discord.RunDiscordBot(dbot)
	}

	TeleBot = b
	DCBot = dbot

}

func ServerAPI(cfg *config.Config) {
	gin.SetMode(cfg.Mode)
	r := gin.Default()
	// r.Use(Cors())
	r.Use(RequestLoggerMiddleware())

	apiV1 := r.Group("/api/v1")
	authMiddleware, err := jwtGinMiddleware(cfg.SecretKey)
	if err != nil {
		log.Fatalf("jwt auth middleware: %v", err)
	}

	err = authMiddleware.MiddlewareInit()
	if err != nil {
		log.Fatalf("authMiddleware.MiddlewareInit: %v", err)
	}

	apiV1.GET("/twitter/callback", TwitterCallBackHandler)
	apiV1.GET("/discord/callback", DiscordCallBackHandler)
	apiV1.GET("/telegram/callback", TelegramCallback)
	apiV1.POST("/brows_official_website/callback", BrowsOfficialWebsiteCallback)

	apiV1.GET("/kol_referral_list", GetUserCreditsHandler)
	apiV1.GET("/credits/list", creditsListHandler)

	user := apiV1.Group("/user")
	user.GET("/login_before", GetNonceStringHandler)
	user.POST("/login", authMiddleware.LoginHandler)
	user.GET("/verify_code", GetNumericVerifyCodeHandler)
	user.POST("/logout", authMiddleware.LogoutHandler)
	user.GET("/refresh_token", authMiddleware.RefreshHandler)

	user.Use(authMiddleware.MiddlewareFunc())
	user.GET("/info", GetUserInfoHandler)
	user.GET("/twitter/auth", TwitterOAuthHandler)
	user.GET("/discord/auth", DiscordOAuthHandler)
	user.GET("/telegram/auth", TelegramOAuthHandler)
	user.POST("/wallet/bind", BindWalletHandler)

	quest := apiV1.Group("/quest")
	quest.GET("/query_missions", QueryMissionHandler)
	quest.Use(authMiddleware.MiddlewareFunc())
	quest.GET("/query_user_credits", QueryUserCreditsHandler)
	quest.GET("/check", CheckQuestHandler)
	quest.POST("/twitter_link", PostTwitterLinkHandler)
	quest.POST("/kol_referral_code", BindingKOLReferralCodeHandler)
	quest.GET("/official_website/brows", BrowsOfficialWebsite)
	quest.GET("/official_website/verify", VerifyBrowsOfficialWebsite)
	quest.GET("/become_volunteer", GetBecomeVolunteerURL)
	quest.GET("/become_volunteer/verify", VerifyBecomeVolunteer)

	quest.GET("/invite/logs", GetInviteLogs)
	quest.GET("/mission/logs", GetMissionLogs)

	if err := r.Run(cfg.ApiListen); err != nil {
		log.Fatalf("starting server: %v\n", err)
	}
}
