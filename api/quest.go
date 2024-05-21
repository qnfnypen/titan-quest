package api

import (
	"context"
	"crypto/hmac"
	"crypto/sha256"
	"database/sql"
	"encoding/base64"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"time"

	"github.com/antihax/optional"
	jwt "github.com/appleboy/gin-jwt/v2"
	"github.com/bwmarrin/discordgo"
	"github.com/gin-gonic/gin"
	"github.com/gnasnik/titan-quest/config"
	"github.com/gnasnik/titan-quest/core/dao"
	errorsx "github.com/gnasnik/titan-quest/core/errors"
	"github.com/gnasnik/titan-quest/core/generated/model"
	"github.com/gnasnik/titan-quest/core/opcrypt"
	swagger "github.com/gnasnik/titan-quest/go-client-generated"
	"github.com/gnasnik/titan-quest/pkg/random"
	"github.com/golang-module/carbon/v2"
	"github.com/mrjones/oauth"
	"github.com/valyala/fastjson"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"

	glog "log"
)

func TwitterOAuthHandler(c *gin.Context) {
	claims := jwt.ExtractClaims(c)
	username := claims[identityKey].(string)

	redirectURI := c.Query("redirect_uri")
	if redirectURI == "" {
		redirectURI = config.Cfg.RedirectURI
	}

	consumer := oauth.NewConsumer(
		config.Cfg.TwitterAPIKey,
		config.Cfg.TwitterAPIKeySecret,
		oauth.ServiceProvider{
			RequestTokenUrl:   "https://api.twitter.com/oauth/request_token",
			AuthorizeTokenUrl: "https://api.twitter.com/oauth/authenticate",
			AccessTokenUrl:    "https://api.twitter.com/oauth/access_token",
		})

	// Step 1. Obtain request token for the authorization popup.
	requestToken, loginUrl, err := consumer.GetRequestTokenAndUrl(redirectURI)
	if err != nil {
		log.Errorf("GetRequestTokenAndUrl: %s", err)
		c.JSON(http.StatusOK, respErrorCode(errorsx.InternalServer, c))
		return
	}

	err = dao.AddTwitterOAuth(c.Request.Context(), &model.TwitterOauth{Username: username, RequestToken: requestToken.Token, RedirectUri: redirectURI, CreatedAt: time.Now()})
	if err != nil {
		log.Errorf("AddTwitterOAuth: %v", err)
		c.JSON(http.StatusOK, respErrorCode(errorsx.InternalServer, c))
		return
	}

	// Step 2. Redirect to the authorization screen.
	c.JSON(http.StatusOK, respJSON(JsonObject{
		"url": loginUrl,
	}))
}

func DiscordOAuthHandler(c *gin.Context) {
	claims := jwt.ExtractClaims(c)
	username := claims[identityKey].(string)

	redirectURI := c.Query("redirect_uri")
	if redirectURI == "" {
		redirectURI = config.Cfg.RedirectURI
	}

	oauthProvider := &oauth2.Config{
		ClientID:     config.Cfg.DiscordClientId,
		ClientSecret: config.Cfg.DiscordClientSecret,
		RedirectURL:  redirectURI,
		Endpoint: oauth2.Endpoint{
			AuthURL:  "https://discord.com/oauth2/authorize",
			TokenURL: "https://discord.com/api/oauth2/token",
		},
		Scopes: []string{"identify", "email"},
	}

	state := random.GenerateRandomString(12)
	err := dao.AddDiscordOAuth(c.Request.Context(), &model.DiscordOauth{Username: username, State: state, RedirectUri: redirectURI, CreatedAt: time.Now()})
	if err != nil {
		log.Errorf("AddTwitterOAuth: %v", err)
		c.JSON(http.StatusOK, respErrorCode(errorsx.InternalServer, c))
		return
	}

	// Step 2. Redirect to the authorization screen.
	c.JSON(http.StatusOK, respJSON(JsonObject{
		"url": oauthProvider.AuthCodeURL(state),
	}))
}

func GoogleOAuthHandler(c *gin.Context) {
	//claims := jwt.ExtractClaims(c)
	//username := claims[identityKey].(string)
	//
	//redirectURI := c.Query("redirect_uri")
	//if redirectURI == "" {
	//	redirectURI = config.Cfg.RedirectURI
	//}
	//
	//oauthProvider := &oauth2.Config{
	//	ClientID:     config.Cfg.DiscordClientId,
	//	ClientSecret: config.Cfg.DiscordClientSecret,
	//	RedirectURL:  redirectURI,
	//	Endpoint:     google.Endpoint,
	//}

	//state := random.GenerateRandomString(12)
	//err := dao.AddDiscordOAuth(c.Request.Context(), &model.DiscordOauth{Username: username, State: state})
	//if err != nil {
	//	log.Errorf("AddTwitterOAuth: %v", err)
	//	c.JSON(http.StatusOK, respErrorCode(errorsx.InternalServer, c))
	//	return
	//}

	// Step 2. Redirect to the authorization screen.
	c.JSON(http.StatusOK, respJSON(nil))
}

func TelegramOAuthHandler(c *gin.Context) {

	redirectURI := c.Query("redirect_uri")
	if redirectURI == "" {
		redirectURI = config.Cfg.RedirectURI
	}

	//https://oauth.telegram.org/embed/titannettest_bot?origin=http%3A%2F%2F127.0.0.1%3A5500&return_to=http%3A%2F%2F127.0.0.1%3A5500%2Findex.html&size=large&request_access=write
	//// https://oauth.telegram.org/auth?bot_id=531675494&origin=https%3A%2F%2Ftelegram.org&embed=1&request_access=write&return_to=https%3A%2F%2Ftelegram.org%2Fblog%2Flogin%23
	origin, _ := url.Parse(redirectURI)

	endpoint := "https://oauth.telegram.org/auth"
	telegramAuthURL, _ := url.Parse(endpoint)

	code := random.GenerateRandomString(12)

	values := telegramAuthURL.Query()
	values.Add("bot_id", config.Cfg.TelegramBotID)
	values.Add("origin", origin.Scheme+"://"+origin.Host)
	values.Add("embed", "0")
	//values.Add("request_access", "write")
	values.Add("return_to", redirectURI+"?code="+code)

	telegramAuthURL.RawQuery = values.Encode()

	c.JSON(http.StatusOK, respJSON(JsonObject{
		"url": telegramAuthURL.String(),
	}))
}

func DiscordCallBackHandler(c *gin.Context) {
	code := c.Query("code")
	state := c.Query("state")

	if code == "" || state == "" {
		c.JSON(http.StatusOK, respErrorCode(errorsx.InvalidParams, c))
		return
	}

	da, err := dao.GetDiscordOAuthByState(c.Request.Context(), state)
	if err != nil {
		log.Errorf("GetDiscordOAuthByState: %v", err)
		c.JSON(http.StatusOK, respErrorCode(errorsx.InternalServer, c))
		return
	}

	if da.RedirectUri == "" {
		da.RedirectUri = config.Cfg.RedirectURI
	}

	oauthProvider := &oauth2.Config{
		ClientID:     config.Cfg.DiscordClientId,
		ClientSecret: config.Cfg.DiscordClientSecret,
		RedirectURL:  da.RedirectUri,
		Endpoint: oauth2.Endpoint{
			AuthURL:  "https://discord.com/oauth2/authorize",
			TokenURL: "https://discord.com/api/oauth2/token",
		},
		Scopes: []string{"identify", "email"},
	}

	tokens, err := oauthProvider.Exchange(c.Request.Context(), code)
	if err != nil {
		log.Errorf("Exchange: %v", err)
		c.JSON(http.StatusOK, respErrorCode(errorsx.InternalServer, c))
		return
	}

	ts, _ := discordgo.New("Bearer " + tokens.AccessToken)

	// Retrive the user data.
	u, err := ts.User("@me")
	if err != nil {
		log.Errorf("Get User: %v", err)
		c.JSON(http.StatusOK, respErrorCode(errorsx.InternalServer, c))
	}

	dau, err := dao.GetDiscordOAuth(c.Request.Context(), u.ID)
	if dau != nil && dau.Username != "" {
		c.JSON(http.StatusOK, respErrorCode(errorsx.SocialMediaAccountIsAlreadyInUse, c))
		return
	}

	err = dao.UpdateDiscordUserInfo(c.Request.Context(), state, u.ID, u.Email)
	if err != nil {
		log.Errorf("UpdateDiscordUserInfo: %v", err)
		c.JSON(http.StatusOK, respErrorCode(errorsx.InternalServer, c))
		return
	}

	c.JSON(http.StatusOK, respJSON(nil))
}

func GoogleCallBackHandler(c *gin.Context) {
	code := c.Query("code")
	state := c.Query("state")

	if code == "" || state == "" {
		c.JSON(http.StatusOK, respErrorCode(errorsx.InvalidParams, c))
		return
	}

	oauthProvider := &oauth2.Config{
		ClientID:     config.Cfg.DiscordClientId,
		ClientSecret: config.Cfg.DiscordClientSecret,
		RedirectURL:  config.Cfg.RedirectURI,
		Endpoint:     google.Endpoint,
	}

	tokens, err := oauthProvider.Exchange(c.Request.Context(), code)
	if err != nil {
		log.Errorf("Exchange: %v", err)
		c.JSON(http.StatusOK, respErrorCode(errorsx.InternalServer, c))
		return
	}

	fmt.Println("==>", tokens)

	c.JSON(http.StatusOK, respJSON(nil))
}

func TwitterCallBackHandler(c *gin.Context) {
	oauthToken := c.Query("oauth_token")
	oauthVerify := c.Query("oauth_verifier")

	if oauthToken == "" || oauthVerify == "" {
		c.JSON(http.StatusOK, respErrorCode(errorsx.InvalidParams, c))
		return
	}

	consumer := oauth.NewConsumer(
		config.Cfg.TwitterAPIKey,
		config.Cfg.TwitterAPIKeySecret,
		oauth.ServiceProvider{
			RequestTokenUrl:   "https://api.twitter.com/oauth/request_token",
			AuthorizeTokenUrl: "https://api.twitter.com/oauth/authenticate",
			AccessTokenUrl:    "https://api.twitter.com/oauth/access_token",
		})

	// Part 2/2: Second request after Authorize app is clicked.
	requestToken := &oauth.RequestToken{oauthToken, oauthVerify}

	// Step 3. Exchange oauth token and oauth verifier for access token.
	accessToken, err := consumer.AuthorizeToken(requestToken, oauthVerify)
	if err != nil {
		log.Errorf("AuthorizeToken: %v", err)
		c.JSON(http.StatusOK, respErrorCode(errorsx.InternalServer, c))
		return
	}

	userId := accessToken.AdditionalData["user_id"]
	screenName := accessToken.AdditionalData["screen_name"]

	ta, err := dao.GetTwitterOauth(c.Request.Context(), userId)
	if ta != nil && ta.Username != "" {
		c.JSON(http.StatusOK, respErrorCode(errorsx.SocialMediaAccountIsAlreadyInUse, c))
		return
	}

	err = dao.UpdateTwitterUserInfo(c.Request.Context(), oauthToken, userId, screenName)
	if err != nil {
		log.Errorf("UpdateTwitterUserInfo: %v", err)
		c.JSON(http.StatusOK, respErrorCode(errorsx.InternalServer, c))
		return
	}

	c.JSON(http.StatusOK, respJSON(nil))
}

func TelegramCallback(c *gin.Context) {
	type User struct {
		Id        int64  `json:"id"`
		FirstName string `json:"first_name"`
		LastName  string `json:"last_name"`
		Username  string `json:"username"`
	}

	type Payload struct {
		User User `json:"user"`
	}

	hash := c.Query("hash")
	payloadB64 := c.Query("payload")

	fmt.Println("code", c.Query("code"))
	fmt.Println("hash", hash)
	fmt.Println("pay", payloadB64)

	payloadBytes, err := base64.StdEncoding.DecodeString(payloadB64)
	if err != nil {
		c.JSON(http.StatusOK, respErrorCode(errorsx.InvalidParams, c))
		return
	}

	var payload Payload
	err = json.Unmarshal(payloadBytes, &payload)
	if err != nil {
		c.JSON(http.StatusOK, respErrorCode(errorsx.InvalidParams, c))
		return
	}

	h := hmac.New(sha256.New, []byte(config.Cfg.TelegramBotToken))
	h.Write([]byte(payloadB64))
	checkHash := hex.EncodeToString(h.Sum(nil))

	if hash != checkHash {
		c.JSON(http.StatusOK, respErrorCode(errorsx.InvalidParams, c))
		return
	}

	user := payload.User
	userId := user.Id
	firstName := user.FirstName
	lastName := user.LastName
	username := user.Username

	fmt.Println("userId", userId)
	fmt.Println("firstName", firstName)
	fmt.Println("lastName", lastName)
	fmt.Println("username", username)
}

type RespMission struct {
	*model.Mission
	SubMission []*model.Mission `json:"sub_mission"`
}

func QueryMissionHandler(c *gin.Context) {
	missions, err := dao.GetMissions(c.Request.Context())
	if err != nil {
		log.Errorf("GetMissions: %v", err)
		c.JSON(http.StatusOK, respErrorCode(errorsx.InternalServer, c))
		return
	}

	var (
		basicMission  []*RespMission
		dailyMission  []*RespMission
		weeklyMission []*RespMission
	)

	for _, mission := range missions {
		subMission, err := dao.GetSubMissions(c.Request.Context(), mission.ID)
		if err != nil {
			log.Errorf("GetSubMissions: %v", err)
			c.JSON(http.StatusOK, respErrorCode(errorsx.InternalServer, c))
			return
		}

		mi := &RespMission{
			Mission:    mission,
			SubMission: subMission,
		}

		// 处理浏览官网跳转

		switch mission.Type {
		case MissionTypeBasic:
			basicMission = append(basicMission, mi)
		case MissionTypeDaily:
			dailyMission = append(dailyMission, mi)
		case MissionTypeWeekly:
			weeklyMission = append(weeklyMission, mi)
		}
	}

	c.JSON(http.StatusOK, respJSON(JsonObject{
		"basic_mission":  basicMission,
		"daily_mission":  dailyMission,
		"weekly_mission": weeklyMission,
	}))
}

func QueryUserCreditsHandler(c *gin.Context) {
	claims := jwt.ExtractClaims(c)
	username := claims[identityKey].(string)

	// 获取已完成的基础任务
	completeBasicMission, err := dao.GetUserMissionByUser(c.Request.Context(), username, 1, dao.QueryOption{})
	if err != nil && !errors.Is(err, sql.ErrNoRows) {
		log.Errorf("GetUserMissionByUser: %v", err)
	}

	var basicMission []*model.UserMission
	for _, bs := range completeBasicMission {
		if bs.MissionID == MissionIdLikeTwitter || bs.MissionID == MissionIdRetweet {
			mission, err := dao.GetMissionById(c.Request.Context(), bs.MissionID)
			if err != nil {
				log.Errorf("GetMissionById: %v", err)
				c.JSON(http.StatusOK, respErrorCode(errorsx.InternalServer, c))
				return
			}

			if mission.OpenUrl != bs.Content {
				continue
			}
		}

		basicMission = append(basicMission, bs)
	}

	// 获取已完成的每日任务
	dailyMission, err := dao.GetUserMissionByUser(c.Request.Context(), username, 2, dao.QueryOption{StartTime: carbon.Now().StartOfDay().String()})
	if err != nil && !errors.Is(err, sql.ErrNoRows) {
		log.Errorf("GetUserMissionByUser: %v", err)
	}

	// 获取已完成的每周任务
	weeklyMission, err := dao.GetUserMissionByUser(c.Request.Context(), username, 3, dao.QueryOption{StartTime: carbon.Now().StartOfWeek().String()})
	if err != nil && !errors.Is(err, sql.ErrNoRows) {
		log.Errorf("GetUserMissionByUser: %v", err)
	}

	credits, err := dao.SumUserCredits(c.Request.Context(), username)
	if err != nil && !errors.Is(err, sql.ErrNoRows) {
		log.Errorf("SumUserCredits: %v", err)
	}
	icredits, err := dao.SumInviteCredits(c.Request.Context(), username)
	if err != nil && !errors.Is(err, sql.ErrNoRows) {
		log.Errorf("SumUserCredits: %v", err)
	}

	var (
		twitterUserId string
		discordUserId string
	)
	twitterUser, err := dao.GetTwitterOauthByUsername(c.Request.Context(), username)
	if twitterUser != nil {
		twitterUserId = twitterUser.TwitterUserID
	}

	discordUser, err := dao.GetDiscordOAuthByUsername(c.Request.Context(), username)
	if discordUser != nil {
		discordUserId = discordUser.DiscordUserID
	}

	c.JSON(http.StatusOK, respJSON(JsonObject{
		"address":         username,
		"credits":         credits + icredits,
		"today_credits":   0,
		"week_credits":    0,
		"twitter_user_id": twitterUserId,
		"discord_user_id": discordUserId,
		"missions": JsonObject{
			"basic_mission":  basicMission,
			"daily_mission":  dailyMission,
			"weekly_mission": weeklyMission,
		},
	}))
}

func CheckQuestHandler(c *gin.Context) {
	claims := jwt.ExtractClaims(c)
	username := claims[identityKey].(string)
	missionId, _ := strconv.ParseInt(c.Query("mission_id"), 10, 64)

	mission, err := dao.GetMissionById(c.Request.Context(), missionId)
	if err != nil {
		log.Errorf("GetMissionById: %v", err)
		c.JSON(http.StatusOK, respErrorCode(errorsx.InternalServer, c))
		return
	}

	if mission == nil {
		log.Errorf("mission not found")
		c.JSON(http.StatusOK, respErrorCode(errorsx.InternalServer, c))
		return
	}

	option := dao.QueryOption{}
	if missionId == MissionIdRetweet || missionId == MissionIdLikeTwitter {
		option.Content = mission.OpenUrl
	}

	if mission.Type == MissionTypeDaily {
		option.StartTime = carbon.Now().StartOfDay().String()
	}

	if mission.Type == MissionTypeWeekly {
		option.StartTime = carbon.Now().StartOfWeek().String()
	}

	ums, err := dao.GetUserMissionByMissionId(c.Request.Context(), username, missionId, option)
	if err != nil {
		log.Errorf("GetUserMissionByMissionId: %v", err)
		c.JSON(http.StatusOK, respErrorCode(errorsx.InternalServer, c))
		return
	}

	expectedCount := 1

	if missionId == MissionIdInviteFriendsToDiscord {
		subMissions, err := dao.GetSubMissions(c.Request.Context(), missionId)
		if err != nil {
			log.Errorf("GetUserMissionByMissionId: %v", err)
			c.JSON(http.StatusOK, respErrorCode(errorsx.InternalServer, c))
			return
		}
		expectedCount = len(subMissions)
	}

	if len(ums) >= expectedCount {
		c.JSON(http.StatusOK, respJSON(JsonObject{
			"missions": ums,
		}))
		return
	}

	switch missionId {
	case MissionIdFollowTwitter:
		err = checkFollowTwitter(c.Request.Context(), mission, username, option)
	case MissionIdRetweet:
		err = checkReTweet(c.Request.Context(), mission, username, option)
	case MissionIdLikeTwitter:
		err = checkLikeTweet(c.Request.Context(), mission, username, option)
	case MissionIdJoinDiscord:
		err = checkJoinDiscord(c.Request.Context(), mission, username, option)
	case MissionIdBindingKOL:
		err = checkBindingKOL(c.Request.Context(), mission, username, option)
	case MissionIdVisitOfficialWebsite:
		err = checkVisitOfficialWebsite(c.Request.Context(), mission, username, option)
	case MissionIdVisitReferrerPage:
		err = checkVisitReferrerPage(c.Request.Context(), mission, username, option)
	case MissionIdJoinTelegram:
		err = checkJoinTelegram(c.Request.Context(), mission, username, option)
	case MissionIdQuoteTweet:
		err = checkQuoteTweet(c.Request.Context(), mission, username, option)
	case MissionIdPostTweet:
		err = checkPostTweet(c.Request.Context(), mission, username, option)
	case MissionIdInviteFriendsToDiscord:
		err = checkInviteFriendsToDiscord(c.Request.Context(), mission, username, option)
	default:
		c.JSON(http.StatusOK, respErrorCode(errorsx.NoImplement, c))
		return
	}

	if err != nil {
		log.Errorf("check mission: %v", err)
		c.JSON(http.StatusOK, respErrorCode(errorsx.MissionUnComplete, c))
		return
	}

	ums, err = dao.GetUserMissionByMissionId(c.Request.Context(), username, missionId, option)
	if err != nil {
		log.Errorf("GetUserMissionByMissionId: %v", err)
		return
	}

	if len(ums) > 0 {
		c.JSON(http.StatusOK, respJSON(JsonObject{
			"missions": ums,
		}))
		return
	}

	c.JSON(http.StatusOK, respJSON(nil))
}

var globalCounter int

func GetUToolKeyByRoundRobin() string {
	keys := config.Cfg.UToolAPIKeys
	if len(keys) == 1 {
		return keys[0]
	}

	globalCounter++

	return keys[globalCounter%len(keys)]
}

func completeConnectWalletMission(ctx context.Context, address string) error {
	mission, err := dao.GetMissionById(ctx, MissionIdConnectWallet)
	if err != nil {
		log.Errorf("GetMissionById: %v", err)
		return err
	}

	ums, err := dao.GetUserMissionByMissionId(ctx, address, mission.ID, dao.QueryOption{})
	if err != nil {
		log.Errorf("GetUserMissionByMissionId: %v", err)
		return err
	}

	if len(ums) == 0 {
		return dao.AddUserMissionAndInviteLog(ctx, &model.UserMission{
			Username:  address,
			MissionID: mission.ID,
			Type:      mission.Type,
			Credit:    mission.Credit,
			Content:   address,
			CreatedAt: time.Now(),
		})
		// return dao.AddUserMission(ctx, &model.UserMission{
		// 	Username:  address,
		// 	MissionID: mission.ID,
		// 	Type:      mission.Type,
		// 	Credit:    mission.Credit,
		// 	Content:   address,
		// 	CreatedAt: time.Now(),
		// })
	}

	return nil
}

func checkFollowTwitter(ctx context.Context, mission *model.Mission, username string, queryOpt dao.QueryOption) error {
	twitterUser, err := dao.GetTwitterOauthByUsername(ctx, username)
	if err != nil {
		log.Errorf("GetTwitterOauthByUsername: %v", err)
		return err
	}

	client := swagger.NewAPIClient(swagger.NewConfiguration())

	apiKey := GetUToolKeyByRoundRobin()

	option := &swagger.TwitterFollowsApiToolsApiFollowingsIdsUsingGETOpts{
		UserId: optional.NewString(twitterUser.TwitterUserID),
	}
	result, _, err := client.TwitterFollowsApiToolsApi.FollowingsIdsUsingGET(ctx, apiKey, option)
	if err != nil {
		log.Errorf("FollowingsIdsUsingGET: %v", err)
		return err
	}

	if result.Code != 1 {
		log.Errorf("code: %d %s", result.Code, result.Msg)
		return errors.New(result.Msg)
	}

	type followIdsResp struct {
		PreviousCursor       int64   `json:"previous_cursor"`
		Ids                  []int64 `json:"ids"`
		PreviousCursorString string  `json:"previous_cursor_str"`
		NextCursor           int64   `json:"next_cursor"`
		NextCursorStr        string  `json:"next_cursor_str"`
	}

	bytesString, ok := result.Data.(string)
	if !ok {
		return errors.New("response not string")
	}

	var resp followIdsResp
	if err := json.Unmarshal([]byte(bytesString), &resp); err != nil {
		return err
	}

	var followed bool

	for _, id := range resp.Ids {
		if id != config.Cfg.OfficialTwitterUserId {
			continue
		}

		followed = true
		break
	}

	if !followed {
		return errors.New("user unfollow")
	}

	ums, err := dao.GetUserMissionByMissionId(ctx, username, mission.ID, queryOpt)
	if err != nil {
		log.Errorf("GetUserMissionByUser: %v", err)
		return err
	}

	if len(ums) == 0 {
		err = dao.AddUserMissionAndInviteLog(ctx, &model.UserMission{
			Username:  username,
			MissionID: mission.ID,
			Type:      mission.Type,
			Credit:    mission.Credit,
			Content:   twitterUser.TwitterUserID,
			CreatedAt: time.Now(),
		})
		// err = dao.AddUserMission(ctx, &model.UserMission{
		// 	Username:  username,
		// 	MissionID: mission.ID,
		// 	Type:      mission.Type,
		// 	Credit:    mission.Credit,
		// 	Content:   twitterUser.TwitterUserID,
		// 	CreatedAt: time.Now(),
		// })

		if err != nil {
			log.Errorf("AddUserMission: %v", err)
			return err
		}
	}

	return nil
}

func checkLikeTweet(ctx context.Context, mission *model.Mission, username string, queryOpt dao.QueryOption) error {
	twitterUser, err := dao.GetTwitterOauthByUsername(ctx, username)
	if err != nil {
		log.Errorf("GetTwitterOauthByUsername: %v", err)
		return err
	}

	openURL, err := url.Parse(mission.OpenUrl)
	if err != nil {
		log.Errorf("Parse OPEN URL: %v", err)
		return err
	}

	tweetId := openURL.Query().Get("tweet_id")
	client := swagger.NewAPIClient(swagger.NewConfiguration())

	apiKey := GetUToolKeyByRoundRobin()
	option := &swagger.TwitterGetTweesApiToolsApiFavoritersV2UsingGETOpts{}

	result, _, err := client.TwitterGetTweesApiToolsApi.FavoritersV2UsingGET(ctx, apiKey, tweetId, option)
	if err != nil {
		log.Errorf("GetListByUserIdOrScreenNameUsingGET: %v", err)
		return err
	}

	if result.Code != 1 {
		log.Errorf("code: %d %s", result.Code, result.Msg)
		return errors.New(result.Msg)
	}

	v, err := fastjson.Parse(result.Data.(string))
	if err != nil {
		return errors.New(err.Error())
	}

	var liked bool

	entries := v.Get("data").Get("favoriters_timeline").Get("timeline").Get("instructions", "0").GetArray("entries")

	for _, e := range entries {
		twitterUserId := e.Get("content").Get("itemContent").Get("user_results").Get("result").GetStringBytes("rest_id")

		if string(twitterUserId) != twitterUser.TwitterUserID {
			continue
		}

		liked = true
		break
	}

	if !liked {
		return errors.New("mission uncompleted")
	}

	ums, err := dao.GetUserMissionByMissionId(ctx, username, mission.ID, queryOpt)
	if err != nil {
		log.Errorf("GetUserMissionByUser: %v", err)
		return err
	}

	if len(ums) == 0 {
		err = dao.AddUserMissionAndInviteLog(ctx, &model.UserMission{
			Username:  username,
			MissionID: mission.ID,
			Type:      mission.Type,
			Credit:    mission.Credit,
			Content:   mission.OpenUrl,
			CreatedAt: time.Now(),
		})
		// err = dao.AddUserMission(ctx, &model.UserMission{
		// 	Username:  username,
		// 	MissionID: mission.ID,
		// 	Type:      mission.Type,
		// 	Credit:    mission.Credit,
		// 	Content:   mission.OpenUrl,
		// 	CreatedAt: time.Now(),
		// })

		if err != nil {
			log.Errorf("AddUserMission: %v", err)
			return err
		}
	}

	return nil
}

func checkReTweet(ctx context.Context, mission *model.Mission, username string, queryOpt dao.QueryOption) error {
	twitterUser, err := dao.GetTwitterOauthByUsername(ctx, username)
	if err != nil {
		log.Errorf("GetTwitterOauthByUsername: %v", err)
		return err
	}

	client := swagger.NewAPIClient(swagger.NewConfiguration())

	apiKey := GetUToolKeyByRoundRobin()
	option := &swagger.TwitterGetTweesApiToolsApiRetweetersV2UsingGETOpts{}

	openURL, err := url.Parse(mission.OpenUrl)
	if err != nil {
		log.Errorf("Parse OPEN URL: %v", err)
		return err
	}

	tweetId := openURL.Query().Get("tweet_id")
	result, _, err := client.TwitterGetTweesApiToolsApi.RetweetersV2UsingGET(ctx, apiKey, tweetId, option)
	if err != nil {
		log.Errorf("GetListByUserIdOrScreenNameUsingGET: %v", err)
		return err
	}

	if result.Code != 1 {
		log.Errorf("code: %d %s", result.Code, result.Msg)
		return errors.New(result.Msg)
	}

	v, err := fastjson.Parse(result.Data.(string))
	if err != nil {
		return errors.New(err.Error())
	}

	var liked bool

	entries := v.Get("data").Get("retweeters_timeline").Get("timeline").Get("instructions", "0").GetArray("entries")

	for _, e := range entries {
		twitterUserId := e.Get("content").Get("itemContent").Get("user_results").Get("result").GetStringBytes("rest_id")

		if string(twitterUserId) != twitterUser.TwitterUserID {
			continue
		}

		liked = true
		break
	}

	if !liked {
		return errors.New("mission uncompleted")
	}

	ums, err := dao.GetUserMissionByMissionId(ctx, username, mission.ID, queryOpt)
	if err != nil {
		log.Errorf("GetUserMissionByUser: %v", err)
		return err
	}

	if len(ums) == 0 {
		err = dao.AddUserMissionAndInviteLog(ctx, &model.UserMission{
			Username:  username,
			MissionID: mission.ID,
			Type:      mission.Type,
			Credit:    mission.Credit,
			Content:   mission.OpenUrl,
			CreatedAt: time.Now(),
		})

		if err != nil {
			log.Errorf("AddUserMission: %v", err)
			return err
		}
	}

	return nil
}

func checkQuoteTweet(ctx context.Context, mission *model.Mission, username string, queryOpt dao.QueryOption) error {
	twitterUser, err := dao.GetTwitterOauthByUsername(ctx, username)
	if err != nil {
		log.Errorf("GetTwitterOauthByUsername: %v", err)
		return err
	}

	//openURL, err := url.Parse(mission.OpenUrl)
	//if err != nil {
	//	log.Errorf("Parse OPEN URL: %v", err)
	//	return err
	//}

	//fmt.Println(openURL)
	//openURLPaths := strings.Split(openURL.Path, "/")
	//
	//if len(openURLPaths) != 4 {
	//	log.Errorf("Invalid URL: %v", err)
	//	return errors.New("invalid post url")
	//}
	//tweetId := openURLPaths[len(openURLPaths)-1]

	client := swagger.NewAPIClient(swagger.NewConfiguration())
	twitterLink, err := dao.GetUserTwitterLink(ctx, username, mission.ID, carbon.Now().StartOfDay().String())
	if err != nil {
		log.Errorf("GetUserTwitterLink: %v", err)
		return err
	}

	replyURL, err := url.Parse(twitterLink.Link)
	if err != nil {
		log.Errorf("Parse twitter Link URL: %v", err)
		return err
	}

	paths := strings.Split(replyURL.Path, "/")

	if len(paths) != 4 {
		log.Errorf("Invalid URL: %v", err)
		return errors.New("invalid url")
	}

	replyId := strings.TrimSpace(paths[len(paths)-1])

	apiKey := GetUToolKeyByRoundRobin()
	option := &swagger.TwitterGetTweesApiToolsApiTweetTimelineUsingGETOpts{}
	result, _, err := client.TwitterGetTweesApiToolsApi.TweetTimelineUsingGET(ctx, apiKey, replyId, option)
	if err != nil {
		log.Errorf("GetListByUserIdOrScreenNameUsingGET: %v", err)
		return err
	}

	if result.Code != 1 {
		log.Errorf("code: %d %s", result.Code, result.Msg)
		return errors.New(result.Msg)
	}

	v, err := fastjson.Parse(result.Data.(string))
	if err != nil {
		return errors.New(err.Error())
	}

	entries := v.Get("data").Get("threaded_conversation_with_injections_v2").Get("instructions", "0").Get("entries", "0")
	res := entries.Get("content").Get("itemContent").Get("tweet_results").Get("result")
	postUserId := res.Get("core").Get("user_results").Get("result").GetStringBytes("rest_id")
	sourceTweetId := res.Get("quoted_status_result").Get("result").GetStringBytes("rest_id")
	mentions := res.Get("legacy").Get("entities").GetArray("user_mentions")
	postCreatedAt := res.Get("legacy").GetStringBytes("created_at")

	if !strings.Contains(mission.OpenUrl, string(sourceTweetId)) {
		return errors.New("not allowed source tweet id")
	}

	if string(postUserId) != twitterUser.TwitterUserID {
		return errors.New(fmt.Sprintf("invalid link, expected twitter user id: %s got: %s", twitterUser.TwitterUserID, string(postUserId)))
	}

	if carbon.Parse(string(postCreatedAt)).Lt(carbon.Now().StartOfDay()) {
		return errors.New("post expiration")
	}

	if len(mentions) < 3 {
		return errors.New("tag users not enough")
	}

	ums, err := dao.GetUserMissionByMissionId(ctx, username, mission.ID, queryOpt)
	if err != nil {
		log.Errorf("GetUserMissionByUser: %v", err)
		return err
	}

	if len(ums) == 0 {
		err = dao.AddUserMissionAndInviteLog(ctx, &model.UserMission{
			Username:  username,
			MissionID: mission.ID,
			Type:      mission.Type,
			Credit:    mission.Credit,
			Content:   twitterUser.TwitterUserID,
			CreatedAt: time.Now(),
		})

		if err != nil {
			log.Errorf("AddUserMission: %v", err)
			return err
		}
	}

	return nil
}

func checkPostTweet(ctx context.Context, mission *model.Mission, username string, queryOpt dao.QueryOption) error {
	twitterUser, err := dao.GetTwitterOauthByUsername(ctx, username)
	if err != nil {
		log.Errorf("GetTwitterOauthByUsername: %v", err)
		return err
	}

	openURL, err := url.Parse(mission.OpenUrl)
	if err != nil {
		log.Errorf("Parse OPEN URL: %v", err)
		return err
	}

	copyContent := openURL.Query().Get("text")

	client := swagger.NewAPIClient(swagger.NewConfiguration())
	twitterLink, err := dao.GetUserTwitterLink(ctx, username, mission.ID, carbon.Now().StartOfDay().String())
	if err != nil {
		log.Errorf("GetUserTwitterLink: %v", err)
		return err
	}

	replyURL, err := url.Parse(twitterLink.Link)
	if err != nil {
		log.Errorf("Parse twitter Link URL: %v", err)
		return err
	}

	paths := strings.Split(replyURL.Path, "/")

	if len(paths) != 4 {
		log.Errorf("Invalid URL: %v", err)
		return errors.New("invalid url")
	}

	replyId := strings.TrimSpace(paths[len(paths)-1])

	apiKey := GetUToolKeyByRoundRobin()
	option := &swagger.TwitterGetTweesApiToolsApiTweetTimelineUsingGETOpts{}
	result, _, err := client.TwitterGetTweesApiToolsApi.TweetTimelineUsingGET(ctx, apiKey, replyId, option)
	if err != nil {
		log.Errorf("GetListByUserIdOrScreenNameUsingGET: %v", err)
		return err
	}

	if result.Code != 1 {
		log.Errorf("code: %d %s", result.Code, result.Msg)
		return errors.New(result.Msg)
	}

	v, err := fastjson.Parse(result.Data.(string))
	if err != nil {
		log.Errorf("parse data: %v", err)
		return err
	}

	entries := v.Get("data").Get("threaded_conversation_with_injections_v2").Get("instructions", "0").Get("entries", "0")
	res := entries.Get("content").Get("itemContent").Get("tweet_results").Get("result")
	postUserId := res.Get("core").Get("user_results").Get("result").GetStringBytes("rest_id")
	fullText := res.Get("legacy").GetStringBytes("full_text")
	postCreatedAt := res.Get("legacy").GetStringBytes("created_at")

	FirstLine := strings.Split(string(fullText), "\n")[0]

	if !strings.Contains(copyContent, FirstLine) {
		return errors.New("invalid content")
	}

	if carbon.Parse(string(postCreatedAt)).Lt(carbon.Now().StartOfDay()) {
		return errors.New("post expiration")
	}

	if string(postUserId) != twitterUser.TwitterUserID {
		return errors.New(fmt.Sprintf("invalid link, expected twitter user id: %s got: %s", twitterUser.TwitterUserID, string(postUserId)))
	}

	ums, err := dao.GetUserMissionByMissionId(ctx, username, mission.ID, queryOpt)
	if err != nil {
		log.Errorf("GetUserMissionByUser: %v", err)
		return err
	}

	if len(ums) == 0 {
		err = dao.AddUserMissionAndInviteLog(ctx, &model.UserMission{
			Username:  username,
			MissionID: mission.ID,
			Type:      mission.Type,
			Credit:    mission.Credit,
			Content:   twitterUser.TwitterUserID,
			CreatedAt: time.Now(),
		})

		if err != nil {
			log.Errorf("AddUserMission: %v", err)
			return err
		}
	}

	return nil
}

func checkJoinTelegram(ctx context.Context, mission *model.Mission, username string, queryOpt dao.QueryOption) error {
	// todo: get telegram id by username

	return errors.New("not implement")

	//telegramUserId := int64(6177508547)
	//
	//fmt.Println(telegramUserId, config.Cfg.OfficialTelegramGroupId)
	//
	//_, err := Bot.ChatMemberOf(&tele.Chat{ID: config.Cfg.OfficialTelegramGroupId}, &tele.Chat{ID: telegramUserId})
	//if err != nil {
	//	fmt.Println("chat member of: ", err)
	//	return err
	//}
	//
	//ums, err := dao.GetUserMissionByMissionId(ctx, username, mission.ID, dao.QueryOption{})
	//if err != nil {
	//	log.Errorf("GetUserMissionByUser: %v", err)
	//	return err
	//}
	//
	//if len(ums) == 0 {
	//	err = dao.AddUserMission(ctx, &model.UserMission{
	//		Username:  username,
	//		MissionID: mission.ID,
	//		Type:      mission.Type,
	//		Credit:    mission.Credit,
	//		Content:   strconv.FormatInt(telegramUserId, 10),
	//		CreatedAt: time.Now(),
	//	})
	//
	//	if err != nil {
	//		log.Errorf("AddUserMission: %v", err)
	//		return err
	//	}
	//}
	//
	//return nil
}

func checkJoinDiscord(ctx context.Context, mission *model.Mission, username string, queryOpt dao.QueryOption) error {
	discordUser, err := dao.GetDiscordOAuthByUsername(ctx, username)
	if err != nil {
		log.Errorf("GetDiscordOAuthByUsername: %v", err)
		return err
	}

	key := "gm::discord::members"
	existing, err := dao.RedisCache.SIsMember(ctx, key, discordUser.DiscordUserID).Result()
	if err != nil {
		log.Errorf("SIsMember: %v", err)
		return err
	}

	if !existing {
		return errors.New("user not join discord")
	}

	ums, err := dao.GetUserMissionByMissionId(ctx, username, mission.ID, queryOpt)
	if err != nil {
		log.Errorf("GetUserMissionByUser: %v", err)
		return err
	}

	if len(ums) == 0 {
		err = dao.AddUserMissionAndInviteLog(ctx, &model.UserMission{
			Username:  username,
			MissionID: mission.ID,
			Type:      mission.Type,
			Credit:    mission.Credit,
			Content:   discordUser.DiscordUserID,
			CreatedAt: time.Now(),
		})

		if err != nil {
			log.Errorf("AddUserMission: %v", err)
			return err
		}
	}

	return nil
}

func checkBindingKOL(ctx context.Context, mission *model.Mission, username string, queryOpt dao.QueryOption) error {
	userInfo, err := dao.GetUserByUsername(ctx, username)
	if err != nil {
		return err
	}

	if userInfo.FromKolRefCode == "" {
		return errors.New("please complete mission first")
	}

	ums, err := dao.GetUserMissionByMissionId(ctx, username, mission.ID, queryOpt)
	if err != nil {
		log.Errorf("GetUserMissionByUser: %v", err)
		return err
	}

	if len(ums) == 0 {
		err = dao.AddUserMissionAndInviteLog(ctx, &model.UserMission{
			Username:  username,
			MissionID: mission.ID,
			Type:      mission.Type,
			Credit:    mission.Credit,
			Content:   userInfo.FromKolRefCode,
			CreatedAt: time.Now(),
		})

		if err != nil {
			log.Errorf("AddUserMission: %v", err)
			return err
		}
	}

	return nil
}

func checkVisitOfficialWebsite(ctx context.Context, mission *model.Mission, username string, queryOpt dao.QueryOption) error {
	return errors.New("not implement")
}

func checkVisitReferrerPage(ctx context.Context, mission *model.Mission, username string, queryOpt dao.QueryOption) error {
	return errors.New("not implement")
}

func checkInviteFriendsToDiscord(ctx context.Context, mission *model.Mission, username string, queryOpt dao.QueryOption) error {
	discordUser, err := dao.GetDiscordOAuthByUsername(ctx, username)
	if err != nil {
		log.Errorf("GetDiscordOAuthByUsername: %v", err)
		return err
	}

	subMission, err := dao.GetSubMissions(ctx, mission.ID)
	if err != nil {
		log.Errorf("GetSubMissions: %v", err)
		return err
	}

	startOfWeek := carbon.Now().StartOfWeek().String()
	key := fmt.Sprintf("gm::discord::invitecounter::%s::%s", discordUser.DiscordUserID, startOfWeek)
	result, err := dao.RedisCache.Get(ctx, key).Result()
	if err != nil {
		log.Errorf("Get invite counter: %v", err)
		return err
	}

	count, _ := strconv.ParseInt(result, 10, 64)

	log.Infof("%s invite friends count: %d", discordUser.DiscordUserID, count)
	fmt.Println("invite friends count:", discordUser.DiscordUserID, count)

	if count <= 0 {
		return errors.New("complete mission first")
	}

	ums, err := dao.GetUserMissionByMissionId(ctx, username, mission.ID, queryOpt)
	if err != nil {
		log.Errorf("GetUserMissionByUser: %v", err)
		return err
	}

	completedSubMissionId := make(map[int64]struct{})
	for _, m := range ums {
		completedSubMissionId[m.SubMissionID] = struct{}{}
	}

	completedCount := len(completedSubMissionId)

	for _, sm := range subMission {
		expectedCount, err := strconv.ParseInt(sm.Title, 10, 64)
		if err != nil {
			log.Errorf("Parse invite friends count expect value: %v", err)
			return err
		}

		if _, ok := completedSubMissionId[sm.ID]; ok {
			continue
		}

		if count < expectedCount {
			continue
		}

		completedCount++
		err = dao.AddUserMissionAndInviteLog(ctx, &model.UserMission{
			Username:     username,
			MissionID:    mission.ID,
			SubMissionID: sm.ID,
			Type:         sm.Type,
			Credit:       sm.Credit,
			Content:      discordUser.DiscordUserID,
			CreatedAt:    time.Now(),
		})

		if err != nil {
			log.Errorf("AddUserMission: %v", err)
			return err
		}
	}

	if completedCount <= 0 {
		return errors.New("complete mission first")
	}

	return nil
}

func PostTwitterLinkHandler(c *gin.Context) {
	claims := jwt.ExtractClaims(c)
	username := claims[identityKey].(string)

	var params model.UserTwitterLink
	if err := c.BindJSON(&params); err != nil {
		c.JSON(http.StatusOK, respErrorCode(errorsx.InvalidParams, c))
		return
	}

	if params.Link == "" || params.MissionID == 0 {
		c.JSON(http.StatusOK, respErrorCode(errorsx.InvalidParams, c))
		return
	}

	params.Username = username
	params.CreatedAt = time.Now()
	err := dao.AddUserTwitterLink(c.Request.Context(), &params)
	if err != nil {
		log.Errorf("AddUserTwitterLink: %v", err)
		c.JSON(http.StatusOK, respErrorCode(errorsx.InternalServer, c))
		return
	}

	c.JSON(http.StatusOK, respJSON(nil))
}

func BindingKOLReferralCodeHandler(c *gin.Context) {
	claims := jwt.ExtractClaims(c)
	username := claims[identityKey].(string)

	var params = struct {
		Code string `json:"code"`
	}{}

	if err := c.BindJSON(&params); err != nil {
		c.JSON(http.StatusOK, respErrorCode(errorsx.InvalidParams, c))
		return
	}

	if params.Code == "" {
		c.JSON(http.StatusOK, respErrorCode(errorsx.InvalidReferralCode, c))
		return
	}

	userInfo, err := dao.GetUserByUsername(c.Request.Context(), username)
	if err != nil {
		log.Errorf("GetUserByUsername: %v", err)
		c.JSON(http.StatusOK, respErrorCode(errorsx.InternalServer, c))
		return
	}

	if userInfo.FromKolRefCode != "" {
		c.JSON(http.StatusOK, respErrorCode(errorsx.ReferralCodeBound, c))
		return
	}

	kolUserId, err := GetKOLUserId(c.Request.Context(), params.Code)
	if err != nil || kolUserId == "" {
		c.JSON(http.StatusOK, respErrorCode(errorsx.InvalidVerifyCode, c))
		return
	}

	err = dao.UpdateUserKOLReferralCode(c.Request.Context(), username, params.Code, kolUserId)
	if err != nil {
		log.Errorf("UpdateUserKOLReferralCode: %v", err)
		c.JSON(http.StatusOK, respErrorCode(errorsx.InternalServer, c))
		return
	}

	c.JSON(http.StatusOK, respJSON(nil))
}

func GetKOLUserId(ctx context.Context, code string) (string, error) {
	client := http.DefaultClient

	requestURL := config.Cfg.TitanAPI.BasePath + "/v1/kol/code?code=" + code
	request, err := http.NewRequest(http.MethodGet, requestURL, nil)
	if err != nil {
		log.Errorf("create request: %v", err)
		return "", err
	}

	request.Header.Add("Authorization", "Bearer "+config.Cfg.TitanAPI.Key)
	resp, err := client.Do(request)
	if err != nil {
		log.Errorf("get response: %v", err)
		return "", err
	}

	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return "", errors.New(resp.Status)
	}

	result, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	v, err := fastjson.Parse(string(result))
	if err != nil {
		return "", errors.New(err.Error())
	}

	userId := v.Get("data").GetStringBytes("kol_user_id")

	return string(userId), nil
}

func GetUserCreditsHandler(c *gin.Context) {
	userId := c.Query("user_id")
	if userId == "" {
		c.JSON(http.StatusOK, respErrorCode(errorsx.InvalidParams, c))
		return
	}

	pageSize, _ := strconv.Atoi(c.Query("page_size"))
	page, _ := strconv.Atoi(c.Query("page"))
	order := c.Query("order")
	orderField := c.Query("order_field")
	option := dao.QueryOption{
		Page:       page,
		PageSize:   pageSize,
		Order:      order,
		OrderField: orderField,
	}

	total, referralList, err := dao.GetUserCreditsByKOLReferralCode(c.Request.Context(), userId, option)
	if err != nil {
		log.Errorf("GetUserCreditsByKOLReferralCode: %v", err)
		c.JSON(http.StatusOK, respErrorCode(errorsx.InternalServer, c))
		return
	}

	commission, err := dao.GetKOLCommissionCredits(c.Request.Context(), userId)
	if err != nil && !errors.Is(err, sql.ErrNoRows) {
		log.Errorf("GetKOLCommissionCredits: %v", err)
		c.JSON(http.StatusOK, respErrorCode(errorsx.InternalServer, c))
		return
	}

	c.JSON(http.StatusOK, respJSON(JsonObject{
		"kol_commission_credits": commission,
		"list":                   referralList,
		"total":                  total,
	}))
}

// BrowsOfficialWebsite 浏览官网
func BrowsOfficialWebsite(c *gin.Context) {
	claims := jwt.ExtractClaims(c)
	username := claims[identityKey].(string)

	tnow := time.Now().Unix()
	value := fmt.Sprintf("%s:%d", username, tnow)
	code, err := opcrypt.AesEncryptCBC([]byte(value), []byte(config.Cfg.AesKey))
	if err != nil {
		log.Errorf("generate code of brows official website error: %v", err)
		c.JSON(http.StatusOK, respErrorCode(errorsx.InternalServer, c))
		return
	}

	c.JSON(http.StatusOK, respJSON(JsonObject{
		"code": code,
		"uri":  config.Cfg.OfficialWebsiteURI,
	}))
}

// BrowsOfficialWebsiteCallback 浏览官网回调
func BrowsOfficialWebsiteCallback(c *gin.Context) {
	var params = struct {
		Code string `json:"code"`
	}{}

	if err := c.BindJSON(&params); err != nil {
		c.JSON(http.StatusOK, respErrorCode(errorsx.InvalidParams, c))
		return
	}

	code := strings.TrimSpace(params.Code)

	if code == "" {
		c.JSON(http.StatusOK, respErrorCode(errorsx.InvalidParams, c))
		return
	}

	tnow := time.Now().Unix()

	value, err := opcrypt.AesDecryptCBC(code, []byte(config.Cfg.AesKey))
	if err != nil {
		log.Errorf("get code of brows official website error: %v", err)
		c.JSON(http.StatusOK, respErrorCode(errorsx.InternalServer, c))
		return
	}

	values := strings.Split(string(value), ":")
	if len(values) < 2 {
		log.Errorf("split code of brows official website error: %v", err)
		c.JSON(http.StatusOK, respErrorCode(errorsx.InternalServer, c))
		return
	}
	username := values[0]
	pnow, _ := strconv.ParseInt(values[1], 10, 64)
	if pnow == 0 {
		log.Errorf("time of code error: %v", err)
		c.JSON(http.StatusOK, respErrorCode(errorsx.InternalServer, c))
		return
	}

	// 如果浏览官网到达了规定时间，则发放积分奖励
	if tnow-pnow >= config.Cfg.BrowsOfficialWebsiteTime {
		err = completeMission(c.Request.Context(), username, MissionIdBrowsOfficialWebSite)
		if err != nil {
			log.Errorf("complete brows official website error: %v", err)
			c.JSON(http.StatusOK, respErrorCode(errorsx.InternalServer, c))
			return
		}
	}

	c.JSON(http.StatusOK, respJSON(nil))
}

// VerifyBrowsOfficialWebsite 验证浏览官网是否完成
func VerifyBrowsOfficialWebsite(c *gin.Context) {
	var msg string

	claims := jwt.ExtractClaims(c)
	username := claims[identityKey].(string)
	lang := c.GetHeader("Lang")

	complete, err := getMission(c.Request.Context(), username, MissionIdBrowsOfficialWebSite)
	if err != nil {
		log.Errorf("get user mission error: %v", err)
		c.JSON(http.StatusOK, respErrorCode(errorsx.InternalServer, c))
		return
	}

	switch strings.ToLower(lang) {
	case "cn":
		if !complete {
			msg = "请先完成任务"
		}
	default:
		if !complete {
			msg = "Please complete the task first"
		}
	}

	c.JSON(http.StatusOK, respJSON(JsonObject{
		"verified": complete,
		"msg":      msg,
	}))
}

// GetInviteLogs 获取邀请记录
func GetInviteLogs(c *gin.Context) {
	claims := jwt.ExtractClaims(c)
	username := claims[identityKey].(string)

	page, _ := c.GetQuery("page")
	size, _ := c.GetQuery("size")
	pageInt, _ := strconv.Atoi(page)
	sizeInt, _ := strconv.Atoi(size)

	out, total, err := dao.GetUserInviteLogs(c.Request.Context(), username, dao.QueryOption{Page: pageInt, PageSize: sizeInt})
	if err != nil {
		log.Errorf("get user invite_log error: %v", err)
		c.JSON(http.StatusOK, respErrorCode(errorsx.InternalServer, c))
		return
	}

	c.JSON(http.StatusOK, respJSON(JsonObject{
		"total": total,
		"list":  out,
	}))
}

// GetMissionLogs 获取任务完成记录
func GetMissionLogs(c *gin.Context) {
	claims := jwt.ExtractClaims(c)
	username := claims[identityKey].(string)

	page, _ := c.GetQuery("page")
	size, _ := c.GetQuery("size")
	pageInt, _ := strconv.Atoi(page)
	sizeInt, _ := strconv.Atoi(size)

	out, total, err := dao.GetMissionLogs(c.Request.Context(), username, dao.QueryOption{Page: pageInt, PageSize: sizeInt})
	if err != nil {
		log.Errorf("get user mission_log error: %v", err)
		c.JSON(http.StatusOK, respErrorCode(errorsx.InternalServer, c))
		return
	}

	c.JSON(http.StatusOK, respJSON(JsonObject{
		"total": total,
		"list":  out,
	}))
}

// GetBecomeVolunteerURL 获取跳转链接
func GetBecomeVolunteerURL(c *gin.Context) {
	var url string
	lang := c.GetHeader("Lang")

	switch strings.ToLower(lang) {
	case "cn":
		url = config.Cfg.GoogleDoc.CnURI
	default:
		url = config.Cfg.GoogleDoc.EnURI
	}

	c.JSON(http.StatusOK, respJSON(JsonObject{
		"url": url,
	}))
}

// VerifyBecomeVolunteer 验证是否完成报表填写
func VerifyBecomeVolunteer(c *gin.Context) {
	var (
		msg, speedID string
		complete     bool
		err          error
	)

	claims := jwt.ExtractClaims(c)
	username := claims[identityKey].(string)
	glog.Println(username)

	lang := c.GetHeader("Lang")

	// 调用谷歌文档接口进行查询
	switch strings.ToLower(lang) {
	case "cn":
		msg = "请先完成任务"
		speedID = config.Cfg.GoogleDoc.CnDocID
	default:
		msg = "Please complete the task first"
		speedID = config.Cfg.GoogleDoc.EnDocID
	}

	// 判断是否完成验证
	complete, err = checkBVComplete(username, speedID)
	if err != nil {
		log.Errorf("check complete error: %v", err)
		c.JSON(http.StatusOK, respErrorCode(errorsx.InternalServer, c))
		return
	}
	// 查询到则添加任务完成记录
	if complete {
		err = completeMission(c.Request.Context(), username, MissionIdBecomeVolunteer)
		if err != nil {
			log.Errorf("get user mission_log error: %v", err)
			c.JSON(http.StatusOK, respErrorCode(errorsx.InternalServer, c))
			return
		}
		msg = ""
	}

	c.JSON(http.StatusOK, respJSON(JsonObject{
		"verified": complete,
		"msg":      msg,
	}))
}
