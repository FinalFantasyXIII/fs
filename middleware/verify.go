package middleware

import (
	"github.com/gin-gonic/gin"
	"strings"
)

func CheckForbiddenPath(forbidden []string) gin.HandlerFunc {
	return func(c *gin.Context) {
		var isForbidden bool
		for _,v := range forbidden{
			isForbidden = strings.Contains(c.Request.URL.Path, v)
			if isForbidden {
				c.HTML(200,"pig.html",nil)
				c.Abort()
			}
		}
		c.Next()
	}
}
