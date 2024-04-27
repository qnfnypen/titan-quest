package api

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/gnasnik/titan-quest/config"
	"github.com/gnasnik/titan-quest/core/bot/discord"
	logging "github.com/ipfs/go-log/v2"
	tele "gopkg.in/telebot.v3"
	"time"
)

var (
	Bot *tele.Bot
)

var log = logging.Logger("api")

func InitBot() {
	pref := tele.Settings{
		Token:  config.Cfg.TelegramBotToken,
		Poller: &tele.LongPoller{Timeout: 10 * time.Second},
	}

	//groupId := int64(-1002050663753)

	b, err := tele.NewBot(pref)
	if err != nil {
		log.Fatalf("create telegram cmd: %v", err)
		return
	}

	if !config.Cfg.DisableDiscordBot {
		go discord.RunDiscordBot(config.Cfg.DiscordBotToken)
	}

	fmt.Println("telegram cmd id:", b.Me.ID)

	Bot = b

}

func ServerAPI(cfg *config.Config) {
	gin.SetMode(cfg.Mode)
	r := gin.Default()
	r.Use(Cors())
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
	apiV1.GET("/google/callback", GoogleCallBackHandler)
	apiV1.GET("/telegram/callback", TelegramCallback)

	user := apiV1.Group("/user")
	user.GET("/login_before", GetNonceStringHandler)
	user.POST("/login", authMiddleware.LoginHandler)
	user.GET("/verify_code", GetNumericVerifyCodeHandler)
	user.POST("/logout", authMiddleware.LogoutHandler)
	user.GET("/refresh_token", authMiddleware.RefreshHandler)

	user.Use(authMiddleware.MiddlewareFunc())
	user.POST("/info", GetUserInfoHandler)
	user.GET("/twitter/auth", TwitterOAuthHandler)
	user.GET("/discord/auth", DiscordOAuthHandler)
	user.GET("/google/auth", GoogleOAuthHandler)
	user.GET("/telegram/auth", TelegramOAuthHandler)

	quest := apiV1.Group("/quest")
	//quest.POST("/x/follow", XFollowHandler)
	quest.GET("/query_missions", QueryMissionHandler)
	quest.Use(authMiddleware.MiddlewareFunc())
	quest.GET("/query_user_credits", QueryUserCreditsHandler)
	quest.GET("/check", CheckQuestHandler)
	quest.POST("/twitter_link", PostTwitterLinkHandler)

	if err := r.Run(cfg.ApiListen); err != nil {
		log.Fatalf("starting server: %v\n", err)
	}
}
