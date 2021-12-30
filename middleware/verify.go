package middleware

import (
	"github.com/gin-gonic/gin"
	"regexp"
)

func CheckForbiddenPath(forbidden []string) gin.HandlerFunc {
	return func(c *gin.Context) {
		for _,v := range forbidden{
			flag , _ := regexp.MatchString(v,c.Request.URL.Path)
			if flag {
				c.HTML(200,"pig.html",nil)
				c.Abort()
			}
		}
		c.Next()
	}
}
