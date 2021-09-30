package rebar

import (
	"github.com/gin-gonic/gin"
)

const (
	I18nKey      = "i18n"
	LoggerKey    = "rebarLogger"
	RequestIDKey = "requestID"
	TxKey        = "tx"
)

type BuffaloValidateError interface {
	Error() string
	String() string
	HasAny() bool
}

func AbortWithError(c *gin.Context, code int, err error) *gin.Error {
	hashmap := gin.H{
		"request_id": RequestIDFrom(c),
	}
	if _, ok := err.(BuffaloValidateError); ok {
		hashmap["validate"] = err
	} else {
		hashmap["error"] = err.Error()
	}
	c.AbortWithStatusJSON(code, hashmap)
	return c.Error(err)
}
