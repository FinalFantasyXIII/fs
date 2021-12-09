package middleware

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"time"
)

func LogPrintter(ch chan<- map[string]string) gin.HandlerFunc{
	return func(c *gin.Context) {
		vs := make(map[string]string)
		vs["access_time"] = time.Now().Format("2006-01-02 15:04:05")
		vs["client_ip"] = c.Request.RemoteAddr
		vs["method"] = c.Request.Method
		vs["url"] = fmt.Sprintf("%#v",c.Request.URL.Path)
		ch <- vs
		c.Next()
	}
}