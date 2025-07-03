package middlewares

import (
	"messenging_test/common"

	"github.com/gin-gonic/gin"
)

// LanguageMiddleware handles language code validation and setting
func LanguageMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		languageCode := c.GetHeader("language_code")
		if languageCode == "" || (common.LanguageCode(languageCode) != common.En && common.LanguageCode(languageCode) != common.Es) {
			c.Set("language_code", common.En)
		} else {
			c.Set("language_code", languageCode)
		}
		c.Next()
	}
}
