package middleware

import (
	"github.com/gin-gonic/gin"
	"github.com/masonhubco/rebar/v2"
)

// I18n returns a middleware that detects language from client side and use that for
// selecting i18n language files
func I18n() gin.HandlerFunc {
	return func(c *gin.Context) {
		accept := c.Request.Header.Get("Accept-Language")
		c.Set(rebar.I18nKey, rebar.WithLanguage(accept))
		c.Next()
	}
}
