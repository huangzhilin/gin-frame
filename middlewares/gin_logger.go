package middlewares

import (
	"bytes"
	"io/ioutil"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

// GinLogger 接收gin框架默认的日志
func GinLogger() gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()
		url := c.Request.URL
		path := url.Path // 请求路径

		body, _ := c.GetRawData()                                // gin请求body默认只能被解析一次，解析后重新设置Body
		c.Request.Body = ioutil.NopCloser(bytes.NewBuffer(body)) // 重新设置request的Body

		c.Next()

		if !strings.HasPrefix(path, "/swagger/") {
			zap.L().Info(path,
				zap.Int("status", c.Writer.Status()),
				zap.Any("request", map[string]string{
					"url":        c.Request.Method + " " + url.String(),
					"body":       string(body),
					"ip":         c.ClientIP(),
					"user-agent": c.Request.UserAgent(),
				}),
				zap.String("errors", c.Errors.ByType(gin.ErrorTypePrivate).String()),
				zap.Duration("cost", time.Since(start)),
			)
		}

	}
}
