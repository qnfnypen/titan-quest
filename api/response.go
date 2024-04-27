package api

import (
	"github.com/gin-gonic/gin"
	err "github.com/gnasnik/titan-quest/core/errors"
	"github.com/gnasnik/titan-quest/core/generated/model"
	"github.com/pkg/errors"
	"strings"
)

type JsonObject map[string]interface{}

func respJSON(v interface{}) gin.H {
	return gin.H{
		"success": true,
		"data":    v,
		"code":    0,
	}
}

func respError(e error) gin.H {
	var apiError err.ApiError
	if !errors.As(e, &apiError) {
		apiError = err.ErrUnknown
	}

	return gin.H{
		"success": false,
		"code":    apiError.Code(),
		"message": apiError.Error(),
	}
}

func respErrorCode(code int, c *gin.Context) gin.H {
	lang := c.GetHeader("Lang")

	var msg string

	messages := strings.Split(err.ErrMap[code], ":")
	if len(messages) == 0 {
		msg = err.ErrMap[code]
	} else {
		if lang == model.LanguageCN {
			msg = messages[1]
		} else {
			msg = messages[0]
		}
	}

	return gin.H{
		"code": -1,
		"err":  code,
		"msg":  msg,
	}
}
