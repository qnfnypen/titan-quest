package api

import (
	"context"
	"database/sql"
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"time"

	jwt "github.com/appleboy/gin-jwt/v2"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/gin-gonic/gin"
	"github.com/gnasnik/titan-quest/core/dao"
	"github.com/gnasnik/titan-quest/core/errors"
	"github.com/gnasnik/titan-quest/core/generated/model"
	"github.com/gnasnik/titan-quest/pkg/random"
	"github.com/rs/xid"
)

const (
	loginStatusFailure = iota
	loginStatusSuccess
)

type login struct {
	Username   string `form:"username" json:"username"`
	Password   string `form:"password" json:"password"`
	VerifyCode string `form:"verify_code" json:"verify_code"`
	Sign       string `form:"sign" json:"sign"`
	Address    string `form:"address" json:"address"`
	Code       string `form:"code" json:"code"` // 邀请码
}

type loginResponse struct {
	Token  string `json:"token"`
	Expire string `json:"expire"`
}

var identityKey = "id"

func jwtGinMiddleware(secretKey string) (*jwt.GinJWTMiddleware, error) {
	return jwt.New(&jwt.GinJWTMiddleware{
		Realm:             "User",
		Key:               []byte(secretKey),
		Timeout:           time.Hour * 8,
		MaxRefresh:        24 * time.Hour,
		IdentityKey:       identityKey,
		SendAuthorization: true,
		PayloadFunc: func(data interface{}) jwt.MapClaims {
			if v, ok := data.(*model.User); ok {
				return jwt.MapClaims{
					identityKey: v.Username,
				}
			}
			return jwt.MapClaims{}
		},
		IdentityHandler: func(c *gin.Context) interface{} {
			claims := jwt.ExtractClaims(c)
			return &model.User{
				Username: claims[identityKey].(string),
			}
		},
		LoginResponse: func(c *gin.Context, code int, token string, expire time.Time) {
			c.JSON(http.StatusOK, gin.H{
				"code": 0,
				"data": loginResponse{
					Token:  token,
					Expire: expire.Format(time.RFC3339),
				},
			})
		},
		LogoutResponse: func(c *gin.Context, code int) {
			c.JSON(http.StatusOK, gin.H{
				"code": 0,
			})
		},
		Authenticator: func(c *gin.Context) (interface{}, error) {
			var loginParams login
			if err := c.BindJSON(&loginParams); err != nil {
				return "", fmt.Errorf("invalid input params")
			}

			if loginParams.Address == "" && loginParams.Username != "" {
				loginParams.Address = loginParams.Username
			}

			var (
				user interface{}
				err  error
			)

			if loginParams.Sign != "" {
				user, err = loginBySignature(c, loginParams.Address, loginParams.Sign, loginParams.Code)
			} else if loginParams.VerifyCode != "" {
				user, err = loginByVerifyCode(c, loginParams.Username, loginParams.VerifyCode, loginParams.Code)
			} else {
				return nil, errors.New("invalid login params")
			}

			if err != nil {
				log.Errorf("login: %v", err)
				return nil, err
			}

			if err := completeConnectWalletMission(c.Request.Context(), loginParams.Address); err != nil {
				log.Errorf("completeConnectWalletMission: %v", err)
			}

			return user, nil
		},
		Unauthorized: func(c *gin.Context, code int, message string) {
			c.JSON(http.StatusUnauthorized, gin.H{
				"code":    code,
				"message": message,
				"success": false,
			})
		},
		// TokenLookup is a string in the form of "<source>:<name>" that is used
		// to extract token from the request.
		// Optional. Default value "header:Authorization".
		// Possible values:
		// - "header:<name>"
		// - "query:<name>"
		// - "cookie:<name>"
		// - "param:<name>"
		//TokenLookup: "header: Authorization, query: token, cookie: jwt",
		TokenLookup: "header: JwtAuthorization",
		// TokenLookup: "query:token",
		// TokenLookup: "cookie:token",

		// TokenHeadName is a string in the header. Default value is "Bearer"
		TokenHeadName: "Bearer",

		// TimeFunc provides the current time. You can override it to use another time value. This is useful for testing or if your server uses a different time zone than your tokens.
		TimeFunc: time.Now,

		RefreshResponse: func(c *gin.Context, code int, token string, t time.Time) {
			c.Next()
		},
	})
}

func loginBySignature(c *gin.Context, address, msg, inviteCode string) (interface{}, error) {
	nonce, err := getNonceFromCache(c.Request.Context(), address, NonceStringTypeSignature)
	if err != nil {
		return nil, errors.NewErrorCode(errors.InvalidParams, c)
	}
	if nonce == "" {
		return nil, errors.NewErrorCode(errors.VerifyCodeExpired, c)
	}
	recoverAddress, err := VerifyMessage(nonce, msg)
	if strings.ToUpper(recoverAddress) != strings.ToUpper(address) {
		return nil, errors.NewErrorCode(errors.PassWordNotAllowed, c)
	}
	err = AddUserInfo(c.Request.Context(), address, inviteCode)
	if err != nil {
		c.JSON(http.StatusOK, respErrorCode(errors.InternalServer, c))
	}
	return &model.User{Username: address, Role: 0}, nil
}

func loginByVerifyCode(c *gin.Context, username, inputCode, inviteCode string) (interface{}, error) {
	code, err := getNonceFromCache(c.Request.Context(), username, NonceStringTypeLogin)
	if err != nil {
		log.Errorf("get user by verify code: %v", err)
		return nil, errors.NewErrorCode(errors.InvalidParams, c)
	}

	if code == "" {
		return nil, errors.NewErrorCode(errors.VerifyCodeExpired, c)
	}

	if code != inputCode {
		return nil, errors.NewErrorCode(errors.InvalidVerifyCode, c)
	}

	// _, err = dao.GetUserByUsername(c.Request.Context(), username)
	// if err == sql.ErrNoRows {
	// 	user := &model.User{
	// 		Username:     username,
	// 		CreatedAt:    time.Now(),
	// 		UpdatedAt:    time.Now(),
	// 		ReferralCode: random.GenerateRandomString(6),
	// 	}
	// 	err = dao.CreateUser(c.Request.Context(), user)
	// 	if err != nil {
	// 		c.JSON(http.StatusOK, respErrorCode(errors.InternalServer, c))
	// 	}
	// }
	err = AddUserInfo(c.Request.Context(), username, inviteCode)
	if err != nil {
		c.JSON(http.StatusOK, respErrorCode(errors.InternalServer, c))
	}

	return &model.User{Username: username, Role: 0}, nil
}

func VerifyMessage(message string, signedMessage string) (string, error) {
	// Hash the unsigned message using EIP-191
	hashedMessage := []byte("\x19Ethereum Signed Message:\n" + strconv.Itoa(len(message)) + message)
	hash := crypto.Keccak256Hash(hashedMessage)
	// Get the bytes of the signed message
	decodedMessage := hexutil.MustDecode(signedMessage)
	// Handles cases where EIP-115 is not implemented (most wallets don't implement it)
	if decodedMessage[64] == 27 || decodedMessage[64] == 28 {
		decodedMessage[64] -= 27
	}
	// Recover a public key from the signed message
	sigPublicKeyECDSA, err := crypto.SigToPub(hash.Bytes(), decodedMessage)
	if sigPublicKeyECDSA == nil {
		log.Errorf("Could not get a public get from the message signature")
	}
	if err != nil {
		return "", err
	}

	return crypto.PubkeyToAddress(*sigPublicKeyECDSA).String(), nil
}

// AddUserInfo 新增用户信息
func AddUserInfo(ctx context.Context, username, code string) error {
	userExt := &model.UsersExt{
		Username:    username,
		InviteCode:  xid.New().String(),
		InvitedCode: code,
	}

	// 判断用户是否存在，不存在则增加
	_, err := dao.GetUserByUsername(ctx, username)
	switch err {
	case sql.ErrNoRows:
		user := &model.User{
			Username:     username,
			CreatedAt:    time.Now(),
			UpdatedAt:    time.Now(),
			ReferralCode: random.GenerateRandomString(6),
		}

		err = dao.CreateUserInfo(ctx, user, userExt)
		if err != nil {
			return err
		}
	case nil:
	default:
		return err
	}

	// 判断用户附属表是否存在
	_, err = dao.GetUserExt(ctx, username)
	switch err {
	case sql.ErrNoRows:
		err = dao.CreateUserExt(ctx, userExt)
		if err != nil {
			return err
		}
	case nil:
	default:
		return err
	}

	return nil
}
