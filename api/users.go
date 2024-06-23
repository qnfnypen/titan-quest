package api

import (
	"context"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"time"

	constant "github.com/TestsLing/aj-captcha-go/const"
	"github.com/cosmos/cosmos-sdk/crypto/keys/secp256k1"

	jwt "github.com/appleboy/gin-jwt/v2"
	"github.com/gin-gonic/gin"
	"github.com/gnasnik/titan-quest/core/dao"
	"github.com/gnasnik/titan-quest/core/errors"
	"github.com/gnasnik/titan-quest/core/generated/model"
	"github.com/gnasnik/titan-quest/pkg/random"
	"github.com/go-redis/redis/v9"
)

type (
	// VerifyCodeReq 获取邮箱验证码
	VerifyCodeReq struct {
		Username  string `json:"username"`
		Token     string `json:"token"`
		PointJSON string `json:"pointJson"`
		Type      string `json:"type"`
	}
)

func GetUserInfoHandler(c *gin.Context) {
	claims := jwt.ExtractClaims(c)
	username := claims[identityKey].(string)
	// user, err := dao.GetUserByUsername(c.Request.Context(), username)
	// if err != nil {
	// 	c.JSON(http.StatusOK, respError(errors.ErrNotFound))
	// 	return
	// }

	resp, err := dao.GetUserResponse(c.Request.Context(), username)
	if err != nil {
		c.JSON(http.StatusOK, respError(errors.ErrNotFound))
		return
	}

	// c.JSON(http.StatusOK, respJSON(user.ToResponseUser()))
	c.JSON(http.StatusOK, respJSON(resp))
}

type NonceStringType string

const (
	NonceStringTypeRegister  NonceStringType = "1"
	NonceStringTypeLogin     NonceStringType = "2"
	NonceStringTypeReset     NonceStringType = "3"
	NonceStringTypeSignature NonceStringType = "4"
)

const titanWalletPrefix = "titan"

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
	// 获取邀请码
	// code := strings.TrimSpace(c.Query("code"))
	if username == "" {
		c.JSON(http.StatusOK, respErrorCode(errors.InvalidParams, c))
		return
	}

	nonce, err := generateNonceString(c.Request.Context(), getRedisNonceSignatureKey(username))
	if err != nil {
		c.JSON(http.StatusOK, respErrorCode(errors.InternalServer, c))
		return
	}

	// inviteCode := xid.New().String()
	// userExt := &model.UsersExt{
	// 	Username:    username,
	// 	InviteCode:  inviteCode,
	// 	InvitedCode: code,
	// }

	// _, err = dao.GetUserByUsername(c.Request.Context(), username)
	// if err == sql.ErrNoRows {
	// 	//c.JSON(http.StatusOK, respErrorCode(errors.UserNotFound, c))
	// 	//return
	// 	user := &model.User{
	// 		Username:     username,
	// 		CreatedAt:    time.Now(),
	// 		UpdatedAt:    time.Now(),
	// 		ReferralCode: random.GenerateRandomString(6),
	// 	}

	// 	err = dao.CreateUserInfo(c.Request.Context(), user, userExt)
	// 	if err != nil {
	// 		c.JSON(http.StatusOK, respErrorCode(errors.InternalServer, c))
	// 		return
	// 	}
	// }

	// if err != nil {
	// 	c.JSON(http.StatusOK, respErrorCode(errors.InternalServer, c))
	// 	return
	// }

	// // 老用户则增加邀请码
	// _, err = dao.GetUserExt(c.Request.Context(), username)
	// if err == sql.ErrNoRows {
	// 	err = dao.CreateUserExt(c.Request.Context(), userExt)
	// 	if err != nil {
	// 		c.JSON(http.StatusOK, respErrorCode(errors.InternalServer, c))
	// 		return
	// 	}
	// }

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
	req := &VerifyCodeReq{}
	err := c.BindJSON(req)
	if err != nil {
		c.JSON(http.StatusOK, respErrorCode(errors.InvalidParams, c))
		return
	}
	userInfo.Username = req.Username
	verifyType := req.Type
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

	// 滑块校验
	ser := factory.GetService(constant.BlockPuzzleCaptcha)
	err = ser.Check(req.Token, req.PointJSON)
	if err != nil {
		c.JSON(http.StatusOK, respErrorCode(errors.CaptchaError, c))
		return
	}

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

func BindWalletHandler(c *gin.Context) {
	claims := jwt.ExtractClaims(c)
	username := claims[identityKey].(string)

	type bindParams struct {
		VerifyCode string `json:"verify_code"`
		Sign       string `json:"sign"`
		PublicKey  string `json:"public_key"`
		Address    string `json:"address"`
	}

	var param bindParams
	if err := c.BindJSON(&param); err != nil {
		c.JSON(http.StatusOK, respErrorCode(errors.InvalidParams, c))
		return
	}

	nonce, err := getNonceFromCache(c.Request.Context(), param.Address, NonceStringTypeSignature)
	if err != nil {
		log.Errorf("query nonce string: %v", err)
		c.JSON(http.StatusOK, respErrorCode(errors.InternalServer, c))
		return
	}

	if nonce == "" {
		c.JSON(http.StatusOK, respErrorCode(errors.VerifyCodeExpired, c))
		return
	}

	success, err := VerifyCosmosAddr(param.Address, param.PublicKey, titanWalletPrefix)
	if err != nil || !success {
		c.JSON(http.StatusOK, respErrorCode(errors.InvalidWalletAddress, c))
		return
	}

	isWalletBound, err := dao.IsWalletAddressExists(c.Request.Context(), param.Address)
	if err != nil {
		c.JSON(http.StatusOK, respErrorCode(errors.InternalServer, c))
		return
	}

	if isWalletBound {
		c.JSON(http.StatusOK, respErrorCode(errors.WalletBound, c))
		return
	}

	bytePubKey, err := hex.DecodeString(param.PublicKey)
	if err != nil {
		c.JSON(http.StatusOK, respErrorCode(errors.InvalidPublicKey, c))
		return
	}

	byteSignature, err := hex.DecodeString(param.Sign)
	if err != nil {
		c.JSON(http.StatusOK, respErrorCode(errors.InvalidSignature, c))
		return
	}

	pubKey := secp256k1.PubKey{Key: bytePubKey}

	success, err = VerifyArbitraryMsg(param.Address, nonce, byteSignature, pubKey)
	if err != nil || !success {
		c.JSON(http.StatusOK, respErrorCode(errors.InvalidSignature, c))
		return
	}

	user, err := dao.GetUserByUsername(c.Request.Context(), username)
	if err != nil || user == nil {
		c.JSON(http.StatusOK, respErrorCode(errors.UserNotFound, c))
		return
	}

	if user.WalletAddress != "" {
		c.JSON(http.StatusOK, respErrorCode(errors.WalletBound, c))
		return
	}

	if err := dao.UpdateUserWalletAddress(context.Background(), username, param.Address); err != nil {
		log.Errorf("update user wallet address: %v", err)
		c.JSON(http.StatusOK, respErrorCode(errors.InternalServer, c))
		return
	}

	err = completeMission(c.Request.Context(), username, MissionIdBindTitanWallet)
	if err != nil {
		log.Errorf("complete brows official website error: %v", err)
		c.JSON(http.StatusOK, respErrorCode(errors.InternalServer, c))
		return
	}

	c.JSON(http.StatusOK, respJSON(nil))
}

// GetBlockCaptcha 滑块验证
func GetBlockCaptcha(c *gin.Context) {
	data, _ := factory.GetService(constant.BlockPuzzleCaptcha).Get()
	//输出json结果给调用方
	c.JSON(200, data)
}
