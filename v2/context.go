package rebar

import (
	"github.com/gin-gonic/gin"
	"github.com/jmoiron/sqlx"
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

func LoggerFrom(c *gin.Context) Logger {
	if maybeALogger, exists := c.Get(LoggerKey); exists {
		if logger, ok := maybeALogger.(Logger); ok {
			return logger
		}
	}
	defaultLogger, _ := NewStandardLogger()
	return defaultLogger
}

func RequestIDFrom(c *gin.Context) string {
	return c.GetString(RequestIDKey)
}

func I18nFrom(c *gin.Context) (lang LanguageScoped, ok bool) {
	if maybeI18n, exists := c.Get(I18nKey); exists {
		lang, ok = maybeI18n.(LanguageScoped)
	}
	return
}

func I18nMustFrom(c *gin.Context) LanguageScoped {
	if lang, exists := I18nFrom(c); exists {
		return lang
	}
	panic(`"` + I18nKey + `" does not exist in context`)
}

func TxFrom(c *gin.Context) (tx *sqlx.Tx, ok bool) {
	if maybeTx, exists := c.Get(TxKey); exists {
		tx, ok = maybeTx.(*sqlx.Tx)
	}
	return
}

func TxMustFrom(c *gin.Context) *sqlx.Tx {
	if tx, exists := TxFrom(c); exists {
		return tx
	}
	panic(`"` + TxKey + `" does not exist in context`)
}
