package api

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	jwt "github.com/appleboy/gin-jwt/v2"
	"github.com/gin-gonic/gin"
	"github.com/gnasnik/titan-quest/core/dao"
	"github.com/gnasnik/titan-quest/core/errors"
	"github.com/gnasnik/titan-quest/core/generated/model"
	"github.com/gnasnik/titan-quest/pkg/random"
	"github.com/go-redis/redis/v9"
	"net/http"
	"strings"
	"time"
)

func GetUserInfoHandler(c *gin.Context) {
	claims := jwt.ExtractClaims(c)
	username := claims[identityKey].(string)
	user, err := dao.GetUserByUsername(c.Request.Context(), username)
	if err != nil {
		c.JSON(http.StatusOK, respError(errors.ErrNotFound))
		return
	}

	c.JSON(http.StatusOK, respJSON(user))
}

type NonceStringType string

const (
	NonceStringTypeRegister  NonceStringType = "1"
	NonceStringTypeLogin     NonceStringType = "2"
	NonceStringTypeReset     NonceStringType = "3"
	NonceStringTypeSignature NonceStringType = "4"
)

var defaultNonceExpiration = 5 * time.Minute

func getRedisNonceSignatureKey(username string) string {
	return fmt.Sprintf("TITAN::QUEST::SIGN::%s", username)
}

func getRedisNonceRegisterKey(username string) string {
	return fmt.Sprintf("TITAN::QUEST::REG::%s", username)
}

func getRedisNonceLoginKey(username string) string {
	return fmt.Sprintf("TITAN::QUEST::LOGIN::%s", username)
}

func getRedisNonceResetKey(username string) string {
	return fmt.Sprintf("TITAN::QUEST::RESET::%s", username)
}

func getNonceFromCache(ctx context.Context, username string, t NonceStringType) (string, error) {
	var key string

	switch t {
	case NonceStringTypeRegister:
		key = getRedisNonceRegisterKey(username)
	case NonceStringTypeLogin:
		key = getRedisNonceLoginKey(username)
	case NonceStringTypeReset:
		key = getRedisNonceResetKey(username)
	case NonceStringTypeSignature:
		key = getRedisNonceSignatureKey(username)
	default:
		return "", fmt.Errorf("unsupported nonce string type")
	}

	bytes, err := dao.RedisCache.Get(ctx, key).Bytes()
	if err == redis.Nil {
		return "", nil
	}

	if err != nil {
		return "", err
	}

	var verifyCode string
	err = json.Unmarshal(bytes, &verifyCode)
	if err != nil {
		return "", err
	}

	return verifyCode, nil
}

func GetNonceStringHandler(c *gin.Context) {
	username := c.Query("username")
	if username == "" {
		c.JSON(http.StatusOK, respErrorCode(errors.InvalidParams, c))
		return
	}

	nonce, err := generateNonceString(c.Request.Context(), getRedisNonceSignatureKey(username))
	if err != nil {
		c.JSON(http.StatusOK, respErrorCode(errors.InternalServer, c))
		return
	}

	_, err = dao.GetUserByUsername(c.Request.Context(), username)
	if err == sql.ErrNoRows {
		//c.JSON(http.StatusOK, respErrorCode(errors.UserNotFound, c))
		//return
		user := &model.User{
			Username:     username,
			CreatedAt:    time.Now(),
			UpdatedAt:    time.Now(),
			ReferralCode: random.GenerateRandomString(6),
		}
		err = dao.CreateUser(c.Request.Context(), user)
		if err != nil {
			c.JSON(http.StatusOK, respErrorCode(errors.InternalServer, c))
			return
		}
	}

	if err != nil {
		c.JSON(http.StatusOK, respErrorCode(errors.InternalServer, c))
		return
	}

	c.JSON(http.StatusOK, respJSON(JsonObject{
		"code": nonce,
	}))
}

func generateNonceString(ctx context.Context, key string) (string, error) {
	rand := random.GenerateRandomNumber(6)
	verifyCode := "TitanNetWork(" + rand + ")"
	bytes, err := json.Marshal(verifyCode)
	if err != nil {
		return "", err
	}

	_, err = dao.RedisCache.Set(ctx, key, bytes, defaultNonceExpiration).Result()
	if err != nil {
		log.Errorf("%v:", err)
		return "", err
	}

	return verifyCode, nil
}

func GetNumericVerifyCodeHandler(c *gin.Context) {
	userInfo := &model.User{}
	userInfo.Username = c.Query("username")
	verifyType := c.Query("type")
	lang := c.GetHeader("Lang")
	userInfo.UserEmail = userInfo.Username

	var key string
	switch NonceStringType(verifyType) {
	case NonceStringTypeRegister:
		key = getRedisNonceRegisterKey(userInfo.Username)
	case NonceStringTypeLogin:
		key = getRedisNonceLoginKey(userInfo.Username)
	case NonceStringTypeReset:
		key = getRedisNonceResetKey(userInfo.Username)
	case NonceStringTypeSignature:
		key = getRedisNonceSignatureKey(userInfo.Username)
	default:
		c.JSON(http.StatusOK, respErrorCode(errors.UnsupportedVerifyCodeType, c))
		return
	}

	nonce, err := getNonceFromCache(c.Request.Context(), userInfo.Username, NonceStringType(verifyType))
	if err != nil {
		c.JSON(http.StatusOK, respErrorCode(errors.InternalServer, c))
		return
	}

	if nonce != "" {
		c.JSON(http.StatusOK, respErrorCode(errors.GetVCFrequently, c))
		return
	}

	verifyCode := random.GenerateRandomNumber(6)

	if err = sendEmail(userInfo.Username, verifyCode, lang); err != nil {
		log.Errorf("send email: %v", err)
		if strings.Contains(err.Error(), "timed out") {
			c.JSON(http.StatusOK, respErrorCode(errors.TimeoutCode, c))
			return
		}
		c.JSON(http.StatusOK, respErrorCode(errors.InternalServer, c))
		return
	}

	if err = cacheVerifyCode(c.Request.Context(), key, verifyCode); err != nil {
		c.JSON(http.StatusOK, respErrorCode(errors.InternalServer, c))
		return
	}

	c.JSON(http.StatusOK, respJSON(JsonObject{
		"msg": "success",
	}))
}

func cacheVerifyCode(ctx context.Context, key, verifyCode string) error {
	bytes, err := json.Marshal(verifyCode)
	if err != nil {
		return err
	}

	_, err = dao.RedisCache.Set(ctx, key, bytes, defaultNonceExpiration).Result()
	if err != nil {
		return err
	}

	return nil
}

const (
	MissionIdConnectWallet int64 = iota + 1001
	MissionIdFollowTwitter
	MissionIdRetweet
	MissionIdLikeTwitter
	MissionIdJoinDiscord
	MissionIdJoinTelegram
)

const (
	MissionIdQuoteTweet int64 = iota + 1106
	MissionIdPostTweet
	MissionIdInviteFriendsToDiscord
)
